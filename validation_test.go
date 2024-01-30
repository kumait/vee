package vee

import (
	"errors"
	"fmt"
	"testing"
	"time"
)

const SpecialCharacters = `!@#$%^&*()_-+=[]{},;:.?/~\"'`

func TestIfNotBlank(t *testing.T) {
	var tests = []struct {
		name            string
		input           string
		shouldHaveError bool
	}{
		{`IfNotBlank("")`, "", false},
		{`IfNotBlank("t")`, "t", true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := Value(test.input, IfNotBlank(StrMinLen(2))).Check()
			if err == nil && test.shouldHaveError {
				t.Errorf("%q => %v, should get an error but got nil", test.input, err)
				return
			}

			if err != nil && !test.shouldHaveError {
				t.Errorf("%q => %v, should get nil but got error", test.input, err)
				return
			}
		})
	}
}

func TestValidatePassword(t *testing.T) {
	passwordConstraints := Constraints(DefaultCheckSemantic,
		NotBlank(),
		StrMinLen(8),
		StrMaxLen(24),
		ContainsNumber(),
		ContainsLower(),
		ContainsUpper(),
		ContainsAny(SpecialCharacters),
		Contains("???"),
	)

	var tests = []struct {
		name            string
		input           string
		shouldHaveError bool
	}{
		{`valid password`, "test1#Password???", false},
		{`invalid password`, "t", true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := Value(test.input, passwordConstraints).Check()
			if err == nil && test.shouldHaveError {
				t.Errorf("%q => error(%v), should get an error but got nil", test.input, err)
				return
			}

			if err != nil && !test.shouldHaveError {
				t.Errorf("%q => error(%v), should get nil but got error", test.input, err)
				return
			}
		})
	}
}

func TestSlice(t *testing.T) {
	DefaultCheckSemantic = CheckSemanticAll
	tags := []string{"tag1", "dd", "loooong"}

	cons := Field[[]string]("tags", tags,
		Len[[]string](2, 12),
		Each[[]string, string](NotBlank(), StrLen(2, 8)),
	)

	err := cons.Check()
	if err != nil {
		t.Fail()
	}
}

func TestIn(t *testing.T) {
	var languages = map[string]bool{
		"en": true,
		"de": true,
		"fr": true,
	}

	var tests = []struct {
		name            string
		input           string
		shouldHaveError bool
	}{
		{`valid language`, "en", false},
		{`invalid language`, "t", true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := Value(test.input, In(languages)).Check()
			if err == nil && test.shouldHaveError {
				t.Errorf("%q => error(%v), should get an error but got nil", test.input, err)
				return
			}

			if err != nil && !test.shouldHaveError {
				t.Errorf("%q => error(%v), should get nil but got error", test.input, err)
				return
			}
		})
	}
}

func TestRange(t *testing.T) {
	var tests = []struct {
		name            string
		input           int
		shouldHaveError bool
	}{
		{"valid value", 19, false},
		{"invalid value", 21, true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := Value(test.input, Range(2, 20)).Check()
			if err == nil && test.shouldHaveError {
				t.Errorf("%q => error(%v), should get an error but got nil", test.input, err)
				return
			}

			if err != nil && !test.shouldHaveError {
				t.Errorf("%q => error(%v), should get nil but got error", test.input, err)
				return
			}
		})
	}
}

func TestRequired(t *testing.T) {
	s := "test"

	var tests = []struct {
		name            string
		input           *string
		shouldHaveError bool
	}{
		{"valid value", &s, false},
		{"invalid value", nil, true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := Value(test.input, Required[*string]()).Check()
			if err == nil && test.shouldHaveError {
				t.Errorf("%v => error(%v), should get an error but got nil", test.input, err)
				return
			}

			if err != nil && !test.shouldHaveError {
				t.Errorf("%v => error(%v), should get nil but got error", test.input, err)
				return
			}
		})
	}
}

type (
	LoginType int

	User struct {
		Email     string
		Name      string
		Password  string
		Language  string
		Tags      []string
		Status    *int
		LoginType LoginType
	}
)

func (lt LoginType) Validate() error {
	if lt != 25 {
		return errors.New("recursive validation")
	}
	return nil
}

func TestSchema(t *testing.T) {
	var languages = map[string]bool{
		"en": true,
		"de": true,
		"fr": true,
	}

	status := 20
	user := User{
		Email:     "test@test.com",
		Name:      "~",
		Password:  "GG@Test123",
		Language:  "en",
		Tags:      []string{"A1"},
		Status:    &status,
		LoginType: 25,
	}

	DefaultCheckSemantic = CheckSemanticAll

	languageConstraints := Constraints(DefaultCheckSemantic, NotBlank(), StrLen(2, 2), In(languages))
	passwordConstraints := Constraints(DefaultCheckSemantic, NotBlank(),
		StrLen(8, 25),
		ContainsUpper(),
		ContainsLower(),
		ContainsNumber(),
		ContainsAny(SpecialCharacters),
	)

	cons := Schema(
		Field("email", user.Email, NotBlank(), Email()),
		Field("name", user.Name, If(func() bool { return user.Name != "~" }, NotBlank(), StrLen(2, 50))),
		Field("password", user.Password, passwordConstraints),
		Field("language", user.Language, languageConstraints),
		Field("tags", user.Tags, Len[[]string](0, 6), Each[[]string, string](NotBlank(), StrLen(1, 12))),
		Field("status", user.Status, IfNotNil[*int, int](Min(0), Max(20))),
		Field("LoginType", user.LoginType, Range[LoginType](20, 30)),
		// custom validation func
		Func(func() error {
			if len(user.Tags) > 3 {
				return errors.New("test error")
			}
			return nil
		}),
	)

	start := time.Now()
	err := cons.Check()
	fmt.Println(time.Since(start))
	if err != nil {
		t.Fail()
	}
}
