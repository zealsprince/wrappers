// Package enum provides wrappers for string-backed enumerated types.
//
// The permitted values live on a zero-sized type that implements Values, so the
// rule reaches them through its type parameter. Nothing has to be registered or
// initialized, and a set with no Values method fails to compile rather than
// accepting everything at runtime.
package enum

import (
	"slices"

	wrappers "github.com/zealsprince/wrappers/v2"
)

// Values supplies the name and the permitted values for an enumerated type.
// Implement it on an empty struct:
//
//	type Direction string
//
//	const (
//		North Direction = "North"
//		South Direction = "South"
//	)
//
//	type directions struct{}
//
//	func (directions) Name() wrappers.Name    { return "CardinalDirection" }
//	func (directions) Values() []Direction    { return []Direction{North, South} }
//
//	type CardinalDirection = enum.Wrapper[Direction, directions]
type Values[T ~string] interface {
	Name() wrappers.Name
	Values() []T
}

// Rule validates that a value is one of V's permitted values.
type Rule[T ~string, V Values[T]] struct{}

func (Rule[T, V]) Name() wrappers.Name {
	var values V

	return values.Name()
}

// Parse accepts the enum's own type and plain strings, so a value arriving from
// JSON as an untyped string converts without the caller restating the type.
func (r Rule[T, V]) Parse(value any) (T, error) {
	switch v := value.(type) {
	case T:
		return v, nil

	case string:
		return T(v), nil

	case []byte:
		return T(v), nil

	default:
		return "", wrappers.ErrValuef(r.Name(), value, "expected a string")
	}
}

func (r Rule[T, V]) Validate(value T) error {
	var values V

	permitted := values.Values()
	if slices.Contains(permitted, value) {
		return nil
	}

	return wrappers.ErrValuef(r.Name(), string(value), "one of %v", permitted)
}

// Wrapper is the strict enum wrapper. Invalid input fails the unmarshal.
type Wrapper[T ~string, V Values[T]] = wrappers.Wrapper[T, Rule[T, V]]

// Lenient discards invalid input instead of failing the unmarshal.
type Lenient[T ~string, V Values[T]] = wrappers.Lenient[T, Rule[T, V]]
