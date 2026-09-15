package wrappers_test

import (
	"encoding/json"
	"errors"
	"math"
	"strings"
	"testing"
	"time"

	wrappers "github.com/zealsprince/wrappers/v2"
)

func TestMustOfPanicsOnInvalidInput(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Error("MustOf() did not panic on invalid input")
		}
	}()

	wrappers.MustOf[wrappers.Int]("not a number")
}

func TestMustOfReturnsTheValueWhenValid(t *testing.T) {
	if got := wrappers.MustOf[wrappers.Int](7).Get(); got != 7 {
		t.Errorf("Get() = %d, want 7", got)
	}
}

func TestMarshalDiscardedWrapperAsNull(t *testing.T) {
	w := wrappers.OfDiscard[wrappers.Int]("not a number")

	out, err := json.Marshal(w)
	if err != nil {
		t.Fatalf("Marshal() error = %v", err)
	}

	if string(out) != "null" {
		t.Errorf("Marshal() = %s, want null", out)
	}
}

func TestUnmarshalMalformedJSON(t *testing.T) {
	var w wrappers.Int

	err := w.UnmarshalJSON([]byte(`{`))
	if err == nil {
		t.Fatal("UnmarshalJSON() error = nil, want a decode error")
	}

	if !w.IsDiscarded() {
		t.Error("IsDiscarded() = false, want true after a decode failure")
	}

	if !w.IsPresent() {
		t.Error("IsPresent() = false, want true, the field was in the payload")
	}
}

func TestErrorMessages(t *testing.T) {
	var w wrappers.Int

	// A rejection carrying a value prints it.
	err := w.Wrap("not a number")
	if err == nil {
		t.Fatal("Wrap() error = nil")
	}

	if !strings.Contains(err.Error(), "Int") || !strings.Contains(err.Error(), "not a number") {
		t.Errorf("Error() = %q, want the wrapper name and the value", err)
	}

	// A nil rejection has no value to print.
	err = w.Wrap(nil)
	if err == nil {
		t.Fatal("Wrap(nil) error = nil")
	}

	if strings.Contains(err.Error(), "%!") || !strings.Contains(err.Error(), "nil") {
		t.Errorf("Error() = %q, want a clean nil message", err)
	}
}

func TestFieldErrorUnwrapsToItsCause(t *testing.T) {
	var data struct {
		Value wrappers.Int `json:"value"`
	}

	err := wrappers.Check(&data)
	if err == nil {
		t.Fatal("Check() = nil, want an error")
	}

	var ferr *wrappers.FieldError
	if !errors.As(err, &ferr) {
		t.Fatalf("Check() = %v, want a *FieldError", err)
	}

	if ferr.Unwrap() == nil {
		t.Error("FieldError.Unwrap() = nil, want the underlying cause")
	}

	if !strings.Contains(ferr.Error(), "Value") {
		t.Errorf("FieldError.Error() = %q, want the field named", ferr)
	}
}

func TestCheckHandlesEveryContainerShape(t *testing.T) {
	t.Run("nil input", func(t *testing.T) {
		if err := wrappers.Check(nil); err != nil {
			t.Errorf("Check(nil) = %v, want nil", err)
		}
	})

	t.Run("map of structs", func(t *testing.T) {
		type row struct {
			Code wrappers.NonEmptyString `json:"code"`
		}

		data := struct {
			Rows map[string]row `json:"rows"`
		}{Rows: map[string]row{"a": {}}}

		err := wrappers.Check(&data)
		if err == nil || !strings.Contains(err.Error(), "Rows.a.Code") {
			t.Errorf("Check() = %v, want the map key in the path", err)
		}
	})

	t.Run("wrapper behind an interface field", func(t *testing.T) {
		data := struct {
			Any any
		}{Any: wrappers.Int{}}

		if err := wrappers.Check(&data); err == nil {
			t.Error("Check() = nil, want the wrapper behind the interface reported")
		}
	})

	t.Run("populated pointer field", func(t *testing.T) {
		value := wrappers.MustOf[wrappers.Int](1)

		data := struct {
			Value *wrappers.Int
		}{Value: &value}

		if err := wrappers.Check(&data); err != nil {
			t.Errorf("Check() = %v, want nil", err)
		}
	})

	t.Run("non-struct input", func(t *testing.T) {
		if err := wrappers.Check("just a string"); err != nil {
			t.Errorf("Check(string) = %v, want nil", err)
		}
	})
}

