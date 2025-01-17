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
