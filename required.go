package vee

import (
	"errors"
	"reflect"
)

type RequiredConstraint[T any] struct {
	Value T
}

func Required[T any]() CheckableValue[T] {
	return &RequiredConstraint[T]{}
}

func (c *RequiredConstraint[T]) Check() error {
	if reflect.ValueOf(c.Value).IsNil() {
		return errors.New("is required")
	}
	return nil
}

func (c *RequiredConstraint[T]) SetValue(value T) {
	c.Value = value
	return
}
