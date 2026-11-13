// Package pressure defines the pressure unit symbols supported by the units package
// and their conversion factors to the base unit, the pascal.
package pressure

import (
	"github.com/asaphin/go-physics-units/conversion"
	"github.com/asaphin/go-physics-units/internal/rates"
)

// BaseUnit is the SI unit of pressure that all pressure factors are relative to.
const BaseUnit = Pascal

// Pressure unit symbols.
const (
	Pascal              = "Pa"
	Hectopascal         = "hPa"
	Millibar            = "mbar"
	Kilopascal          = "kPa"
	Torr                = "Torr" // 1/760 standard atmosphere
	MillimeterOfMercury = "mmHg" // conventional: 1 mm of mercury of density 13595.1 kg/m³ under standard gravity
	InchOfMercury       = "inHg" // conventional, as for MillimeterOfMercury
	PoundPerSquareInch  = "psi"  // pound-force per square inch
	Bar                 = "bar"
	Atmosphere          = "atm" // standard atmosphere
	Megapascal          = "MPa"
	Gigapascal          = "GPa"
)

// The factors are exact: defined values, written as constant expressions that Go
// evaluates exactly before rounding once.
//
// Source: NIST SP 811, appendix B.
var conversionFactors = conversion.Factors{
	Pascal:              1,
	Hectopascal:         100,
	Millibar:            100,
	Kilopascal:          1e3,
	Torr:                101325.0 / 760,
	MillimeterOfMercury: 13595.1 * 9.80665 / 1000,
	InchOfMercury:       13595.1 * 9.80665 * 0.0254,
	PoundPerSquareInch:  0.45359237 * 9.80665 / (0.0254 * 0.0254),
	Bar:                 1e5,
	Atmosphere:          101325,
	Megapascal:          1e6,
	Gigapascal:          1e9,
}

// ConversionFactors returns a copy of the pressure conversion factors: how many
// base units (Pa) there are in one unit of each symbol.
func ConversionFactors() conversion.Factors {
	return conversion.CopyConversionFactors(conversionFactors)
}

// Rate returns the multiplier that converts a value in unit from into unit to,
// and whether both are pressure units. Both forms of the micro sign are accepted.
func Rate(from, to string) (float64, bool) {
	return rates.FactorRate(conversionFactors, from, to)
}
