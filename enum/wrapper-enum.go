package enum

import (
	"fmt"

	"github.com/zealsprince/wrappers"
)

// Note: This wrapper should not be used directly, especially not as a struct field type. It is a reusable wrapper that should be embedded in a struct.
// See the implementations within this sub-package for examples of how to use this wrapper.

// WrapperEnum is a reusable string wrapper that validates its value against a provided list of valid values.
type WrapperEnum[T ~string] struct {
	wrappers.Wrapper[T, string]
	name        wrappers.Name
	validValues []T
}

var _ wrappers.WrapperProvider = (*WrapperEnum[string])(nil) // Ensure that WrapperEnum implements WrapperProvider.

func (wrapper *WrapperEnum[T]) SetValidValues(values []T) {
	wrapper.validValues = values
}

func (wrapper *WrapperEnum[T]) Get() *T {
	return &wrapper.Value
}

func (wrapper *WrapperEnum[T]) GetAny() any {
	return wrapper.Get()
}

func (wrapper *WrapperEnum[T]) validateAndSet(value T) bool {
	for _, validValue := range wrapper.validValues {
		if value == validValue {
			wrapper.Value = value
			return true
		}
	}

	wrapper.Discard()
	return false
}

func (wrapper *WrapperEnum[T]) Wrap(value any, discard bool) error {
	switch v := value.(type) {
	case nil:
		wrapper.Discard()

	case wrappers.WrapperProvider:
		if v.IsDiscarded() {
			wrapper.Discard()
			return nil
		}

		return wrapper.Wrap(v.UnwrapAny(), discard)

	case T:
		if ok := wrapper.validateAndSet(T(v)); !ok {
			wrapper.Discard()
			if !discard {
				return fmt.Errorf("invalid value %v, valid values are %+v", v, wrapper.validValues)
			}
		}

	case string:
		if v == "" {
			wrapper.Discard()
			return nil
		}

		if ok := wrapper.validateAndSet(T(v)); !ok {
			wrapper.Discard()
			if !discard {
				return fmt.Errorf("invalid value %v, valid values are %+v", v, wrapper.validValues)
			}
		}

	default:
		wrapper.Discard()
		if !discard {
			return wrappers.ErrorType(wrapper.name, value)
		}
	}

	return nil
}

func (wrapper *WrapperEnum[T]) Unwrap() string {
	if wrapper.IsDiscarded() {
		return ""
	}

	return string(wrapper.Value)
}

func (wrapper *WrapperEnum[T]) UnwrapAny() any {
	return wrapper.Unwrap()
}

func (wrapper *WrapperEnum[T]) MarshalJSON() ([]byte, error) {
	return wrappers.MarshalJSON(wrapper)
}

func (wrapper *WrapperEnum[T]) UnmarshalJSON(data []byte) error {
	if wrapper == nil {
		return fmt.Errorf("unmarshal into nil wrapper")
	}

	return wrappers.UnmarshalJSON(data, wrapper)
}
