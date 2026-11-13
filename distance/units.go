// Package distance defines the distance (length) unit symbols supported by the
// units package and their conversion factors to the base unit, the metre.
package distance

import (
	"math"

	"github.com/asaphin/go-physics-units/conversion"
	"github.com/asaphin/go-physics-units/internal/rates"
)

// BaseUnit is the SI unit of length that all distance factors are relative to.
const BaseUnit = Meter

// Distance unit symbols.
const (
	PlanckLength     = "lp"
	Femtometer       = "fm"
	Picometer        = "pm"
	Angstrom         = "Å"
	Nanometer        = "nm"
	Micrometer       = "μm" // U+03BC GREEK SMALL LETTER MU; "µm" (U+00B5 MICRO SIGN) is accepted too
	Millimeter       = "mm"
	Centimeter       = "cm"
	Decimeter        = "dm"
	Meter            = "m"
	Kilometer        = "km"
	Inch             = "in"
	Feet             = "ft"
	Yard             = "yd"
	Mile             = "mi"
	NauticalMile     = "nmi"
	AstronomicalUnit = "AU"
	LightYear        = "ly"
	Parsec           = "pc"
)

// PlankLength is a misspelled alias of PlanckLength.
//
// Deprecated: Use PlanckLength.
const PlankLength = PlanckLength

const (
	astronomicalUnit = 1.495978707e11 // m, exact by definition (IAU 2012)
	speedOfLight     = 299792458      // m/s, exact by definition
	julianYear       = 31557600       // s, 365.25 days
)

var conversionFactors = conversion.Factors{
	PlanckLength:     1.616255e-35, // CODATA 2018
	Femtometer:       1e-15,
	Picometer:        1e-12,
	Angstrom:         1e-10,
	Nanometer:        1e-9,
	Micrometer:       1e-6,
	Millimeter:       0.001,
	Centimeter:       0.01,
	Decimeter:        0.1,
	Meter:            1,
	Kilometer:        1000,
	Inch:             0.0254,
	Feet:             0.3048,
	Yard:             0.9144,
	Mile:             1609.344,
	NauticalMile:     1852,
	AstronomicalUnit: astronomicalUnit,
	LightYear:        speedOfLight * julianYear,           // exact (IAU)
	Parsec:           astronomicalUnit * 648000 / math.Pi, // exact by definition (IAU 2015)
}

// ConversionFactors returns a copy of the distance conversion factors: how many
// base units (m) there are in one unit of each symbol.
func ConversionFactors() conversion.Factors {
	return conversion.CopyConversionFactors(conversionFactors)
}

// Rate returns the multiplier that converts a value in unit from into unit to,
// and whether both are distance units. For example, Rate(Kilometer, Meter) is 1000.
// Both forms of the micro sign are accepted.
func Rate(from, to string) (float64, bool) {
	return rates.FactorRate(conversionFactors, from, to)
}
