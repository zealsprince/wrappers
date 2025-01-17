package wrappers

import (
	"fmt"
	"reflect"
)

// ValidationError represents an error during validation.
type ValidationError struct {
	WrapperName string
	Value       any
	Reason      string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("invalid value %+v for wrapper %q: %s", e.Value, e.WrapperName, e.Reason)
}

func ErrorNil(name Name) error {
	return &ValidationError{
		WrapperName: string(name),
		Value:       nil,
		Reason:      "value is nil",
	}
}

func ErrorType(name Name, value any) error {
	if value == nil {
		return ErrorNil(name)
	}

	// Check if the value is a pointer type.
	extra := ""
	if reflect.TypeOf(value).Kind() != reflect.Ptr {
		extra = "(if you're passing a wrapper, make sure NOT to dereference it to a value)"
	}

	return &ValidationError{
		WrapperName: string(name),
		Value:       value,
		Reason:      fmt.Sprintf("type mismatch, got %T %s", value, extra),
	}
}

func ErrorValue(name Name, value any, expected string) error {
	return &ValidationError{
		WrapperName: string(name),
		Value:       value,
		Reason:      fmt.Sprintf("expected %s", expected),
	}
}

func ErrorParse(name Name, value any, err error) error {
	return &ValidationError{
		WrapperName: string(name),
		Value:       value,
		Reason:      fmt.Sprintf("failed to parse: %v", err),
	}
}
