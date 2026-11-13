// Package frequency defines the frequency unit symbols supported by the units package
// and their conversion factors to the base unit, the hertz.
package frequency

import (
	"github.com/asaphin/go-physics-units/conversion"
	"github.com/asaphin/go-physics-units/internal/rates"
)

// BaseUnit is the SI unit of frequency that all frequency factors are relative to.
const BaseUnit = Hertz

// Frequency unit symbols.
const (
	RevolutionPerMinute = "rpm" // one revolution per minute, as a rotational frequency
	Hertz               = "Hz"
	Kilohertz           = "kHz"
	Megahertz           = "MHz"
	Gigahertz           = "GHz"
	Terahertz           = "THz"
)

// The factors are exact: defined values, written as constant expressions that Go
// evaluates exactly before rounding once.
//
// Source: the SI Brochure (9th edition).
var conversionFactors = conversion.Factors{
	RevolutionPerMinute: 1.0 / 60,
	Hertz:               1,
	Kilohertz:           1e3,
	Megahertz:           1e6,
	Gigahertz:           1e9,
	Terahertz:           1e12,
}

// ConversionFactors returns a copy of the frequency conversion factors: how many
// base units (Hz) there are in one unit of each symbol.
func ConversionFactors() conversion.Factors {
	return conversion.CopyConversionFactors(conversionFactors)
}

// Rate returns the multiplier that converts a value in unit from into unit to,
// and whether both are frequency units. Both forms of the micro sign are accepted.
func Rate(from, to string) (float64, bool) {
	return rates.FactorRate(conversionFactors, from, to)
}
