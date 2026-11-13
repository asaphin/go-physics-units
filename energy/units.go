// Package energy defines the energy unit symbols supported by the units package
// and their conversion factors to the base unit, the joule.
package energy

import (
	"github.com/asaphin/go-physics-units/conversion"
	"github.com/asaphin/go-physics-units/internal/rates"
)

// BaseUnit is the SI unit of energy that all energy factors are relative to.
const BaseUnit = Joule

// Energy unit symbols.
const (
	Electronvolt = "eV" // exact since the 2019 redefinition of the SI
	Erg          = "erg"
	Joule        = "J"
	FootPound    = "ft·lbf" // foot-pound force, with U+00B7 MIDDLE DOT
	Calorie      = "cal"    // thermochemical calorie
	Kilojoule    = "kJ"
	BTU          = "Btu"  // International Table British thermal unit
	Kilocalorie  = "kcal" // thermochemical kilocalorie, the dietary Calorie
	WattHour     = "Wh"
	Megajoule    = "MJ"
	KilowattHour = "kWh"
	Gigajoule    = "GJ"
	MegawattHour = "MWh"
)

// The factors are exact: defined values, written as constant expressions that Go
// evaluates exactly before rounding once.
//
// Source: NIST SP 811, appendix B; the electronvolt from the SI Brochure (9th
// edition).
var conversionFactors = conversion.Factors{
	Electronvolt: 1.602176634e-19,
	Erg:          1e-7,
	Joule:        1,
	FootPound:    0.3048 * 0.45359237 * 9.80665,
	Calorie:      4.184,
	Kilojoule:    1e3,
	BTU:          1055.05585262,
	Kilocalorie:  4184,
	WattHour:     3600,
	Megajoule:    1e6,
	KilowattHour: 3.6e6,
	Gigajoule:    1e9,
	MegawattHour: 3.6e9,
}

// ConversionFactors returns a copy of the energy conversion factors: how many
// base units (J) there are in one unit of each symbol.
func ConversionFactors() conversion.Factors {
	return conversion.CopyConversionFactors(conversionFactors)
}

// Rate returns the multiplier that converts a value in unit from into unit to,
// and whether both are energy units. Both forms of the micro sign are accepted.
func Rate(from, to string) (float64, bool) {
	return rates.FactorRate(conversionFactors, from, to)
}
