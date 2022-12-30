package behaviors

import (
	"context"

	"github.com/mitz-it/golang-cqrs/commands/v2"
	"github.com/mitz-it/golang-cqrs/queries/v2"
)

type Action func(ctx context.Context, command commands.ICommand) (commands.IResponse, error)
type Request func(ctx context.Context, query queries.IQuery) (queries.IResponse, error)

type IBehavior interface {
	SetNextAction(next Action)
	SetNextRequest(next Request)
	HandleCommand(ctx context.Context, command commands.ICommand) (commands.IResponse, error)
	HandleQuery(ctx context.Context, query queries.IQuery) (queries.IResponse, error)
}

type Behavior struct {
	NextAction  Action
	NextRequest Request
}
