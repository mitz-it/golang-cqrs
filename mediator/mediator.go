package cqrs

import (
	"context"
	"fmt"
	"strings"

	"go.uber.org/dig"
	"golang.org/x/exp/slices"

	cqrs_behaviors "github.com/mitz-it/golang-cqrs/behaviors"
	cqrs_commands "github.com/mitz-it/golang-cqrs/commands"
	cqrs_queries "github.com/mitz-it/golang-cqrs/queries"
)

type MediatorParams struct {
	dig.In

	CommandBehaviors []cqrs_behaviors.IBehavior      `group:"CommandBehaviors"`
	QueryBehaviors   []cqrs_behaviors.IBehavior      `group:"QueryBehaviors"`
	Handlers         []cqrs_commands.ICommandHandler `group:"CommandHandlers"`
	QueryHandlers    []cqrs_queries.IQueryHandler    `group:"QueryHandlers"`
}

type IMediator interface {
	Send(ctx context.Context, command cqrs_commands.ICommand) (cqrs_commands.IResponse, error)
	Request(ctx context.Context, query cqrs_queries.IQuery) (cqrs_queries.IResponse, error)
}

type Mediator struct {
	commandBehaviors []cqrs_behaviors.IBehavior
	queryBehaviors   []cqrs_behaviors.IBehavior
	handlers         []cqrs_commands.ICommandHandler
	queryHandlers    []cqrs_queries.IQueryHandler
}

func (mediator Mediator) Send(ctx context.Context, command cqrs_commands.ICommand) (cqrs_commands.IResponse, error) {
	position := slices.IndexFunc(mediator.handlers, func(handler cqrs_commands.ICommandHandler) bool {
		handlerName := fmt.Sprintf("%T", handler)
		commandName := fmt.Sprintf("%T", command)
		return strings.Contains(handlerName, commandName)
	})

	handler := mediator.handlers[position]

	mediator.commandBehaviors[len(mediator.commandBehaviors)-1].SetNextAction(handler.HandleCommand)
	return mediator.commandBehaviors[0].HandleCommand(ctx, command)
}

func (mediator Mediator) Request(ctx context.Context, query cqrs_queries.IQuery) (cqrs_queries.IResponse, error) {
	position := slices.IndexFunc(mediator.queryHandlers, func(handler cqrs_queries.IQueryHandler) bool {
		handlerName := fmt.Sprintf("%T", handler)
		commandName := fmt.Sprintf("%T", query)
		return strings.Contains(handlerName, commandName)
	})

	handler := mediator.queryHandlers[position]

	mediator.queryBehaviors[len(mediator.queryBehaviors)-1].SetNextRequest(handler.HandleQuery)

	return mediator.queryBehaviors[0].HandleQuery(ctx, query)
}

func sortCommandBehaviors(behaviors []cqrs_behaviors.IBehavior) []cqrs_behaviors.IBehavior {
	firstBehaviorName := fmt.Sprintf("%T", behaviors[0])
	lastBehaviorName := fmt.Sprintf("%T", behaviors[len(behaviors)-1])
	validationBehaviorName := fmt.Sprintf("%T", &cqrs_behaviors.ValidationBehavior{})
	eventDispatchBehaviorName := fmt.Sprintf("%T", &cqrs_behaviors.EventDispatchBehavior{})

	if firstBehaviorName != validationBehaviorName {
		validationBehaviorPosition := slices.IndexFunc(behaviors, func(behavior cqrs_behaviors.IBehavior) bool {
			behaviorName := fmt.Sprintf("%T", behavior)
			return strings.Contains(behaviorName, validationBehaviorName)
		})

		firstBehavior := behaviors[0]
		behaviors[0] = behaviors[validationBehaviorPosition]
		behaviors[validationBehaviorPosition] = firstBehavior
	}

	if lastBehaviorName != eventDispatchBehaviorName {
		eventDispatchBehaviorPosition := slices.IndexFunc(behaviors, func(behavior cqrs_behaviors.IBehavior) bool {
			behaviorName := fmt.Sprintf("%T", behavior)
			return strings.Contains(behaviorName, eventDispatchBehaviorName)
		})

		lastBehavior := behaviors[len(behaviors)-1]
		behaviors[len(behaviors)-1] = behaviors[eventDispatchBehaviorPosition]
		behaviors[eventDispatchBehaviorPosition] = lastBehavior

	}

	return behaviors
}

func NewMediator(params MediatorParams) IMediator {
	for index, behavior := range sortCommandBehaviors(params.CommandBehaviors) {
		if index < len(params.CommandBehaviors)-1 {
			behavior.SetNextAction(params.CommandBehaviors[index+1].HandleCommand)
		} else {
			behavior.SetNextAction(nil)
		}
	}

	for index, behavior := range params.QueryBehaviors {
		if index < len(params.QueryBehaviors)-1 {
			behavior.SetNextRequest(params.QueryBehaviors[index+1].HandleQuery)
		} else {
			behavior.SetNextRequest(nil)
		}
	}
	return Mediator{
		commandBehaviors: params.CommandBehaviors,
		queryBehaviors:   params.QueryBehaviors,
		handlers:         params.Handlers,
		queryHandlers:    params.QueryHandlers,
	}
}
