package cqrs

import (
	"context"
	"errors"
	"fmt"
	"reflect"

	"github.com/ahmetb/go-linq/v3"
)

type IQueryHandler[TQuery any, TResponse any] interface {
	Handle(ctx context.Context, query TQuery) (TResponse, error)
}

var queryHandlers map[reflect.Type]interface{}

func init() {
	queryHandlers = make(map[reflect.Type]interface{})
}

func RegisterQueryHandler[TQuery any, TResponse any](handler IQueryHandler[TQuery, TResponse]) error {
	var query TQuery
	queryType := reflect.TypeOf(query)

	_, found := queryHandlers[queryType]

	if found {
		msg := fmt.Sprintf("handler for query of type %s is already registered", queryType.String())
		return errors.New(msg)
	}

	queryHandlers[queryType] = handler

	return nil
}

func Request[TQuery any, TResponse any](ctx context.Context, query TQuery) (TResponse, error) {
	queryType := reflect.TypeOf(query)

	h, found := queryHandlers[queryType]

	if !found {
		msg := fmt.Sprintf("no handler registered for query %T", query)
		return *new(TResponse), errors.New(msg)
	}

	handler, casted := h.(IQueryHandler[TQuery, TResponse])

	if !casted {
		msg := fmt.Sprintf("handler of type %T is not assignable for query of type %T and response of type %T", handler, query, *new(TResponse))
		return *new(TResponse), errors.New(msg)
	}

	if len(queryBehaviors) <= 0 {
		return handler.Handle(ctx, query)
	}

	sortedBehaviors := sortBehaviors(queryBehaviors)

	queryHandle := func() (interface{}, error) {
		return handler.Handle(ctx, query)
	}

	aggregatedPipeline := linq.From(sortedBehaviors).AggregateWithSeedT(queryHandle, func(next NextFunc, b IBehavior) NextFunc {
		var nextFunc NextFunc = func() (interface{}, error) {
			return b.Handle(ctx, query, next)
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
