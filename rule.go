package wrappers

import (
	"encoding/json"
	"fmt"
	"math"
	"strconv"
)

// The parsers below are the shared coercion table. A rule embeds the one matching
// its underlying type and supplies Name plus, optionally, Validate. Embedding a
// parser cannot be forgotten: without it the rule does not satisfy Rule[T] and
// the package does not compile.

// StringParser coerces input to string. Numbers and bools are formatted rather
// than rejected, since upstreams flip a field between "42" and 42 freely.
type StringParser struct{}

func (StringParser) Parse(value any) (string, error) {
	switch v := value.(type) {
	case string:
		return v, nil

	case json.Number:
		return v.String(), nil

	case bool:
		return strconv.FormatBool(v), nil

	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		return fmt.Sprintf("%d", v), nil

	case float32:
		return strconv.FormatFloat(float64(v), 'f', -1, 32), nil

	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64), nil

	case []byte:
		return string(v), nil

	default:
		return "", errType("", value)
	}
}

func (StringParser) Validate(string) error { return nil }

// Int64Parser coerces input to int64. Fractional and out-of-range input is
// rejected rather than truncated, which is the v1 behaviour this replaces.
type Int64Parser struct{}

func (Int64Parser) Parse(value any) (int64, error) {
	switch v := value.(type) {
	case int64:
		return v, nil

	case int:
		return int64(v), nil

	case int8:
		return int64(v), nil

	case int16:
		return int64(v), nil

	case int32:
		return int64(v), nil

	case uint64:
		if v > math.MaxInt64 {
			return 0, ValueErrorf("", v, "exceeds int64")
		}

		return int64(v), nil

	case uint:
		return Int64Parser{}.Parse(uint64(v))

	case uint8:
		return int64(v), nil

	case uint16:
		return int64(v), nil

	case uint32:
		return int64(v), nil

	case json.Number:
		return parseIntString(v.String())

	case string:
		return parseIntString(v)

	case float32:
		return floatToInt(float64(v))

	case float64:
		return floatToInt(v)

	default:
		return 0, errType("", value)
	}
}

func (Int64Parser) Validate(int64) error { return nil }

func parseIntString(s string) (int64, error) {
	n, err := strconv.ParseInt(s, 10, 64)
	if err == nil {
		return n, nil
	}

	// A JSON number like 1.0 is an integer in every sense that matters, so fall
	// back to the float path rather than rejecting it outright.
	f, ferr := strconv.ParseFloat(s, 64)
	if ferr != nil {
		return 0, ValueErrorf("", s, "not an integer")
	}

	return floatToInt(f)
}

func floatToInt(f float64) (int64, error) {
	if math.IsNaN(f) || math.IsInf(f, 0) {
		return 0, ValueErrorf("", f, "not a finite number")
	}

	if f != math.Trunc(f) {
		return 0, ValueErrorf("", f, "has a fractional part")
	}

	// The bounds are checked before conversion because float-to-int conversion of
	// an out-of-range value is undefined in Go, not merely lossy.
	if f < math.MinInt64 || f >= math.MaxInt64 {
		return 0, ValueErrorf("", f, "exceeds int64")
	}

	return int64(f), nil
}

// Float64Parser coerces input to float64.
type Float64Parser struct{}

func (Float64Parser) Parse(value any) (float64, error) {
	switch v := value.(type) {
	case float64:
		return v, nil

	case float32:
		return float64(v), nil

	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		return strconv.ParseFloat(fmt.Sprintf("%d", v), 64)

	case json.Number:
		return parseFloatString(v.String())

	case string:
		return parseFloatString(v)

	default:
		return 0, errType("", value)
	}
}

func (Float64Parser) Validate(float64) error { return nil }

func parseFloatString(s string) (float64, error) {
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0, ValueErrorf("", s, "not a number")
	}

	if math.IsNaN(f) || math.IsInf(f, 0) {
		return 0, ValueErrorf("", s, "not a finite number")
	}

	return f, nil
}

// BoolParser coerces input to bool, accepting the string and numeric spellings
// that show up in query strings and loosely typed payloads.
type BoolParser struct{}

func (BoolParser) Parse(value any) (bool, error) {
	switch v := value.(type) {
	case bool:
		return v, nil

	case string:
		b, err := strconv.ParseBool(v)
		if err != nil {
			return false, ValueErrorf("", v, "not a boolean")
		}

		return b, nil

	case json.Number:
		n, err := parseIntString(v.String())
		if err != nil {
			return false, err
		}

		return n != 0, nil

	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64:
		n, err := Int64Parser{}.Parse(v)
		if err != nil {
			return false, err
		}

		return n != 0, nil

	default:
		return false, errType("", value)
	}
}

func (BoolParser) Validate(bool) error { return nil }
