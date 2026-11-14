// Package density defines the density units of the units package
// and their conversion factors to the base unit, the kilogram per cubic metre.
//
// Use the units with units.NewDensity, from the root package
// github.com/asaphin/go-physics-units:
//
//	d := units.NewDensity(1.2, density.KilogramPerCubicMeter)
//
// ParseUnit turns a symbol such as "kg/m³" into a Unit, and String turns a
// Unit into its symbol.
package density

// Unit is a density unit: one of the constants below. String returns its
// symbol and ParseUnit finds the unit of a symbol.
//
// The zero Unit is the base unit, KilogramPerCubicMeter. The numbers behind the constants may
// change between versions: store and exchange units by their symbols.
type Unit uint8

// BaseUnit is the SI unit of density that all density factors are relative to.
const BaseUnit = KilogramPerCubicMeter

// Density units, with their symbols.
const (
	KilogramPerCubicMeter  Unit = iota // kg/m³
	GramPerLiter                       // g/L: equal to the base unit
	PoundPerCubicFoot                  // lb/ft³
	PoundPerGallon                     // lb/gal: per US liquid gallon
	KilogramPerLiter                   // kg/L
	GramPerMilliliter                  // g/mL
	GramPerCubicCentimeter             // g/cm³

	unitCount // the number of units; keep it last
)

// definitions holds the symbol of each unit and its factor: how many base
// units (kg/m³) there are in one unit.
//
// The factors are exact: defined values, written as constant expressions that Go
// evaluates exactly before rounding once.
//
// Source: NIST SP 811, appendix B.
var definitions = [unitCount]definition{
	KilogramPerCubicMeter:  {"kg/m³", 1},
	GramPerLiter:           {"g/L", 1},
	PoundPerCubicFoot:      {"lb/ft³", 0.45359237 / (0.3048 * 0.3048 * 0.3048)},
	PoundPerGallon:         {"lb/gal", 0.45359237 / (231 * 0.0254 * 0.0254 * 0.0254)},
	KilogramPerLiter:       {"kg/L", 1000},
	GramPerMilliliter:      {"g/mL", 1000},
	GramPerCubicCentimeter: {"g/cm³", 1000},
}
