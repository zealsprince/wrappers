package wrappers

// The elementary wrappers. Each is an alias to an instantiated Wrapper, so the
// zero value is usable and `*String` satisfies json.Unmarshaler with no setup.

type stringRule struct{ StringParser }

func (stringRule) Name() Name { return "String" }

// NonEmptyString rejects the empty string. Plain String accepts it, which is the
// v1 behaviour reversed: v1 treated "" as nil and errored, with no way to opt out.
type nonEmptyStringRule struct{ StringParser }

func (nonEmptyStringRule) Name() Name { return "NonEmptyString" }

func (r nonEmptyStringRule) Validate(s string) error {
	if s == "" {
		return ValueErrorf(r.Name(), s, "must not be empty")
	}

	return nil
}

type intRule struct{ Int64Parser }

func (intRule) Name() Name { return "Int" }

type floatRule struct{ Float64Parser }

func (floatRule) Name() Name { return "Float" }

type boolRule struct{ BoolParser }

func (boolRule) Name() Name { return "Bool" }

type (
	String         = Wrapper[string, stringRule]
	NonEmptyString = Wrapper[string, nonEmptyStringRule]
	Int            = Wrapper[int64, intRule]
	Float          = Wrapper[float64, floatRule]
	Bool           = Wrapper[bool, boolRule]
)

type (
	LenientString         = Lenient[string, stringRule]
	LenientNonEmptyString = Lenient[string, nonEmptyStringRule]
	LenientInt            = Lenient[int64, intRule]
	LenientFloat          = Lenient[float64, floatRule]
	LenientBool           = Lenient[bool, boolRule]
)
