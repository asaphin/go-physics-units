// Package resistance defines the electrical resistance units of the units package and their
// conversion factors to the base unit, the ohm.
//
// Use the units with units.NewResistance, from the root package
// github.com/asaphin/go-physics-units:
//
//	r := units.NewResistance(4.7, resistance.Kiloohm)
//
// ParseUnit turns a symbol such as "kΩ" into a Unit, and String turns a
// Unit into its symbol.
package resistance

// Unit is an electrical resistance unit: one of the constants below. String returns its
// symbol and ParseUnit finds the unit of a symbol.
//
// The zero Unit is the base unit, Ohm. The numbers behind the constants may
// change between versions: store and exchange units by their symbols.
type Unit uint8

// BaseUnit is the SI unit of electrical resistance that all factors are relative to.
const BaseUnit = Ohm

// Resistance units, with their symbols.
const (
	Ohm      Unit = iota // Ω, with U+03A9 GREEK CAPITAL LETTER OMEGA; "Ω" (U+2126 OHM SIGN) is accepted too
	Milliohm             // mΩ
	Kiloohm              // kΩ
	Megaohm              // MΩ
	Gigaohm              // GΩ

	unitCount // the number of units; keep it last
)

// definitions holds the symbol of each unit and its factor: how many base
// units (Ω) there are in one unit.
//
// Source: the SI Brochure (9th edition).
var definitions = [unitCount]definition{
	Ohm:      {"Ω", 1},
	Milliohm: {"mΩ", 1e-3},
	Kiloohm:  {"kΩ", 1e3},
	Megaohm:  {"MΩ", 1e6},
	Gigaohm:  {"GΩ", 1e9},
}
