package regex

import (
	"github.com/zealsprince/wrappers"
)

const (
	WrapperRegexEmailName    wrappers.Name = "WrapperRegexEmail"
	WrapperRegexEmailPattern string        = `^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`
)

// WrapperRegexEmail is a specialized wrapper for validating email addresses.
type WrapperRegexEmail struct {
	WrapperRegex
}

func (wrapper *WrapperRegexEmail) Initialize() {
	wrapper.WrapperRegex.SetPattern(WrapperRegexEmailName, WrapperRegexEmailPattern)
	wrapper.WrapperBase.Initialize()
}

// UnmarshalJSON ensures the wrapper is initialized before unmarshalling and proxies the call.
func (wrapper *WrapperRegexEmail) UnmarshalJSON(data []byte) error {
	if !wrapper.IsInitialized() {
		wrapper.Initialize()
	}
	return wrapper.WrapperRegex.UnmarshalJSON(data)
}
