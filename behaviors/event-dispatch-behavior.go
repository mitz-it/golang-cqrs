package cqrs_behaviors

import (
	cqrs_commands "gitlab.internal.cloud.payly.com.br/microservices/chassis/cqrs/commands"

	cqrs_queries "gitlab.internal.cloud.payly.com.br/microservices/chassis/cqrs/queries"

	events "gitlab.internal.cloud.payly.com.br/microservices/chassis/events"
	"gitlab.internal.cloud.payly.com.br/microservices/chassis/logging"
)

type EventDispatchBehavior struct {
	Behavior
	eventDispatcher events.IEventDispatcher
	logger          *logging.Logger
}

func (behavior *EventDispatchBehavior) SetNext(next Action) {
	behavior.Next = next
}

func (behavior *EventDispatchBehavior) SetNextRequest(next Request) {
	behavior.NextRequest = next
}

func (behavior *EventDispatchBehavior) Handle(command cqrs_commands.ICommand) (cqrs_commands.IResponse, error) {
	response, err := behavior.Next(command)
	behavior.logger.Standard.Info().Msgf("dispatching domain events")
	behavior.eventDispatcher.CommitDomainEventsStack()
	return response, err
}

func (behavior *EventDispatchBehavior) HandleQuery(query cqrs_queries.IQuery) cqrs_queries.IResponse {
	return behavior.NextRequest(query)
}

func NewEventDispatchBehavior(eventDispatcher events.IEventDispatcher, logger *logging.Logger) IBehavior {
	eventDispatchBehavior := &EventDispatchBehavior{}
	eventDispatchBehavior.eventDispatcher = eventDispatcher
	eventDispatchBehavior.logger = logger
	return eventDispatchBehavior
}
