# Wrappers

![Go](https://img.shields.io/badge/Go-1.24.1%2B-blue.svg)
![License](https://img.shields.io/badge/License-MIT-green.svg)
[![Build Go application](https://github.com/zealsprince/wrappers/actions/workflows/go-build.yml/badge.svg)](https://github.com/zealsprince/wrappers/actions/workflows/go-build.yml)
[![Test Go application](https://github.com/zealsprince/wrappers/actions/workflows/go-test.yml/badge.svg)](https://github.com/zealsprince/wrappers/actions/workflows/go-test.yml)

Wrappers are struct field types that validate and convert their own value while
`encoding/json` is still unmarshalling. A field typed `regex.Email` accepts an email
address and rejects anything else, and a field typed `wrappers.Int` takes `1`, `"1"`
or `1.0` and gives you an `int64` either way.

The conversion is the part a struct tag validator can't do for you. Those run after
`encoding/json` has already failed on the type mismatch, so an upstream that changes
a field from a number to a string breaks you before validation gets a turn.

## The problem

Unmarshal, then check every field by hand:

```go
type User struct {
    ID    string `json:"id"`
    Name  string `json:"name"`
    Email string `json:"email"`
}

var user User
if err := json.Unmarshal(payload, &user); err != nil {
    return err
}

if !isValidEmail(user.Email) {
    return errors.New("invalid email address")
}

id, err := strconv.ParseInt(user.ID, 10, 64)
if err != nil || id <= 0 {
    return errors.New("invalid user ID")
}

if user.Name == "" {
    return errors.New("invalid user name")
}
```

Or declare what each field accepts and let the unmarshal enforce it:

```go
type User struct {
    ID    wrappers.Int            `json:"id"`
    Name  wrappers.NonEmptyString `json:"name"`
    Email regex.Email             `json:"email"`
}

var user User
if err := json.Unmarshal(payload, &user); err != nil {
    return err // a bad email or an unparseable ID arrives here
}

if err := wrappers.Check(&user); err != nil {
    return err // a field the payload never included arrives here
}
```

## Install

```
go get github.com/zealsprince/wrappers/v2
```

```go
import (
    wrappers "github.com/zealsprince/wrappers/v2"
    "github.com/zealsprince/wrappers/v2/regex"
)
```

## How it works

A wrapper has three operations.

**Wrapping** converts an input value to the wrapper's type and runs its rule against
the result. `wrappers.Int` accepts every integer width, a float with no fractional
part, and a numeric string. It rejects `1.9` and anything that overflows `int64`
rather than truncating. Wrapping happens automatically during unmarshalling, and you
can call `Wrap` directly.

**Unwrapping** reads the value back. `Get` returns the wrapped Go type, so
`wrappers.Time.Get()` is a `time.Time`. `Unwrap` returns the serialized form, which
for `Time` is an RFC3339 string.

**Discarding** is what a wrapper does with a value that fails its rule. It keeps the
zero value and records the rejection, which you read with `IsDiscarded`. A later
successful `Wrap` clears it.

The zero wrapper is usable. You don't need pointer fields and there's nothing to
initialize. Marshalling doesn't mutate, so marshalling the same struct from several
goroutines is safe.

### Absent fields

A field missing from the payload is never handed to `UnmarshalJSON`, so no validation
inside a wrapper can see it. `Check` walks a struct and reports fields that never
received input, along with any that were discarded:

```go
var user User
json.Unmarshal([]byte(`{"name":"Andrew"}`), &user)

err := wrappers.Check(&user)
// ID: Int: absent from input
// Email: Email: absent from input
```

It descends into nested structs, slices and maps, and joins every problem into one
error so a bad payload takes one round trip to diagnose rather than one per field.
Tag a field `wrappers:"optional"` to allow its absence.

### Strict and lenient

Each wrapper has a lenient counterpart that discards bad input instead of failing the
unmarshal. Use it when you want the rest of the payload even if one field is garbage:

```go
type Contact struct {
    Email        regex.Email        `json:"email"`
    BillingEmail regex.LenientEmail `json:"billingEmail"`
}

var c Contact
err := json.Unmarshal(payload, &c)
// A bad "email" fails here. A bad "billingEmail" does not.

if c.BillingEmail.IsValid() {
    send(c.BillingEmail.Get())
}
```

### Omitting fields

Tag a field `omitzero` and absent or discarded values drop out of the marshalled
output:

```go
type Contact struct {
    Email regex.Email `json:"email"`
    Phone regex.Phone `json:"phone,omitzero"`
}
```

## Packages

| Package | Provides |
| --- | --- |
| `wrappers` | `String`, `NonEmptyString`, `Int`, `Float`, `Bool`, `Time`, `TimeUnix`, `TimeISO8601`, `Duration` |
| `wrappers/regex` | `Email`, `Phone`, `URL`, `SepaIBAN`, `SepaBIC`, `VIN` |
| `wrappers/enum` | `Wrapper[T, V]` for string-backed enumerated types, plus `CardinalDirections` as a worked example |
| `wrappers/country` | `Country`, ISO 3166 codes via `github.com/biter777/countries` |

Country is separate so that importing the core doesn't link the countries dataset.

`Time` accepts RFC3339 and refuses bare numbers, because an epoch value is ambiguous
between seconds and milliseconds. `TimeUnix` is the explicit way to ask for seconds.
`Duration` reads bare numbers as nanoseconds, matching `time.Duration`.

## Writing a rule

Validation is a type parameter. A rule is an empty struct with a name, a parser and a
check. The wrapper finds it through its own type, so there's no registration or setup
step:

```go
package myrules

import (
    "regexp"

    wrappers "github.com/zealsprince/wrappers/v2"
)

var slugPattern = regexp.MustCompile(`^[a-z0-9]+(-[a-z0-9]+)*$`)

type slugRule struct{ wrappers.StringParser }

func (slugRule) Name() wrappers.Name { return "Slug" }

func (r slugRule) Validate(v string) error {
    if !slugPattern.MatchString(v) {
        return wrappers.ValueErrorf(r.Name(), v, "expected a lowercase slug")
    }

    return nil
}

type (
    Slug        = wrappers.Wrapper[string, slugRule]
    LenientSlug = wrappers.Lenient[string, slugRule]
)
```

Embedding `StringParser` supplies the conversion from loose input. `Int64Parser`,
`Float64Parser` and `BoolParser` do the same for the other underlying types. A rule
that leaves one out doesn't satisfy `Rule[T]`, so the package won't compile.

For a type the parsers don't cover, write `Parse` yourself and optionally an
`Unwrap(T) any` to control the serialized form. `country` and the time rules are
worked examples.

## Migrating from v1

v1 stays available at `github.com/zealsprince/wrappers` and its tags are unaffected.

| v1 | v2 |
| --- | --- |
| `*wrappers.WrapperString` | `wrappers.String` (no pointer) |
| `wrappers.New[*T]()` | the zero value |
| `wrappers.NewWithValue[*T](v)` | `wrappers.Of[T](v)` |
| `wrappers.NewWithValueDiscard[*T](v)` | `wrappers.OfDiscard[T](v)` |
| `Wrap(v, false)` / `Wrap(v, true)` | `Wrap(v)` / `WrapDiscard(v)` |
| `Discarder[W]` with `.Proxy` | `Lenient[T, R]`, methods promoted |
| `Initialize` and `SetPattern` | nothing, rules are reachable from the type |
| nilling a field to omit it | `omitzero` |

Two behaviours changed. `Country` unwraps to the alpha-2 code, so `"DE"` round trips
instead of coming back as `"Germany"`. Integer wrappers reject fractional and
overflowing JSON numbers instead of truncating them.

## Motivations

This came out of building the [Data Layer](https://data-layer.com/) at
[Savages Corp](https://github.com/savages-corp). Writing integrations meant matching
external API structures and then validating every field after unmarshalling, because
somebody on the other end, often in a customer-managed environment, would change a
field's type or format without telling anyone. Catching that is a cat-and-mouse game,
and the code you write to play it is the same twenty lines in every handler.

Moving the check into the field type puts it where the shape is declared. The struct
says what it accepts, and a payload that doesn't match produces an error naming the
field instead of a panic three functions later.

## License

Released under the [MIT License](LICENSE).
