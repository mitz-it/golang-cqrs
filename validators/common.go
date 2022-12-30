package cqrs_validators

import (
	"reflect"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

func Nested(target interface{}, fieldRules ...*validation.FieldRules) *validation.FieldRules {
	return validation.Field(target, validation.By(func(field interface{}) error {
		value := reflect.Indirect(reflect.ValueOf(field))

		if value.CanAddr() {
			addr := value.Addr().Interface()
			return validation.ValidateStruct(addr, fieldRules...)
		}

		return validation.ValidateStruct(target, fieldRules...)
	}))
}

var StringRule []validation.Rule

func init() {
	StringRule = []validation.Rule{
		validation.Required,
		validation.Length(1, 200),
	}
}
