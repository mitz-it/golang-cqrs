package cqrs_behaviors

import (
	"fmt"
	"strings"

	cqrs_commands "github.com/mitz-it/golang-cqrs/commands"
	cqrs_queries "github.com/mitz-it/golang-cqrs/queries"
	cqrs_validators "github.com/mitz-it/golang-cqrs/validators"

	"go.uber.org/dig"
	"golang.org/x/exp/slices"
)

type ValidationBehavior struct {
	Behavior
	commandValidators []cqrs_validators.ICommandValidator
}

type ValidationBehaviorParams struct {
	dig.In

	CommandValidators []cqrs_validators.ICommandValidator `group:"CommandValidators"`
}

func (behavior *ValidationBehavior) SetNext(next Action) {
	behavior.Next = next
}

func (behavior *ValidationBehavior) SetNextRequest(next Request) {
	behavior.NextRequest = next
}

func (behavior *ValidationBehavior) Handle(command cqrs_commands.ICommand) (cqrs_commands.IResponse, error) {
	index := slices.IndexFunc(behavior.commandValidators, func(validator cqrs_validators.ICommandValidator) bool {
		validatorSplittedName := strings.Split(fmt.Sprintf("%T", validator), ".")
		validatorName := validatorSplittedName[len(validatorSplittedName)-1]
		commandSplittedName := strings.Split(fmt.Sprintf("%T", command), ".")
		commandName := commandSplittedName[len(commandSplittedName)-1]
		return strings.Contains(validatorName, commandName)
	})

	if index > -1 {
		err := behavior.commandValidators[index].Validate(command)

		if err != nil {
			return nil, err
		}
	}

	return behavior.Next(command)
}

func (behavior *ValidationBehavior) HandleQuery(query cqrs_queries.IQuery) cqrs_queries.IResponse {
	return nil
}

func NewValidationBehavior(params ValidationBehaviorParams) IBehavior {
	return &ValidationBehavior{commandValidators: params.CommandValidators}
}
