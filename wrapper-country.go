package wrappers

import (
	"fmt"

	"github.com/biter777/countries"
)

const (
	WrapperCountryName Name = "WrapperCountry"
)

type WrapperCountry Wrapper[countries.CountryCode, string]

var _ WrapperProvider = (*WrapperCountry)(nil) // Ensure that WrapperCountry implements WrapperProvider.

func (wrapper *WrapperCountry) Get() countries.CountryCode {
	return wrapper.Value
}

func (wrapper *WrapperCountry) GetAny() any {
	return wrapper.Get()
}

func (wrapper *WrapperCountry) Wrap(value any, discard bool) error {
	switch v := value.(type) {
	case nil:
		wrapper.Discard()

	case WrapperProvider:
		if v.IsDiscarded() {
			wrapper.Discard()
			return nil
		}

		return wrapper.Wrap(v.UnwrapAny(), discard)

	case countries.CountryCode:
		if v == countries.Unknown {
			wrapper.Discard()
			if !discard {
				return ErrorValue(WrapperCountryName, value, "DE")
			}
		}
		wrapper.Value = v

	case countries.Country:
		if v.Code == countries.Unknown {
			wrapper.Discard()
			if !discard {
				return ErrorValue(WrapperCountryName, value, "DE")
			}
		}
		wrapper.Value = v.Code

	case string:
		if v == "Unknown" {
			wrapper.Discard()
			return nil
		}

		code := countries.ByName(v)
		if code != countries.Unknown {
			wrapper.Value = code

		} else {
			wrapper.Discard()
			if !discard {
				return ErrorValue(WrapperCountryName, value, "DE")
			}
		}

	default:
		wrapper.Discard()
		if !discard {
			return ErrorType(WrapperCountryName, value)
		}
	}

	return nil
}

func (wrapper *WrapperCountry) Unwrap() string {
	if wrapper.IsDiscarded() {
		return countries.Unknown.String()
	}

	return wrapper.Value.String()
}

func (wrapper *WrapperCountry) UnwrapAny() any {
	return wrapper.Unwrap()
}

func (wrapper *WrapperCountry) MarshalJSON() ([]byte, error) {
	return MarshalJSON(wrapper)
}

func (wrapper *WrapperCountry) UnmarshalJSON(data []byte) error {
	if wrapper == nil {
		return fmt.Errorf("unmarshal into nil wrapper")
	}

	return UnmarshalJSON(data, wrapper)
}
