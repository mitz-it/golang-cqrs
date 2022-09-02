package cqrs_behaviors

import (
	cqrs_commands "gitlab.internal.cloud.payly.com.br/microservices/chassis/cqrs/commands"
	cqrs_queries "gitlab.internal.cloud.payly.com.br/microservices/chassis/cqrs/queries"

	"gitlab.internal.cloud.payly.com.br/microservices/chassis/logging"
)

type LoggingBehavior struct {
	Behavior
	logger *logging.Logger
}

func (behavior *LoggingBehavior) SetNext(next Action) {
	behavior.Next = next
}

func (behavior *LoggingBehavior) SetNextRequest(next Request) {
	behavior.NextRequest = next
}

func (behavior *LoggingBehavior) Handle(command cqrs_commands.ICommand) (cqrs_commands.IResponse, error) {
	behavior.logger.Standard.Info().Interface("serialized-command", command)
	defer behavior.logger.Standard.Info().Msgf("Log end")
	return behavior.Next(command)
}

func (behavior *LoggingBehavior) HandleQuery(query cqrs_queries.IQuery) cqrs_queries.IResponse {
	behavior.logger.Standard.Info().Interface("serialized-query", query)
	response := behavior.NextRequest(query)
	defer behavior.logger.Standard.Info().Interface("query-return", response)
	return response
}

func NewLoggingBehavior(logger *logging.Logger) IBehavior {
	loggingBehavior := &LoggingBehavior{}
	loggingBehavior.logger = logger
	return loggingBehavior
}
