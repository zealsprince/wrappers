package wrappers_test

import (
	"errors"
	"math"
	"testing"

	wrappers "github.com/zealsprince/wrappers/v2"
)

// The coercion table is the reason this library exists over a struct tag
// validator, so it gets exercised input type by input type.

func TestStringCoercion(t *testing.T) {
	tests := []struct {
		name  string
		input any
		want  string
	}{
		{"string", "andrew", "andrew"},
		{"int", 42, "42"},
		{"int64", int64(-7), "-7"},
		{"uint8", uint8(3), "3"},
		{"float64 integral", float64(1), "1"},
		{"float64 fractional", 1.5, "1.5"},
		{"bool", true, "true"},
		{"bytes", []byte("raw"), "raw"},
		{"empty string is allowed", "", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var w wrappers.String

			if err := w.Wrap(tt.input); err != nil {
				t.Fatalf("Wrap(%v) error = %v, want nil", tt.input, err)
			}

			if got := w.Get(); got != tt.want {
				t.Errorf("Get() = %q, want %q", got, tt.want)
			}
		})
	}
}

// A float that formats as 1e+21 through %v would land in the payload in exponent
// form, so the parser has to use a plain decimal formatting.
func TestStringCoercionOfLargeFloats(t *testing.T) {
	var w wrappers.String

	if err := w.Wrap(1e21); err != nil {
		t.Fatalf("Wrap(1e21) error = %v", err)
	}

	if got := w.Get(); got != "1000000000000000000000" {
		t.Errorf("Get() = %q, want a plain decimal", got)
	}
}

func TestStringRejectsUnconvertibleInput(t *testing.T) {
	var w wrappers.String

	err := w.Wrap(struct{ A int }{1})
	if err == nil {
		t.Fatal("Wrap(struct) error = nil, want a type error")
	}

	if !errors.Is(err, wrappers.ErrType) {
		t.Errorf("errors.Is(err, ErrType) = false, want true")
	}

	// The shared parsers cannot name themselves, so the core stamps the rule name on.
	var verr *wrappers.ValidationError
	if errors.As(err, &verr) && verr.Wrapper != "String" {
		t.Errorf("ValidationError.Wrapper = %q, want %q", verr.Wrapper, "String")
	}
}

func TestNonEmptyStringRejectsTheEmptyString(t *testing.T) {
	var w wrappers.NonEmptyString

	if err := w.Wrap(""); err == nil {
		t.Fatal("Wrap(\"\") error = nil, want a rejection")
	}

	if err := w.Wrap("andrew"); err != nil {
		t.Errorf("Wrap(\"andrew\") error = %v, want nil", err)
	}
}

func TestIntCoercion(t *testing.T) {
	tests := []struct {
		name  string
		input any
		want  int64
	}{
		{"int", 42, 42},
		{"int32", int32(-7), -7},
		{"uint32", uint32(9), 9},
		{"float integral", float64(3), 3},
		{"string", "42", 42},
		{"negative string", "-7", -7},
		{"MaxInt64", int64(math.MaxInt64), math.MaxInt64},
		{"MinInt64", int64(math.MinInt64), math.MinInt64},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var w wrappers.Int

			if err := w.Wrap(tt.input); err != nil {
				t.Fatalf("Wrap(%v) error = %v, want nil", tt.input, err)
			}

			if got := w.Get(); got != tt.want {
				t.Errorf("Get() = %d, want %d", got, tt.want)
			}
		})
	}
}

func TestIntRejections(t *testing.T) {
	tests := []struct {
		name  string
		input any
	}{
		{"fractional float", 1.9},
		{"NaN", math.NaN()},
		{"positive infinity", math.Inf(1)},
		{"negative infinity", math.Inf(-1)},
		{"overflowing float", 1e30},
		{"uint64 past MaxInt64", uint64(math.MaxInt64) + 1},
		{"non numeric string", "andrew"},
		{"bool", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var w wrappers.Int

			if err := w.Wrap(tt.input); err == nil {
				t.Fatalf("Wrap(%v) error = nil, want a rejection (got %d)", tt.input, w.Get())
			}

			if !w.IsDiscarded() {
				t.Error("IsDiscarded() = false, want true after a rejection")
			}
		})
	}
}

func TestFloatCoercion(t *testing.T) {
	tests := []struct {
		name  string
		input any
		want  float64
	}{
		{"float64", 1.5, 1.5},
		{"float32", float32(0.5), 0.5},
		{"int", 3, 3},
		{"string", "1.5", 1.5},
		{"exponent string", "1e3", 1000},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var w wrappers.Float

			if err := w.Wrap(tt.input); err != nil {
				t.Fatalf("Wrap(%v) error = %v, want nil", tt.input, err)
			}

			if got := w.Get(); got != tt.want {
				t.Errorf("Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

// A non-finite float would marshal to invalid JSON, so it is rejected on the way in.
func TestFloatRejectsNonFinite(t *testing.T) {
	for _, input := range []string{"NaN", "Inf", "-Inf"} {
		t.Run(input, func(t *testing.T) {
			var w wrappers.Float

			if err := w.Wrap(input); err == nil {
				t.Errorf("Wrap(%q) error = nil, want a rejection", input)
			}
		})
	}
}

func TestBoolCoercion(t *testing.T) {
	tests := []struct {
		name  string
		input any
		want  bool
	}{
		{"bool", true, true},
		{"string true", "true", true},
		{"string false", "false", false},
		{"string 1", "1", true},
		{"int 1", 1, true},
		{"int 0", 0, false},
		{"float 1", 1.0, true},
		// v1 grouped these into one case clause, so v kept the interface type and
		// int64(0) compared unequal to the untyped 0 that defaults to int. Every
		// width gets its own case here, and every width gets its own test.
		{"int64 zero", int64(0), false},
		{"int32 zero", int32(0), false},
		{"int8 zero", int8(0), false},
		{"float32 zero", float32(0), false},
		{"float64 zero", float64(0), false},
		{"int64 one", int64(1), true},
		{"float32 one", float32(1), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var w wrappers.Bool

			if err := w.Wrap(tt.input); err != nil {
				t.Fatalf("Wrap(%v) error = %v, want nil", tt.input, err)
			}

			if got := w.Get(); got != tt.want {
				t.Errorf("Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBoolRejectsNonsense(t *testing.T) {
	var w wrappers.Bool

	if err := w.Wrap("maybe"); err == nil {
		t.Error("Wrap(\"maybe\") error = nil, want a rejection")
	}
}
