// Package radioactivity defines the radioactivity (activity) units of the units
// package and their conversion factors to the base unit, the becquerel.
//
// Use the units with units.NewRadioactivity, from the root package
// github.com/asaphin/go-physics-units:
//
//	a := units.NewRadioactivity(3.7, radioactivity.Gigabecquerel)
//
// ParseUnit turns a symbol such as "GBq" into a Unit, and String turns a
// Unit into its symbol.
package radioactivity

// Unit is a radioactivity unit: one of the constants below. String returns its
// symbol and ParseUnit finds the unit of a symbol.
//
// The zero Unit is the base unit, Becquerel. The numbers behind the constants
// may change between versions: store and exchange units by their symbols.
type Unit uint8

// BaseUnit is the SI unit of radioactivity that all factors are relative to.
const BaseUnit = Becquerel

// Radioactivity units, with their symbols.
const (
	Becquerel     Unit = iota // Bq
	Kilobecquerel             // kBq
	Megabecquerel             // MBq
	Gigabecquerel             // GBq
	Terabecquerel             // TBq
	Curie                     // Ci
	Millicurie                // mCi
	Microcurie                // μCi: with U+03BC GREEK SMALL LETTER MU
	Nanocurie                 // nCi

	unitCount // the number of units; keep it last
)

// definitions holds the symbol of each unit and its factor: how many base
// units (Bq) there are in one unit.
//
// Source: the SI Brochure (9th edition); the curie is 3.7e10 Bq exactly (NIST
// SP 811).
var definitions = [unitCount]definition{
	Becquerel:     {"Bq", 1},
	Kilobecquerel: {"kBq", 1e3},
	Megabecquerel: {"MBq", 1e6},
	Gigabecquerel: {"GBq", 1e9},
	Terabecquerel: {"TBq", 1e12},
	Curie:         {"Ci", 3.7e10},
	Millicurie:    {"mCi", 3.7e7},
	Microcurie:    {"μCi", 3.7e4},
	Nanocurie:     {"nCi", 37},
}
