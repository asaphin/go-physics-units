package velocity

import (
	"github.com/asaphin/go-physics-units/conversion"
	"github.com/asaphin/go-physics-units/internal/immutable"
	"github.com/asaphin/go-physics-units/internal/rates"
	"sync"
)

const BaseUnit = MeterPerSecond

const (
	MeterPerSecond   = "m/s"
	KilometerPerHour = "km/h"
	MilePerHour      = "mph"
)

var conversionFactors = conversion.Factors{
	MeterPerSecond:   1,
	KilometerPerHour: 0.277778,
	MilePerHour:      0.44704,
}

// ConversionFactors shows how many base units (m/s) in specified unit.
func ConversionFactors() conversion.Factors {
	return conversion.CopyConversionFactors(conversionFactors)
}

var conversionRates immutable.Float64Map
var conversionRatesSync sync.Once

// ConversionRates returns pointer to conversion rates storage.
// Rates stored by composite keys unitFrom + unitTo
func ConversionRates() immutable.Float64Map {
	conversionRatesSync.Do(func() {
		conversionRates = immutable.MakeImmutable(rates.FactorsToRates(conversionFactors))
	})

	return conversionRates
}
