package cqrs_behaviors

import (
	"fmt"
	"strings"

	commands "github.com/mitz-it/golang-cqrs/commands"
	queries "github.com/mitz-it/golang-cqrs/queries"
	validators "github.com/mitz-it/golang-cqrs/validators"

	"go.uber.org/dig"
	"golang.org/x/exp/slices"
)

type ValidationBehavior struct {
	Behavior
	commandValidators []validators.ICommandValidator
	queryValidators   []validators.IQueryValidator
}

type ValidationBehaviorParams struct {
	dig.In

	CommandValidators []validators.ICommandValidator `group:"CommandValidators"`
	QueryValidators   []validators.IQueryValidator   `group:"QueryValidators"`
}

func (behavior *ValidationBehavior) SetNext(next Action) {
	behavior.Next = next
}

func (behavior *ValidationBehavior) SetNextRequest(next Request) {
	behavior.NextRequest = next
}

func (behavior *ValidationBehavior) Handle(command commands.ICommand) (commands.IResponse, error) {
	if behavior.commandValidators == nil || len(behavior.commandValidators) <= 0 {
		return behavior.Next(command)
	}

	index := slices.IndexFunc(behavior.commandValidators, func(validator validators.ICommandValidator) bool {
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

func (behavior *ValidationBehavior) HandleQuery(query queries.IQuery) (queries.IResponse, error) {
	if behavior.queryValidators == nil || len(behavior.queryValidators) <= 0 {
		return behavior.NextRequest(query)
	}

	index := slices.IndexFunc(behavior.queryValidators, func(validator validators.IQueryValidator) bool {
		validatorSplittedName := strings.Split(fmt.Sprintf("%T", validator), ".")
		validatorName := validatorSplittedName[len(validatorSplittedName)-1]
		querySplittedName := strings.Split(fmt.Sprintf("%T", query), ".")
		queryName := querySplittedName[len(querySplittedName)-1]
		return strings.Contains(validatorName, queryName)
	})

	if index > -1 {
		err := behavior.queryValidators[index].Validate(query)

		if err != nil {
			return nil, err
		}
	}

	return behavior.Next(query)
}

func NewValidationBehavior(params ValidationBehaviorParams) IBehavior {
	return &ValidationBehavior{
		commandValidators: params.CommandValidators,
		queryValidators:   params.QueryValidators,
	}
}
