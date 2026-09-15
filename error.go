package wrappers

import (
	"errors"
	"fmt"
)

// Name identifies a wrapper in error messages. Rules supply it.
type Name string

// Sentinel causes. Callers match these with errors.Is rather than comparing strings.
var (
	ErrNil     = errors.New("value is nil")
	ErrType    = errors.New("type mismatch")
	ErrValue   = errors.New("value out of range or malformed")
	ErrMissing = errors.New("field absent from input")
)

// ValidationError carries which wrapper rejected what, and why.
type ValidationError struct {
	Wrapper Name
	Value   any
	Reason  string
	Cause   error
}

func (e *ValidationError) Error() string {
	if e.Value == nil {
		return fmt.Sprintf("%s: %s", e.Wrapper, e.Reason)
	}

	return fmt.Sprintf("%s: invalid value %#v: %s", e.Wrapper, e.Value, e.Reason)
}

func (e *ValidationError) Unwrap() error { return e.Cause }

// FieldError names the struct field a ValidationError came from. Check produces these.
type FieldError struct {
	Field string
	Cause error
}

func (e *FieldError) Error() string { return fmt.Sprintf("%s: %v", e.Field, e.Cause) }

func (e *FieldError) Unwrap() error { return e.Cause }

func errNil(name Name) error {
	return &ValidationError{Wrapper: name, Reason: "value is nil", Cause: ErrNil}
}

func errType(name Name, value any) error {
	return &ValidationError{
		Wrapper: name,
		Value:   value,
		Reason:  fmt.Sprintf("cannot convert %T", value),
		Cause:   ErrType,
	}
}

// ValueErrorf builds a rejection with a formatted reason. Rules call this from
// Validate. TypeError covers the other case, input that could not be converted
// at all. Both are exported because rules live in other packages.
func ValueErrorf(name Name, value any, format string, args ...any) error {
	return &ValidationError{
		Wrapper: name,
		Value:   value,
		Reason:  fmt.Sprintf(format, args...),
		Cause:   ErrValue,
	}
}

// named stamps the rule's name onto errors that parsers raised before they knew
// it. Parsers are shared across rules, so they cannot name themselves.
func named(name Name, err error) error {
	var verr *ValidationError
	if errors.As(err, &verr) && verr.Wrapper == "" {
		verr.Wrapper = name
	}

	return err
}

// TypeError builds a rejection for input a rule cannot convert.
func TypeError(name Name, value any) error {
	return errType(name, value)
}
