package cqrs_dig_test

import (
	"context"
	"fmt"
	"testing"

	cqrs_dig "github.com/mitz-it/golang-cqrs-dig"

	cqrs "github.com/mitz-it/golang-cqrs"
	"github.com/stretchr/testify/assert"
	"go.uber.org/dig"
)

type PongService struct {
}

func NewPingService() *PongService {
	return &PongService{}
}

func (s *PongService) SetPong(ping string, pong string) string {
	return fmt.Sprintf("%s %s", ping, pong)
}

type PingResponse struct {
	Pong string
}

type PingCommand struct {
	Ping string
}

type PingQuery struct {
	Ping string
}

type PingEvent struct {
	Ping string
}

type CommandHandler struct {
	service *PongService
}

func (h *CommandHandler) Handle(ctx context.Context, command *PingCommand) (*PingResponse, error) {
	pong := h.service.SetPong(command.Ping, "pong")

	response := &PingResponse{
		Pong: pong,
	}

	return response, nil
}

func NewCommandHandler(service *PongService) cqrs.ICommandHandler[*PingCommand, *PingResponse] {
	return &CommandHandler{
		service: service,
	}
}

type QueryHandler struct {
	service *PongService
}

func (h *QueryHandler) Handle(ctx context.Context, query *PingQuery) (*PingResponse, error) {
	pong := h.service.SetPong(query.Ping, "pong")

	response := &PingResponse{
		Pong: pong,
	}

	return response, nil
}

func NewQueryHandler(service *PongService) cqrs.IQueryHandler[*PingQuery, *PingResponse] {
	return &QueryHandler{
		service: service,
	}
}

type PingBehavior struct {
	service *PongService
}

func (b *PingBehavior) Handle(ctx context.Context, request interface{}, next cqrs.NextFunc) (interface{}, error) {
	res, err := next()

	if err != nil {
		return nil, err
	}

	response := res.(*PingResponse)

	response.Pong = b.service.SetPong(response.Pong, "behavior also says pong")

	return response, nil
}

func NewPingBehavior(service *PongService) *PingBehavior {
	return &PingBehavior{
		service: service,
	}
}

type PingEventHandler struct {
	service *PongService
}

func (h *PingEventHandler) Handle(ctx context.Context, event *PingEvent) error {
	event.Ping = h.service.SetPong(event.Ping, "and event says pong!")
	return nil
}

func NewPingEventHandler(service *PongService) cqrs.IEvenHandler[*PingEvent] {
	return &PingEventHandler{
		service: service,
	}
}

func Test_ProvideCommandHandler_WhenHasInjectedService_ShouldInvokeAllDependencies(t *testing.T) {
	// arrange
	container := dig.New()
	container.Provide(NewPingService)

	command := &PingCommand{
		Ping: "ping",
	}

	// act
	cqrs_dig.ProvideCommandHandler[*PingCommand, *PingResponse](container, NewCommandHandler)

	response, _ := cqrs.Send[*PingCommand, *PingResponse](context.TODO(), command)

	// assert
	assert.Equal(t, "ping pong", response.Pong)
}

func Test_ProvideCommandHandler_WhenProvideCommandBehavior_ShouldInvokeAllDependencies(t *testing.T) {
	// arrange
	container := dig.New()
	container.Provide(NewPingService)
	cqrs_dig.ProvideCommandHandler[*PingCommand, *PingResponse](container, NewCommandHandler)

	command := &PingCommand{
		Ping: "ping",
	}

	// act
	err := cqrs_dig.ProvideCommandBehavior[*PingBehavior](container, 0, NewPingBehavior)
	response, _ := cqrs.Send[*PingCommand, *PingResponse](context.TODO(), command)

	// assert
	assert.Equal(t, "ping pong behavior also says pong", response.Pong)
	assert.Nil(t, err)
}

func Test_ProvideQueryHandler_WhenHasInjectedService_ShouldInvokeAllDependencies(t *testing.T) {
	// arrange
	container := dig.New()
	container.Provide(NewPingService)

	query := &PingQuery{
		Ping: "ping",
	}

	// act
	cqrs_dig.ProvideQueryHandler[*PingQuery, *PingResponse](container, NewQueryHandler)

	response, _ := cqrs.Request[*PingQuery, *PingResponse](context.TODO(), query)

	// assert
	assert.Equal(t, "ping pong", response.Pong)
}

func Test_ProvideQueryHandler_WhenProvideCommandBehavior_ShouldInvokeAllDependencies(t *testing.T) {
	// arrange
	container := dig.New()
	container.Provide(NewPingService)
	cqrs_dig.ProvideQueryHandler[*PingQuery, *PingResponse](container, NewQueryHandler)

	query := &PingQuery{
		Ping: "ping",
	}

	// act
	err := cqrs_dig.ProvideQueryBehavior[*PingBehavior](container, 0, NewPingBehavior)
	response, _ := cqrs.Request[*PingQuery, *PingResponse](context.TODO(), query)

	// assert
	assert.Equal(t, "ping pong behavior also says pong", response.Pong)
	assert.Nil(t, err)
}

func Test_ProvideEventHandler_WhenHasInjectedDependencies_ShouldInvokeAllDependencies(t *testing.T) {
	// arrange
	container := dig.New()
	container.Provide(NewPingService)

	event := &PingEvent{
		Ping: "ping",
	}

	// act
	cqrs_dig.ProvideEventHandler[*PingEvent](container, NewPingEventHandler)

	err := cqrs.PublishEvent[*PingEvent](context.TODO(), event)

	// assert
	assert.Equal(t, "ping and event says pong!", event.Ping)
	assert.Nil(t, err)
}
