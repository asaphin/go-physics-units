// Package fuelconsumption defines the fuel consumption (volume per distance)
// units of the units package and their conversion factors to the base unit, the
// cubic metre per metre.
//
// Use the units with units.NewFuelConsumption, from the root package
// github.com/asaphin/go-physics-units:
//
//	fc := units.NewFuelConsumption(6.5, fuelconsumption.LiterPer100Kilometers)
//
// ParseUnit turns a symbol such as "L/100km" into a Unit, and String turns a
// Unit into its symbol.
package fuelconsumption

// Unit is a fuel consumption unit: one of the constants below. String returns
// its symbol and ParseUnit finds the unit of a symbol.
//
// The zero Unit is the base unit, CubicMeterPerMeter. The numbers behind the
// constants may change between versions: store and exchange units by their
// symbols.
type Unit uint8

// BaseUnit is the SI unit of fuel consumption that all factors are relative to.
const BaseUnit = CubicMeterPerMeter

// Fuel consumption units, with their symbols.
const (
	CubicMeterPerMeter    Unit = iota // m³/m
	LiterPer100Kilometers             // L/100km
	LiterPerKilometer                 // L/km
	GallonPer100Miles                 // gal/100mi: US liquid gallons per 100 miles

	unitCount // the number of units; keep it last
)

// definitions holds the symbol of each unit and its factor: how many base
// units (m³/m) there are in one unit.
//
// Source: NIST SP 811, appendix B; the US liquid gallon is 231 in³.
var definitions = [unitCount]definition{
	CubicMeterPerMeter:    {"m³/m", 1},
	LiterPer100Kilometers: {"L/100km", 1e-3 / 100e3},
	LiterPerKilometer:     {"L/km", 1e-3 / 1e3},
	GallonPer100Miles:     {"gal/100mi", 231 * 0.0254 * 0.0254 * 0.0254 / (100 * 1609.344)},
}
