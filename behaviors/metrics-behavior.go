package cqrs_behaviors

import (
	"fmt"
	"strings"
	"time"

	cqrs_commands "gitlab.internal.cloud.payly.com.br/microservices/chassis/cqrs/commands"
	cqrs_queries "gitlab.internal.cloud.payly.com.br/microservices/chassis/cqrs/queries"

	"gitlab.internal.cloud.payly.com.br/microservices/chassis/logging"
)

type MetricsBehavior struct {
	Behavior
	logger *logging.Logger
}

func (behavior *MetricsBehavior) SetNext(next Action) {
	behavior.Next = next
}

func (behavior *MetricsBehavior) SetNextRequest(next Request) {
	behavior.NextRequest = next
}

func (behavior *MetricsBehavior) Handle(command cqrs_commands.ICommand) (cqrs_commands.IResponse, error) {
	start := time.Now()

	actionName := generateActionName(command)

	defer func() {
		end := time.Since(start)

		behavior.logger.Standard.Info().Msgf("command %s duration: %d ms", actionName, int(end.Milliseconds()))
	}()

	return behavior.Next(command)
}

func (behavior *MetricsBehavior) HandleQuery(query cqrs_queries.IQuery) cqrs_queries.IResponse {
	start := time.Now()

	actionName := generateActionName(query)

	defer func() {
		end := time.Since(start)

		behavior.logger.Standard.Info().Msgf("query %s duration: %d ms", actionName, int(end.Milliseconds()))
	}()

	return behavior.NextRequest(query)
}

func generateActionName(handler any) string {
	return strings.Split(fmt.Sprintf("%T", handler), ".")[1]
}

func NewMetricsBehavior(logger *logging.Logger) IBehavior {
	metricsBehavior := &MetricsBehavior{}
	metricsBehavior.logger = logger
	return metricsBehavior
}
