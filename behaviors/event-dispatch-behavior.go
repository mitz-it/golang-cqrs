package cqrs_behaviors

import (
	cqrs_commands "github.com/mitz-it/golang-cqrs/commands"
	logging "github.com/mitz-it/golang-logging"

	cqrs_queries "github.com/mitz-it/golang-cqrs/queries"

	events "github.com/mitz-it/golang-events"
)

type EventDispatchBehavior struct {
	Behavior
	eventDispatcher events.IEventDispatcher
	logger          *logging.Logger
}

func (behavior *EventDispatchBehavior) SetNextAction(next Action) {
	behavior.NextAction = next
}

func (behavior *EventDispatchBehavior) SetNextRequest(next Request) {
	behavior.NextRequest = next
}

func (behavior *EventDispatchBehavior) HandleCommand(command cqrs_commands.ICommand) (cqrs_commands.IResponse, error) {
	response, err := behavior.NextAction(command)
	behavior.logger.Standard.Info().Msgf("dispatching domain events")
	behavior.eventDispatcher.CommitDomainEventsStack()
	return response, err
}

func (behavior *EventDispatchBehavior) HandleQuery(query cqrs_queries.IQuery) (cqrs_queries.IResponse, error) {
	return behavior.NextRequest(query)
}

func NewEventDispatchBehavior(eventDispatcher events.IEventDispatcher, logger *logging.Logger) IBehavior {
	eventDispatchBehavior := &EventDispatchBehavior{}
	eventDispatchBehavior.eventDispatcher = eventDispatcher
	eventDispatchBehavior.logger = logger
	return eventDispatchBehavior
}
