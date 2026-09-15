package wrappers_test

import (
	"encoding/json"
	"errors"
	"math"
	"sync"
	"testing"

	wrappers "github.com/zealsprince/wrappers/v2"
)

// Each test below pins a defect that v1 actually shipped. They are named after
// the defect so a failure says what regressed rather than which assertion moved.

// v1: w.Wrap("bad", true) then w.Wrap(42, false) left value=42 but Unwrap()=0 and
// IsDiscarded()=true forever, because nothing ever cleared the flag.
func TestSuccessfulWrapClearsAnEarlierDiscard(t *testing.T) {
	var w wrappers.Int

	if w.WrapDiscard("not a number") {
		t.Fatal("WrapDiscard(\"not a number\") = true, want false")
	}

	if !w.IsDiscarded() {
		t.Fatal("IsDiscarded() = false, want true after a rejected value")
	}

	if err := w.Wrap(42); err != nil {
		t.Fatalf("Wrap(42) error = %v, want nil", err)
	}

	if w.IsDiscarded() {
		t.Error("IsDiscarded() = true, want false after a successful wrap")
	}

	if got := w.Get(); got != 42 {
		t.Errorf("Get() = %d, want 42", got)
	}
}

// v1 decoded JSON numbers through float64, so {"id":1.9} became 1 and {"id":1e30}
// became -9223372036854775808. Out-of-range float to int conversion is undefined
// in Go, not merely lossy, so that second one was garbage rather than truncation.
func TestIntegerRulesRejectFractionalAndOutOfRange(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"fractional", `1.9`},
		{"overflow", `1e30`},
		{"negative overflow", `-1e30`},
		{"far past MaxInt64", `92233720368547758079`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var w wrappers.Int

			if err := json.Unmarshal([]byte(tt.input), &w); err == nil {
				t.Fatalf("Unmarshal(%s) error = nil, want rejection (got %d)", tt.input, w.Get())
			}
		})
	}
}

// Lenient coercion is the reason this library exists, so the fix above must not
// have tightened the cases that legitimately worked.
func TestIntegerRulesStillCoerceLooseInput(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  int64
	}{
		{"json integer", `42`, 42},
		{"integral float", `1.0`, 1},
		{"numeric string", `"42"`, 42},
		{"negative string", `"-7"`, -7},
		{"exponent form", `1e3`, 1000},
		{"MaxInt64", `9223372036854775807`, math.MaxInt64},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var w wrappers.Int

			if err := json.Unmarshal([]byte(tt.input), &w); err != nil {
				t.Fatalf("Unmarshal(%s) error = %v, want nil", tt.input, err)
			}

			if got := w.Get(); got != tt.want {
				t.Errorf("Get() = %d, want %d", got, tt.want)
			}
		})
	}
}

// v1's NewWithValue was typed on the wrapped type, so NewWithValue[*WrapperInt]("1")
// did not compile even though Wrap("1", false) worked. The constructor and the
// wrap path now accept the same input.
func TestConstructorAcceptsTheSameLooseInputAsWrap(t *testing.T) {
	w, err := wrappers.Of[wrappers.Int]("1")
	if err != nil {
		t.Fatalf("Of(\"1\") error = %v, want nil", err)
	}

	if got := w.Get(); got != 1 {
		t.Errorf("Get() = %d, want 1", got)
	}
}

// v1's NewWithValueDiscard returned the zero T on error, so an invalid value gave
// back a nil pointer and the first method call panicked.
func TestOfDiscardAlwaysReturnsAUsableWrapper(t *testing.T) {
	w := wrappers.OfDiscard[wrappers.NonEmptyString]("")

	if !w.IsDiscarded() {
		t.Error("IsDiscarded() = false, want true for a rejected value")
	}

	if w.IsValid() {
		t.Error("IsValid() = true, want false")
	}

	// Would have panicked on v1's nil pointer.
	if got := w.Get(); got != "" {
		t.Errorf("Get() = %q, want the empty default", got)
	}
}

// v1 required every field to be a pointer, so an absent field stayed nil, produced
// no error, and panicked on Unwrap. The zero wrapper now reads safely.
func TestZeroWrapperIsUsable(t *testing.T) {
	var w wrappers.String

	if w.IsPresent() {
		t.Error("IsPresent() = true, want false on a zero wrapper")
	}

	if w.IsValid() {
		t.Error("IsValid() = true, want false on a zero wrapper")
	}

	if got := w.Get(); got != "" {
		t.Errorf("Get() = %q, want the empty default", got)
	}
}

