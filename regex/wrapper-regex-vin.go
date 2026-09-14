package regex

import "github.com/zealsprince/wrappers"

const (
	WrapperRegexVinName    wrappers.Name = "WrapperRegexVin"
	WrapperRegexVinPattern string        = `^[A-HJ-NPR-Z0-9]{17}$`
)

type WrapperRegexVin struct {
	WrapperRegex
}

func (wrapper *WrapperRegexVin) Initialize() {
	wrapper.WrapperRegex.SetPattern(WrapperRegexVinName, WrapperRegexVinPattern)
	wrapper.WrapperBase.Initialize()
}

// UnmarshalJSON ensures the wrapper is initialized before unmarshalling and proxies the call.
// Without this, promotion hands the embedded WrapperRegex to the unmarshaller and
// SetPattern never runs, so every value fails with "Regex not set".
func (wrapper *WrapperRegexVin) UnmarshalJSON(data []byte) error {
	if !wrapper.IsInitialized() {
		wrapper.Initialize()
	}
	return wrapper.WrapperRegex.UnmarshalJSON(data)
}
