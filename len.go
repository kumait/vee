package vee

import (
	"fmt"
	"reflect"
)

type LenConstraint[T any] struct {
	Value T
	Min   int
	Max   int
}

func Len[T any](min, max int) CheckableValue[T] {
	return &LenConstraint[T]{
		Min: min,
		Max: max,
	}
}

func (c *LenConstraint[T]) SetValue(value T) {
	c.Value = value
}

func (c *LenConstraint[T]) Check() error {
	var l int
	v := reflect.ValueOf(c.Value)
	switch v.Kind() {
	case reflect.Slice, reflect.Array, reflect.Map:
		l = v.Len()
	default:
		panic("invalid value")
	}

	if c.Min == c.Max {
		if l != c.Min {
			return fmt.Errorf("must have %d items", c.Max)
		}
	}

	if l > c.Max {
		return fmt.Errorf("must have %d items at most", c.Max)
	} else if l < c.Min {
		return fmt.Errorf("must have %d items at least", c.Min)
	}

	return nil
}
