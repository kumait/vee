package vee

import (
	"fmt"
)

type (
	CheckSemantic int

	Checkable interface {
		Check() error
	}

	CheckableValue[T any] interface {
		Checkable
		SetValue(value T) // TODO: consider returning CheckableValue[T]
	}

	ValueConstraint[T any] struct {
		Value      T
		Constraint CheckableValue[T]
	}

	ValueConstraints[T any] struct {
		Value       T
		Constraints []CheckableValue[T]
		Semantic    CheckSemantic
	}

	FieldConstraint[T any] struct {
		Constraint CheckableValue[T]
		FieldName  string
	}

	EachConstraint[T ~[]E, E any] struct {
		Value      T
		Constraint CheckableValue[E]
	}

	IfConstraint[T any] struct {
		Value      T
		Predicate  func() bool
		Constraint CheckableValue[T]
	}

	IfNotNilConstraint[T ~*E, E any] struct {
		Value      T
		Constraint CheckableValue[E]
	}

	FuncConstraint struct {
		Func func() error
	}

	SchemaConstraint struct {
		Constraints []Checkable
	}
)

const (
	CheckSemanticFirst CheckSemantic = iota
	CheckSemanticAll
)

var (
	DefaultCheckSemantic = CheckSemanticFirst
)

func Value[T any](value T, cons ...CheckableValue[T]) CheckableValue[T] {
	c := Constraints(DefaultCheckSemantic, cons...)
	c.SetValue(value)
	return &ValueConstraint[T]{
		Value:      value,
		Constraint: c,
	}
}

func Field[T any](name string, value T, cons ...CheckableValue[T]) CheckableValue[T] {
	fc := &FieldConstraint[T]{
		Constraint: Value(value, cons...),
		FieldName:  name,
	}

	return fc
}

func Each[T ~[]E, E any](cons ...CheckableValue[E]) CheckableValue[T] {
	return &EachConstraint[T, E]{
		Constraint: Constraints(DefaultCheckSemantic, cons...),
	}
}

func If[T any](predicate func() bool, cons ...CheckableValue[T]) CheckableValue[T] {
	return &IfConstraint[T]{
		Predicate:  predicate,
		Constraint: Constraints(DefaultCheckSemantic, cons...),
	}
}

func IfNotNil[T ~*E, E any](cons ...CheckableValue[E]) CheckableValue[T] {
	return &IfNotNilConstraint[T, E]{
		Constraint: Constraints(DefaultCheckSemantic, cons...),
	}
}

func Constraints[T any](semantic CheckSemantic, cons ...CheckableValue[T]) CheckableValue[T] {
	return &ValueConstraints[T]{
		Constraints: cons,
		Semantic:    semantic,
	}
}

func First[T any](cons ...CheckableValue[T]) CheckableValue[T] {
	return Constraints(CheckSemanticFirst, cons...)
}

func All[T any](cons ...CheckableValue[T]) CheckableValue[T] {
	return Constraints(CheckSemanticAll, cons...)
}

func Func(f func() error) Checkable {
	return &FuncConstraint{Func: f}
}

func Schema(cons ...Checkable) Checkable {
	return &SchemaConstraint{Constraints: cons}
}

func (c *ValueConstraint[T]) SetValue(value T) {
	c.Value = value
}

func (c *ValueConstraint[T]) Check() error {
	if t, ok := any(c.Value).(Validatable); ok {
		err := t.Validate()
		if err != nil {
			return err
		}
	}

	err := c.Constraint.Check()
	if err != nil {
		return err
	}

	return nil

}

func (c *ValueConstraints[T]) SetValue(value T) {
	c.Value = value
	for _, con := range c.Constraints {
		con.SetValue(value)
	}
}

func (c *ValueConstraints[T]) Check() error {
	if c.Semantic == CheckSemanticAll {
		var e ErrList
		for _, con := range c.Constraints {
			err := con.Check()
			if err != nil {
				e = append(e, err)
			}
		}

		if e != nil {
			return e
		}
	} else {
		for _, con := range c.Constraints {
			err := con.Check()
			if err != nil {
				return err
			}
		}
	}

	if t, ok := any(c.Value).(Validatable); ok {
		err := t.Validate()
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *FieldConstraint[T]) SetValue(value T) {
	c.Constraint.SetValue(value)
}

func (c *FieldConstraint[T]) Check() error {
	err := c.Constraint.Check()
	if err != nil {
		return FieldError(c.FieldName, err)
	}
	return nil
}

func (c *EachConstraint[T, E]) SetValue(value T) {
	c.Value = value
}

func (c *EachConstraint[T, E]) Check() error {
	if DefaultCheckSemantic == CheckSemanticAll {
		var e ErrList
		for i, item := range c.Value {
			var ie ErrList
			c.Constraint.SetValue(item)
			err := c.Constraint.Check()
			if err != nil {
				ie = append(ie, err)
			}

			if len(ie) > 0 {
				e = append(e, FieldError(fmt.Sprintf("#%d", i), ie))
			}
		}
		if e != nil {
			return e
		}
	} else {
		for i, item := range c.Value {
			c.Constraint.SetValue(item)
			err := c.Constraint.Check()
			if err != nil {
				return FieldError(fmt.Sprintf("#%d", i), err)
			}
		}
	}

	return nil
}

func (c *IfConstraint[T]) SetValue(value T) {
	c.Value = value
}

func (c *IfConstraint[T]) Check() error {
	if c.Predicate() {
		c.Constraint.SetValue(c.Value)
		return c.Constraint.Check()
	}

	return nil
}

func (c *IfNotNilConstraint[T, E]) SetValue(value T) {
	c.Value = value
}

func (c *IfNotNilConstraint[T, E]) Check() error {
	if c.Value != nil {
		c.Constraint.SetValue(*c.Value)
		return c.Constraint.Check()
	}

	return nil
}

func (c *SchemaConstraint) Check() error {
	if DefaultCheckSemantic == CheckSemanticAll {
		var e ErrList
		for _, con := range c.Constraints {
			err := con.Check()
			if err != nil {
				e = append(e, err)
			}
		}
		if e != nil {
			return e
		}
	} else {
		for _, con := range c.Constraints {
			err := con.Check()
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (c *FuncConstraint) Check() error {
	return c.Func()
}
