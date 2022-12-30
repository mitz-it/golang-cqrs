package cqrs

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"go.uber.org/dig"
	"golang.org/x/exp/slices"

	"github.com/mitz-it/golang-cqrs/behaviors/v2"
	"github.com/mitz-it/golang-cqrs/commands/v2"
	"github.com/mitz-it/golang-cqrs/queries/v2"
)

type MediatorParams struct {
	dig.In

	CommandBehaviors []behaviors.IBehavior      `group:"CommandBehaviors"`
	QueryBehaviors   []behaviors.IBehavior      `group:"QueryBehaviors"`
	Handlers         []commands.ICommandHandler `group:"CommandHandlers"`
	QueryHandlers    []queries.IQueryHandler    `group:"QueryHandlers"`
}

type IMediator interface {
	Send(ctx context.Context, command commands.ICommand) (commands.IResponse, error)
	Request(ctx context.Context, query queries.IQuery) (queries.IResponse, error)
}

type Mediator struct {
	commandBehaviors []behaviors.IBehavior
	queryBehaviors   []behaviors.IBehavior
	handlers         []commands.ICommandHandler
	queryHandlers    []queries.IQueryHandler
}

func (mediator Mediator) Send(ctx context.Context, command commands.ICommand) (commands.IResponse, error) {
	position := slices.IndexFunc(mediator.handlers, func(handler commands.ICommandHandler) bool {
		handlerName := fmt.Sprintf("%T", handler)
		commandName := fmt.Sprintf("%T", command)
		return strings.Contains(handlerName, commandName)
	})

	if position <= -1 {
		message := fmt.Sprintf("command handler not found for command of type: %T", command)
		err := errors.New(message)
		panic(err)
	}

	handler := mediator.handlers[position]

	if len(mediator.commandBehaviors) <= 0 {
		return handler.HandleCommand(ctx, command)
	}

	mediator.commandBehaviors[len(mediator.commandBehaviors)-1].SetNextAction(handler.HandleCommand)

	return mediator.commandBehaviors[0].HandleCommand(ctx, command)
}

func (mediator Mediator) Request(ctx context.Context, query queries.IQuery) (queries.IResponse, error) {
	position := slices.IndexFunc(mediator.queryHandlers, func(handler queries.IQueryHandler) bool {
		handlerName := fmt.Sprintf("%T", handler)
		commandName := fmt.Sprintf("%T", query)
		return strings.Contains(handlerName, commandName)
	})

	if position <= -1 {
		message := fmt.Sprintf("query handler not found for query of type: %T", query)
		err := errors.New(message)
		panic(err)
	}

	handler := mediator.queryHandlers[position]

	if len(mediator.queryBehaviors) <= 0 {
		return handler.HandleQuery(ctx, query)
	}

	mediator.queryBehaviors[len(mediator.queryBehaviors)-1].SetNextRequest(handler.HandleQuery)

	return mediator.queryBehaviors[0].HandleQuery(ctx, query)
}

func sortCommandBehaviors(behaviors []behaviors.IBehavior) []behaviors.IBehavior {
	firstBehaviorName := fmt.Sprintf("%T", behaviors[0])
	lastBehaviorName := fmt.Sprintf("%T", behaviors[len(behaviors)-1])
	validationBehaviorName := fmt.Sprintf("%T", &behaviors.ValidationBehavior{})
	eventDispatchBehaviorName := fmt.Sprintf("%T", &behaviors.EventDispatchBehavior{})

	if firstBehaviorName != validationBehaviorName {
		validationBehaviorPosition := slices.IndexFunc(behaviors, func(behavior behaviors.IBehavior) bool {
			behaviorName := fmt.Sprintf("%T", behavior)
			return strings.Contains(behaviorName, validationBehaviorName)
		})

		firstBehavior := behaviors[0]
		behaviors[0] = behaviors[validationBehaviorPosition]
		behaviors[validationBehaviorPosition] = firstBehavior
	}

	if lastBehaviorName != eventDispatchBehaviorName {
		eventDispatchBehaviorPosition := slices.IndexFunc(behaviors, func(behavior behaviors.IBehavior) bool {
			behaviorName := fmt.Sprintf("%T", behavior)
			return strings.Contains(behaviorName, eventDispatchBehaviorName)
		})

		lastBehavior := behaviors[len(behaviors)-1]
		behaviors[len(behaviors)-1] = behaviors[eventDispatchBehaviorPosition]
		behaviors[eventDispatchBehaviorPosition] = lastBehavior

	}

	return behaviors
}

func configureQueryBehaviors(params MediatorParams) {
	for index, behavior := range params.QueryBehaviors {
		if index < len(params.QueryBehaviors)-1 {
			behavior.SetNextRequest(params.QueryBehaviors[index+1].HandleQuery)
		} else {
			behavior.SetNextRequest(nil)
		}
	}
}

func configureCommandBehaviors(params MediatorParams) {
	if len(params.CommandBehaviors) <= 0 {
		return
	}

	for index, behavior := range sortCommandBehaviors(params.CommandBehaviors) {
		if index < len(params.CommandBehaviors)-1 {
			behavior.SetNextAction(params.CommandBehaviors[index+1].HandleCommand)
		} else {
			behavior.SetNextAction(nil)
		}
	}
}

func NewMediator(params MediatorParams) IMediator {
	configureCommandBehaviors(params)

	configureQueryBehaviors(params)

	return Mediator{
		commandBehaviors: params.CommandBehaviors,
		queryBehaviors:   params.QueryBehaviors,
		handlers:         params.Handlers,
		queryHandlers:    params.QueryHandlers,
	}
}
