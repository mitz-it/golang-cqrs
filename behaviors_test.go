package cqrs

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

type Behavior1 struct {
}

func (b *Behavior1) Handle(ctx context.Context, request interface{}, next NextFunc) (interface{}, error) {
	return next()
}

type Behavior2 struct {
}

func (b *Behavior2) Handle(ctx context.Context, request interface{}, next NextFunc) (interface{}, error) {
	return next()
}

type Behavior3 struct {
}

func (b *Behavior3) Handle(ctx context.Context, request interface{}, next NextFunc) (interface{}, error) {
	next()
	var res = struct {
		Foo string
	}{
		Foo: "test",
	}
	return &res, nil
}

func behaviors_cleanup(t *testing.T) {
	commandBehaviors = make(map[int]interface{})
	queryBehaviors = make(map[int]interface{})
}

func TestRegisterCommandBehavior_WhenPositionIsNotTaken_ShouldRegisterBehavior(t *testing.T) {
	// arrange
	defer behaviors_cleanup(t)
	behavior := &Behavior1{}

	// act
	RegisterCommandBehavior(0, behavior)

	// assert
	assert.Equal(t, commandBehaviors[0], behavior)
}

func TestRegisterCommandBehavior_WhenPositionIsTaken_ShouldReturnError(t *testing.T) {
	// arrange
	defer behaviors_cleanup(t)
	behavior1 := &Behavior1{}
	behavior2 := &Behavior1{}

	// act
	RegisterCommandBehavior(0, behavior1)
	err := RegisterCommandBehavior(0, behavior2)

	// assert
	assert.Error(t, err)
}

func TestRegisterQueryBehavior_WhenPositionIsNotTaken_ShouldRegisterBehavior(t *testing.T) {
	// arrange
	defer behaviors_cleanup(t)
	behavior := &Behavior1{}

	// act
	RegisterQueryBehavior(0, behavior)

	// assert
	assert.Equal(t, queryBehaviors[0], behavior)
}

func TestRegisterQueryBehavior_WhenPositionIsTaken_ShouldReturnError(t *testing.T) {
	// arrange
	defer behaviors_cleanup(t)
	behavior1 := &Behavior1{}
	behavior2 := &Behavior1{}

	// act
	RegisterQueryBehavior(0, behavior1)
	err := RegisterQueryBehavior(0, behavior2)

	// assert
	assert.Error(t, err)
}

func TestSortBehaviors_GivenBehaviorMap_ShouldSortBehaviors(t *testing.T) {
	// arrange
	behaviorsMap := map[int]interface{}{
		2: &Behavior3{},
		0: &Behavior1{},
		1: &Behavior2{},
	}

	// act
	sorted := sortBehaviors(behaviorsMap)

	// assert
	assert.Equal(t, sorted[0], &Behavior3{})
	assert.Equal(t, sorted[1], &Behavior2{})
	assert.Equal(t, sorted[2], &Behavior1{})
}
