package wrappers

import (
	"fmt"
	"time"
)

const (
	WrapperTimeDurationName Name = "WrapperTimeDuration"
)

type WrapperTimeDuration Wrapper[time.Duration, string]

var _ WrapperProvider = (*WrapperTimeDuration)(nil) // Ensure that WrapperTimeDuration implements WrapperProvider.

func (wrapper *WrapperTimeDuration) Get() time.Duration {
	return wrapper.Value
}

func (wrapper *WrapperTimeDuration) GetAny() any {
	return wrapper.Get()
}

func (wrapper *WrapperTimeDuration) Wrap(value any, discard bool) error {
	fmt.Printf("v: %v\n", value)
	switch v := value.(type) {
	case nil:
		wrapper.Discard()

	case WrapperProvider:
		if v.IsDiscarded() {
			wrapper.Discard()
			return nil
		}

		return wrapper.Wrap(v.UnwrapAny(), discard)

	case time.Duration:
		wrapper.Value = v

	case int, int8, int16, int32, int64:
		wrapper.Value = time.Duration(v.(int))

	case float32, float64:
		wrapper.Value = time.Duration(v.(float64))

	case string:
		converted, err := time.ParseDuration(v)
		if err != nil {
			wrapper.Discard()
			if !discard {
				return ErrorValue(WrapperTimeDurationName, value, "time.Duration")
			}
			return nil
		}

		wrapper.Value = converted

	default:
		wrapper.Discard()
		if !discard {
			return ErrorType(WrapperTimeDurationName, value)
		}
	}

	return nil
}

func (wrapper *WrapperTimeDuration) Unwrap() string {
	if wrapper.IsDiscarded() {
		return "0s"
	}

	return wrapper.Value.String()
}

func (wrapper *WrapperTimeDuration) UnwrapAny() any {
	return wrapper.Unwrap()
}

func (wrapper *WrapperTimeDuration) MarshalJSON() ([]byte, error) {
	return MarshalJSON(wrapper)
}

func (wrapper *WrapperTimeDuration) UnmarshalJSON(data []byte) error {
	if wrapper == nil {
		return fmt.Errorf("unmarshal into nil wrapper")
	}

	return UnmarshalJSON(data, wrapper)
}
