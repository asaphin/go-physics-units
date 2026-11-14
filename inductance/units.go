// Package inductance defines the inductance units of the units package and
// their conversion factors to the base unit, the henry.
//
// Use the units with units.NewInductance, from the root package
// github.com/asaphin/go-physics-units:
//
//	l := units.NewInductance(10, inductance.Millihenry)
//
// ParseUnit turns a symbol such as "mH" into a Unit, and String turns a
// Unit into its symbol.
package inductance

// Unit is an inductance unit: one of the constants below. String returns its
// symbol and ParseUnit finds the unit of a symbol.
//
// The zero Unit is the base unit, Henry. The numbers behind the constants may
// change between versions: store and exchange units by their symbols.
type Unit uint8

// BaseUnit is the SI unit of inductance that all factors are relative to.
const BaseUnit = Henry

// Inductance units, with their symbols.
const (
	Henry      Unit = iota // H
	Millihenry             // mH
	Microhenry             // μH: with U+03BC GREEK SMALL LETTER MU
	Nanohenry              // nH

	unitCount // the number of units; keep it last
)

// definitions holds the symbol of each unit and its factor: how many base
// units (H) there are in one unit.
//
// Source: the SI Brochure (9th edition).
var definitions = [unitCount]definition{
	Henry:      {"H", 1},
	Millihenry: {"mH", 1e-3},
	Microhenry: {"μH", 1e-6},
	Nanohenry:  {"nH", 1e-9},
}
