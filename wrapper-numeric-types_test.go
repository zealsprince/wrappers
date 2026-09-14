package wrappers

import (
	"testing"
	"time"
)

// Both wrappers below used to group several numeric types into one case clause,
// which leaves v with the interface type rather than the concrete one. The
// duration wrapper then asserted v.(int) and panicked, and the bool wrapper
// compared an interface against an untyped constant and got the wrong answer.
// Nesting reaches both, since WrapperInt.UnwrapAny hands over an int64.

func TestWrapperTimeDurationAcceptsEveryNumericType(t *testing.T) {
	tests := []struct {
		name  string
		value any
		want  time.Duration
	}{
		{"int", int(5), 5},
		{"int8", int8(5), 5},
		{"int16", int16(5), 5},
		{"int32", int32(5), 5},
		{"int64", int64(5), 5},
		{"float32", float32(5), 5},
		{"float64", float64(5), 5},
		{"time.Duration", 5 * time.Second, 5 * time.Second},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wrapper := New[*WrapperTimeDuration]()

			// This panicked for every type except int, float64 and time.Duration.
			if err := wrapper.Wrap(tt.value, false); err != nil {
				t.Fatalf("Wrap(%v) error = %v, want nil", tt.value, err)
			}

			if got := wrapper.Get(); got != tt.want {
				t.Errorf("Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWrapperTimeDurationAcceptsNestedWrapper(t *testing.T) {
	count := New[*WrapperInt]()
	if err := count.Wrap(5, false); err != nil {
		t.Fatalf("Wrap(5) error = %v", err)
	}

	wrapper := New[*WrapperTimeDuration]()

	// UnwrapAny returns an int64 here, which is what used to panic.
	if err := wrapper.Wrap(count, false); err != nil {
		t.Fatalf("Wrap(WrapperInt) error = %v, want nil", err)
	}

	if got := wrapper.Get(); got != time.Duration(5) {
		t.Errorf("Get() = %v, want %v", got, time.Duration(5))
	}
}

func TestWrapperBoolTreatsEveryNumericZeroAsFalse(t *testing.T) {
	tests := []struct {
		name  string
		value any
		want  bool
	}{
		{"int zero", int(0), false},
		{"int8 zero", int8(0), false},
		{"int16 zero", int16(0), false},
		{"int32 zero", int32(0), false},
		{"int64 zero", int64(0), false},
		{"float32 zero", float32(0), false},
		{"float64 zero", float64(0), false},
		{"int64 one", int64(1), true},
		{"float32 one", float32(1), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wrapper := New[*WrapperBool]()

			if err := wrapper.Wrap(tt.value, false); err != nil {
				t.Fatalf("Wrap(%v) error = %v, want nil", tt.value, err)
			}

			// int64(0), int32(0) and float32(0) all came back as true.
			if got := wrapper.Unwrap(); got != tt.want {
				t.Errorf("Unwrap() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWrapperBoolAcceptsNestedWrapper(t *testing.T) {
	count := New[*WrapperInt]()
	if err := count.Wrap(0, false); err != nil {
		t.Fatalf("Wrap(0) error = %v", err)
	}

	wrapper := New[*WrapperBool]()
	if err := wrapper.Wrap(count, false); err != nil {
		t.Fatalf("Wrap(WrapperInt) error = %v, want nil", err)
	}

	if wrapper.Unwrap() {
		t.Error("Unwrap() = true, want false for a wrapped zero")
	}
}
