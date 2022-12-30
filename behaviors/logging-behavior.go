package behaviors

import (
	"context"

	"github.com/mitz-it/golang-cqrs/commands"
	"github.com/mitz-it/golang-cqrs/queries"
	logging "github.com/mitz-it/golang-logging"
)

type LoggingBehavior struct {
	Behavior
	logger *logging.Logger
}

func (behavior *LoggingBehavior) SetNextAction(next Action) {
	behavior.NextAction = next
}

func (behavior *LoggingBehavior) SetNextRequest(next Request) {
	behavior.NextRequest = next
}

func (behavior *LoggingBehavior) HandleCommand(ctx context.Context, command commands.ICommand) (commands.IResponse, error) {
	behavior.logger.Standard.Info().Interface("serialized-command", command)
	response, err := behavior.NextAction(ctx, command)
	defer behavior.logger.Standard.Info().Interface("command-return", response)
	return response, err
}

func (behavior *LoggingBehavior) HandleQuery(ctx context.Context, query queries.IQuery) (queries.IResponse, error) {
	behavior.logger.Standard.Info().Interface("serialized-query", query)
	response, err := behavior.NextRequest(ctx, query)
	defer behavior.logger.Standard.Info().Interface("query-return", response)
	return response, err
}

func NewLoggingBehavior(logger *logging.Logger) IBehavior {
	loggingBehavior := &LoggingBehavior{}
	loggingBehavior.logger = logger
	return loggingBehavior
}
