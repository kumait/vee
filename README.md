# vee: Simple Go Validator

## What is vee
Use Vee for common validation scenarios, such as required, string length, ranges, regexes and so on.

## Features
* Simple to use and flexible.
* Conditional validations.
* Recursive validations.
* No reflection.

## Usage

### Value Validation
To validate a value, use the `Value` function as below.

```go
package examples

import (
	"fmt"
	v "github.com/kumait/vee"
)

func validateName(name string) {
	err := v.Value(name, v.NotBlank()).Check()
	if err != nil {
		fmt.Printf("Validation error: %v\n", err)
	}
}

func validatePercentage(number int) {
	err := v.Value(number, v.Range(0, 100)).Check()
	if err != nil {
		fmt.Printf("Validation error: %v\n", err)
	}
}
```

### Schema Validation
Usually, validations are done on structs, you can validate a struct using the `Schema` function.

```go
package dto

import (
	v "github.com/kumait/vee"
)

type (
	CreateFeedbackRequest struct {
		Type  string `json:"type"`
		Email string `json:"email"`
		Name  string `json:"name"`
		Body  string `json:"body"`
	}
)

var FeedbackTypes = map[string]bool{
	"general":     true,
	"bug":         true,
	"enhancement": true,
}

func (d *CreateFeedbackRequest) Validate() error {
	return v.Schema(
		v.Field("type", d.Type, v.In(FeedbackTypes)),
		v.Field("email", d.Email, v.IfNotBlank(v.Email(), v.StrMaxLen(255))),
		v.Field("name", d.Name, v.IfNotBlank(v.StrMaxLen(255))),
		v.Field("body", d.Body, v.NotBlank(), v.StrLen(1, 500)),
	).Check()
}
```
The example above shows a case for conditional validation wherein a field is only validated if it is not blank.

### Recursive Validation

Recursive validation is supported when a type conforms to the `Validatble` interface which defines a single method `Validate() error`.

```go
package main

import (
	v "github.com/kumait/vee"
)

type (
	User struct {
		Name    string
		Address Address
	}
	Address struct {
		Street string
		City   string
	}
)

func (r User) Validate() error {
	return v.Schema(
		v.Field("Name", r.Name, v.NotBlank(), v.StrMaxLen(100)),
		v.Field("Address", r.Address),
	).Check()
}

func (a Address) Validate() error {
	return v.Schema(
		v.Field("Street", a.Street, v.NotBlank()),
		v.Field("Address", a.City, v.NotBlank()),
	).Check()
}

func main() {
	request := User{
		Name: "test",
		Address: Address{
			Street: "",
			City:   "Ottawa",
		},
	}

	err := request.Validate()
	if err != nil {
		println(err.Error())
	}
}
```
In the above example `Address` conform to the `Validable` interface, vee will recursively validate `Address` because of that.

### Slice and Array Validation

Slice and array validation is supported through the `Each` constraint.

```go
package examples

import v "github.com/kumait/vee"

type (
	Attribute struct {
		Name  string `json:"name"`
		Value string `json:"value"`
	}

	CreateVocRequest struct {
		Term        string       `json:"term"`
		POS         string       `json:"pos"`
		Attributes  []Attribute  `json:"attributes"`
		Tags        []string     `json:"tags"`
	}
)

var tagsConstraint = v.Constraints(v.DefaultCheckSemantic, v.Len[[]string](0, 20), v.Each[[]string, string](v.NotBlank(), v.StrMaxLen(20)))

func (d *CreateVocRequest) Validate() error {
	return v.Schema(
		v.Field("term", d.Term, v.StrLen(1, 100)),
		v.Field("attributes", d.Attributes, v.Len[[]Attribute](0, 50), v.Each[[]Attribute, Attribute]()),
		v.Field("tags", d.Tags, tagsConstraint),
	).Check()
}

func (d Attribute) Validate() error {
	return v.Schema(
		v.Field("name", d.Name, v.NotBlank(), v.StrMaxLen(255)),
		v.Field("value", d.Value, v.NotBlank(), v.StrMaxLen(255)),
	).Check()
}
```
In the above example the slice `[]Attributes` is validated for its length and each item in it is validated.

Constraints can also be saved to variables and used in several validations, the `DefaultCheckSemantic` is set by default to `CheckSemanticFirst` which returns the first error only,
this can be changed to `CheckSemanticAll` to return all errors.
