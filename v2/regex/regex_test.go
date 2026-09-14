package regex_test

import (
	"encoding/json"
	"errors"
	"testing"

	wrappers "github.com/zealsprince/wrappers/v2"
	"github.com/zealsprince/wrappers/v2/regex"
)

// Every field here is a plain value on a struct nobody touched first. v1's suite
// assigned an initialized wrapper to the field before unmarshalling, which no
// caller does, and that is precisely how four broken wrappers passed their tests.
type payload struct {
	Email    regex.Email    `json:"email"`
	Phone    regex.Phone    `json:"phone"`
	URL      regex.URL      `json:"url"`
	SepaIBAN regex.SepaIBAN `json:"sepaIban"`
	SepaBIC  regex.SepaBIC  `json:"sepaBic"`
	VIN      regex.VIN      `json:"vin"`
}

func TestAcceptsValidValues(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
		get   func(payload) string
	}{
		{"email", `{"email":"andrew@example.com"}`, "andrew@example.com", func(p payload) string { return p.Email.Get() }},
		{"phone", `{"phone":"+4915112345678"}`, "+4915112345678", func(p payload) string { return p.Phone.Get() }},
		{"url", `{"url":"https://example.com/a?b=c"}`, "https://example.com/a?b=c", func(p payload) string { return p.URL.Get() }},
		{"sepaIban", `{"sepaIban":"DE89370400440532013000"}`, "DE89370400440532013000", func(p payload) string { return p.SepaIBAN.Get() }},
		{"sepaBic", `{"sepaBic":"DEUTDEFF500"}`, "DEUTDEFF500", func(p payload) string { return p.SepaBIC.Get() }},
		{"vin", `{"vin":"1HGCM82633A004352"}`, "1HGCM82633A004352", func(p payload) string { return p.VIN.Get() }},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var data payload

			if err := json.Unmarshal([]byte(tt.input), &data); err != nil {
				t.Fatalf("Unmarshal(%s) error = %v, want nil", tt.input, err)
			}

			if got := tt.get(data); got != tt.want {
				t.Errorf("Get() = %q, want %q", got, tt.want)
			}
		})
	}
}

// A wrapper whose pattern failed to load would pass the test above by accepting
// everything, so the rejection path has to be pinned on the same code path.
func TestRejectsInvalidValues(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"email", `{"email":"not an email"}`},
		{"phone", `{"phone":"not a phone"}`},
		{"url", `{"url":"not a url"}`},
		{"sepaIban", `{"sepaIban":"nope"}`},
		{"sepaBic", `{"sepaBic":"nope"}`},
		{"vin", `{"vin":"TOOSHORT"}`},
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

			if !errors.Is(err, wrappers.ErrValue) {
				t.Errorf("errors.Is(err, ErrValue) = false, want true")
			}
		})
	}
}

// The lenient variants take the same input without failing and report it instead.
func TestLenientDiscardsInsteadOfFailing(t *testing.T) {
	var data struct {
		Email regex.LenientEmail `json:"email"`
	}

	if err := json.Unmarshal([]byte(`{"email":"not an email"}`), &data); err != nil {
		t.Fatalf("Unmarshal() error = %v, want nil from a lenient wrapper", err)
	}

	if data.Email.IsValid() {
		t.Error("IsValid() = true, want false for a rejected value")
	}

	if !data.Email.IsDiscarded() {
		t.Error("IsDiscarded() = false, want true")
	}

	if got := data.Email.Get(); got != "" {
		t.Errorf("Get() = %q, want the empty default", got)
	}
}

func TestRoundTrip(t *testing.T) {
	input := `{"email":"andrew@example.com","phone":"+4915112345678","url":"https://example.com","sepaIban":"DE89370400440532013000","sepaBic":"DEUTDEFF500","vin":"1HGCM82633A004352"}`

	var data payload
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
