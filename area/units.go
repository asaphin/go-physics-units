package area

import (
	"github.com/asaphin/go-physics-units/conversion"
	"github.com/asaphin/go-physics-units/internal/immutable"
	"github.com/asaphin/go-physics-units/internal/rates"
)

const BaseUnit = SquareMeter

const (
	SquareMillimeter = "mm²"
	SquareCentimeter = "cm²"
	SquareInch       = "in²"
	SquareFoot       = "ft²"
	SquareMeter      = "m²"
	SquareKilometer  = "km²"
	SquareMile       = "mi²"
)

var conversionFactors = conversion.Factors{
	SquareMillimeter: 1e-6,
	SquareCentimeter: 1e-4,
	SquareInch:       0.00064516,
	SquareFoot:       0.09290304,
	SquareMeter:      1,
	SquareKilometer:  1e6,
	SquareMile:       2589988.1103360,
}

// ConversionFactors shows how many base units (m²) in specified unit.
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
