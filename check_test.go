package wrappers_test

import (
	"encoding/json"
	"errors"
	"strings"
	"testing"
	"time"

	wrappers "github.com/zealsprince/wrappers/v2"
)

type account struct {
	Name    wrappers.NonEmptyString `json:"name"`
	Balance wrappers.Int            `json:"balance"`
	Note    wrappers.String         `json:"note" wrappers:"optional"`
}

// The hole v1 never closed: a field absent from the payload is never handed to
// UnmarshalJSON, so no validation inside the wrapper can ever see it.
func TestCheckReportsAbsentFields(t *testing.T) {
	var data account

	if err := json.Unmarshal([]byte(`{"name":"andrew"}`), &data); err != nil {
		t.Fatalf("Unmarshal() error = %v", err)
	}

	err := wrappers.Check(&data)
	if err == nil {
		t.Fatal("Check() = nil, want an error for the absent balance")
	}

	if !errors.Is(err, wrappers.ErrMissing) {
		t.Errorf("errors.Is(err, ErrMissing) = false, want true")
	}

	if !strings.Contains(err.Error(), "Balance") {
		t.Errorf("Check() = %v, want the field named", err)
	}

	// Note is tagged optional, so its absence is not a problem.
	if strings.Contains(err.Error(), "Note") {
		t.Errorf("Check() = %v, want the optional field left alone", err)
	}
}

func TestCheckPassesWhenEverythingIsPresent(t *testing.T) {
	var data account

	if err := json.Unmarshal([]byte(`{"name":"andrew","balance":10,"note":"hi"}`), &data); err != nil {
		t.Fatalf("Unmarshal() error = %v", err)
	}

	if err := wrappers.Check(&data); err != nil {
		t.Errorf("Check() = %v, want nil", err)
	}
}

// One call has to report every bad field, not just the first, or fixing a payload
// turns into a round trip per field.
func TestCheckJoinsEveryProblem(t *testing.T) {
	var data account

	if err := json.Unmarshal([]byte(`{}`), &data); err != nil {
		t.Fatalf("Unmarshal() error = %v", err)
	}

	err := wrappers.Check(&data)
	if err == nil {
		t.Fatal("Check() = nil, want errors")
	}

	for _, field := range []string{"Name", "Balance"} {
		if !strings.Contains(err.Error(), field) {
			t.Errorf("Check() = %v, want %s reported", err, field)
		}
	}
}

// A field that arrived but failed its rule is a different problem from an absent
// one, and Check has to distinguish them.
func TestCheckReportsDiscardedFields(t *testing.T) {
	var data struct {
		Name wrappers.LenientNonEmptyString `json:"name"`
	}

	if err := json.Unmarshal([]byte(`{"name":""}`), &data); err != nil {
		t.Fatalf("Unmarshal() error = %v, want the lenient wrapper to swallow it", err)
	}

	err := wrappers.Check(&data)
	if err == nil {
		t.Fatal("Check() = nil, want the discarded field reported")
	}

	if !errors.Is(err, wrappers.ErrValue) {
		t.Errorf("errors.Is(err, ErrValue) = false, want true for a discarded value")
	}

	if errors.Is(err, wrappers.ErrMissing) {
		t.Error("errors.Is(err, ErrMissing) = true, want a discard reported as ErrValue")
	}
}

func TestCheckWalksNestedStructs(t *testing.T) {
	type inner struct {
		Code wrappers.NonEmptyString `json:"code"`
	}

	var data struct {
		Inner inner `json:"inner"`
	}

	if err := json.Unmarshal([]byte(`{"inner":{}}`), &data); err != nil {
		t.Fatalf("Unmarshal() error = %v", err)
	}

	err := wrappers.Check(&data)
	if err == nil {
		t.Fatal("Check() = nil, want the nested field reported")
	}

	if !strings.Contains(err.Error(), "Inner.Code") {
		t.Errorf("Check() = %v, want the dotted path Inner.Code", err)
	}
}

