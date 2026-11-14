// Package charge defines the electric charge units of the units package and their
// conversion factors to the base unit, the coulomb.
package charge

// Unit is an electric charge unit: one of the constants below. String returns its
// symbol and ParseUnit finds the unit of a symbol.
//
// The zero Unit is the base unit, Coulomb. The numbers behind the constants may
// change between versions: store and exchange units by their symbols.
type Unit uint8

// BaseUnit is the SI unit of electric charge that all factors are relative to.
const BaseUnit = Coulomb

// Charge units, with their symbols.
const (
	Coulomb         Unit = iota // C
	Microcoulomb                // μC, with U+03BC GREEK SMALL LETTER MU; "µC" (U+00B5 MICRO SIGN) is accepted too
	Millicoulomb                // mC
	MilliampereHour             // mAh: the battery capacity unit, 3.6 C
	AmpereHour                  // Ah

	unitCount // the number of units; keep it last
)

// definitions holds the symbol of each unit and its factor: how many base
// units (C) there are in one unit.
//
// Source: the SI Brochure (9th edition). The ampere-hour is 3600 C exactly.
var definitions = [unitCount]definition{
	Coulomb:         {"C", 1},
	Microcoulomb:    {"μC", 1e-6},
	Millicoulomb:    {"mC", 1e-3},
	MilliampereHour: {"mAh", 3.6},
	AmpereHour:      {"Ah", 3600},
}
