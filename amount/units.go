// Package amount defines the amount of substance units of the units package and
// their conversion factors to the base unit, the mole.
package amount

// Unit is an amount of substance unit: one of the constants below. String
// returns its symbol and ParseUnit finds the unit of a symbol.
//
// The zero Unit is the base unit, Mole. The numbers behind the constants may
// change between versions: store and exchange units by their symbols.
type Unit uint8

// BaseUnit is the SI unit of amount of substance that all factors are relative
// to.
const BaseUnit = Mole

// Amount of substance units, with their symbols.
const (
	Mole      Unit = iota // mol
	Millimole             // mmol
	Micromole             // μmol: with U+03BC GREEK SMALL LETTER MU; "µmol" (U+00B5 MICRO SIGN) is accepted too
	Nanomole              // nmol
	Picomole              // pmol
	Kilomole              // kmol

	unitCount // the number of units; keep it last
)

// definitions holds the symbol of each unit and its factor: how many base
// units (mol) there are in one unit.
//
// Source: the SI Brochure (9th edition).
var definitions = [unitCount]definition{
	Mole:      {"mol", 1},
	Millimole: {"mmol", 1e-3},
	Micromole: {"μmol", 1e-6},
	Nanomole:  {"nmol", 1e-9},
	Picomole:  {"pmol", 1e-12},
	Kilomole:  {"kmol", 1e3},
}
