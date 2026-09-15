package enum

import (
	wrappers "github.com/zealsprince/wrappers/v2"
)

// CardinalDirection ships as the worked example of the pattern above. The whole
// declaration is the type, its constants, and a Values implementation.
type CardinalDirection string

const (
	DirectionNorth CardinalDirection = "north"
	DirectionEast  CardinalDirection = "east"
	DirectionSouth CardinalDirection = "south"
	DirectionWest  CardinalDirection = "west"
)

type cardinalDirections struct{}

func (cardinalDirections) Name() wrappers.Name { return "CardinalDirection" }

func (cardinalDirections) Values() []CardinalDirection {
	return []CardinalDirection{
		DirectionNorth,
		DirectionEast,
		DirectionSouth,
		DirectionWest,
	}
}

type (
	CardinalDirections        = Wrapper[CardinalDirection, cardinalDirections]
	LenientCardinalDirections = Lenient[CardinalDirection, cardinalDirections]
)
