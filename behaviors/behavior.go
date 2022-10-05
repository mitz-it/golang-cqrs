package cqrs_behaviors

import (
	"context"

	commands "github.com/mitz-it/golang-cqrs/commands"

	queries "github.com/mitz-it/golang-cqrs/queries"
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
