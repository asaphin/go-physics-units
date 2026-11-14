// Package magneticfield defines the magnetic flux density units of the units
// package and their conversion factors to the base unit, the tesla.
package magneticfield

// Unit is a magnetic flux density unit: one of the constants below. String
// returns its symbol and ParseUnit finds the unit of a symbol.
//
// The zero Unit is the base unit, Tesla. The numbers behind the constants may
// change between versions: store and exchange units by their symbols.
type Unit uint8

// BaseUnit is the SI unit of magnetic flux density that all factors are
// relative to.
const BaseUnit = Tesla

// Magnetic flux density units, with their symbols.
const (
	Tesla      Unit = iota // T
	Millitesla             // mT
	Microtesla             // μT: with U+03BC GREEK SMALL LETTER MU
	Nanotesla              // nT
	Gauss                  // G
	Milligauss             // mG

	unitCount // the number of units; keep it last
)

// definitions holds the symbol of each unit and its factor: how many base
// units (T) there are in one unit.
//
// Source: the SI Brochure (9th edition); the gauss is 1e-4 T (NIST SP 811).
var definitions = [unitCount]definition{
	Tesla:      {"T", 1},
	Millitesla: {"mT", 1e-3},
	Microtesla: {"μT", 1e-6},
	Nanotesla:  {"nT", 1e-9},
	Gauss:      {"G", 1e-4},
	Milligauss: {"mG", 1e-7},
}
