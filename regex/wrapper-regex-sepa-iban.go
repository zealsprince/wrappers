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
