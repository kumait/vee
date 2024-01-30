package vee

import (
	"errors"
)

type (
	InConstraint[T comparable] struct {
		Value       T
		ValidValues map[T]bool
	}

	NotInConstraint[T comparable] struct {
		Value         T
		InvalidValues map[T]bool
	}
)

func In[T comparable](values map[T]bool) CheckableValue[T] {
	return &InConstraint[T]{
		ValidValues: values,
	}
}

func NotIn[T comparable](values map[T]bool) CheckableValue[T] {
	return &NotInConstraint[T]{
		InvalidValues: values,
	}
}

func (c *InConstraint[T]) SetValue(value T) {
	c.Value = value
}

func (c *InConstraint[T]) Check() error {
	if _, ok := c.ValidValues[c.Value]; !ok {
		return errors.New("is not in valid values")
	}
	return nil
}

func (c *NotInConstraint[T]) SetValue(value T) {
	c.Value = value
}

func (c *NotInConstraint[T]) Check() error {
	if _, ok := c.InvalidValues[c.Value]; ok {
		return errors.New("is in invalid values")
	}
	return nil
}
