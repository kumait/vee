package vee

import (
	"errors"
	"fmt"
	"regexp"
)

var EmailRegex = regexp.MustCompile(`^\w+([.-]?\w+)*@\w+([.-]?\w+)*(\.\w{2,3})+$`)

type RegexConstraint struct {
	Value string
	Error string
	Regex *regexp.Regexp
}

func Email() CheckableValue[string] {
	return &RegexConstraint{
		Error: "invalid email",
		Regex: EmailRegex,
	}
}

func (c *RegexConstraint) SetValue(value string) {
	c.Value = value
	return
}

func (c *RegexConstraint) Check() (err error) {
	if !c.Regex.Match([]byte(c.Value)) {
		err = errors.New(c.Error)
	}
	return
}

func Regex(regex *regexp.Regexp) CheckableValue[string] {
	return &RegexConstraint{
		Error: fmt.Sprintf("value does not match regex %s", regex.String()),
		Regex: regex,
	}
}
