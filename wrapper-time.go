package wrappers

import (
	"fmt"
	"time"
)

const (
	WrapperTimeName Name = "WrapperTime"
)

type WrapperTime Wrapper[time.Time, string]

var _ WrapperProvider = (*WrapperTime)(nil) // Ensure that WrapperTime implements WrapperProvider.

func (wrapper *WrapperTime) Get() time.Time {
	return wrapper.Value
}

func (wrapper *WrapperTime) GetAny() any {
	return wrapper.Get()
}

func (wrapper *WrapperTime) Wrap(value any, discard bool) error {
	switch v := value.(type) {
	case nil:
		wrapper.Discard()
		if !discard {
			return ErrorNil(WrapperTimeName)
		}

	case WrapperProvider:
		if v.IsDiscarded() {
			wrapper.Discard()
			return nil
		}

		return wrapper.Wrap(v.UnwrapAny(), discard)

	case time.Time:
		wrapper.Value = v

	case string:
		converted, err := time.Parse(time.RFC3339, v)
		if err != nil {
			wrapper.Discard()
			if !discard {
				return ErrorValue(WrapperTimeName, value, "RFC3339")
			}
		}

		wrapper.Value = converted

	default:
		wrapper.Discard()
		if !discard {
			return ErrorType(WrapperTimeName, value)
		}
	}

	return nil
}

func (wrapper *WrapperTime) Unwrap() string {
	if wrapper.IsDiscarded() {
		return ""
	}

	return wrapper.Value.Format(time.RFC3339)
}

func (wrapper *WrapperTime) UnwrapAny() any {
	return wrapper.Unwrap()
}

func (wrapper *WrapperTime) MarshalJSON() ([]byte, error) {
	return MarshalJSON(wrapper)
}

func (wrapper *WrapperTime) UnmarshalJSON(data []byte) error {
	if wrapper == nil {
		return fmt.Errorf("unmarshal into nil wrapper")
	}

	return UnmarshalJSON(data, wrapper)
}
