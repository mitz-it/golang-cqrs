package cqrs_validators

import cqrs_commands "github.com/mitz-it/golang-cqrs/commands"

type ICommandValidator interface {
	Validate(command cqrs_commands.ICommand) error
}
