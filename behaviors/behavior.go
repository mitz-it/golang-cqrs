package cqrs_behaviors

import (
	commands "github.com/mitz-it/golang-cqrs/commands"

	queries "github.com/mitz-it/golang-cqrs/queries"
)

type Action func(command commands.ICommand) (commands.IResponse, error)
type Request func(query queries.IQuery) (queries.IResponse, error)

type IBehavior interface {
	SetNext(next Action)
	SetNextRequest(next Request)
	Handle(command commands.ICommand) (commands.IResponse, error)
	HandleQuery(query queries.IQuery) (queries.IResponse, error)
}

type Behavior struct {
	Next        Action
	NextRequest Request
}
