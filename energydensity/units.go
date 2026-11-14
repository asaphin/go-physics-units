// Package energydensity defines the energy density (energy per volume) units of
// the units package and their conversion factors to the base unit, the joule
// per cubic metre.
package energydensity

// Unit is an energy density unit: one of the constants below. String returns
// its symbol and ParseUnit finds the unit of a symbol.
//
// The zero Unit is the base unit, JoulePerCubicMeter. The numbers behind the
// constants may change between versions: store and exchange units by their
// symbols.
type Unit uint8

// BaseUnit is the SI unit of energy density that all factors are relative to.
const BaseUnit = JoulePerCubicMeter

// Energy density units, with their symbols.
const (
	JoulePerCubicMeter        Unit = iota // J/m³
	KilojoulePerCubicMeter                // kJ/m³
	MegajoulePerCubicMeter                // MJ/m³
	MegajoulePerLiter                     // MJ/L
	WattHourPerLiter                      // Wh/L
	KilowattHourPerCubicMeter             // kWh/m³
	KilowattHourPerLiter                  // kWh/L

	unitCount // the number of units; keep it last
)

// definitions holds the symbol of each unit and its factor: how many base
// units (J/m³) there are in one unit.
//
// Source: the SI Brochure (9th edition).
var definitions = [unitCount]definition{
	JoulePerCubicMeter:        {"J/m³", 1},
	KilojoulePerCubicMeter:    {"kJ/m³", 1e3},
	MegajoulePerCubicMeter:    {"MJ/m³", 1e6},
	MegajoulePerLiter:         {"MJ/L", 1e6 / 1e-3},
	WattHourPerLiter:          {"Wh/L", 3600 / 1e-3},
	KilowattHourPerCubicMeter: {"kWh/m³", 3.6e6},
	KilowattHourPerLiter:      {"kWh/L", 3.6e6 / 1e-3},
}
