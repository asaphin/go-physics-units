// Package mass defines the mass unit symbols supported by the units package and
// their conversion factors to the base unit, the kilogram.
package mass

import (
	"github.com/asaphin/go-physics-units/conversion"
	"github.com/asaphin/go-physics-units/internal/rates"
)

// BaseUnit is the SI unit of mass that all mass factors are relative to.
const BaseUnit = Kilogram

// Mass unit symbols.
const (
	Dalton    = "Da"
	Microgram = "μg" // U+03BC GREEK SMALL LETTER MU; "µg" (U+00B5 MICRO SIGN) is accepted too
	Milligram = "mg"
	Gram      = "g"
	Kilogram  = "kg"
	Ounce     = "oz"
	Pound     = "lb"
	Quintal   = "q" // metric quintal, 100 kg
	Tonne     = "t"
	EarthMass = "M⊕"
	SolarMass = "M☉"
)

var conversionFactors = conversion.Factors{
	Dalton:    1.6605390666e-27, // CODATA 2018
	Microgram: 1e-9,
	Milligram: 1e-6,
	Gram:      0.001,
	Kilogram:  1,
	Ounce:     0.028349523125, // avoirdupois ounce, exactly 1/16 lb
	Pound:     0.45359237,     // international avoirdupois pound, exact
	Quintal:   100,
	Tonne:     1000,
	EarthMass: 5.9722e24,
	SolarMass: 1.98847e30,
}

var conversionRates = rates.MustNewTable("mass", conversionFactors)

// ConversionFactors returns a copy of the mass conversion factors: how many
// base units (kg) there are in one unit of each symbol.
func ConversionFactors() conversion.Factors {
	return conversion.CopyConversionFactors(conversionFactors)
}

// ConversionRates returns the shared, immutable table of conversion rates
// between every pair of mass units.
func ConversionRates() rates.Table {
	return conversionRates
}
