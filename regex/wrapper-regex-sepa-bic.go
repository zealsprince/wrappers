package regex

import "github.com/zealsprince/wrappers"

const (
	WrapperRegexSepaBicName    wrappers.Name = "WrapperRegexSepaBic"
	WrapperRegexSepaBicPattern string        = `^[A-Z]{6,6}[A-Z2-9][A-NP-Z0-9]([A-Z0-9]{3,3}){0,1}$`
)

type WrapperRegexSepaBic struct {
	WrapperRegex
}

func (wrapper *WrapperRegexSepaBic) Initialize() {
	wrapper.WrapperRegex.SetPattern(WrapperRegexSepaBicName, WrapperRegexSepaBicPattern)
	wrapper.WrapperBase.Initialize()
}
