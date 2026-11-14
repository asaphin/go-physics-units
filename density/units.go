// Package density defines the density unit symbols supported by the units package
// and their conversion factors to the base unit, the kilogram per cubic metre.
package density

import (
	"github.com/asaphin/go-physics-units/conversion"
	"github.com/asaphin/go-physics-units/internal/rates"
)

// BaseUnit is the SI unit of density that all density factors are relative to.
const BaseUnit = GramPerLiter

// Density unit symbols.
const (
	GramPerLiter           = "g/L" // equal to the base unit
	PoundPerCubicFoot      = "lb/ft³"
	KilogramPerCubicMeter  = "kg/m³"
	PoundPerGallon         = "lb/gal" // per US liquid gallon
	KilogramPerLiter       = "kg/L"
	GramPerMilliliter      = "g/mL"
	GramPerCubicCentimeter = "g/cm³"
)

// The factors are exact: defined values, written as constant expressions that Go
// evaluates exactly before rounding once.
//
// Source: NIST SP 811, appendix B.
var conversionFactors = conversion.Factors{
	GramPerLiter:           1,
	PoundPerCubicFoot:      0.45359237 / (0.3048 * 0.3048 * 0.3048),
	KilogramPerCubicMeter:  1,
	PoundPerGallon:         0.45359237 / (231 * 0.0254 * 0.0254 * 0.0254),
	KilogramPerLiter:       1000,
	GramPerMilliliter:      1000,
	GramPerCubicCentimeter: 1000,
}

// ConversionFactors returns a copy of the density conversion factors: how many
// base units (kg/m³) there are in one unit of each symbol.
func ConversionFactors() conversion.Factors {
	return conversion.CopyConversionFactors(conversionFactors)
}

// Rate returns the multiplier that converts a value in unit from into unit to,
// and whether both are density units. Both forms of the micro sign are accepted.
func Rate(from, to string) (float64, bool) {
	return rates.FactorRate(conversionFactors, from, to)
}
