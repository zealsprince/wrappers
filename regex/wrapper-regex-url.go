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
