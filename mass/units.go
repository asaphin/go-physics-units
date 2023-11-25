package mass

import (
	"github.com/asaphin/go-physics-units/conversion"
	"github.com/asaphin/go-physics-units/internal/immutable"
	"github.com/asaphin/go-physics-units/internal/rates"
	"sync"
)

const BaseUnit = Kilogram

const (
	Gram     = "g"
	Kilogram = "kg"
	Pound    = "lb"
	Ounce    = "oz"
)

var conversionFactors = conversion.Factors{
	Gram:     0.001,
	Kilogram: 1,
	Pound:    0.453592,
	Ounce:    0.0283495,
}

// ConversionFactors shows how many base units (kg) in specified unit.
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
