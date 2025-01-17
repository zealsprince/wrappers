package wrappers

import (
	"fmt"
	"strings"
	"time"
)

const (
	WrapperTimeISO8601Name Name = "WrapperTimeISO8601"
)

type WrapperTimeISO8601 Wrapper[time.Time, string]

var _ WrapperProvider = (*WrapperTimeISO8601)(nil) // Ensure that WrapperTimeISO8601 implements WrapperProvider.

func (wrapper *WrapperTimeISO8601) Get() time.Time {
	return wrapper.Value
}

func (wrapper *WrapperTimeISO8601) GetAny() any {
	return wrapper.Get()
}

func (wrapper *WrapperTimeISO8601) Wrap(value any, discard bool) error {
	switch v := value.(type) {
	case nil:
		wrapper.Discard()
		if !discard {
			return ErrorNil(WrapperTimeISO8601Name)
		}

	case WrapperTime:
		if v.IsDiscarded() {
			wrapper.Discard()
			return nil
		}

		wrapper.Value = v.Get()

	case WrapperProvider:
		if v.IsDiscarded() {
			wrapper.Discard()
			return nil
		}

		return wrapper.Wrap(v.UnwrapAny(), discard)

	case time.Time:
		wrapper.Value = v

	case string:
		if v == "" {
			wrapper.Discard()
			return nil
		}

		v = strings.TrimSuffix(v, "+0000")
		if !strings.HasSuffix(v, "Z") {
			v += "Z"
		}

		parsed, err := time.Parse(time.RFC3339, v)
		if err != nil {
			wrapper.Discard()
			if !discard {
				return ErrorValue(WrapperTimeISO8601Name, value, err.Error())
			} else {
				return nil
			}
		}

		wrapper.Value = parsed

	default:
		wrapper.Discard()
		if !discard {
			return ErrorType(WrapperTimeISO8601Name, value)
		}
	}

	return nil
}

func (wrapper *WrapperTimeISO8601) Unwrap() string {
	if wrapper.IsDiscarded() {
		return new(time.Time).Format(time.RFC3339)
	}

	return wrapper.Value.Format(time.RFC3339)
}

func (wrapper *WrapperTimeISO8601) UnwrapAny() any {
	return wrapper.Unwrap()
}

func (wrapper *WrapperTimeISO8601) MarshalJSON() ([]byte, error) {
	return MarshalJSON(wrapper)
}

func (wrapper *WrapperTimeISO8601) UnmarshalJSON(data []byte) error {
	if wrapper == nil {
		return fmt.Errorf("unmarshal into nil wrapper")
	}

	return UnmarshalJSON(data, wrapper)
}
