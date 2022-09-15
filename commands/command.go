package cqrs_commands

type ICommand interface {
}
type IResponse interface {
}

type ICommandHandler interface {
	HandleCommand(command ICommand) (IResponse, error)
}