func TestCheckWalksSlices(t *testing.T) {
	type row struct {
		Code wrappers.NonEmptyString `json:"code"`
	}

	var data struct {
		Rows []row `json:"rows"`
	}

	if err := json.Unmarshal([]byte(`{"rows":[{"code":"a"},{}]}`), &data); err != nil {
		t.Fatalf("Unmarshal() error = %v", err)
	}

	err := wrappers.Check(&data)
	if err == nil {
		t.Fatal("Check() = nil, want the second row reported")
	}

	if !strings.Contains(err.Error(), "Rows[1].Code") {
		t.Errorf("Check() = %v, want the indexed path Rows[1].Code", err)
	}
}

// Pointer fields are still legal, and a nil one is the v1 panic waiting to happen.
func TestCheckReportsNilWrapperPointers(t *testing.T) {
	var data struct {
		Name *wrappers.NonEmptyString `json:"name"`
	}

	if err := json.Unmarshal([]byte(`{}`), &data); err != nil {
		t.Fatalf("Unmarshal() error = %v", err)
	}

	err := wrappers.Check(&data)
	if err == nil {
		t.Fatal("Check() = nil, want the nil pointer reported")
	}

	if !errors.Is(err, wrappers.ErrMissing) {
		t.Errorf("errors.Is(err, ErrMissing) = false, want true")
	}
}

func TestCheckAcceptsAValueAsWellAsAPointer(t *testing.T) {
	var data account
	if err := json.Unmarshal([]byte(`{"name":"a","balance":1}`), &data); err != nil {
		t.Fatalf("Unmarshal() error = %v", err)
	}

	if err := wrappers.Check(data); err != nil {
		t.Errorf("Check(value) = %v, want nil", err)
	}
}

func TestCheckIgnoresNonWrapperFields(t *testing.T) {
	var data struct {
		Plain  string
		Number int
		hidden string
	}

	if err := wrappers.Check(&data); err != nil {
		t.Errorf("Check() = %v, want nil for a struct with no wrappers", err)
	}

	_ = data.hidden
}

type node struct {
	Name wrappers.NonEmptyString
	Next *node
}

// Check follows pointers, so a structure that points back at itself must not
// send the walk around the loop until the stack runs out.
func TestCheckTerminatesOnCycles(t *testing.T) {
	t.Run("two node cycle", func(t *testing.T) {
		first := &node{}
		second := &node{}
		first.Next = second
		second.Next = first

		err := runWithTimeout(t, func() error { return wrappers.Check(first) })
		if err == nil {
			t.Error("Check() = nil, want both empty names reported")
		}
	})

	t.Run("self reference", func(t *testing.T) {
		self := &node{}
		self.Next = self

		if err := runWithTimeout(t, func() error { return wrappers.Check(self) }); err == nil {
			t.Error("Check() = nil, want the empty name reported")
		}
	})

	// A diamond is not a cycle. The same node reached down two paths is visited
	// once, so its problem is reported once rather than twice, but it must not be
	// dropped entirely just because the walk has seen that address before.
	t.Run("diamond reports the shared node", func(t *testing.T) {
		shared := &node{}

		data := struct {
			Left  *node
			Right *node
		}{Left: shared, Right: shared}

		err := runWithTimeout(t, func() error { return wrappers.Check(&data) })
		if err == nil {
			t.Fatal("Check() = nil, want the shared node's empty name reported")
		}

		if got := strings.Count(err.Error(), "absent from input"); got != 1 {
			t.Errorf("reported %d times, want exactly 1 for a shared node", got)
		}
	})
}

func runWithTimeout(t *testing.T, f func() error) error {
	t.Helper()

	type result struct{ err error }
	done := make(chan result, 1)

	go func() {
		done <- result{err: f()}
	}()

	select {
	case r := <-done:
		return r.err

	case <-time.After(5 * time.Second):
		t.Fatal("Check() did not return, the walk is not terminating")

		return nil
	}
}
