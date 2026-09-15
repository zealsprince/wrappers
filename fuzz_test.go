package wrappers_test

import (
	"encoding/json"
	"testing"

	wrappers "github.com/zealsprince/wrappers/v2"
	"github.com/zealsprince/wrappers/v2/regex"
)

// Arbitrary bytes reach UnmarshalJSON straight off the wire, so no input may
// panic. A rejection is a fine outcome, a crash is not.
func FuzzUnmarshalNeverPanics(f *testing.F) {
	for _, seed := range []string{
		`"a@b.com"`, `1`, `1.9`, `1e309`, `-1e309`, `true`, `null`, `""`,
		`"2026-03-06T12:30:00Z"`, `"1h30m"`, `{}`, `[]`, `"\ud800"`,
		`99999999999999999999999999`, `"DE"`, `"north"`,
	} {
		f.Add([]byte(seed))
	}

	f.Fuzz(func(t *testing.T, data []byte) {
		targets := []json.Unmarshaler{
			&wrappers.String{}, &wrappers.NonEmptyString{}, &wrappers.Int{},
			&wrappers.Float{}, &wrappers.Bool{}, &wrappers.Time{},
			&wrappers.TimeUnix{}, &wrappers.TimeISO8601{}, &wrappers.Duration{},
			&regex.Email{}, &regex.VIN{},
			&wrappers.LenientInt{}, &wrappers.LenientString{},
		}

		for _, target := range targets {
			// An error is expected for most input. A panic is the failure.
			_ = target.UnmarshalJSON(data)
		}
	})
}

// The same for Wrap, which takes already-decoded Go values.
func FuzzWrapNeverPanics(f *testing.F) {
	f.Add("abc", int64(0), 0.0, false)
	f.Add("", int64(-1), 1e308, true)

	f.Fuzz(func(t *testing.T, s string, i int64, fl float64, b bool) {
		for _, input := range []any{s, i, fl, b, []byte(s), nil} {
			var (
				str wrappers.String
				num wrappers.Int
				flt wrappers.Float
				bl  wrappers.Bool
				dur wrappers.Duration
			)

			_ = str.Wrap(input)
			_ = num.Wrap(input)
			_ = flt.Wrap(input)
			_ = bl.Wrap(input)
			_ = dur.Wrap(input)
		}
	})
}