// v1 lazily wrote the initialized flag inside MarshalJSON, so marshalling one
// wrapper from several goroutines raced. Run this with -race.
func TestMarshalDoesNotMutate(t *testing.T) {
	w := wrappers.MustOf[wrappers.String]("value")

	var wg sync.WaitGroup
	for range 8 {
		wg.Add(1)

		go func() {
			defer wg.Done()

			if _, err := json.Marshal(w); err != nil {
				t.Errorf("Marshal() error = %v", err)
			}
		}()
	}

	wg.Wait()
}

func TestNestedWrappers(t *testing.T) {
	inner := wrappers.MustOf[wrappers.Int](42)

	var outer wrappers.String
	if err := outer.Wrap(inner); err != nil {
		t.Fatalf("Wrap(inner) error = %v, want nil", err)
	}

	if got := outer.Get(); got != "42" {
		t.Errorf("Get() = %q, want %q", got, "42")
	}
}

// A discarded wrapper handed to another wrapper propagates the discard rather
// than quietly contributing its zero value.
func TestNestedDiscardPropagates(t *testing.T) {
	inner := wrappers.OfDiscard[wrappers.Int]("not a number")

	var outer wrappers.String
	if err := outer.Wrap(inner); err != nil {
		t.Fatalf("Wrap(discarded) error = %v, want nil", err)
	}

	if !outer.IsDiscarded() {
		t.Error("IsDiscarded() = false, want the discard to propagate")
	}
}

func TestNullDiscardsAndReportsErrNil(t *testing.T) {
	var w wrappers.String

	err := json.Unmarshal([]byte(`null`), &w)
	if err == nil {
		t.Fatal("Unmarshal(null) error = nil, want ErrNil")
	}

	if !errors.Is(err, wrappers.ErrNil) {
		t.Errorf("errors.Is(err, ErrNil) = false, want true (got %v)", err)
	}

	if !w.IsDiscarded() {
		t.Error("IsDiscarded() = false, want true after null")
	}
}

func TestLenientSwallowsWhatWrapperRejects(t *testing.T) {
	strict := struct {
		Value wrappers.Int `json:"value"`
	}{}

	if err := json.Unmarshal([]byte(`{"value":"nope"}`), &strict); err == nil {
		t.Error("strict Unmarshal() error = nil, want a rejection")
	}

	lenient := struct {
		Value wrappers.LenientInt `json:"value"`
	}{}

	if err := json.Unmarshal([]byte(`{"value":"nope"}`), &lenient); err != nil {
		t.Errorf("lenient Unmarshal() error = %v, want nil", err)
	}

	if lenient.Value.IsValid() {
		t.Error("IsValid() = true, want false after a discarded value")
	}
}

func TestRoundTripPreservesPayload(t *testing.T) {
	type record struct {
		Name  wrappers.NonEmptyString `json:"name"`
		Count wrappers.Int            `json:"count"`
		Ratio wrappers.Float          `json:"ratio"`
		Fresh wrappers.Bool           `json:"fresh"`
	}

	input := `{"name":"andrew","count":3,"ratio":1.5,"fresh":true}`

	var data record
	if err := json.Unmarshal([]byte(input), &data); err != nil {
		t.Fatalf("Unmarshal() error = %v", err)
	}

	out, err := json.Marshal(data)
	if err != nil {
		t.Fatalf("Marshal() error = %v", err)
	}

	if string(out) != input {
		t.Errorf("round trip changed the payload:\n got %s\nwant %s", out, input)
	}
}

// omitzero replaces v1's instruction to nil the field out by hand before marshalling.
func TestOmitZeroDropsAbsentAndDiscardedFields(t *testing.T) {
	type record struct {
		Name wrappers.String `json:"name"`
		Note wrappers.String `json:"note,omitzero"`
	}

	var data record
	if err := json.Unmarshal([]byte(`{"name":"andrew"}`), &data); err != nil {
		t.Fatalf("Unmarshal() error = %v", err)
	}

	out, err := json.Marshal(data)
	if err != nil {
		t.Fatalf("Marshal() error = %v", err)
	}

	if string(out) != `{"name":"andrew"}` {
		t.Errorf("Marshal() = %s, want the absent field omitted", out)
	}
}
