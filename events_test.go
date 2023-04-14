package cqrs

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

type FakeEvent struct {
	Message string
}

type FakeEventHandler1 struct {
}

func (h *FakeEventHandler1) Handle(ctx context.Context, event *FakeEvent) error {
	return nil
}

type FakeEventHandler2 struct {
}

func (h *FakeEventHandler2) Handle(ctx context.Context, event *FakeEvent) error {
	return errors.New("something unexpected happened")
}

type FakeEventHandler3 struct {
}

func (h *FakeEventHandler3) Handle(ctx context.Context, event *FakeEvent) error {
	err := errors.New("something unexpected happened")
	panic(err)
}

func events_cleanup(t *testing.T) {
	t.Cleanup(func() {
		eventHandlers = make(map[reflect.Type][]interface{})
		eventListener = make(chan *EventDelivery)
	})
}

func TestRegisterEventSubscriber_WhenFirstHandler_ShouldAddHandlerToMap(t *testing.T) {
	// arrange
	defer events_cleanup(t)
	var event *FakeEvent
	eventType := reflect.TypeOf(event)
	handler := &FakeEventHandler1{}

	// act
	err := RegisterEventSubscriber[*FakeEvent](handler)

	// assert
	assert.Nil(t, err)
	assert.Contains(t, eventHandlers[eventType], handler)
	assert.Len(t, eventHandlers[eventType], 1)

}

func TestRegisterEventSubscriber_WhenNotFirstHandler_ShouldAddHandlersToMap(t *testing.T) {
	// arrange
	defer events_cleanup(t)
	var event *FakeEvent
	eventType := reflect.TypeOf(event)
	handler1 := &FakeEventHandler1{}
	handler2 := &FakeEventHandler2{}

	// act
	RegisterEventSubscriber[*FakeEvent](handler1)
	err := RegisterEventSubscriber[*FakeEvent](handler2)

	// assert
	assert.Nil(t, err)
	assert.Contains(t, eventHandlers[eventType], handler1)
	assert.Contains(t, eventHandlers[eventType], handler2)
	assert.Len(t, eventHandlers[eventType], 2)
}

func TestRegisterEventSubscribers_WhenMultipleHadlers_ShouldAddHandlersToMap(t *testing.T) {
	// arrange
	defer events_cleanup(t)
	var event *FakeEvent
	eventType := reflect.TypeOf(event)
	defer t.Cleanup(func() {

	})
	handler1 := &FakeEventHandler1{}
	handler2 := &FakeEventHandler2{}

	// act
	err := RegisterEventSubscribers[*FakeEvent](handler1, handler2)

	// assert
	assert.Nil(t, err)
	assert.Contains(t, eventHandlers[eventType], handler1)
	assert.Contains(t, eventHandlers[eventType], handler2)
	assert.Len(t, eventHandlers[eventType], 2)
}

func TestRegisterEventSubscribers_WhenNotAnyHandler_ShouldReturnError(t *testing.T) {
	// arrange
	defer events_cleanup(t)
	handlers := []IEventHandler[*FakeEvent]{}
	// act
	err := RegisterEventSubscribers(handlers...)

	// assert
	assert.NotNil(t, err)
}

func TestPublishEvent_WhenEventHandlerNotFound_ShouldReturn(t *testing.T) {
	// arrange
	defer events_cleanup(t)
	event := FakeEvent{
		Message: "test",
	}
	// act
	err := PublishEvent(context.TODO(), event)

	// assert
	assert.NotNil(t, err)
}

func TestPublishEvent_WhenEvenHandlerCanBeCastedAndHandle_ShouldNotReturnError(t *testing.T) {
	// arrange
	defer events_cleanup(t)
	event := &FakeEvent{
		Message: "test",
	}
	handler := &FakeEventHandler1{}
	RegisterEventSubscriber[*FakeEvent](handler)

	// act
	err := PublishEvent(context.TODO(), event)
	// assert
	assert.Nil(t, err)
}

func TestPublishEvent_WhenEvenHandlerCanBeCastedAndNotHandle_ShouldReturnError(t *testing.T) {
	// arrange
	defer events_cleanup(t)
	event := &FakeEvent{
		Message: "test",
	}
	handler := &FakeEventHandler2{}
	RegisterEventSubscriber[*FakeEvent](handler)

	// act
	err := PublishEvent(context.TODO(), event)

	// assert
	assert.NotNil(t, err)
}

func TestPublishEvent_WhenEvenHandlerCantBeCastedAndHandle_ShouldNotReturnError(t *testing.T) {
	// arrange
	defer events_cleanup(t)
	var event interface{} = &FakeEvent{
		Message: "test",
	}
	handler := &FakeEventHandler1{}

	RegisterEventSubscriber[*FakeEvent](handler)

	// act
	err := PublishEvent(context.TODO(), event)

	// assert
	assert.Nil(t, err)
}

func TestPublishEvent_WhenEvenHandlerCantBeCastedAndNotHandle_ShouldReturnError(t *testing.T) {
	defer events_cleanup(t)
	var event interface{} = &FakeEvent{
		Message: "test",
	}
	handler := &FakeEventHandler2{}

	RegisterEventSubscriber[*FakeEvent](handler)

	// act
	err := PublishEvent(context.TODO(), event)

	// assert
	assert.NotNil(t, err)
}

func TestPublishEventAsync_WhenListening_ShouldHandleEvent(t *testing.T) {
	// arrange
	defer events_cleanup(t)
	event := &FakeEvent{
		Message: "test",
	}
	handler := &FakeEventHandler1{}
	RegisterEventSubscriber[*FakeEvent](handler)
	Listen()

	// act
	publish := func() {
		PublishEventAsync(context.TODO(), event)
	}

	// assert
	assert.NotPanics(t, publish)
}

func TestPublishEventAsync_WhenHandlerNotFound_ShouldReturnEarlyOnDelivery(t *testing.T) {
	// arrange
	defer events_cleanup(t)
	event := struct {
		Message string
	}{
		Message: "test",
	}
	Listen()

	// act
	publish := func() {
		PublishEventAsync(context.TODO(), event)
	}

	// assert
	assert.NotPanics(t, publish)
}

func TestPublishEventAsync_WhenHandlerPanic_ShouldRecover(t *testing.T) {
	// arrange
	defer events_cleanup(t)
	event := &FakeEvent{
		Message: "test",
	}
	handler := &FakeEventHandler3{}
	RegisterEventSubscriber[*FakeEvent](handler)
	Listen()

	// act
	publish := func() {
		PublishEventAsync(context.TODO(), event)
	}

	// assert
	assert.NotPanics(t, publish)
}