func TestParsersHandleJSONNumbers(t *testing.T) {
	tests := []struct {
		name   string
		target any
		input  string
		want   string
	}{
		{"string from number", new(wrappers.String), `42`, `"42"`},
		{"float from number", new(wrappers.Float), `1.5`, `1.5`},
		{"bool from number", new(wrappers.Bool), `1`, `true`},
		{"bool from zero", new(wrappers.Bool), `0`, `false`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := json.Unmarshal([]byte(tt.input), tt.target); err != nil {
				t.Fatalf("Unmarshal(%s) error = %v", tt.input, err)
			}

			out, err := json.Marshal(tt.target)
			if err != nil {
				t.Fatalf("Marshal() error = %v", err)
			}

			if string(out) != tt.want {
				t.Errorf("round trip = %s, want %s", out, tt.want)
			}
		})
	}
}

func TestParserRejectionPaths(t *testing.T) {
	t.Run("string from unconvertible", func(t *testing.T) {
		var w wrappers.String

		if err := w.Wrap(map[string]int{}); err == nil {
			t.Error("Wrap(map) error = nil, want a type error")
		}
	})

	t.Run("float from unconvertible", func(t *testing.T) {
		var w wrappers.Float

		if err := w.Wrap(map[string]int{}); err == nil {
			t.Error("Wrap(map) error = nil, want a type error")
		}
	})

	t.Run("float from nonsense string", func(t *testing.T) {
		var w wrappers.Float

		if err := w.Wrap("andrew"); err == nil {
			t.Error("Wrap(\"andrew\") error = nil, want a rejection")
		}
	})

	t.Run("bool from unconvertible", func(t *testing.T) {
		var w wrappers.Bool

		if err := w.Wrap(map[string]int{}); err == nil {
			t.Error("Wrap(map) error = nil, want a type error")
		}
	})

	t.Run("bool from a fractional number", func(t *testing.T) {
		var w wrappers.Bool

		if err := w.Wrap(1.5); err == nil {
			t.Error("Wrap(1.5) error = nil, want a rejection")
		}
	})

	t.Run("int from uint64 past MaxInt64", func(t *testing.T) {
		var w wrappers.Int

		if err := w.Wrap(uint64(math.MaxInt64) + 1); err == nil {
			t.Error("Wrap(huge uint64) error = nil, want a rejection")
		}
	})

	t.Run("int from uint", func(t *testing.T) {
		var w wrappers.Int

		if err := w.Wrap(uint(9)); err != nil {
			t.Errorf("Wrap(uint(9)) error = %v, want nil", err)
		}
	})

	t.Run("string from uint", func(t *testing.T) {
		var w wrappers.String

		if err := w.Wrap(uint16(9)); err != nil {
			t.Errorf("Wrap(uint16(9)) error = %v, want nil", err)
		}

		if got := w.Get(); got != "9" {
			t.Errorf("Get() = %q, want %q", got, "9")
		}
	})

	t.Run("string from float32", func(t *testing.T) {
		var w wrappers.String

		if err := w.Wrap(float32(1.5)); err != nil {
			t.Errorf("Wrap(float32) error = %v, want nil", err)
		}
	})

	t.Run("float from int", func(t *testing.T) {
		var w wrappers.Float

		if err := w.Wrap(3); err != nil {
			t.Errorf("Wrap(3) error = %v, want nil", err)
		}
	})
}

