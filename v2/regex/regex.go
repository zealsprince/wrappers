// Package regex provides string wrappers validated against a regular expression.
//
// Every pattern here compiles once, at package initialization, into a package
// level variable. There is no SetPattern call, no compiled-pattern cache and no
// Initialize step, which is deliberate: in v1 a derived wrapper had to remember
// to override UnmarshalJSON so its Initialize would run, and four of the six
// shipped wrappers did not, so they rejected every value with "Regex not set".
// A rule reaches its pattern through package scope instead, so there is nothing
// left to forget.
package regex

import (
	"regexp"

	wrappers "github.com/zealsprince/wrappers/v2"
)

// MustCompile panics on a malformed pattern at import time rather than on the
// first payload that reaches the wrapper.
var (
	emailPattern    = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@([a-zA-Z0-9]+\-?)*[a-zA-Z0-9]+\.[a-zA-Z]{2,}$`)
	phonePattern    = regexp.MustCompile(`^(?:\+?[1-9]\d{1,14}|0\d{1,14})$`)
	urlPattern      = regexp.MustCompile(`^(https?|ftp)://[^\s/$.?#].[^\s]*$`)
	sepaIBANPattern = regexp.MustCompile(`^[A-Z]{2}[0-9]{2}[A-Z0-9]{1,30}$`)
	sepaBICPattern  = regexp.MustCompile(`^[A-Z]{6}[A-Z0-9]{2}([A-Z0-9]{3})?$`)
	vinPattern      = regexp.MustCompile(`^[A-HJ-NPR-Z0-9]{17}$`)
)

// match is the shared body of every Validate below. Each rule is otherwise three
// lines, and six copies of the same guard is what the helper exists to avoid.
func match(re *regexp.Regexp, name wrappers.Name, value, expected string) error {
	if re.MatchString(value) {
		return nil
	}

	return wrappers.ErrValuef(name, value, "expected %s", expected)
}

type emailRule struct{ wrappers.StringParser }

func (emailRule) Name() wrappers.Name { return "Email" }

func (r emailRule) Validate(v string) error {
	return match(emailPattern, r.Name(), v, "an email address")
}

type phoneRule struct{ wrappers.StringParser }

func (phoneRule) Name() wrappers.Name { return "Phone" }

func (r phoneRule) Validate(v string) error {
	return match(phonePattern, r.Name(), v, "an E.164 phone number")
}

type urlRule struct{ wrappers.StringParser }

func (urlRule) Name() wrappers.Name { return "URL" }

func (r urlRule) Validate(v string) error {
	return match(urlPattern, r.Name(), v, "an http, https or ftp URL")
}

type sepaIBANRule struct{ wrappers.StringParser }

func (sepaIBANRule) Name() wrappers.Name { return "SepaIBAN" }

func (r sepaIBANRule) Validate(v string) error {
	return match(sepaIBANPattern, r.Name(), v, "a SEPA IBAN")
}

type sepaBICRule struct{ wrappers.StringParser }

func (sepaBICRule) Name() wrappers.Name { return "SepaBIC" }

func (r sepaBICRule) Validate(v string) error {
	return match(sepaBICPattern, r.Name(), v, "a SEPA BIC")
}

type vinRule struct{ wrappers.StringParser }

func (vinRule) Name() wrappers.Name { return "VIN" }

func (r vinRule) Validate(v string) error {
	return match(vinPattern, r.Name(), v, "a 17 character vehicle identification number")
}

// The strict wrappers. Invalid input fails the unmarshal.
type (
	Email    = wrappers.Wrapper[string, emailRule]
	Phone    = wrappers.Wrapper[string, phoneRule]
	URL      = wrappers.Wrapper[string, urlRule]
	SepaIBAN = wrappers.Wrapper[string, sepaIBANRule]
	SepaBIC  = wrappers.Wrapper[string, sepaBICRule]
	VIN      = wrappers.Wrapper[string, vinRule]
)

// The lenient wrappers. Invalid input is discarded and the caller checks IsValid.
type (
	LenientEmail    = wrappers.Lenient[string, emailRule]
	LenientPhone    = wrappers.Lenient[string, phoneRule]
	LenientURL      = wrappers.Lenient[string, urlRule]
	LenientSepaIBAN = wrappers.Lenient[string, sepaIBANRule]
	LenientSepaBIC  = wrappers.Lenient[string, sepaBICRule]
	LenientVIN      = wrappers.Lenient[string, vinRule]
)
