package velocity

import "github.com/asaphin/go-physics-units/conversion-factors"

const BaseUnit = MeterPerSecond

const (
	MeterPerSecond   = "m/s"
	KilometerPerHour = "km/h"
	MilePerHour      = "mph"
)

var velocityConversionFactors = conversion.Factors{
	MeterPerSecond:   1,
	KilometerPerHour: 0.277778,
	MilePerHour:      0.44704,
}

// ConversionFactors shows how many base velocity units (m/s) in specified velocity unit.
func ConversionFactors() conversion.Factors {
	return conversion.CopyConversionFactors(velocityConversionFactors)
}
