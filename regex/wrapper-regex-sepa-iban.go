package regex

import "github.com/zealsprince/wrappers"

const (
	WrapperRegexSepaIbanName    wrappers.Name = "WrapperRegexSepaIban"
	WrapperRegexSepaIbanPattern string        = `^[A-Z]{2}[0-9]{2}[A-Z0-9]{1,30}$`
)

type WrapperRegexSepaIban struct {
	WrapperRegex
}

func (wrapper *WrapperRegexSepaIban) Initialize() {
	wrapper.WrapperRegex.SetPattern(WrapperRegexSepaIbanName, WrapperRegexSepaIbanPattern)
	wrapper.WrapperBase.Initialize()
}

// UnmarshalJSON ensures the wrapper is initialized before unmarshalling and proxies the call.
// Without this, promotion hands the embedded WrapperRegex to the unmarshaller and
// SetPattern never runs, so every value fails with "Regex not set".
func (wrapper *WrapperRegexSepaIban) UnmarshalJSON(data []byte) error {
	if !wrapper.IsInitialized() {
		wrapper.Initialize()
	}
	return wrapper.WrapperRegex.UnmarshalJSON(data)
}
