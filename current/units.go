// Package current defines the electric current units of the units package and their
// conversion factors to the base unit, the ampere.
package current

// Unit is an electric current unit: one of the constants below. String returns its
// symbol and ParseUnit finds the unit of a symbol.
//
// The zero Unit is the base unit, Ampere. The numbers behind the constants may
// change between versions: store and exchange units by their symbols.
type Unit uint8

// BaseUnit is the SI unit of electric current that all factors are relative to.
const BaseUnit = Ampere

// Current units, with their symbols.
const (
	Ampere      Unit = iota // A
	Microampere             // μA, with U+03BC GREEK SMALL LETTER MU; "µA" (U+00B5 MICRO SIGN) is accepted too
	Milliampere             // mA
	Kiloampere              // kA

	unitCount // the number of units; keep it last
)

// definitions holds the symbol of each unit and its factor: how many base
// units (A) there are in one unit.
//
// Source: the SI Brochure (9th edition); the ampere is an SI base unit.
var definitions = [unitCount]definition{
	Ampere:      {"A", 1},
	Microampere: {"μA", 1e-6},
	Milliampere: {"mA", 1e-3},
	Kiloampere:  {"kA", 1e3},
}
