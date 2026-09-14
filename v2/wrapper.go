package wrappers

import (
	"bytes"
	"encoding/json"
)

// Rule is the validation policy for a wrapper, carried as a type parameter rather
// than embedded. A rule is expected to be zero-sized: the core reaches it with
// `var r R`, so there is no construction step and nothing to forget to initialize.
type Rule[T any] interface {
	// Name identifies the rule in error messages.
	Name() Name

	// Parse converts loosely typed input into T. Nested wrappers and nil are
	// resolved by the core before Parse is called, so implementations only see
	// concrete values. Reuse the built-in parsers (StringParser, Int64Parser and
	// friends) by embedding them instead of writing this from scratch.
	Parse(any) (T, error)

	// Validate applies the semantic check on top of a parsed value. Plain types
	// embed a parser that supplies a no-op.
	Validate(T) error
}

// unwrapper is the shape every wrapper satisfies. It lets one wrapper accept
// another as input without the core knowing the concrete types involved.
type unwrapper interface {
	UnwrapAny() any
	IsDiscarded() bool
}

// Wrapper holds a validated value of type T under rule R.
//
// The zero Wrapper is usable. It carries no initialized flag, no lazy setup and
// no mutation during marshalling, so marshalling the same value from several
// goroutines is safe.
type Wrapper[T any, R Rule[T]] struct {
	value     T
	discarded bool
	present   bool
}

// Of builds a wrapper from a value, validating it.
func Of[T any, R Rule[T]](value any) (Wrapper[T, R], error) {
	var w Wrapper[T, R]
	err := w.Wrap(value)

	return w, err
}

// MustOf builds a wrapper from a value and panics if it does not validate.
// Intended for package-level constants and tests, not for request data.
func MustOf[T any, R Rule[T]](value any) Wrapper[T, R] {
	w, err := Of[T, R](value)
	if err != nil {
		panic(err)
	}

	return w
}

// OfDiscard builds a wrapper from a value, discarding it if invalid. Unlike v1
// this always returns a usable wrapper, never a nil pointer.
func OfDiscard[T any, R Rule[T]](value any) Wrapper[T, R] {
	var w Wrapper[T, R]
	w.WrapDiscard(value)

	return w
}

// Wrap validates a value and stores it. A successful Wrap clears any prior
// discard, so a wrapper can be reused without carrying a stale flag.
func (w *Wrapper[T, R]) Wrap(value any) error {
	var rule R

	w.present = true

	// A wrapper handed to another wrapper contributes its unwrapped value, and a
	// discarded one propagates the discard rather than its zero value.
	if nested, ok := value.(unwrapper); ok {
		if nested.IsDiscarded() {
			w.discard()

			return nil
		}

		value = nested.UnwrapAny()
	}

	if value == nil {
		w.discard()

		return errNil(rule.Name())
	}

	parsed, err := rule.Parse(value)
	if err != nil {
		w.discard()

		return named(rule.Name(), err)
	}

	if err := rule.Validate(parsed); err != nil {
		w.discard()

		return named(rule.Name(), err)
	}

	w.value = parsed
	w.discarded = false

	return nil
}

// WrapDiscard validates a value and discards it if invalid, reporting whether it
// was kept. The v1 boolean parameter is gone: the two behaviours are two methods.
func (w *Wrapper[T, R]) WrapDiscard(value any) bool {
	return w.Wrap(value) == nil
}

// Get returns the stored value in its wrapped type, zero if discarded.
func (w Wrapper[T, R]) Get() T {
	if w.discarded {
		var zero T

		return zero
	}

	return w.value
}

// Unwrap returns the value in the form the rule serializes to.
func (w Wrapper[T, R]) Unwrap() any {
	var rule R
	if u, ok := any(rule).(interface{ Unwrap(T) any }); ok {
		return u.Unwrap(w.Get())
	}

	return w.Get()
}

// UnwrapAny satisfies the unwrapper interface so wrappers can nest.
func (w Wrapper[T, R]) UnwrapAny() any { return w.Unwrap() }

// IsDiscarded reports whether the last wrap rejected its input.
func (w Wrapper[T, R]) IsDiscarded() bool { return w.discarded }

// IsPresent reports whether the wrapper was ever given input. A field absent from
// a JSON payload leaves this false, which is what Check keys off.
func (w Wrapper[T, R]) IsPresent() bool { return w.present }

// IsValid reports whether the wrapper holds a value that passed its rule.
func (w Wrapper[T, R]) IsValid() bool { return w.present && !w.discarded }

// Name returns the rule's name.
func (w Wrapper[T, R]) Name() Name {
	var rule R

	return rule.Name()
}

func (w *Wrapper[T, R]) discard() {
	var zero T

	w.value = zero
	w.discarded = true
}

// Discard clears the value and flags the wrapper as discarded.
func (w *Wrapper[T, R]) Discard() { w.discard() }

// MarshalJSON writes the unwrapped value, or null if discarded. It does not
// mutate the wrapper.
func (w Wrapper[T, R]) MarshalJSON() ([]byte, error) {
	if w.discarded {
		return []byte("null"), nil
	}

	return json.Marshal(w.Unwrap())
}

// UnmarshalJSON decodes and validates in one step. Numbers are decoded as
// json.Number so integer rules can reject fractional and out-of-range input
// instead of silently truncating through float64.
func (w *Wrapper[T, R]) UnmarshalJSON(data []byte) error {
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.UseNumber()

	var raw any
	if err := decoder.Decode(&raw); err != nil {
		w.present = true
		w.discard()

		return err
	}

	return w.Wrap(raw)
}

// IsZero lets `json:",omitzero"` drop absent and discarded fields, which removes
// the v1 step of nilling the field out by hand to omit it.
func (w Wrapper[T, R]) IsZero() bool { return !w.present || w.discarded }

// Lenient is a Wrapper that never fails unmarshalling. Invalid input is discarded
// silently and the caller checks IsValid. It replaces v1's Discarder, with no
// Proxy field to reach through: the wrapper's own methods are promoted.
//
// Lenient is about bad values. Absent fields are a separate axis, handled by the
// `wrappers:"optional"` tag that Check reads.
type Lenient[T any, R Rule[T]] struct {
	Wrapper[T, R]
}

// UnmarshalJSON decodes into the underlying wrapper and swallows rejection.
func (o *Lenient[T, R]) UnmarshalJSON(data []byte) error {
	if err := o.Wrapper.UnmarshalJSON(data); err != nil {
		o.Wrapper.discard()
	}

	return nil
}
