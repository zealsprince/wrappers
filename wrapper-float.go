package wrappers

import (
	"fmt"
	"strconv"
)

const (
	WrapperFloatName Name = "WrapperFloat"
)

type WrapperFloat Wrapper[float64, float64]

var _ WrapperProvider = (*WrapperFloat)(nil) // Ensure that WrapperFloat implements WrapperProvider.

func (wrapper *WrapperFloat) Get() float64 {
	return wrapper.Value
}

func (wrapper *WrapperFloat) GetAny() any {
	return wrapper.Get()
}

func (wrapper *WrapperFloat) Wrap(value any, discard bool) error {
	switch v := value.(type) {
	case nil:
		wrapper.Discard()
		if !discard {
			return ErrorNil(WrapperFloatName)
		}

	case WrapperProvider:
		if v.IsDiscarded() {
			wrapper.Discard()
			return nil
		}

		return wrapper.Wrap(v.UnwrapAny(), discard)

	case string:
		converted, err := strconv.ParseFloat(v, 64)
		if err != nil {
			wrapper.Discard()
			if !discard {
				return ErrorValue(WrapperFloatName, value, "float")
			}
			return nil
		}

		wrapper.Value = converted

	case int:
		wrapper.Value = float64(v)

	case int16:
		wrapper.Value = float64(v)

	case int32:
		wrapper.Value = float64(v)

	case int64:
		wrapper.Value = float64(v)

	case float32:
		wrapper.Value = float64(v)

	case float64:
		wrapper.Value = v

	default:
		wrapper.Discard()
		if !discard {
			return ErrorType(WrapperFloatName, value)
		}
	}

	return nil
}

func (wrapper *WrapperFloat) Unwrap() float64 {
	if wrapper.IsDiscarded() {
		return 0
	}

	return wrapper.Value
}

func (wrapper *WrapperFloat) UnwrapAny() any {
	return wrapper.Unwrap()
}

func (wrapper *WrapperFloat) MarshalJSON() ([]byte, error) {
	return MarshalJSON(wrapper)
}

func (wrapper *WrapperFloat) UnmarshalJSON(data []byte) error {
	if wrapper == nil {
		return fmt.Errorf("unmarshal into nil wrapper")
	}

	return UnmarshalJSON(data, wrapper)
}
