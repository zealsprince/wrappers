package wrappers

import (
	"testing"

	"github.com/biter777/countries"
)

// NewWithValueDiscard promises a discarded wrapper, so it has to hand one back
// even when the value is rejected. Up to v1.1.6 it delegated to NewWithValue,
// which returns the zero T on error, so an invalid value produced a nil pointer
// and the first method call on it panicked.
func TestNewWithValueDiscardReturnsUsableWrapper(t *testing.T) {
	wrapper := NewWithValueDiscard[*WrapperCountry](countries.Unknown)

	if wrapper == nil {
		t.Fatal("NewWithValueDiscard() = nil, want a discarded wrapper")
	}

	// These would panic on a nil pointer, which is the regression being pinned.
	if !wrapper.IsDiscarded() {
		t.Error("IsDiscarded() = false, want true for a rejected value")
	}

	// WrapperCountry unwraps to the country name, and a discarded one reports
	// "Unknown" rather than an empty string.
	if got := wrapper.Unwrap(); got != countries.Unknown.String() {
		t.Errorf("Unwrap() = %q, want %q", got, countries.Unknown.String())
	}
}

// The valid path has to keep working unchanged.
func TestNewWithValueDiscardKeepsValidValues(t *testing.T) {
	wrapper := NewWithValueDiscard[*WrapperCountry](countries.DEU)

	if wrapper == nil {
		t.Fatal("NewWithValueDiscard() = nil, want a populated wrapper")
	}

	if wrapper.IsDiscarded() {
		t.Error("IsDiscarded() = true, want false for a valid value")
	}

	if got := wrapper.Unwrap(); got != countries.DEU.String() {
		t.Errorf("Unwrap() = %q, want %q", got, countries.DEU.String())
	}
}
