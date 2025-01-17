package regex

import (
	"fmt"
	"regexp"
	"sync"

	"github.com/zealsprince/wrappers"
)

// Note: This wrapper should not be used directly, especially not as a struct field type. It is a reusable wrapper that should be embedded in a struct.
// See the implementations within this sub-package for examples of how to use this wrapper.

// WrapperRegex is a reusable string wrapper that validates its value against a provided regex.
type WrapperRegex struct {
	wrappers.Wrapper[string, string]
	name    wrappers.Name
	regex   *regexp.Regexp
	pattern string
}

var _ wrappers.WrapperProvider = (*WrapperRegex)(nil) // Ensure that WrapperRegex implements WrapperProvider.

var (
	regexCache   = make(map[string]*regexp.Regexp)
	regexCacheMu sync.RWMutex
)

// NewWrapperRegexRegex creates a new WrapperRegex with the given name and regex pattern.
// It caches the compiled regex to optimize performance.
func (wrapper *WrapperRegex) SetPattern(name wrappers.Name, pattern string) error {
	regexCacheMu.RLock()
	regex, exists := regexCache[pattern]
	regexCacheMu.RUnlock()

	if !exists {
		var err error
		regex, err = regexp.Compile(pattern)
		if err != nil {
			return fmt.Errorf("failed to compile regex for %s: %w", name, err)
		}
		wrapper.regex = regex

		regexCacheMu.Lock()
		regexCache[pattern] = regex
		regexCacheMu.Unlock()
	}

	// Assign the values to the wrapper
	wrapper.name = name
	wrapper.pattern = pattern
	wrapper.regex = regex

	return nil
}

// Get returns the wrapped string.
func (wrapper *WrapperRegex) Get() string {
	return wrapper.Value
}

// GetAny returns the wrapped string.
func (wrapper *WrapperRegex) GetAny() any {
	return wrapper.Get()
}

// Wrap validates and wraps the input value using the regex.
func (wrapper *WrapperRegex) Wrap(value any, discard bool) error {
	var str string
	switch v := value.(type) {
	case nil:
		wrapper.Discard()
		return nil

	case wrappers.WrapperProvider:
		if v.IsDiscarded() {
			wrapper.Discard()
			return nil
		}

		return wrapper.Wrap(v.UnwrapAny(), discard)

	case string:
		str = v

	default:
		wrapper.Discard()
		if !discard {
			return wrappers.ErrorType(wrapper.name, value)
		}
		return nil
	}

	if wrapper.regex == nil {
		return fmt.Errorf("Regex not set - if you are embedding this wrapper, make sure the implementation calls SetPattern during the Initialize method and initializes during UnmarshalJSON")
	}

	if !wrapper.regex.MatchString(str) {
		wrapper.Discard()
		if !discard {
			return wrappers.ErrorValue(wrapper.name, str, wrapper.pattern)
		}
		return nil
	}

	wrapper.Value = str

	return nil
}

// Unwrap returns the wrapped string pointer if not discarded.
func (wrapper *WrapperRegex) Unwrap() string {
	if wrapper.IsDiscarded() {
		return ""
	}

	return wrapper.Value
}

func (wrapper *WrapperRegex) UnwrapAny() any {
	return wrapper.Unwrap()
}

// MarshalJSON marshals the wrapped value to JSON, handling discards.
func (wrapper *WrapperRegex) MarshalJSON() ([]byte, error) {
	return wrappers.MarshalJSON(wrapper)
}

// UnmarshalJSON unmarshals JSON data into the wrapper.
func (wrapper *WrapperRegex) UnmarshalJSON(data []byte) error {
	if wrapper == nil {
		return fmt.Errorf("unmarshal into nil wrapper")
	}

	return wrappers.UnmarshalJSON(data, wrapper)
}
