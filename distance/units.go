// Package distance defines the distance (length) units of the units package and
// their conversion factors to the base unit, the metre.
package distance

import "math"

// Unit is a distance unit: one of the constants below. String returns its
// symbol and ParseUnit finds the unit of a symbol.
//
// The zero Unit is the base unit, Meter. The numbers behind the constants may
// change between versions: store and exchange units by their symbols.
type Unit uint8

// BaseUnit is the SI unit of length that all distance factors are relative to.
const BaseUnit = Meter

// Distance units, with their symbols.
const (
	Meter            Unit = iota // m
	PlanckLength                 // lp
	Femtometer                   // fm
	Picometer                    // pm
	Angstrom                     // Å
	Nanometer                    // nm
	Micrometer                   // μm, with U+03BC GREEK SMALL LETTER MU; "µm" (U+00B5 MICRO SIGN) is accepted too
	Millimeter                   // mm
	Centimeter                   // cm
	Decimeter                    // dm
	Kilometer                    // km
	Inch                         // in
	Feet                         // ft
	Yard                         // yd
	Mile                         // mi
	NauticalMile                 // nmi
	AstronomicalUnit             // AU
	LightYear                    // ly
	Parsec                       // pc
	Kiloparsec                   // kpc
	Megaparsec                   // Mpc
	Gigaparsec                   // Gpc
	Mil                          // mil: a thousandth of an inch, also called a thou
	Fathom                       // ftm: 6 ft

	unitCount // the number of units; keep it last
)

// PlankLength is a misspelled alias of PlanckLength.
//
// Deprecated: Use PlanckLength.
//
//go:fix inline
const PlankLength = PlanckLength

const (
	astronomicalUnit = 1.495978707e11 // m, exact by definition (IAU 2012)
	speedOfLight     = 299792458      // m/s, exact by definition
	julianYear       = 31557600       // s, 365.25 days
)

// definitions holds the symbol of each unit and its factor: how many base
// units (m) there are in one unit.
var definitions = [unitCount]definition{
	Meter:            {"m", 1},
	PlanckLength:     {"lp", 1.616255e-35}, // CODATA 2022
	Femtometer:       {"fm", 1e-15},
	Picometer:        {"pm", 1e-12},
	Angstrom:         {"Å", 1e-10},
	Nanometer:        {"nm", 1e-9},
	Micrometer:       {"μm", 1e-6},
	Millimeter:       {"mm", 0.001},
	Centimeter:       {"cm", 0.01},
	Decimeter:        {"dm", 0.1},
	Kilometer:        {"km", 1000},
	Inch:             {"in", 0.0254},
	Feet:             {"ft", 0.3048},
	Yard:             {"yd", 0.9144},
	Mile:             {"mi", 1609.344},
	NauticalMile:     {"nmi", 1852},
	AstronomicalUnit: {"AU", astronomicalUnit},
	LightYear:        {"ly", speedOfLight * julianYear},           // exact (IAU)
	Parsec:           {"pc", astronomicalUnit * 648000 / math.Pi}, // exact by definition (IAU 2015)
	Kiloparsec:       {"kpc", astronomicalUnit * 648000 / math.Pi * 1e3},
	Megaparsec:       {"Mpc", astronomicalUnit * 648000 / math.Pi * 1e6},
	Gigaparsec:       {"Gpc", astronomicalUnit * 648000 / math.Pi * 1e9},
	Mil:              {"mil", 0.0254 / 1000},
	Fathom:           {"ftm", 6 * 0.3048},
}
