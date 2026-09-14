package regex

import "github.com/zealsprince/wrappers"

const (
	WrapperRegexUrlName    wrappers.Name = "WrapperRegexUrl"
	WrapperRegexUrlPattern string        = `^(https?|ftp)://[^\s/$.?#].[^\s]*$`
)

type WrapperRegexUrl struct {
	WrapperRegex
}

func (wrapper *WrapperRegexUrl) Initialize() {
	wrapper.WrapperRegex.SetPattern(WrapperRegexUrlName, WrapperRegexUrlPattern)
	wrapper.WrapperBase.Initialize()
}

// UnmarshalJSON ensures the wrapper is initialized before unmarshalling and proxies the call.
// Without this, promotion hands the embedded WrapperRegex to the unmarshaller and
// SetPattern never runs, so every value fails with "Regex not set".
func (wrapper *WrapperRegexUrl) UnmarshalJSON(data []byte) error {
	if !wrapper.IsInitialized() {
		wrapper.Initialize()
	}
	return wrapper.WrapperRegex.UnmarshalJSON(data)
}
