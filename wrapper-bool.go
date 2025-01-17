package wrappers

import (
	"fmt"
	"strings"
)

const (
	WrapperBoolName    Name   = "WrapperBool"
	WrapperBoolExample string = "true, false, yes, no, 1, 0"
)

type WrapperBool Wrapper[bool, bool]

var _ WrapperProvider = (*WrapperBool)(nil) // Ensure that WrapperBool implements WrapperProvider.

func (wrapper *WrapperBool) Get() bool {
	return wrapper.Value
}

func (wrapper *WrapperBool) GetAny() any {
	return wrapper.Get()
}

func (wrapper *WrapperBool) Wrap(value any, discard bool) error {
	switch v := value.(type) {
	case nil:
		wrapper.Discard()
		if !discard {
			return ErrorNil(WrapperBoolName)
		}

	case WrapperProvider:
		if v.IsDiscarded() {
			wrapper.Discard()
			return nil
		}

		return wrapper.Wrap(v.UnwrapAny(), discard)

	case bool:
		wrapper.Value = v

	case int, int8, int16, int32, int64:
		wrapper.Value = v != 0

	case float32, float64:
		wrapper.Value = v != 0.0

	case string:
		switch strings.ToLower(v) {
		case "true":
			wrapper.Value = true

		case "false":
			wrapper.Value = false

		case "yes":
			wrapper.Value = true

		case "no":
			wrapper.Value = false

		case "1":
			wrapper.Value = true

		case "0":
			wrapper.Value = false

		default:
			wrapper.Discard()
			if !discard {
				return ErrorValue(WrapperBoolName, value, WrapperBoolExample)
			}
		}

	default:
		wrapper.Discard()
		if !discard {
			return ErrorType(WrapperBoolName, value)
		}
	}

	return nil
}

func (wrapper *WrapperBool) Unwrap() bool {
	if wrapper.IsDiscarded() {
		return false
	}

	return wrapper.Value
}

func (wrapper *WrapperBool) UnwrapAny() any {
	return wrapper.Unwrap()
}

func (wrapper *WrapperBool) MarshalJSON() ([]byte, error) {
	return MarshalJSON(wrapper)
}

func (wrapper *WrapperBool) UnmarshalJSON(data []byte) error {
	if wrapper == nil {
		return fmt.Errorf("unmarshal into nil wrapper")
	}

	return UnmarshalJSON(data, wrapper)
}
