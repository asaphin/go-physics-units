// Package time defines the time (duration) unit symbols supported by the units
// package and their conversion factors to the base unit, the second.
//
// It shares its name with the standard library package; import one of them
// under an alias when both are needed.
package time

import (
	"github.com/asaphin/go-physics-units/conversion"
	"github.com/asaphin/go-physics-units/internal/rates"
)

// BaseUnit is the SI unit of time that all time factors are relative to.
const BaseUnit = Second

// Time unit symbols.
const (
	PlanckTime  = "tp"
	Nanosecond  = "ns"
	Microsecond = "μs" // U+03BC GREEK SMALL LETTER MU; "µs" (U+00B5 MICRO SIGN) is accepted too
	Millisecond = "ms"
	Second      = "s"
	Minute      = "min"
	Hour        = "h"
	Day         = "d"
	Week        = "wk"
	Year        = "yr" // Julian year of 365.25 days, as used in astronomy
)

// PlankTime is a misspelled alias of PlanckTime.
//
// Deprecated: Use PlanckTime.
const PlankTime = PlanckTime

var conversionFactors = conversion.Factors{
	PlanckTime:  5.391247e-44, // CODATA 2018
	Nanosecond:  1e-9,
	Microsecond: 1e-6,
	Millisecond: 0.001,
	Second:      1,
	Minute:      60,
	Hour:        3600,
	Day:         86400,
	Week:        604800,
	Year:        31557600,
}

// ConversionFactors returns a copy of the time conversion factors: how many
// base units (s) there are in one unit of each symbol.
func ConversionFactors() conversion.Factors {
	return conversion.CopyConversionFactors(conversionFactors)
}

// Rate returns the multiplier that converts a value in unit from into unit to,
// and whether both are time units. For example, Rate(Hour, Minute) is 60.
// Both forms of the micro sign are accepted.
func Rate(from, to string) (float64, bool) {
	return rates.FactorRate(conversionFactors, from, to)
}
