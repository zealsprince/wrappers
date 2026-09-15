package wrappers_test

import (
	"encoding/json"
	"testing"
	"time"

	wrappers "github.com/zealsprince/wrappers/v2"
)

func TestTimeAcceptsRFC3339(t *testing.T) {
	var w wrappers.Time

	if err := json.Unmarshal([]byte(`"2026-03-06T12:30:00Z"`), &w); err != nil {
		t.Fatalf("Unmarshal() error = %v, want nil", err)
	}

	want := time.Date(2026, 3, 6, 12, 30, 0, 0, time.UTC)
	if got := w.Get(); !got.Equal(want) {
		t.Errorf("Get() = %v, want %v", got, want)
	}
}

// v1's Wrap assigned the zero time when the parse failed and the discard flag was
// set, so a bad timestamp became a silent 0001-01-01 rather than a rejection.
func TestTimeRejectsBadInputInsteadOfZeroing(t *testing.T) {
	tests := []string{`"not a time"`, `"2026-13-45T00:00:00Z"`, `"06/03/2026"`, `true`}

	for _, input := range tests {
		t.Run(input, func(t *testing.T) {
			var w wrappers.Time

			if err := json.Unmarshal([]byte(input), &w); err == nil {
				t.Fatalf("Unmarshal(%s) error = nil, want a rejection", input)
			}

			if !w.IsDiscarded() {
				t.Error("IsDiscarded() = false, want true")
			}
		})
	}
}

// Epoch input is ambiguous between seconds and milliseconds, so Time refuses it
// and TimeUnix is the explicit way to ask for seconds.
func TestTimeRefusesBareNumbers(t *testing.T) {
	var w wrappers.Time

	if err := json.Unmarshal([]byte(`1772800200`), &w); err == nil {
		t.Error("Unmarshal(number) error = nil, want Time to refuse an ambiguous epoch")
	}
}

func TestTimeUnixAcceptsEpochSeconds(t *testing.T) {
	var w wrappers.TimeUnix

	if err := json.Unmarshal([]byte(`1772800200`), &w); err != nil {
		t.Fatalf("Unmarshal() error = %v, want nil", err)
	}

	if got := w.Get().Unix(); got != 1772800200 {
		t.Errorf("Get().Unix() = %d, want 1772800200", got)
	}

	out, err := json.Marshal(w)
	if err != nil {
		t.Fatalf("Marshal() error = %v", err)
	}

	if string(out) != "1772800200" {
		t.Errorf("Marshal() = %s, want the epoch back", out)
	}
}

func TestTimeRoundTrips(t *testing.T) {
	input := `"2026-03-06T12:30:00Z"`

	var w wrappers.Time
	if err := json.Unmarshal([]byte(input), &w); err != nil {
		t.Fatalf("Unmarshal() error = %v", err)
	}

	out, err := json.Marshal(w)
	if err != nil {
		t.Fatalf("Marshal() error = %v", err)
	}

	if string(out) != input {
		t.Errorf("Marshal() = %s, want %s", out, input)
	}
}

func TestTimeISO8601NormalizesLooseSpellings(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"zulu", `"2026-03-06T12:30:00Z"`},
		{"plus 0000 offset", `"2026-03-06T12:30:00+0000"`},
		{"no zone at all", `"2026-03-06T12:30:00"`},
	}

	want := time.Date(2026, 3, 6, 12, 30, 0, 0, time.UTC)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var w wrappers.TimeISO8601

			if err := json.Unmarshal([]byte(tt.input), &w); err != nil {
				t.Fatalf("Unmarshal(%s) error = %v, want nil", tt.input, err)
			}

			if got := w.Get(); !got.Equal(want) {
				t.Errorf("Get() = %v, want %v", got, want)
			}
		})
	}
}

// A real numeric offset has to survive normalization rather than being given a
// second zone designator.
func TestTimeISO8601KeepsNumericOffsets(t *testing.T) {
	var w wrappers.TimeISO8601

	if err := json.Unmarshal([]byte(`"2026-03-06T12:30:00+02:00"`), &w); err != nil {
		t.Fatalf("Unmarshal() error = %v, want nil", err)
	}

	want := time.Date(2026, 3, 6, 10, 30, 0, 0, time.UTC)
	if got := w.Get(); !got.Equal(want) {
		t.Errorf("Get() = %v, want %v", got, want)
	}
}

func TestDurationCoercion(t *testing.T) {
	tests := []struct {
		name  string
		input any
		want  time.Duration
	}{
		{"duration string", "1h30m", 90 * time.Minute},
		{"time.Duration", 5 * time.Second, 5 * time.Second},
		{"int nanoseconds", 5, 5},
		{"int64 nanoseconds", int64(5), 5},
		{"int32 nanoseconds", int32(5), 5},
		{"float32 nanoseconds", float32(5), 5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var w wrappers.Duration

			// Every one of these but int and float64 panicked on v1.
			if err := w.Wrap(tt.input); err != nil {
				t.Fatalf("Wrap(%v) error = %v, want nil", tt.input, err)
			}

			if got := w.Get(); got != tt.want {
				t.Errorf("Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDurationRejectsNonsense(t *testing.T) {
	for _, input := range []string{`"not a duration"`, `"1 hour"`, `true`} {
		t.Run(input, func(t *testing.T) {
			var w wrappers.Duration

			if err := json.Unmarshal([]byte(input), &w); err == nil {
				t.Errorf("Unmarshal(%s) error = nil, want a rejection", input)
			}
		})
	}
}

func TestDurationRoundTripsAsAString(t *testing.T) {
	input := `"1h30m0s"`

	var w wrappers.Duration
	if err := json.Unmarshal([]byte(input), &w); err != nil {
		t.Fatalf("Unmarshal() error = %v", err)
	}

	out, err := json.Marshal(w)
	if err != nil {
		t.Fatalf("Marshal() error = %v", err)
	}

	if string(out) != input {
		t.Errorf("Marshal() = %s, want %s", out, input)
	}
}
