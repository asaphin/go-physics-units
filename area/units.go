// Package area defines the area unit symbols supported by the units package and
// their conversion factors to the base unit, the square metre.
package area

import (
	"github.com/asaphin/go-physics-units/conversion"
	"github.com/asaphin/go-physics-units/internal/rates"
)

// BaseUnit is the SI unit of area that all area factors are relative to.
const BaseUnit = SquareMeter

// Area unit symbols.
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
	SquareMile:       2589988.110336,
}

var conversionRates = rates.MustNewTable("area", conversionFactors)

// ConversionFactors returns a copy of the area conversion factors: how many
// base units (m²) there are in one unit of each symbol.
func ConversionFactors() conversion.Factors {
	return conversion.CopyConversionFactors(conversionFactors)
}

// ConversionRates returns the shared, immutable table of conversion rates
// between every pair of area units.
func ConversionRates() rates.Table {
	return conversionRates
}
