// Package force defines the force unit symbols supported by the units package
// and their conversion factors to the base unit, the newton.
package force

import (
	"github.com/asaphin/go-physics-units/conversion"
	"github.com/asaphin/go-physics-units/internal/rates"
)

// BaseUnit is the SI unit of force that all force factors are relative to.
const BaseUnit = Newton

// Force unit symbols.
const (
	Dyne          = "dyn"
	Millinewton   = "mN"
	Newton        = "N"
	PoundForce    = "lbf" // the avoirdupois pound under standard gravity
	KilogramForce = "kgf" // the kilogram under standard gravity
	Kilonewton    = "kN"  // not to be confused with velocity.Knot, "kn"
	Meganewton    = "MN"
)

// The factors are exact: defined values, written as constant expressions that Go
// evaluates exactly before rounding once.
//
// Source: NIST SP 811, appendix B.
var conversionFactors = conversion.Factors{
	Dyne:          1e-5,
	Millinewton:   1e-3,
	Newton:        1,
	PoundForce:    0.45359237 * 9.80665,
	KilogramForce: 9.80665,
	Kilonewton:    1e3,
	Meganewton:    1e6,
}

// ConversionFactors returns a copy of the force conversion factors: how many
// base units (N) there are in one unit of each symbol.
func ConversionFactors() conversion.Factors {
	return conversion.CopyConversionFactors(conversionFactors)
}

// Rate returns the multiplier that converts a value in unit from into unit to,
// and whether both are force units. Both forms of the micro sign are accepted.
func Rate(from, to string) (float64, bool) {
	return rates.FactorRate(conversionFactors, from, to)
}
