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
	switch v := value.(type) {
	case nil:
		wrapper.Discard()
		if !discard {
			return ErrorNil(WrapperTimeDurationName)
		}

	case WrapperProvider:
		if v.IsDiscarded() {
			wrapper.Discard()
			return nil
		}

		return wrapper.Wrap(v.UnwrapAny(), discard)

	case time.Duration:
		wrapper.Value = v

	// One case per type. A multi-type case leaves v with the interface type, so
	// the old v.(int) assertion panicked on everything that was not literally an
	// int, including the int64 that WrapperInt.UnwrapAny hands over when wrappers
	// are nested.
	case int:
		wrapper.Value = time.Duration(v)

	case int8:
		wrapper.Value = time.Duration(v)

	case int16:
		wrapper.Value = time.Duration(v)

	case int32:
		wrapper.Value = time.Duration(v)

	case int64:
		wrapper.Value = time.Duration(v)

	case float32:
		wrapper.Value = time.Duration(v)

	case float64:
		wrapper.Value = time.Duration(v)

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