func TestTimeRuleNamesAndEdges(t *testing.T) {
	t.Run("names are reported", func(t *testing.T) {
		var (
			unix wrappers.TimeUnix
			iso  wrappers.TimeISO8601
			dur  wrappers.Duration
		)

		for _, tc := range []struct {
			got  wrappers.Name
			want wrappers.Name
		}{
			{unix.Name(), "TimeUnix"},
			{iso.Name(), "TimeISO8601"},
			{dur.Name(), "Duration"},
		} {
			if tc.got != tc.want {
				t.Errorf("Name() = %q, want %q", tc.got, tc.want)
			}
		}
	})

	t.Run("time.Time passes through every rule", func(t *testing.T) {
		now := time.Date(2026, 3, 6, 12, 0, 0, 0, time.UTC)

		var unix wrappers.TimeUnix
		if err := unix.Wrap(now); err != nil {
			t.Errorf("TimeUnix.Wrap(time.Time) error = %v", err)
		}

		var iso wrappers.TimeISO8601
		if err := iso.Wrap(now); err != nil {
			t.Errorf("TimeISO8601.Wrap(time.Time) error = %v", err)
		}
	})

	t.Run("iso8601 rejects a non-string non-time", func(t *testing.T) {
		var w wrappers.TimeISO8601

		if err := w.Wrap(true); err == nil {
			t.Error("Wrap(true) error = nil, want a rejection")
		}
	})

	t.Run("unix rejects unconvertible input", func(t *testing.T) {
		var w wrappers.TimeUnix

		if err := w.Wrap("not a number"); err == nil {
			t.Error("Wrap(string) error = nil, want a rejection")
		}
	})

	t.Run("iso8601 normalizes a short string without panicking", func(t *testing.T) {
		var w wrappers.TimeISO8601

		if err := w.Wrap("2026"); err == nil {
			t.Error("Wrap(\"2026\") error = nil, want a rejection")
		}
	})

	t.Run("duration rejects unconvertible input", func(t *testing.T) {
		var w wrappers.Duration

		if err := w.Wrap(map[string]int{}); err == nil {
			t.Error("Wrap(map) error = nil, want a type error")
		}
	})

	t.Run("duration rejects a fractional nanosecond count", func(t *testing.T) {
		var w wrappers.Duration

		if err := w.Wrap(1.5); err == nil {
			t.Error("Wrap(1.5) error = nil, want a rejection")
		}
	})
}

func TestRemainingEdges(t *testing.T) {
	t.Run("check on a nil wrapper pointer directly", func(t *testing.T) {
		if err := wrappers.Check((*wrappers.Int)(nil)); err != nil {
			t.Errorf("Check(nil pointer) = %v, want nil at the root", err)
		}
	})

	t.Run("check sees a nil wrapper behind an interface field", func(t *testing.T) {
		data := struct {
			Any any
		}{Any: (*wrappers.Int)(nil)}

		if err := wrappers.Check(&data); err == nil {
			t.Error("Check() = nil, want the nil wrapper reported")
		}
	})

	t.Run("check ignores nil pointers to non-wrappers", func(t *testing.T) {
		data := struct {
			Name *string
		}{}

		if err := wrappers.Check(&data); err != nil {
			t.Errorf("Check() = %v, want nil for a nil *string", err)
		}
	})

	t.Run("int accepts the narrow widths", func(t *testing.T) {
		for _, input := range []any{int16(5), uint8(5), uint16(5)} {
			var w wrappers.Int

			if err := w.Wrap(input); err != nil {
				t.Errorf("Wrap(%T) error = %v, want nil", input, err)
			}

			if got := w.Get(); got != 5 {
				t.Errorf("Get() = %d, want 5", got)
			}
		}
	})

	t.Run("bool rejects a fractional json number", func(t *testing.T) {
		var w wrappers.Bool

		if err := json.Unmarshal([]byte(`1.5`), &w); err == nil {
			t.Error("Unmarshal(1.5) error = nil, want a rejection")
		}
	})

	t.Run("iso8601 marshals back to RFC3339", func(t *testing.T) {
		var w wrappers.TimeISO8601

		if err := json.Unmarshal([]byte(`"2026-03-06T12:30:00+0000"`), &w); err != nil {
			t.Fatalf("Unmarshal() error = %v", err)
		}

		out, err := json.Marshal(w)
		if err != nil {
			t.Fatalf("Marshal() error = %v", err)
		}

		if string(out) != `"2026-03-06T12:30:00Z"` {
			t.Errorf("Marshal() = %s, want the normalized RFC3339 form", out)
		}
	})
}
