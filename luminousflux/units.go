// Package luminousflux defines the luminous flux units of the units package and
// their conversion factors to the base unit, the lumen.
//
// Use the units with units.NewLuminousFlux, from the root package
// github.com/asaphin/go-physics-units:
//
//	f := units.NewLuminousFlux(800, luminousflux.Lumen)
//
// ParseUnit turns a symbol such as "lm" into a Unit, and String turns a
// Unit into its symbol.
package luminousflux

// Unit is a luminous flux unit: one of the constants below. String returns its
// symbol and ParseUnit finds the unit of a symbol.
//
// The zero Unit is the base unit, Lumen. The numbers behind the constants may
// change between versions: store and exchange units by their symbols.
type Unit uint8

// BaseUnit is the SI unit of luminous flux that all factors are relative to.
const BaseUnit = Lumen

// Luminous flux units, with their symbols.
const (
	Lumen     Unit = iota // lm
	Kilolumen             // klm

	unitCount // the number of units; keep it last
)

// definitions holds the symbol of each unit and its factor: how many base
// units (lm) there are in one unit.
//
// Source: the SI Brochure (9th edition).
var definitions = [unitCount]definition{
	Lumen:     {"lm", 1},
	Kilolumen: {"klm", 1e3},
}
