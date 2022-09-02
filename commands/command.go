package cqrs_commands

type ICommand interface {
}
type IResponse interface {
}

type ICommandHandler interface {
	Handle(command ICommand) (IResponse, error)
}
