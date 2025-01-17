package regex

import "github.com/zealsprince/wrappers"

const (
	WrapperRegexPhoneName    wrappers.Name = "WrapperRegexPhone"
	WrapperRegexPhonePattern string        = `^(?:\+?[1-9]\d{1,14}|0\d{1,14})$`
)

type WrapperRegexPhone struct {
	WrapperRegex
}

func (wrapper *WrapperRegexPhone) Initialize() {
	wrapper.WrapperRegex.SetPattern(WrapperRegexPhoneName, WrapperRegexPhonePattern)
	wrapper.WrapperBase.Initialize()
}

func (wrapper *WrapperRegexPhone) UnmarshalJSON(data []byte) error {
	if !wrapper.IsInitialized() {
		wrapper.Initialize()
	}
	return wrapper.WrapperRegex.UnmarshalJSON(data)
}
