package cqrs_validators

import commands "github.com/mitz-it/golang-cqrs/commands"

type ICommandValidator interface {
	Validate(command commands.ICommand) error
}
