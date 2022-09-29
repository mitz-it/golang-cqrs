package cqrs_commands

import "context"

type ICommand interface {
}
type IResponse interface {
}

type ICommandHandler interface {
	HandleCommand(ctx context.Context, command ICommand) (IResponse, error)
}
