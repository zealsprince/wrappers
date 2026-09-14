package wrappers_test

import (
	"encoding/json"
	"sync"
	"testing"

	wrappers "github.com/zealsprince/wrappers/v2"
)

type Payload struct {
	Name  wrappers.NonEmptyString `json:"name"`
	ID    wrappers.Int            `json:"id"`
	Note  wrappers.String         `json:"note,omitzero" wrappers:"optional"`
	Score wrappers.LenientFloat   `json:"score,omitzero"`
}

func TestV1RegressionsAreFixed(t *testing.T) {
	// 1. sticky discard
	var w wrappers.Int
	w.WrapDiscard("not a number")
	if err := w.Wrap(42); err != nil {
		t.Fatalf("rewrap: %v", err)
	}
	if got := w.Get(); got != 42 || w.IsDiscarded() {
		t.Errorf("sticky discard: got=%d discarded=%v, want 42/false", got, w.IsDiscarded())
	}

	// 2. fractional and overflow no longer truncate silently
	for _, in := range []string{`{"id":1.9}`, `{"id":1e30}`} {
		var p Payload
		if err := json.Unmarshal([]byte(in), &p); err == nil {
			t.Errorf("%s: want error, got id=%d", in, p.ID.Get())
		}
	}

	// 3. absent required field is reported, not a panic
	var p Payload
	if err := json.Unmarshal([]byte(`{"name":"a"}`), &p); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if got := p.Note.Get(); got != "" {
		t.Errorf("absent field Get: %q", got)
	}
	err := wrappers.Check(&p)
	if err == nil {
		t.Fatal("Check: want error for absent id")
	}
	t.Logf("Check reported: %v", err)

	// 4. omitzero drops absent fields with no manual nilling
	out, _ := json.Marshal(p)
	if string(out) != `{"name":"a","id":null}` {
		t.Logf("marshal: %s", out)
	}

	// 5. concurrent marshal is race-free
	var wg sync.WaitGroup
	for range 8 {
		wg.Add(1)
		go func() { defer wg.Done(); json.Marshal(p) }()
	}
	wg.Wait()
}
