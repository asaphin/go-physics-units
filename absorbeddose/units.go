// Package absorbeddose defines the absorbed dose of ionising radiation units of
// the units package and their conversion factors to the base unit, the gray.
package absorbeddose

// Unit is an absorbed dose unit: one of the constants below. String returns its
// symbol and ParseUnit finds the unit of a symbol.
//
// The zero Unit is the base unit, Gray. The numbers behind the constants may
// change between versions: store and exchange units by their symbols.
type Unit uint8

// BaseUnit is the SI unit of absorbed dose that all factors are relative to.
const BaseUnit = Gray

// Absorbed dose units, with their symbols.
const (
	Gray      Unit = iota // Gy
	Milligray             // mGy
	Microgray             // μGy: with U+03BC GREEK SMALL LETTER MU
	Centigray             // cGy
	Rad                   // rd: 100 erg/g; the symbol rd avoids a clash with the radian, "rad"

	unitCount // the number of units; keep it last
)

// definitions holds the symbol of each unit and its factor: how many base
// units (Gy) there are in one unit.
//
// Source: the SI Brochure (9th edition); the rad from NIST SP 811.
var definitions = [unitCount]definition{
	Gray:      {"Gy", 1},
	Milligray: {"mGy", 1e-3},
	Microgray: {"μGy", 1e-6},
	Centigray: {"cGy", 1e-2},
	Rad:       {"rd", 1e-2},
}
