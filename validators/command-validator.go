package cqrs_validators

import cqrs_commands "gitlab.internal.cloud.payly.com.br/microservices/chassis/cqrs/commands"

type ICommandValidator interface {
	Validate(command cqrs_commands.ICommand) error
}
