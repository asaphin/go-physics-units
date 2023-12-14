package mass

import (
	"github.com/asaphin/go-physics-units/conversion"
	"github.com/asaphin/go-physics-units/internal/immutable"
	"github.com/asaphin/go-physics-units/internal/rates"
)

const BaseUnit = Kilogram

const (
	Dalton    = "Da"
	Microgram = "µg"
	Milligram = "mg"
	Gram      = "g"
	Kilogram  = "kg"
	Ounce     = "oz"
	Pound     = "lb"
	Quintal   = "Qts"
	Tonne     = "t"
	EarthMass = "M⊕"
	SolarMass = "M☉"
)

var conversionFactors = conversion.Factors{
	Dalton:    1.6605390666e-27,
	Microgram: 1e-9,
	Milligram: 1e-6,
	Gram:      0.001,
	Kilogram:  1,
	Ounce:     0.0283495,
	Pound:     0.453592,
	Quintal:   100,
	Tonne:     1000,
	EarthMass: 5.9722e24,
	SolarMass: 1.98847e30,
}

// ConversionFactors shows how many base units (kg) in specified unit.
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
