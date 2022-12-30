package cqrs_behaviors

import (
	"context"
	"fmt"
	"strings"

	commands "github.com/mitz-it/golang-cqrs/commands"
	queries "github.com/mitz-it/golang-cqrs/queries"
	validation "github.com/mitz-it/golang-validation"

	"go.uber.org/dig"
	"golang.org/x/exp/slices"
)

type ValidationBehavior struct {
	Behavior
	commandValidators []validation.IValidator
	queryValidators   []validation.IValidator
}

type ValidationBehaviorParams struct {
	dig.In

	CommandValidators []validation.IValidator `group:"CommandValidators"`
	QueryValidators   []validation.IValidator `group:"QueryValidators"`
}

func (behavior *ValidationBehavior) SetNextAction(next Action) {
	behavior.NextAction = next
}

func (behavior *ValidationBehavior) SetNextRequest(next Request) {
	behavior.NextRequest = next
}

func (behavior *ValidationBehavior) HandleCommand(ctx context.Context, command commands.ICommand) (commands.IResponse, error) {
	if behavior.commandValidators == nil || len(behavior.commandValidators) <= 0 {
		return behavior.NextAction(ctx, command)
	}

	index := slices.IndexFunc(behavior.commandValidators, func(validator validation.IValidator) bool {
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

	return behavior.NextAction(ctx, command)
}

func (behavior *ValidationBehavior) HandleQuery(ctx context.Context, query queries.IQuery) (queries.IResponse, error) {
	if behavior.queryValidators == nil || len(behavior.queryValidators) <= 0 {
		return behavior.NextRequest(ctx, query)
	}

	index := slices.IndexFunc(behavior.queryValidators, func(validator validation.IValidator) bool {
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

	return behavior.NextRequest(ctx, query)
}

func NewValidationBehavior(params ValidationBehaviorParams) IBehavior {
	return &ValidationBehavior{
		commandValidators: params.CommandValidators,
		queryValidators:   params.QueryValidators,
	}
}
