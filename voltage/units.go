// Package voltage defines the voltage units of the units package and their
// conversion factors to the base unit, the volt.
package voltage

// Unit is a voltage unit: one of the constants below. String returns its
// symbol and ParseUnit finds the unit of a symbol.
//
// The zero Unit is the base unit, Volt. The numbers behind the constants may
// change between versions: store and exchange units by their symbols.
type Unit uint8

// BaseUnit is the SI unit of voltage that all factors are relative to.
const BaseUnit = Volt

// Voltage units, with their symbols.
const (
	Volt      Unit = iota // V
	Microvolt             // μV, with U+03BC GREEK SMALL LETTER MU; "µV" (U+00B5 MICRO SIGN) is accepted too
	Millivolt             // mV
	Kilovolt              // kV
	Megavolt              // MV

	unitCount // the number of units; keep it last
)

// definitions holds the symbol of each unit and its factor: how many base
// units (V) there are in one unit.
//
// Source: the SI Brochure (9th edition).
var definitions = [unitCount]definition{
	Volt:      {"V", 1},
	Microvolt: {"μV", 1e-6},
	Millivolt: {"mV", 1e-3},
	Kilovolt:  {"kV", 1e3},
	Megavolt:  {"MV", 1e6},
}
