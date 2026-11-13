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

// ConversionFactors returns a copy of the area conversion factors: how many
// base units (m²) there are in one unit of each symbol.
func ConversionFactors() conversion.Factors {
	return conversion.CopyConversionFactors(conversionFactors)
}

// Rate returns the multiplier that converts a value in unit from into unit to,
// and whether both are area units. For example, Rate(SquareKilometer, SquareMeter) is 1e6.
// Both forms of the micro sign are accepted.
func Rate(from, to string) (float64, bool) {
	return rates.FactorRate(conversionFactors, from, to)
}
