// Package mass defines the mass units of the units package and
// their conversion factors to the base unit, the kilogram.
package mass

// Unit is a mass unit: one of the constants below. String returns its
// symbol and ParseUnit finds the unit of a symbol.
//
// The zero Unit is the base unit, Kilogram. The numbers behind the constants may
// change between versions: store and exchange units by their symbols.
type Unit uint8

// BaseUnit is the SI unit of mass that all mass factors are relative to.
const BaseUnit = Kilogram

// Mass units, with their symbols.
const (
	Kilogram  Unit = iota // kg
	Dalton                // Da
	Microgram             // μg, with U+03BC GREEK SMALL LETTER MU; "µg" (U+00B5 MICRO SIGN) is accepted too
	Milligram             // mg
	Gram                  // g
	Ounce                 // oz
	Pound                 // lb
	Quintal               // q: metric quintal, 100 kg
	Tonne                 // t
	EarthMass             // M⊕
	SolarMass             // M☉

	unitCount // the number of units; keep it last
)

// definitions holds the symbol of each unit and its factor: how many base
// units (kg) there are in one unit.
var definitions = [unitCount]definition{
	Kilogram:  {"kg", 1},
	Dalton:    {"Da", 1.6605390666e-27}, // CODATA 2018
	Microgram: {"μg", 1e-9},
	Milligram: {"mg", 1e-6},
	Gram:      {"g", 0.001},
	Ounce:     {"oz", 0.028349523125}, // avoirdupois ounce, exactly 1/16 lb
	Pound:     {"lb", 0.45359237},     // international avoirdupois pound, exact
	Quintal:   {"q", 100},
	Tonne:     {"t", 1000},
	EarthMass: {"M⊕", 5.9722e24},
	SolarMass: {"M☉", 1.98847e30},
}
