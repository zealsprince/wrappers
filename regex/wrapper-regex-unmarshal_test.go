package regex

import (
	"encoding/json"
	"testing"
)

// Every wrapper in this package has to survive being unmarshalled into a struct
// that nobody touched first. The rest of the suite in here assigns an initialized
// wrapper to the field before calling json.Unmarshal, which no caller ever does,
// and that is exactly how the missing UnmarshalJSON overrides on Url, SepaIban,
// SepaBic and Vin went unnoticed: their pattern was never set, so every value came
// back as "Regex not set".
func TestUnmarshalIntoZeroStruct(t *testing.T) {
	type payload struct {
		Email    *WrapperRegexEmail    `json:"email"`
		Phone    *WrapperRegexPhone    `json:"phone"`
		Url      *WrapperRegexUrl      `json:"url"`
		SepaIban *WrapperRegexSepaIban `json:"sepaIban"`
		SepaBic  *WrapperRegexSepaBic  `json:"sepaBic"`
		Vin      *WrapperRegexVin      `json:"vin"`
	}

	tests := []struct {
		name  string
		input string
		want  string
		get   func(payload) string
	}{
		{"email", `{"email":"andrew@example.com"}`, "andrew@example.com", func(p payload) string { return p.Email.Unwrap() }},
		{"phone", `{"phone":"+4915112345678"}`, "+4915112345678", func(p payload) string { return p.Phone.Unwrap() }},
		{"url", `{"url":"https://example.com/a?b=c"}`, "https://example.com/a?b=c", func(p payload) string { return p.Url.Unwrap() }},
		{"sepaIban", `{"sepaIban":"DE89370400440532013000"}`, "DE89370400440532013000", func(p payload) string { return p.SepaIban.Unwrap() }},
		{"sepaBic", `{"sepaBic":"DEUTDEFF500"}`, "DEUTDEFF500", func(p payload) string { return p.SepaBic.Unwrap() }},
		{"vin", `{"vin":"1HGCM82633A004352"}`, "1HGCM82633A004352", func(p payload) string { return p.Vin.Unwrap() }},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var data payload

			if err := json.Unmarshal([]byte(tt.input), &data); err != nil {
				t.Fatalf("Unmarshal(%s) error = %v, want nil", tt.input, err)
			}

			if got := tt.get(data); got != tt.want {
				t.Errorf("Unwrap() = %q, want %q", got, tt.want)
			}
		})
	}
}

// The flip side: an invalid value still has to be rejected on the same path. A
// wrapper whose pattern silently failed to load would pass the test above if it
// accepted everything, so pin the rejection too.
func TestUnmarshalIntoZeroStructRejects(t *testing.T) {
	type payload struct {
		Email    *WrapperRegexEmail    `json:"email"`
		Phone    *WrapperRegexPhone    `json:"phone"`
		Url      *WrapperRegexUrl      `json:"url"`
		SepaIban *WrapperRegexSepaIban `json:"sepaIban"`
		SepaBic  *WrapperRegexSepaBic  `json:"sepaBic"`
		Vin      *WrapperRegexVin      `json:"vin"`
	}

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

			if err := json.Unmarshal([]byte(tt.input), &data); err == nil {
				t.Fatalf("Unmarshal(%s) error = nil, want a validation error", tt.input)
			}
		})
	}
}
