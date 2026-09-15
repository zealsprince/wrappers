package country_test

import (
	"encoding/json"
	"testing"

	"github.com/biter777/countries"
	wrappers "github.com/zealsprince/wrappers/v2"
	"github.com/zealsprince/wrappers/v2/country"
)

func TestAcceptsEverySpelling(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"alpha-2", `"DE"`},
		{"alpha-3", `"DEU"`},
		{"english name", `"Germany"`},
		{"lowercase", `"de"`},
		{"numeric code", `276`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var w country.Country

			if err := json.Unmarshal([]byte(tt.input), &w); err != nil {
				t.Fatalf("Unmarshal(%s) error = %v, want nil", tt.input, err)
			}

			if got := w.Get(); got != countries.DEU {
				t.Errorf("Get() = %v, want %v", got, countries.DEU)
			}
		})
	}
}

// v1 emitted the English name, so a payload that arrived as "DE" marshalled back
// out as "Germany" and did not round trip.
func TestRoundTripsAsAlpha2(t *testing.T) {
	input := `"DE"`

	var w country.Country
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

// countries.ByName("XX") returns the None sentinel, which reports IsValid as
// true. v1 only compared against Unknown, so "XX" was accepted and unwrapped
// as "None".
func TestRejectsTheNoneSentinel(t *testing.T) {
	var w country.Country

	err := json.Unmarshal([]byte(`"XX"`), &w)
	if err == nil {
		t.Fatalf("Unmarshal(\"XX\") error = nil, want a rejection (got %q)", w.Get().Alpha2())
	}

	if !w.IsDiscarded() {
		t.Error("IsDiscarded() = false, want true")
	}
}

func TestRejectsUnknownInput(t *testing.T) {
	// 998 is the None sentinel and 999 is International. Both report IsValid as
	// true, and neither is a country.
	for _, input := range []string{`"ZZ"`, `"Banana"`, `""`, `true`, `998`, `999`, `0`, `12345`} {
		t.Run(input, func(t *testing.T) {
			var w country.Country

			if err := json.Unmarshal([]byte(input), &w); err == nil {
				t.Errorf("Unmarshal(%s) error = nil, want a rejection (got %q)", input, w.Get().Alpha2())
			}
		})
	}
}

func TestLenientDiscardsInsteadOfFailing(t *testing.T) {
	var data struct {
		Country country.LenientCountry `json:"country"`
	}

	if err := json.Unmarshal([]byte(`{"country":"XX"}`), &data); err != nil {
		t.Fatalf("Unmarshal() error = %v, want nil from a lenient wrapper", err)
	}

	if data.Country.IsValid() {
		t.Error("IsValid() = true, want false")
	}
}

func TestCheckSeesCountryFields(t *testing.T) {
	var data struct {
		Country country.Country `json:"country"`
	}

	if err := json.Unmarshal([]byte(`{}`), &data); err != nil {
		t.Fatalf("Unmarshal() error = %v", err)
	}

	if err := wrappers.Check(&data); err == nil {
		t.Error("Check() = nil, want the absent country reported")
	}
}

func TestAcceptsCountryStructAndRejectsOtherTypes(t *testing.T) {
	t.Run("countries.Country value", func(t *testing.T) {
		var w country.Country

		if err := w.Wrap(*countries.DEU.Info()); err != nil {
			t.Fatalf("Wrap(countries.Country) error = %v, want nil", err)
		}

		if got := w.Get(); got != countries.DEU {
			t.Errorf("Get() = %v, want %v", got, countries.DEU)
		}
	})

	t.Run("countries.CountryCode value", func(t *testing.T) {
		var w country.Country

		if err := w.Wrap(countries.DEU); err != nil {
			t.Fatalf("Wrap(CountryCode) error = %v, want nil", err)
		}
	})

	t.Run("unconvertible type", func(t *testing.T) {
		var w country.Country

		if err := w.Wrap(map[string]int{}); err == nil {
			t.Error("Wrap(map) error = nil, want a type error")
		}
	})

	t.Run("a numeric code that does not fit an int64", func(t *testing.T) {
		var w country.Country

		if err := w.Wrap(1.5); err == nil {
			t.Error("Wrap(1.5) error = nil, want a rejection")
		}
	})
}

// isAlpha2 has to reject on length and on each character, not just length.
func TestRejectsCodesThatAreNotTwoUppercaseLetters(t *testing.T) {
	var w country.Country

	// 999 is International, whose Alpha2 is a long word rather than a code.
	if err := w.Wrap(999); err == nil {
		t.Error("Wrap(999) error = nil, want International rejected")
	}

	// 998 is None, whose Alpha2 is "None": four characters, mixed case.
	if err := w.Wrap(998); err == nil {
		t.Error("Wrap(998) error = nil, want None rejected")
	}
}
