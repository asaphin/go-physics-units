// Package velocity defines the velocity (speed) unit symbols supported by the
// units package and their conversion factors to the base unit, metres per second.
package velocity

import (
	"github.com/asaphin/go-physics-units/conversion"
	"github.com/asaphin/go-physics-units/internal/rates"
)

// BaseUnit is the SI unit of velocity that all velocity factors are relative to.
const BaseUnit = MeterPerSecond

// Velocity unit symbols.
const (
	KilometerPerHour = "km/h"
	FeetPerSecond    = "ft/s"
	MilePerHour      = "mph"
	Knot             = "kn"
	MeterPerSecond   = "m/s"
	SpeedOfLight     = "c"
)

var conversionFactors = conversion.Factors{
	KilometerPerHour: 1000.0 / 3600,
	FeetPerSecond:    0.3048,
	MilePerHour:      0.44704,
	Knot:             1852.0 / 3600, // one nautical mile per hour
	MeterPerSecond:   1,
	SpeedOfLight:     299792458, // exact by definition
}

var conversionRates = rates.MustNewTable("velocity", conversionFactors)

// ConversionFactors returns a copy of the velocity conversion factors: how many
// base units (m/s) there are in one unit of each symbol.
func ConversionFactors() conversion.Factors {
	return conversion.CopyConversionFactors(conversionFactors)
}

// ConversionRates returns the shared, immutable table of conversion rates
// between every pair of velocity units.
func ConversionRates() rates.Table {
	return conversionRates
}
