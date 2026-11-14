// Package equivalentdose defines the equivalent and effective dose of ionising
// radiation units of the units package and their conversion factors to the base
// unit, the sievert.
//
// Use the units with units.NewEquivalentDose, from the root package
// github.com/asaphin/go-physics-units:
//
//	h := units.NewEquivalentDose(2.4, equivalentdose.Millisievert)
//
// ParseUnit turns a symbol such as "mSv" into a Unit, and String turns a
// Unit into its symbol.
package equivalentdose

// Unit is an equivalent dose unit: one of the constants below. String returns
// its symbol and ParseUnit finds the unit of a symbol.
//
// The zero Unit is the base unit, Sievert. The numbers behind the constants may
// change between versions: store and exchange units by their symbols.
type Unit uint8

// BaseUnit is the SI unit of equivalent dose that all factors are relative to.
const BaseUnit = Sievert

// Equivalent dose units, with their symbols.
const (
	Sievert      Unit = iota // Sv
	Millisievert             // mSv
	Microsievert             // μSv: with U+03BC GREEK SMALL LETTER MU
	Nanosievert              // nSv
	Rem                      // rem
	Millirem                 // mrem

	unitCount // the number of units; keep it last
)

// definitions holds the symbol of each unit and its factor: how many base
// units (Sv) there are in one unit.
//
// Source: the SI Brochure (9th edition); the rem from NIST SP 811.
var definitions = [unitCount]definition{
	Sievert:      {"Sv", 1},
	Millisievert: {"mSv", 1e-3},
	Microsievert: {"μSv", 1e-6},
	Nanosievert:  {"nSv", 1e-9},
	Rem:          {"rem", 1e-2},
	Millirem:     {"mrem", 1e-5},
}
