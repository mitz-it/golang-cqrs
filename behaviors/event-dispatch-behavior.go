package behaviors

import (
	"context"

	"github.com/mitz-it/golang-cqrs/commands/v2"
	"github.com/mitz-it/golang-cqrs/queries/v2"
	logging "github.com/mitz-it/golang-logging"

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

func (behavior *EventDispatchBehavior) HandleCommand(ctx context.Context, command commands.ICommand) (commands.IResponse, error) {
	response, err := behavior.NextAction(ctx, command)
	behavior.logger.Standard.Info().Msgf("dispatching domain events")
	behavior.eventDispatcher.CommitDomainEventsStack(ctx)
	return response, err
}

func (behavior *EventDispatchBehavior) HandleQuery(ctx context.Context, query queries.IQuery) (queries.IResponse, error) {
	return behavior.NextRequest(ctx, query)
}

func NewEventDispatchBehavior(eventDispatcher events.IEventDispatcher, logger *logging.Logger) IBehavior {
	eventDispatchBehavior := &EventDispatchBehavior{}
	eventDispatchBehavior.eventDispatcher = eventDispatcher
	eventDispatchBehavior.logger = logger
	return eventDispatchBehavior
}
