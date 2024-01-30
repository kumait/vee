package vee

import (
	"errors"
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"
)

type (
	StrLenConstraint struct {
		Value     string
		MinLength int
		MaxLength int
	}

	StrMinLenConstraint struct {
		Value     string
		MinLength int
	}

	StrMaxLenConstraint struct {
		Value     string
		MaxLength int
	}

	NotBlankConstraint struct {
		Value string
	}

	ContainsPredicateConstraint struct {
		Value        string
		ErrorMessage string
		Predicate    func(rune) bool
	}

	ContainsConstraint struct {
		Value string
		Str   string
		Any   bool
	}

	IfNotBlankConstraint struct {
		Value      string
		Constraint CheckableValue[string]
	}
)

var (
	containsUpperConstraint = &ContainsPredicateConstraint{
		Predicate: func(r rune) bool {
			return unicode.IsUpper(r)
		},
		ErrorMessage: "must contain one upper case character at least",
	}

	containsLowerConstraint = &ContainsPredicateConstraint{
		Predicate: func(r rune) bool {
			return unicode.IsLower(r)
		},
		ErrorMessage: "must contain one lower case character at least",
	}

	containsNumberConstraint = &ContainsPredicateConstraint{
		Predicate: func(r rune) bool {
			return unicode.IsNumber(r)
		},
		ErrorMessage: "must contain one number character at least",
	}
)

var (
	notBlankError = errors.New("cannot be blank")
)

func StrLen(min int, max int) CheckableValue[string] {
	return &StrLenConstraint{
		MinLength: min,
		MaxLength: max,
	}
}

func StrMinLen(min int) CheckableValue[string] {
	return &StrMinLenConstraint{
		MinLength: min,
	}
}

func StrMaxLen(max int) CheckableValue[string] {
	return &StrMaxLenConstraint{
		MaxLength: max,
	}
}

func NotBlank() CheckableValue[string] {
	return new(NotBlankConstraint)
}

func IfNotBlank(cons ...CheckableValue[string]) CheckableValue[string] {
	return &IfNotBlankConstraint{
		Constraint: Constraints(DefaultCheckSemantic, cons...),
	}
}

func ContainsPredicate(predicate func(rune) bool, message string) CheckableValue[string] {
	return &ContainsPredicateConstraint{
		Predicate:    predicate,
		ErrorMessage: message,
	}
}

func ContainsUpper() CheckableValue[string] {
	return containsUpperConstraint
}

func ContainsLower() CheckableValue[string] {
	return containsLowerConstraint
}

func ContainsNumber() CheckableValue[string] {
	return containsNumberConstraint
}

func Contains(str string) CheckableValue[string] {
	return &ContainsConstraint{
		Str: str,
		Any: false,
	}
}

func ContainsAny(str string) CheckableValue[string] {
	return &ContainsConstraint{
		Str: str,
		Any: true,
	}
}

func (c *StrLenConstraint) SetValue(value string) {
	c.Value = value
}

func (c *StrLenConstraint) Check() error {
	l := utf8.RuneCountInString(c.Value)
	if c.MinLength == c.MaxLength {
		if l != c.MinLength {
			return fmt.Errorf("must have %d characters", c.MaxLength)
		}
	}

	if l > c.MaxLength {
		return fmt.Errorf("must have %d characters at most", c.MaxLength)
	} else if l < c.MinLength {
		return fmt.Errorf("must have %d characters at least", c.MinLength)
	}

	return nil
}

func (c *StrMinLenConstraint) SetValue(value string) {
	c.Value = value
}

func (c *StrMinLenConstraint) Check() error {
	l := utf8.RuneCountInString(c.Value)
	if l < c.MinLength {
		return fmt.Errorf("must have %d characters at least", c.MinLength)
	}

	return nil
}

func (c *StrMaxLenConstraint) SetValue(value string) {
	c.Value = value
}

func (c *StrMaxLenConstraint) Check() error {
	l := utf8.RuneCountInString(c.Value)
	if l > c.MaxLength {
		return fmt.Errorf("must have %d characters at most", c.MaxLength)
	}

	return nil
}

func (c *NotBlankConstraint) SetValue(value string) {
	c.Value = value
}

func (c *NotBlankConstraint) Check() error {
	if c.Value == "" {
		return notBlankError
	}

	return nil
}

func (c *ContainsPredicateConstraint) SetValue(value string) {
	c.Value = value
}

func (c *ContainsPredicateConstraint) Check() error {
	for _, r := range c.Value {
		if c.Predicate(r) {
			return nil
		}
	}
	return errors.New(c.ErrorMessage)
}

func (c *ContainsConstraint) SetValue(value string) {
	c.Value = value
}

func (c *ContainsConstraint) Check() error {
	if c.Any {
		if !strings.ContainsAny(c.Value, c.Str) {
			return errors.New(fmt.Sprintf("must contain any of the characters %s", c.Str))
		}
		return nil
	}

	if !strings.Contains(c.Value, c.Str) {
		return errors.New(fmt.Sprintf(`must contain the string "%s"`, c.Str))
	}
	return nil
}

func (c *IfNotBlankConstraint) SetValue(value string) {
	c.Value = value
}

func (c *IfNotBlankConstraint) Check() error {
	if c.Value != "" {
		c.Constraint.SetValue(c.Value)
		return c.Constraint.Check()
	}

	return nil
}
