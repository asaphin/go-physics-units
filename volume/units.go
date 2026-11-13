// Package volume defines the volume unit symbols supported by the units package
// and their conversion factors to the base unit, the cubic metre.
package volume

import (
	"github.com/asaphin/go-physics-units/conversion"
	"github.com/asaphin/go-physics-units/internal/rates"
)

// BaseUnit is the SI unit of volume that all volume factors are relative to.
const BaseUnit = CubicMeter

// Volume unit symbols.
const (
	CubicMillimeter = "mm³"
	Microliter      = "μL" // U+03BC GREEK SMALL LETTER MU; "µL" (U+00B5 MICRO SIGN) is accepted too
	CubicCentimeter = "cm³"
	Milliliter      = "mL"
	CubicInch       = "in³"
	Liter           = "L"
	Pint            = "pt"  // US liquid pint, 1/8 US gallon
	Quart           = "qt"  // US liquid quart, 1/4 US gallon
	Gallon          = "gal" // US liquid gallon, 231 in³
	CubicFoot       = "ft³"
	Barrel          = "bbl" // oil barrel, 42 US gallons
	CubicYard       = "yd³"
	CubicMeter      = "m³"
	CubicKilometer  = "km³"
)

// The factors are exact: defined values, written as constant expressions that Go
// evaluates exactly before rounding once.
//
// Source: NIST SP 811, appendix B; the inch, foot and yard are the
// international ones (exactly 0.0254, 0.3048 and 0.9144 m).
var conversionFactors = conversion.Factors{
	CubicMillimeter: 1e-9,
	Microliter:      1e-9,
	CubicCentimeter: 1e-6,
	Milliliter:      1e-6,
	CubicInch:       0.0254 * 0.0254 * 0.0254,
	Liter:           1e-3,
	Pint:            231 * 0.0254 * 0.0254 * 0.0254 / 8,
	Quart:           231 * 0.0254 * 0.0254 * 0.0254 / 4,
	Gallon:          231 * 0.0254 * 0.0254 * 0.0254,
	CubicFoot:       0.3048 * 0.3048 * 0.3048,
	Barrel:          42 * 231 * 0.0254 * 0.0254 * 0.0254,
	CubicYard:       0.9144 * 0.9144 * 0.9144,
	CubicMeter:      1,
	CubicKilometer:  1e9,
}

// ConversionFactors returns a copy of the volume conversion factors: how many
// base units (m³) there are in one unit of each symbol.
func ConversionFactors() conversion.Factors {
	return conversion.CopyConversionFactors(conversionFactors)
}

// Rate returns the multiplier that converts a value in unit from into unit to,
// and whether both are volume units. Both forms of the micro sign are accepted.
func Rate(from, to string) (float64, bool) {
	return rates.FactorRate(conversionFactors, from, to)
}
