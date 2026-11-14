// Package acceleration defines the acceleration unit symbols supported by the units package
// and their conversion factors to the base unit, the metre per second squared.
package acceleration

import (
	"github.com/asaphin/go-physics-units/conversion"
	"github.com/asaphin/go-physics-units/internal/rates"
)

// BaseUnit is the SI unit of acceleration that all acceleration factors are relative to.
const BaseUnit = MeterPerSecondSquared

// Acceleration unit symbols.
const (
	CentimeterPerSecondSquared = "cm/s²"
	FootPerSecondSquared       = "ft/s²"
	MeterPerSecondSquared      = "m/s²"
	StandardGravity            = "gₙ" // standard acceleration of gravity, gₙ with U+2099 LATIN SUBSCRIPT SMALL LETTER N; exact by definition (3rd CGPM, 1901)
)

// The factors are exact: defined values, written as constant expressions that Go
// evaluates exactly before rounding once.
//
// Source: NIST SP 811, appendix B.
var conversionFactors = conversion.Factors{
	CentimeterPerSecondSquared: 0.01,
	FootPerSecondSquared:       0.3048,
	MeterPerSecondSquared:      1,
	StandardGravity:            9.80665,
}

// ConversionFactors returns a copy of the acceleration conversion factors: how many
// base units (m/s²) there are in one unit of each symbol.
func ConversionFactors() conversion.Factors {
	return conversion.CopyConversionFactors(conversionFactors)
}

// Rate returns the multiplier that converts a value in unit from into unit to,
// and whether both are acceleration units. Both forms of the micro sign are accepted.
func Rate(from, to string) (float64, bool) {
	return rates.FactorRate(conversionFactors, from, to)
}
