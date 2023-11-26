package velocity

import (
	"github.com/asaphin/go-physics-units/conversion"
	"github.com/asaphin/go-physics-units/internal/immutable"
	"github.com/asaphin/go-physics-units/internal/rates"
)

const BaseUnit = MeterPerSecond

const (
	KilometerPerHour = "km/h"
	FeetPerSecond    = "ft/s"
	MilePerHour      = "mph"
	Knot             = "kn"
	MeterPerSecond   = "m/s"
	SpeedOfLight     = "c"
)

var conversionFactors = conversion.Factors{
	KilometerPerHour: 0.277778,
	FeetPerSecond:    0.3048,
	MilePerHour:      0.44704,
	Knot:             0.514444,
	MeterPerSecond:   1,
	SpeedOfLight:     299792458,
}

// ConversionFactors shows how many base units (m/s) in specified unit.
func ConversionFactors() conversion.Factors {
	return conversion.CopyConversionFactors(conversionFactors)
}

var conversionRates immutable.Float64Map

// ConversionRates returns pointer to conversion rates storage.
// Rates stored by composite keys unitFrom + unitTo.
func ConversionRates() immutable.Float64Map {
	if conversionRates == nil {
		conversionRates = immutable.MakeImmutable(rates.FactorsToRates(conversionFactors))
	}

	return conversionRates
}
