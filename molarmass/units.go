// Package molarmass defines the molar mass units of the units package and their
// conversion factors to the base unit, the kilogram per mole.
//
// Use the units with units.NewMolarMass, from the root package
// github.com/asaphin/go-physics-units:
//
//	mm := units.NewMolarMass(18.015, molarmass.GramPerMole)
//
// ParseUnit turns a symbol such as "g/mol" into a Unit, and String turns a
// Unit into its symbol.
package molarmass

// Unit is a molar mass unit: one of the constants below. String returns its
// symbol and ParseUnit finds the unit of a symbol.
//
// The zero Unit is the base unit, KilogramPerMole. The numbers behind the
// constants may change between versions: store and exchange units by their
// symbols.
type Unit uint8

// BaseUnit is the SI unit of molar mass that all factors are relative to.
const BaseUnit = KilogramPerMole

// Molar mass units, with their symbols.
const (
	KilogramPerMole Unit = iota // kg/mol
	GramPerMole                 // g/mol: numerically equal to the mass of a molecule in daltons

	unitCount // the number of units; keep it last
)

// definitions holds the symbol of each unit and its factor: how many base
// units (kg/mol) there are in one unit.
//
// Source: the SI Brochure (9th edition).
var definitions = [unitCount]definition{
	KilogramPerMole: {"kg/mol", 1},
	GramPerMole:     {"g/mol", 1e-3},
}
