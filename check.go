package wrappers

import (
	"errors"
	"reflect"
	"strconv"
)

const (
	// TagHeader is the struct tag Check reads. v1 declared this and never used it,
	// because a wrapper's UnmarshalJSON cannot see the tag on the field holding it.
	// Check can, because it walks the struct.
	TagHeader = "wrappers"

	// TagOptional marks a field that is allowed to be absent from the input.
	TagOptional = "optional"
)

// checkable is the part of a wrapper Check needs. Every Wrapper and Optional
// satisfies it through their value receivers.
type checkable interface {
	IsPresent() bool
	IsDiscarded() bool
	Name() Name
}

// Check reports fields that never received input or that were discarded. It
// closes the gap that neither v1 nor encoding/json covers: a field simply absent
// from a payload is never unmarshalled, so no amount of validation inside
// UnmarshalJSON will ever see it.
//
// Pass a struct or a pointer to one. Nested structs and slices of structs are
// walked. Fields tagged `wrappers:"optional"` may be absent but are still
// reported if present and discarded. Errors are joined, so one call reports every
// bad field rather than only the first.
func Check(value any) error {
	c := checker{seen: make(map[visit]bool)}

	c.walk(reflect.ValueOf(value), "")

	return errors.Join(c.problems...)
}

// visit identifies somewhere the walk has already been. A struct can point back
// at itself, directly or through a chain, and without this the walk would follow
// the cycle until the stack ran out.
type visit struct {
	pointer uintptr
	typ     reflect.Type
}

type checker struct {
	problems []error
	seen     map[visit]bool
}

func (c *checker) report(path string, cause error) {
	c.problems = append(c.problems, &FieldError{Field: path, Cause: cause})
}

func (c *checker) walk(v reflect.Value, path string) {
	if !v.IsValid() {
		return
	}

	switch v.Kind() {
	case reflect.Pointer, reflect.Interface:
		if v.IsNil() {
			// A nil wrapper pointer is the v1 panic waiting to happen. Report it
			// rather than letting the caller find it by dereferencing.
			if path != "" && implementsCheckable(v.Type()) {
				c.report(path, ErrMissing)
			}

			return
		}

		if v.Kind() == reflect.Pointer {
			here := visit{pointer: v.Pointer(), typ: v.Type()}
			if c.seen[here] {
				return
			}

			c.seen[here] = true
		}

		c.walk(v.Elem(), path)

	case reflect.Struct:
		// A wrapper is itself a struct, so it has to be tested before descending,
		// otherwise the walk would recurse into its unexported fields.
		if c.inspect(v, path, "") {
			return
		}

		c.walkFields(v, path)

	case reflect.Slice, reflect.Array:
		for i := range v.Len() {
			c.walk(v.Index(i), indexPath(path, i))
		}

	case reflect.Map:
		for _, key := range v.MapKeys() {
			c.walk(v.MapIndex(key), joinPath(path, key.String()))
		}
	}
}

func (c *checker) walkFields(v reflect.Value, path string) {
	t := v.Type()

	for i := range t.NumField() {
		field := t.Field(i)
		if !field.IsExported() {
			continue
		}

		child := joinPath(path, field.Name)
		tag := field.Tag.Get(TagHeader)

		// Test the field itself before recursing, so a wrapper field is inspected
		// with its tag in hand rather than as an anonymous struct.
		if c.inspect(v.Field(i), child, tag) {
			continue
		}

		if v.Field(i).Kind() == reflect.Pointer && v.Field(i).IsNil() && implementsCheckable(v.Field(i).Type()) {
			if tag != TagOptional {
				c.report(child, ErrMissing)
			}

			continue
		}

		c.walk(v.Field(i), child)
	}
}

// inspect reports a wrapper field and returns whether the value was a wrapper.
func (c *checker) inspect(v reflect.Value, path string, tag string) bool {
	w, ok := asCheckable(v)
	if !ok {
		return false
	}

	switch {
	case !w.IsPresent():
		if tag != TagOptional {
			c.report(path, &ValidationError{Wrapper: w.Name(), Reason: "absent from input", Cause: ErrMissing})
		}

	case w.IsDiscarded():
		c.report(path, &ValidationError{Wrapper: w.Name(), Reason: "value was discarded", Cause: ErrValue})
	}

	return true
}

func asCheckable(v reflect.Value) (checkable, bool) {
	if v.Kind() == reflect.Pointer {
		if v.IsNil() {
			return nil, false
		}

		v = v.Elem()
	}

	if v.Kind() != reflect.Struct || !v.CanInterface() {
		return nil, false
	}

	w, ok := v.Interface().(checkable)

	return w, ok
}

func implementsCheckable(t reflect.Type) bool {
	for t.Kind() == reflect.Pointer {
		t = t.Elem()
	}

	if t.Kind() != reflect.Struct {
		return false
	}

	return reflect.PointerTo(t).Implements(reflect.TypeOf((*checkable)(nil)).Elem())
}

func joinPath(base, name string) string {
	if base == "" {
		return name
	}

	return base + "." + name
}

func indexPath(base string, i int) string {
	return base + "[" + strconv.Itoa(i) + "]"
}
