package cqrs

import (
	"context"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

type Query1 struct {
}

type QueryHandler1 struct {
}

func (c *QueryHandler1) Handle(ctx context.Context, query *Query1) (*Response, error) {
	return &Response{}, nil
}

func querys_cleanup(t *testing.T) {
	t.Cleanup(func() {
		queryHandlers = make(map[reflect.Type]interface{})
	})
}

func TestRegisterQueryHandler_WhenNotAnyHandlerForQuery_ShouldAddHandlerToMap(t *testing.T) {
	// arrange
	defer querys_cleanup(t)
	var query *Query1
	queryType := reflect.TypeOf(query)
	handler := &QueryHandler1{}
	// act
	err := RegisterQueryHandler[*Query1, *Response](handler)

	// assert
	assert.Nil(t, err)
	assert.Equal(t, queryHandlers[queryType], handler)
}

func TestRegisterQueryHandler_WhenHandlerForQueryAlreadyRegisterd_ShouldReturnError(t *testing.T) {
	// arrange
	defer querys_cleanup(t)
	handler := &QueryHandler1{}

	// act
	RegisterQueryHandler[*Query1, *Response](handler)
	err := RegisterQueryHandler[*Query1, *Response](handler)

	// assert
	assert.Error(t, err)
}

func TestRequest_WhenNoHandlerRegistered_ShouldReturnError(t *testing.T) {
	// arrange
	defer querys_cleanup(t)
	query := &Query1{}
	// act
	_, err := Request[*Query1, *Response](context.TODO(), query)

	// assert
	assert.Error(t, err)
}

func TestRequest_WhenHandlerCantBeCasted_ShouldReturnError(t *testing.T) {
	// arrange
	defer querys_cleanup(t)
	var handler interface{} = func() {

	}
	query := &Query1{}
	queryType := reflect.TypeOf(query)
	queryHandlers[queryType] = handler

	// act
	_, err := Request[*Query1, *Response](context.TODO(), query)

	// assert
	assert.Error(t, err)
}

func TestRequest_WhenNoQueryBehaviors_ShouldCallHandler(t *testing.T) {
	// arrange
	defer querys_cleanup(t)
	query := &Query1{}
	handler := &QueryHandler1{}
	RegisterQueryHandler[*Query1, *Response](handler)

	// act
	res, err := Request[*Query1, *Response](context.TODO(), query)

	// assert
	assert.Nil(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, &Response{}, res)
}

func TestRequest_WhenHaveQueryBehaviors_ShouldCallHandlerThroughPipeline(t *testing.T) {
	// arrange
	defer querys_cleanup(t)
	query := &Query1{}
	handler := &QueryHandler1{}
	RegisterQueryBehavior(0, &Behavior1{})
	RegisterQueryHandler[*Query1, *Response](handler)

	// act
	res, err := Request[*Query1, *Response](context.TODO(), query)

	// assert
	assert.Nil(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, &Response{}, res)
}

func TestRequest_WhenHaveQueryBehaviorsAndResponseCanBeCasted_ShouldDefaultResponse(t *testing.T) {
	// arrange
	defer querys_cleanup(t)
	query := &Query1{}
	handler := &QueryHandler1{}
	RegisterQueryBehavior(0, &Behavior3{})
	RegisterQueryHandler[*Query1, *Response](handler)

	// act
	res, err := Request[*Query1, *Response](context.TODO(), query)

	// assert
	assert.Nil(t, err)
	assert.Equal(t, &Response{}, res)
}
