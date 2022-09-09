package cqrs_behaviors

import (
	cqrs_commands "github.com/mitz-it/golang-cqrs/commands"

	cqrs_queries "github.com/mitz-it/golang-cqrs/queries"
)

type Action func(command cqrs_commands.ICommand) (cqrs_commands.IResponse, error)
type Request func(query cqrs_queries.IQuery) cqrs_queries.IResponse

type IBehavior interface {
	SetNext(next Action)
	SetNextRequest(next Request)
	Handle(command cqrs_commands.ICommand) (cqrs_commands.IResponse, error)
	HandleQuery(query cqrs_queries.IQuery) cqrs_queries.IResponse
}

type Behavior struct {
	Next        Action
	NextRequest Request
}
