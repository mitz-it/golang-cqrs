package cqrs

import (
	"context"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

type Response struct {
}

type Command1 struct {
}

type CommandHandler1 struct {
}

func (c *CommandHandler1) Handle(ctx context.Context, command *Command1) (*Response, error) {
	return &Response{}, nil
}

func commands_cleanup(t *testing.T) {
	t.Cleanup(func() {
		commandHandlers = make(map[reflect.Type]interface{})
	})
}

func TestRegisterCommandHandler_WhenNotAnyHandlerForCommand_ShouldAddHandlerToMap(t *testing.T) {
	// arrange
	defer commands_cleanup(t)
	var command *Command1
	commandType := reflect.TypeOf(command)
	handler := &CommandHandler1{}
	// act
	err := RegisterCommandHandler[*Command1, *Response](handler)

	// assert
	assert.Nil(t, err)
	assert.Equal(t, commandHandlers[commandType], handler)
}

func TestRegisterCommandHandler_WhenHandlerForCommandAlreadyRegisterd_ShouldReturnError(t *testing.T) {
	// arrange
	defer commands_cleanup(t)
	handler := &CommandHandler1{}

	// act
	RegisterCommandHandler[*Command1, *Response](handler)
	err := RegisterCommandHandler[*Command1, *Response](handler)

	// assert
	assert.Error(t, err)
}

func TestSend_WhenNoHandlerRegistered_ShouldReturnError(t *testing.T) {
	// arrange
	defer commands_cleanup(t)
	command := &Command1{}
	// act
	_, err := Send[*Command1, *Response](context.TODO(), command)

	// assert
	assert.Error(t, err)
}

func TestSend_WhenHandlerCantBeCasted_ShouldReturnError(t *testing.T) {
	// arrange
	defer commands_cleanup(t)
	var handler interface{} = func() {

	}
	command := &Command1{}
	commandType := reflect.TypeOf(command)
	commandHandlers[commandType] = handler

	// act
	_, err := Send[*Command1, *Response](context.TODO(), command)

	// assert
	assert.Error(t, err)
}

func TestSend_WhenNoCommandBehaviors_ShouldCallHandler(t *testing.T) {
	// arrange
	defer commands_cleanup(t)
	command := &Command1{}
	handler := &CommandHandler1{}
	RegisterCommandHandler[*Command1, *Response](handler)

	// act
	res, err := Send[*Command1, *Response](context.TODO(), command)

	// assert
	assert.Nil(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, &Response{}, res)
}

func TestSend_WhenHaveCommandBehaviors_ShouldCallHandlerThroughPipeline(t *testing.T) {
	// arrange
	defer commands_cleanup(t)
	command := &Command1{}
	handler := &CommandHandler1{}
	RegisterCommandBehavior(0, &Behavior1{})
	RegisterCommandHandler[*Command1, *Response](handler)

	// act
	res, err := Send[*Command1, *Response](context.TODO(), command)

	// assert
	assert.Nil(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, &Response{}, res)
}

func TestSend_WhenHaveCommandBehaviorsAndResponseCanBeCasted_ShouldDefaultResponse(t *testing.T) {
	// arrange
	defer commands_cleanup(t)
	command := &Command1{}
	handler := &CommandHandler1{}
	RegisterCommandBehavior(0, &Behavior3{})
	RegisterCommandHandler[*Command1, *Response](handler)

	// act
	res, err := Send[*Command1, *Response](context.TODO(), command)

	// assert
	assert.Nil(t, err)
	assert.Equal(t, &Response{}, res)
}
