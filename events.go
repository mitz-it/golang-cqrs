package mediator

import (
	"context"
	"reflect"

	"go.uber.org/multierr"
)

type IEvenHandler[TEvent any] interface {
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

func RegisterEventSubcriber[TEvent any](handler IEvenHandler[TEvent]) error {
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

func PublishEvent[TEvent any](ctx context.Context, event TEvent) error {
	eventType := reflect.TypeOf(event)
	handlers, found := eventHandlers[eventType]

	if !found {
		return nil
	}

	var err error = nil

	for _, h := range handlers {
		handler, ok := h.(IEvenHandler[TEvent])

		if ok {
			handleErr := handler.Handle(ctx, event)

			if handleErr != nil {
				multierr.Append(err, handleErr)
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
			multierr.Append(err, handleErr.(error))
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
