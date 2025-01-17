package wrappers

import (
	"fmt"
)

const (
	WrapperStringName Name = "WrapperString"
)

type WrapperString Wrapper[string, string]

var _ WrapperProvider = (*WrapperString)(nil) // Ensure that WrapperString implements WrapperProvider.

func (wrapper *WrapperString) Get() string {
	return wrapper.Value
}

func (wrapper *WrapperString) GetAny() any {
	return wrapper.Get()
}

func (wrapper *WrapperString) Wrap(value any, discard bool) error {
	switch v := value.(type) {
	case nil:
		wrapper.Discard()

	case WrapperProvider:
		if v.IsDiscarded() {
			wrapper.Discard()
			return nil
		}

		return wrapper.Wrap(v.UnwrapAny(), discard)

	case int, int8, int16, int32, int64, float32, float64:
		wrapper.Value = fmt.Sprintf("%v", v)

	case string:
		if v == "" {
			wrapper.Discard()
			return nil
		}

		wrapper.Value = v

	default:
		wrapper.Discard()
		if !discard {
			return ErrorType(WrapperStringName, value)
		}
	}

	return nil
}

func (wrapper *WrapperString) Unwrap() string {
	if wrapper.IsDiscarded() {
		return ""
	}

	return wrapper.Value
}

func (wrapper *WrapperString) UnwrapAny() any {
	return wrapper.Unwrap()
}

func (wrapper *WrapperString) MarshalJSON() ([]byte, error) {
	return MarshalJSON(wrapper)
}

func (wrapper *WrapperString) UnmarshalJSON(data []byte) error {
	if wrapper == nil {
		return fmt.Errorf("unmarshal into nil wrapper")
	}

	return UnmarshalJSON(data, wrapper)
}
