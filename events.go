package cqrs

import (
	"context"
	"errors"
	"fmt"
	"reflect"

	"go.uber.org/multierr"
)

type IEventHandler[TEvent any] interface {
	Handle(ctx context.Context, event TEvent) error
}

type EventDelivery struct {
	ctx       context.Context
	eventType reflect.Type
	event     interface{}
}

var eventHandlers map[reflect.Type][]interface{}
var eventListener chan *EventDelivery

func init() {
	eventHandlers = make(map[reflect.Type][]interface{})
	eventListener = make(chan *EventDelivery)
}

func RegisterEventSubscriber[TEvent any](handler IEventHandler[TEvent]) error {
	var event TEvent
	eventType := reflect.TypeOf(event)
	handlers, found := eventHandlers[eventType]

	if !found {
		eventHandlers[eventType] = []interface{}{
			handler,
		}
		return nil
	}

	eventHandlers[eventType] = append(handlers, handler)

	return nil
}

func RegisterEventSubscribers[TEvent any](handlers ...IEventHandler[TEvent]) error {
	if len(handlers) <= 0 {
		return errors.New("at least one handler must be provided")
	}

	for _, handler := range handlers {
		RegisterEventSubscriber(handler)
	}

	return nil
}

func PublishEvent[TEvent any](ctx context.Context, event TEvent) error {
	eventType := reflect.TypeOf(event)
	handlers, found := eventHandlers[eventType]

	if !found {
		msg := fmt.Sprintf("no event handler found event of type: %T", event)
		return errors.New(msg)
	}

	var err error = nil

	for _, h := range handlers {
		handler, ok := h.(IEventHandler[TEvent])

		if ok {
			handleErr := handler.Handle(ctx, event)

			if handleErr != nil {
				err = multierr.Append(err, handleErr)
			}

			continue
		}

		args := []reflect.Value{
			reflect.ValueOf(ctx),
			reflect.ValueOf(event),
		}

		r := reflect.ValueOf(h).MethodByName("Handle").Call(args)

		handleErr := r[0].Interface()

		if handleErr != nil {
			err = multierr.Append(err, handleErr.(error))
		}
	}

	return err
}

func PublishEventAsync[TEvent any](ctx context.Context, event TEvent) error {
	eventType := reflect.TypeOf(event)

	delivery := &EventDelivery{
		ctx:       ctx,
		eventType: eventType,
		event:     event,
	}

	eventListener <- delivery

	return nil
}

func Listen() {
	go listen()
}

func handleRecover() {
	if err := recover(); err != nil {
		go listen()
	}
}

func listen() {
	defer handleRecover()
	for delivery := range eventListener {
		event := delivery.event
		eventType := delivery.eventType
		handlers, ok := eventHandlers[eventType]

		if !ok {
			return
		}

		args := []reflect.Value{
			reflect.ValueOf(delivery.ctx),
			reflect.ValueOf(event),
		}

		for _, handler := range handlers {
			reflect.ValueOf(handler).MethodByName("Handle").Call(args)
		}
	}
}
