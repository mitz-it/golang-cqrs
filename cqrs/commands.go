package cqrs

import (
	"context"
	"errors"
	"fmt"
	"reflect"

	"github.com/ahmetb/go-linq/v3"
)

type ICommandHandler[TCommand any, TResponse any] interface {
	Handle(ctx context.Context, command TCommand) (TResponse, error)
}

var commandHandlers map[reflect.Type]interface{}

func init() {
	commandHandlers = make(map[reflect.Type]interface{})
}

func RegisterCommandHandler[TCommand any, TResponse any](handler ICommandHandler[TCommand, TResponse]) error {
	var command TCommand
	commandType := reflect.TypeOf(command)

	_, found := commandHandlers[commandType]

	if found {
		msg := fmt.Sprintf("handler for command of type %s is already registered", commandType.String())
		return errors.New(msg)
	}

	commandHandlers[commandType] = handler

	return nil
}

func Send[TCommand any, TResponse any](ctx context.Context, command TCommand) (TResponse, error) {
	commandType := reflect.TypeOf(command)

	h, found := commandHandlers[commandType]

	if !found {
		msg := fmt.Sprintf("no handler registered for command %T", command)
		return *new(TResponse), errors.New(msg)
	}

	handler, casted := h.(ICommandHandler[TCommand, TResponse])

	if !casted {
		msg := fmt.Sprintf("handler of type %T is not assignable for command of type %T and response of type %T", handler, command, *new(TResponse))
		return *new(TResponse), errors.New(msg)
	}

	if len(commandBehaviors) <= 0 {
		return handler.Handle(ctx, command)
	}

	sortedBehaviors := sortBehaviors(commandBehaviors)

	commandHandle := func() (interface{}, error) {
		return handler.Handle(ctx, command)
	}

	aggregatedPipeline := linq.From(sortedBehaviors).AggregateWithSeedT(commandHandle, func(next NextFunc, b IBehavior) NextFunc {
		var nextFunc NextFunc = func() (interface{}, error) {
			return b.Handle(ctx, command, next)
		}
		return nextFunc
	})

	pipeline := aggregatedPipeline.(NextFunc)

	res, err := pipeline()

	response, casted := res.(TResponse)

	if !casted {
		return *new(TResponse), err
	}

	return response, err
}
