package vee

import (
	"fmt"
	"golang.org/x/exp/constraints"
)

type (
	RangeConstraint[T constraints.Ordered] struct {
		Value T
		Min   T
		Max   T
	}

	MinConstraint[T constraints.Ordered] struct {
		Value T
		Min   T
	}

	MaxConstraint[T constraints.Ordered] struct {
		Value T
		Max   T
	}
)

func Range[T constraints.Ordered](min, max T) CheckableValue[T] {
	return &RangeConstraint[T]{
		Min: min,
		Max: max,
	}
}

func Min[T constraints.Ordered](min T) CheckableValue[T] {
	return &MinConstraint[T]{
		Min: min,
	}
}

func Max[T constraints.Ordered](max T) CheckableValue[T] {
	return &RangeConstraint[T]{
		Max: max,
	}
}

func (c *RangeConstraint[T]) SetValue(value T) {
	c.Value = value
}

func (c *RangeConstraint[T]) Check() error {
	if c.Value > c.Max {
		return fmt.Errorf("is greater than maximum %v", c.Max)
	} else if c.Value < c.Min {
		return fmt.Errorf("is less than minimum %v", c.Min)
	}

	return nil
}

func (c *MinConstraint[T]) SetValue(value T) {
	c.Value = value
}

func (c *MinConstraint[T]) Check() error {
	if c.Value < c.Min {
		return fmt.Errorf("is less than minimum %v", c.Min)
	}

	return nil
}

func (c *MaxConstraint[T]) SetValue(value T) {
	c.Value = value
}

func (c *MaxConstraint[T]) Check() error {
	if c.Value > c.Max {
		return fmt.Errorf("is greater than maximum %v", c.Max)
	}

	return nil
}
