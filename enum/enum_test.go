package enum_test

import (
	"encoding/json"
	"errors"
	"testing"

	wrappers "github.com/zealsprince/wrappers/v2"
	"github.com/zealsprince/wrappers/v2/enum"
)

// Again, a zero struct. A wrapper that needed setup would fail here.
type payload struct {
	Heading enum.CardinalDirections `json:"heading"`
}

func TestAcceptsPermittedValues(t *testing.T) {
	for _, want := range []enum.CardinalDirection{
		enum.DirectionNorth,
		enum.DirectionEast,
		enum.DirectionSouth,
		enum.DirectionWest,
	} {
		t.Run(string(want), func(t *testing.T) {
			var data payload

			input := `{"heading":"` + string(want) + `"}`
			if err := json.Unmarshal([]byte(input), &data); err != nil {
				t.Fatalf("Unmarshal(%s) error = %v, want nil", input, err)
			}

			if got := data.Heading.Get(); got != want {
				t.Errorf("Get() = %q, want %q", got, want)
			}
		})
	}
}

func TestRejectsValuesOutsideTheSet(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"unknown value", `{"heading":"northwest"}`},
		{"wrong case", `{"heading":"North"}`},
		{"empty string", `{"heading":""}`},
		{"wrong type", `{"heading":42}`},
		{"null", `{"heading":null}`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var data payload

			err := json.Unmarshal([]byte(tt.input), &data)
			if err == nil {
				t.Fatalf("Unmarshal(%s) error = nil, want a validation error", tt.input)
			}

			var verr *wrappers.ValidationError
			if !errors.As(err, &verr) {
				t.Fatalf("error = %v, want a *wrappers.ValidationError", err)
			}

			if data.Heading.IsValid() {
				t.Error("IsValid() = true, want false after a rejected value")
			}
		})
	}
}

// v1 discarded the empty string silently while erroring on every other bad value.
// Pin the consistent behaviour so it does not drift back.
func TestEmptyStringIsRejectedLikeAnyOtherBadValue(t *testing.T) {
	var data payload

	err := json.Unmarshal([]byte(`{"heading":""}`), &data)
	if err == nil {
		t.Fatal("Unmarshal() error = nil, want the empty string rejected")
	}

	if !errors.Is(err, wrappers.ErrValue) {
		t.Errorf("errors.Is(err, ErrValue) = false, want true")
	}
}

func TestLenientDiscardsInsteadOfFailing(t *testing.T) {
	var data struct {
		Heading enum.LenientCardinalDirections `json:"heading"`
	}

	if err := json.Unmarshal([]byte(`{"heading":"northwest"}`), &data); err != nil {
		t.Fatalf("Unmarshal() error = %v, want nil from a lenient wrapper", err)
	}

	if data.Heading.IsValid() {
		t.Error("IsValid() = true, want false")
	}
}

func TestRoundTrip(t *testing.T) {
	input := `{"heading":"south"}`

	var data payload
	if err := json.Unmarshal([]byte(input), &data); err != nil {
		t.Fatalf("Unmarshal() error = %v", err)
	}

	out, err := json.Marshal(data)
	if err != nil {
		t.Fatalf("Marshal() error = %v", err)
	}

	if string(out) != input {
		t.Errorf("round trip changed the payload: got %s, want %s", out, input)
	}
}

// Check has to see through to a wrapper nested behind an enum alias.
func TestCheckReportsAbsentField(t *testing.T) {
	var data payload

	if err := json.Unmarshal([]byte(`{}`), &data); err != nil {
		t.Fatalf("Unmarshal() error = %v", err)
	}

	err := wrappers.Check(&data)
	if err == nil {
		t.Fatal("Check() = nil, want an error for the absent field")
	}

	if !errors.Is(err, wrappers.ErrMissing) {
		t.Errorf("errors.Is(err, ErrMissing) = false, want true (got %v)", err)
	}
}

func TestParseAcceptsTheEnumTypeAndBytes(t *testing.T) {
	t.Run("the enum's own type", func(t *testing.T) {
		var w enum.CardinalDirections

		if err := w.Wrap(enum.DirectionNorth); err != nil {
			t.Fatalf("Wrap(CardinalDirection) error = %v, want nil", err)
		}

		if got := w.Get(); got != enum.DirectionNorth {
			t.Errorf("Get() = %q, want %q", got, enum.DirectionNorth)
		}
	})

	t.Run("raw bytes", func(t *testing.T) {
		var w enum.CardinalDirections

		if err := w.Wrap([]byte("south")); err != nil {
			t.Fatalf("Wrap([]byte) error = %v, want nil", err)
		}

		if got := w.Get(); got != enum.DirectionSouth {
			t.Errorf("Get() = %q, want %q", got, enum.DirectionSouth)
		}
	})

	t.Run("unconvertible type", func(t *testing.T) {
		var w enum.CardinalDirections

		if err := w.Wrap(map[string]int{}); err == nil {
			t.Error("Wrap(map) error = nil, want a type error")
		}
	})
}
