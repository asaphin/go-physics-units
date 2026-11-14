// Package fueleconomy defines the fuel economy (distance per volume) units of
// the units package and their conversion factors to the base unit, the metre
// per cubic metre.
package fueleconomy

// Unit is a fuel economy unit: one of the constants below. String returns its
// symbol and ParseUnit finds the unit of a symbol.
//
// The zero Unit is the base unit, MeterPerCubicMeter. The numbers behind the
// constants may change between versions: store and exchange units by their
// symbols.
type Unit uint8

// BaseUnit is the SI unit of fuel economy that all factors are relative to.
const BaseUnit = MeterPerCubicMeter

// Fuel economy units, with their symbols.
const (
	MeterPerCubicMeter    Unit = iota // m/m³
	KilometerPerLiter                 // km/L
	MilePerGallon                     // mpg: miles per US liquid gallon
	MilePerImperialGallon             // mpg(imp): miles per imperial gallon, 4.54609 L

	unitCount // the number of units; keep it last
)

// definitions holds the symbol of each unit and its factor: how many base
// units (m/m³) there are in one unit.
//
// Source: NIST SP 811, appendix B; the US liquid gallon is 231 in³, the
// imperial gallon 4.54609 L.
var definitions = [unitCount]definition{
	MeterPerCubicMeter:    {"m/m³", 1},
	KilometerPerLiter:     {"km/L", 1000 / 1e-3},
	MilePerGallon:         {"mpg", 1609.344 / (231 * 0.0254 * 0.0254 * 0.0254)},
	MilePerImperialGallon: {"mpg(imp)", 1609.344 / 4.54609e-3},
}
