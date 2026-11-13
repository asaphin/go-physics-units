// Package power defines the power unit symbols supported by the units package
// and their conversion factors to the base unit, the watt.
package power

import (
	"github.com/asaphin/go-physics-units/conversion"
	"github.com/asaphin/go-physics-units/internal/rates"
)

// BaseUnit is the SI unit of power that all power factors are relative to.
const BaseUnit = Watt

// Power unit symbols.
const (
	Milliwatt        = "mW"
	Watt             = "W"
	BTUPerHour       = "Btu/h" // International Table British thermal unit per hour
	MetricHorsepower = "PS"    // 75 kgf·m/s
	Horsepower       = "hp"    // mechanical horsepower, 550 ft·lbf/s
	Kilowatt         = "kW"
	Megawatt         = "MW"
	Gigawatt         = "GW"
)

// The factors are exact: defined values, written as constant expressions that Go
// evaluates exactly before rounding once.
//
// Source: NIST SP 811, appendix B.
var conversionFactors = conversion.Factors{
	Milliwatt:        1e-3,
	Watt:             1,
	BTUPerHour:       1055.05585262 / 3600,
	MetricHorsepower: 75 * 9.80665,
	Horsepower:       550 * 0.3048 * 0.45359237 * 9.80665,
	Kilowatt:         1e3,
	Megawatt:         1e6,
	Gigawatt:         1e9,
}

// ConversionFactors returns a copy of the power conversion factors: how many
// base units (W) there are in one unit of each symbol.
func ConversionFactors() conversion.Factors {
	return conversion.CopyConversionFactors(conversionFactors)
}

// Rate returns the multiplier that converts a value in unit from into unit to,
// and whether both are power units. Both forms of the micro sign are accepted.
func Rate(from, to string) (float64, bool) {
	return rates.FactorRate(conversionFactors, from, to)
}
