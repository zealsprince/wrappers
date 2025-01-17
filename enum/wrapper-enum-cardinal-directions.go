package enum

import (
	"github.com/zealsprince/wrappers"
)

type CardinalDirection string

const (
	WrapperEnumCardinalDirectionsName wrappers.Name = "WrapperEnumCardinalDirections"

	DirectionNorth CardinalDirection = "north"
	DirectionEast  CardinalDirection = "east"
	DirectionSouth CardinalDirection = "south"
	DirectionWest  CardinalDirection = "west"
)

type WrapperEnumCardinalDirections struct {
	WrapperEnum[CardinalDirection]
}

func (wrapper *WrapperEnumCardinalDirections) Initialize() {
	wrapper.WrapperEnum.SetValidValues(
		[]CardinalDirection{
			DirectionNorth,
			DirectionEast,
			DirectionSouth,
			DirectionWest,
		})

	wrapper.WrapperBase.Initialize()
}

func (wrapper *WrapperEnumCardinalDirections) UnmarshalJSON(data []byte) error {
	if !wrapper.IsInitialized() {
		wrapper.Initialize()
	}
	return wrapper.WrapperEnum.UnmarshalJSON(data)
}
