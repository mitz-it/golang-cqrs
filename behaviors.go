package mediator

import (
	"context"
	"errors"
	"fmt"
	"sort"
)

type NextFunc func() (interface{}, error)

type IBehavior interface {
	Handle(ctx context.Context, command interface{}, next NextFunc) (interface{}, error)
}

var commandBehaviors map[int]interface{}
var queryBehaviors map[int]interface{}

func init() {
	commandBehaviors = make(map[int]interface{})
	queryBehaviors = make(map[int]interface{})
}

func RegisterCommandBehavior(order int, behavior IBehavior) error {
	_, found := commandBehaviors[order]

	if found {
		msg := fmt.Sprintf("position %d is taken by another command behavior.", order)
		return errors.New(msg)
	}

	commandBehaviors[order] = behavior

	return nil

}

func RegisterQueryBehavior(order int, behavior IBehavior) error {
	_, found := queryBehaviors[order]

	if found {
		msg := fmt.Sprintf("position %d is taken by another query behavior.", order)
		return errors.New(msg)
	}

	queryBehaviors[order] = behavior

	return nil
}

func sortBehaviors(behaviors map[int]interface{}) []interface{} {
	keys := make([]int, len(behaviors)-1)

	for key := range behaviors {
		keys = append(keys, key)
	}

	sort.Ints(keys)

	sorted := make([]interface{}, len(behaviors)-1)

	for _, key := range keys {
		sorted = append(sorted, behaviors[key])
	}

	return sorted
}
