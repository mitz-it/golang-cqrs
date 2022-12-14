package cqrs_behaviors

import (
	"context"
	"fmt"
	"strings"
	"time"

	cqrs_commands "github.com/mitz-it/golang-cqrs/commands"
	cqrs_queries "github.com/mitz-it/golang-cqrs/queries"
	logging "github.com/mitz-it/golang-logging"
)

type MetricsBehavior struct {
	Behavior
	logger *logging.Logger
}

func (behavior *MetricsBehavior) SetNextAction(next Action) {
	behavior.NextAction = next
}

func (behavior *MetricsBehavior) SetNextRequest(next Request) {
	behavior.NextRequest = next
}

func (behavior *MetricsBehavior) HandleCommand(ctx context.Context, command cqrs_commands.ICommand) (cqrs_commands.IResponse, error) {
	start := time.Now()

	actionName := generateActionName(command)

	defer func() {
		end := time.Since(start)

		behavior.logger.Standard.Info().Msgf("command %s duration: %d ms", actionName, int(end.Milliseconds()))
	}()

	return behavior.NextAction(ctx, command)
}

func (behavior *MetricsBehavior) HandleQuery(ctx context.Context, query cqrs_queries.IQuery) (cqrs_queries.IResponse, error) {
	start := time.Now()

	actionName := generateActionName(query)

	defer func() {
		end := time.Since(start)

		behavior.logger.Standard.Info().Msgf("query %s duration: %d ms", actionName, int(end.Milliseconds()))
	}()

	return behavior.NextRequest(ctx, query)
}

func generateActionName(handler any) string {
	return strings.Split(fmt.Sprintf("%T", handler), ".")[1]
}

func NewMetricsBehavior(logger *logging.Logger) IBehavior {
	metricsBehavior := &MetricsBehavior{}
	metricsBehavior.logger = logger
	return metricsBehavior
}
