package behaviors

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/mitz-it/golang-cqrs/commands/v2"
	"github.com/mitz-it/golang-cqrs/queries/v2"
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

func (behavior *MetricsBehavior) HandleCommand(ctx context.Context, command commands.ICommand) (commands.IResponse, error) {
	start := time.Now()

	actionName := generateActionName(command)

	defer func() {
		end := time.Since(start)

		behavior.logger.Standard.Info().Msgf("command %s duration: %d ms", actionName, int(end.Milliseconds()))
	}()

	return behavior.NextAction(ctx, command)
}

func (behavior *MetricsBehavior) HandleQuery(ctx context.Context, query queries.IQuery) (queries.IResponse, error) {
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
