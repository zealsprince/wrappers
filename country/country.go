// Package country provides a wrapper for ISO 3166 country codes.
//
// It lives in its own package so that importing the core does not pull in the
// countries dataset. v1 put this in the root package, so every consumer linked
// it whether they used countries or not.
package country

import (
	"encoding/json"

	"github.com/biter777/countries"
	wrappers "github.com/zealsprince/wrappers/v2"
)

type rule struct{}

func (rule) Name() wrappers.Name { return "Country" }

// Parse accepts a country code value, an alpha-2 or alpha-3 code, an English
// country name, or an ISO 3166 numeric code. Lookup is case insensitive.
func (r rule) Parse(value any) (countries.CountryCode, error) {
	switch v := value.(type) {
	case countries.CountryCode:
		return v, nil

	case countries.Country:
		return v.Code, nil

	case string:
		return countries.ByName(v), nil

	case json.Number, int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64:
		numeric, err := wrappers.Int64Parser{}.Parse(v)
		if err != nil {
			return countries.Unknown, wrappers.ValueErrorf(r.Name(), value, "expected an ISO 3166 numeric code")
		}

		return countries.ByNumeric(int(numeric)), nil

	default:
		return countries.Unknown, wrappers.TypeError(r.Name(), value)
	}
}

// Validate checks the shape of the alpha-2 code rather than listing sentinels.
// The countries package carries pseudo-entries (None at 998, International at
// 999) that report IsValid as true and whose Alpha2 is not a two letter code, so
// validity alone is not enough. v1 only compared against Unknown, which let "XX"
// through and round-tripped it as "None". Every real ISO 3166 alpha-2 is exactly
// two uppercase letters, so that is what gets checked.
func (r rule) Validate(code countries.CountryCode) error {
	if !code.IsValid() || !isAlpha2(code.Alpha2()) {
		return wrappers.ValueErrorf(r.Name(), code.Alpha2(), "a known ISO 3166 country")
	}

	return nil
}

func isAlpha2(code string) bool {
	if len(code) != 2 {
		return false
	}

	for i := range len(code) {
		if code[i] < 'A' || code[i] > 'Z' {
			return false
		}
	}

	return true
}

// Unwrap emits the alpha-2 code. v1 emitted the English name, so a payload that
// arrived as "DE" marshalled back out as "Germany" and did not round trip.
func (rule) Unwrap(code countries.CountryCode) any { return code.Alpha2() }

type (
	Country        = wrappers.Wrapper[countries.CountryCode, rule]
	LenientCountry = wrappers.Lenient[countries.CountryCode, rule]
)
