// Package capacitance defines the capacitance units of the units package and
// their conversion factors to the base unit, the farad.
//
// Use the units with units.NewCapacitance, from the root package
// github.com/asaphin/go-physics-units:
//
//	c := units.NewCapacitance(100, capacitance.Microfarad)
//
// ParseUnit turns a symbol such as "μF" into a Unit, and String turns a
// Unit into its symbol.
package capacitance

// Unit is a capacitance unit: one of the constants below. String returns its
// symbol and ParseUnit finds the unit of a symbol.
//
// The zero Unit is the base unit, Farad. The numbers behind the constants may
// change between versions: store and exchange units by their symbols.
type Unit uint8

// BaseUnit is the SI unit of capacitance that all factors are relative to.
const BaseUnit = Farad

// Capacitance units, with their symbols.
const (
	Farad      Unit = iota // F
	Millifarad             // mF
	Microfarad             // μF: with U+03BC GREEK SMALL LETTER MU
	Nanofarad              // nF
	Picofarad              // pF

	unitCount // the number of units; keep it last
)

// definitions holds the symbol of each unit and its factor: how many base
// units (F) there are in one unit.
//
// Source: the SI Brochure (9th edition).
var definitions = [unitCount]definition{
	Farad:      {"F", 1},
	Millifarad: {"mF", 1e-3},
	Microfarad: {"μF", 1e-6},
	Nanofarad:  {"nF", 1e-9},
	Picofarad:  {"pF", 1e-12},
}
