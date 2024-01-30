package vee

import (
	"fmt"
	"strings"
)

type (
	Validatable interface {
		Validate() error
	}

	ErrList []error

	ErrField struct {
		FieldName string
		Err       error
	}
)

func (el ErrList) Error() string {
	if len(el) == 1 {
		return fmt.Sprintf("%v", el[0])
	}

	var sb strings.Builder
	sb.WriteString("[")
	for i, err := range el {
		if i != 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(fmt.Sprintf("%v", err))
	}
	sb.WriteString("]")

	return sb.String()
}

func (el ErrList) Dto() []map[string]string {
	errMap := make([]map[string]string, 0, len(el))

	for _, err := range el {
		if ef, ok := err.(ErrField); ok {
			errMap = append(errMap, map[string]string{ef.FieldName: ef.Err.Error()})
		} else {
			errMap = append(errMap, map[string]string{ef.FieldName: err.Error()})
		}
	}

	return errMap
}

func (ef ErrField) Error() string {
	return fmt.Sprintf("%s: %v", ef.FieldName, ef.Err)
}

func FieldError(name string, err error) error {
	return ErrField{
		FieldName: name,
		Err:       err,
	}
}
