package wrappers

import (
	"encoding/json"
	"strings"
	"time"
)

// Time accepts RFC3339 timestamps and serializes back to RFC3339.
//
// It deliberately does not accept bare numbers. An epoch value is ambiguous
// between seconds and milliseconds, and guessing from magnitude is the kind of
// silent wrong answer this library exists to prevent. Use TimeUnix, which says
// which unit it means.
type timeRule struct{}

func (timeRule) Name() Name { return "Time" }

func (r timeRule) Parse(value any) (time.Time, error) {
	switch v := value.(type) {
	case time.Time:
		return v, nil

	case string:
		parsed, err := time.Parse(time.RFC3339, v)
		if err != nil {
			return time.Time{}, ValueErrorf(r.Name(), v, "expected an RFC3339 timestamp")
		}

		return parsed, nil

	default:
		return time.Time{}, TypeError(r.Name(), value)
	}
}

func (timeRule) Validate(time.Time) error { return nil }

// Unwrap sends the value out as RFC3339 rather than as Go's time.Time shape.
func (timeRule) Unwrap(v time.Time) any { return v.Format(time.RFC3339) }

// TimeUnix accepts epoch seconds as a number and serializes back to a number.
type timeUnixRule struct{}

func (timeUnixRule) Name() Name { return "TimeUnix" }

func (r timeUnixRule) Parse(value any) (time.Time, error) {
	switch v := value.(type) {
	case time.Time:
		return v, nil

	default:
		seconds, err := Int64Parser{}.Parse(v)
		if err != nil {
			return time.Time{}, ValueErrorf(r.Name(), value, "expected epoch seconds")
		}

		return time.Unix(seconds, 0).UTC(), nil
	}
}

func (timeUnixRule) Validate(time.Time) error { return nil }

func (timeUnixRule) Unwrap(v time.Time) any { return v.Unix() }

// TimeISO8601 is Time with the loosest ISO 8601 spellings normalized first. Some
// upstreams send a +0000 offset or drop the zone designator entirely, and both
// are rejected by a strict RFC3339 parse.
type timeISO8601Rule struct{}

func (timeISO8601Rule) Name() Name { return "TimeISO8601" }

func (r timeISO8601Rule) Parse(value any) (time.Time, error) {
	v, ok := value.(string)
	if !ok {
		return timeRule{}.Parse(value)
	}

	normalized := normalizeISO8601(v)

	parsed, err := time.Parse(time.RFC3339, normalized)
	if err != nil {
		return time.Time{}, ValueErrorf(r.Name(), v, "expected an ISO 8601 timestamp")
	}

	return parsed, nil
}

func (timeISO8601Rule) Validate(time.Time) error { return nil }

func (timeISO8601Rule) Unwrap(v time.Time) any { return v.Format(time.RFC3339) }

func normalizeISO8601(v string) string {
	// A +0000 offset means UTC, which RFC3339 spells Z.
	if after, found := strings.CutSuffix(v, "+0000"); found {
		return after + "Z"
	}

	// A timestamp with no zone at all is treated as UTC, matching v1.
	if !strings.HasSuffix(v, "Z") && !hasNumericZone(v) {
		return v + "Z"
	}

	return v
}

// hasNumericZone reports a trailing +HH:MM or -HH:MM, which must be left alone.
func hasNumericZone(v string) bool {
	if len(v) < 6 {
		return false
	}

	tail := v[len(v)-6:]

	return (tail[0] == '+' || tail[0] == '-') && tail[3] == ':'
}

// Duration accepts Go duration strings such as "1h30m" and bare numbers, which
// are nanoseconds to match time.Duration's own unit. It serializes to the string
// form, so a round trip stays readable.
type durationRule struct{}

func (durationRule) Name() Name { return "Duration" }

func (r durationRule) Parse(value any) (time.Duration, error) {
	switch v := value.(type) {
	case time.Duration:
		return v, nil

	case string:
		parsed, err := time.ParseDuration(v)
		if err != nil {
			return 0, ValueErrorf(r.Name(), v, "expected a duration such as 1h30m")
		}

		return parsed, nil

	case json.Number, int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64:
		nanoseconds, err := Int64Parser{}.Parse(v)
		if err != nil {
			return 0, ValueErrorf(r.Name(), value, "expected whole nanoseconds")
		}

		return time.Duration(nanoseconds), nil

	default:
		return 0, TypeError(r.Name(), value)
	}
}

func (durationRule) Validate(time.Duration) error { return nil }

func (durationRule) Unwrap(v time.Duration) any { return v.String() }

type (
	Time        = Wrapper[time.Time, timeRule]
	TimeUnix    = Wrapper[time.Time, timeUnixRule]
	TimeISO8601 = Wrapper[time.Time, timeISO8601Rule]
	Duration    = Wrapper[time.Duration, durationRule]
)

type (
	LenientTime        = Lenient[time.Time, timeRule]
	LenientTimeUnix    = Lenient[time.Time, timeUnixRule]
	LenientTimeISO8601 = Lenient[time.Time, timeISO8601Rule]
	LenientDuration    = Lenient[time.Duration, durationRule]
)
