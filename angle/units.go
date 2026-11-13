// Package angle defines the angle unit symbols supported by the units package
// and their conversion factors to the base unit, the radian.
package angle

import (
	"math"

	"github.com/asaphin/go-physics-units/conversion"
	"github.com/asaphin/go-physics-units/internal/rates"
)

// BaseUnit is the SI unit of angle that all angle factors are relative to.
const BaseUnit = Radian

// Angle unit symbols.
const (
	Arcsecond   = "arcsec" // 1/3600 degree
	Arcminute   = "arcmin" // 1/60 degree
	Milliradian = "mrad"
	Gradian     = "gon" // 1/400 revolution
	Degree      = "°"   // U+00B0 DEGREE SIGN
	Radian      = "rad"
	Revolution  = "rev" // one full turn
)

// The factors are exact: defined values, written as constant expressions that Go
// evaluates exactly before rounding once.
//
// Source: the SI Brochure (9th edition), table 8. math.Pi is an untyped
// constant with far more precision than a float64, so each factor is rounded
// once, like the others.
var conversionFactors = conversion.Factors{
	Arcsecond:   math.Pi / 648000,
	Arcminute:   math.Pi / 10800,
	Milliradian: 1e-3,
	Gradian:     math.Pi / 200,
	Degree:      math.Pi / 180,
	Radian:      1,
	Revolution:  2 * math.Pi,
}

// ConversionFactors returns a copy of the angle conversion factors: how many
// base units (rad) there are in one unit of each symbol.
func ConversionFactors() conversion.Factors {
	return conversion.CopyConversionFactors(conversionFactors)
}

// Rate returns the multiplier that converts a value in unit from into unit to,
// and whether both are angle units. Both forms of the micro sign are accepted.
func Rate(from, to string) (float64, bool) {
	return rates.FactorRate(conversionFactors, from, to)
}
