package wrappers

import (
	"fmt"
	"strconv"
)

const (
	WrapperIntName Name = "WrapperInt"
)

type WrapperInt Wrapper[int64, int64]

var _ WrapperProvider = (*WrapperInt)(nil) // Ensure that WrapperInt implements WrapperProvider.

func (wrapper *WrapperInt) Get() int64 {
	return wrapper.Value
}

func (wrapper *WrapperInt) GetAny() any {
	return wrapper.Get()
}

func (wrapper *WrapperInt) Wrap(value any, discard bool) error {
	switch v := value.(type) {
	case nil:
		wrapper.Discard()
		if !discard {
			return ErrorNil(WrapperIntName)
		}

	case WrapperProvider:
		if v.IsDiscarded() {
			wrapper.Discard()
			return nil
		}

		return wrapper.Wrap(v.UnwrapAny(), discard)

	case string:
		converted, err := strconv.Atoi(v)
		if err != nil {
			wrapper.Discard()
			if !discard {
				return ErrorValue(WrapperIntName, value, "int")
			}
		}

		wrapper.Value = int64(converted)

	case int:
		wrapper.Value = int64(v)

	case int16:
		wrapper.Value = int64(v)

	case int32:
		wrapper.Value = int64(v)

	case int64:
		wrapper.Value = v

	case float32:
		wrapper.Value = int64(v)

	case float64:
		wrapper.Value = int64(v)

	default:
		wrapper.Discard()
		if !discard {
			return ErrorType(WrapperIntName, value)
		}
	}

	return nil
}

func (wrapper *WrapperInt) Unwrap() int64 {
	if wrapper.IsDiscarded() {
		return 0
	}

	return int64(wrapper.Value)
}

func (wrapper *WrapperInt) UnwrapAny() any {
	return wrapper.Unwrap()
}

func (wrapper *WrapperInt) MarshalJSON() ([]byte, error) {
	return MarshalJSON(wrapper)
}

func (wrapper *WrapperInt) UnmarshalJSON(data []byte) error {
	if wrapper == nil {
		return fmt.Errorf("unmarshal into nil wrapper")
	}

	return UnmarshalJSON(data, wrapper)
}
