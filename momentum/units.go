// Package momentum defines the momentum and impulse units of the units package
// and their conversion factors to the base unit, the kilogram metre per second.
//
// Use the units with units.NewMomentum, from the root package
// github.com/asaphin/go-physics-units:
//
//	p := units.NewMomentum(10, momentum.NewtonSecond)
//
// ParseUnit turns a symbol such as "N·s" into a Unit, and String turns a
// Unit into its symbol.
package momentum

// Unit is a momentum unit: one of the constants below. String returns its
// symbol and ParseUnit finds the unit of a symbol.
//
// The zero Unit is the base unit, KilogramMeterPerSecond. The numbers behind
// the constants may change between versions: store and exchange units by their
// symbols.
type Unit uint8

// BaseUnit is the SI unit of momentum that all factors are relative to.
const BaseUnit = KilogramMeterPerSecond

// Momentum units, with their symbols.
const (
	KilogramMeterPerSecond  Unit = iota // kg·m/s
	NewtonSecond                        // N·s: the unit of impulse, equal to the base unit
	GramCentimeterPerSecond             // g·cm/s
	PoundForceSecond                    // lbf·s
	PoundFootPerSecond                  // lb·ft/s

	unitCount // the number of units; keep it last
)

// definitions holds the symbol of each unit and its factor: how many base
// units (kg·m/s) there are in one unit.
//
// Source: NIST SP 811, appendix B; the pound-force under standard gravity.
var definitions = [unitCount]definition{
	KilogramMeterPerSecond:  {"kg·m/s", 1},
	NewtonSecond:            {"N·s", 1},
	GramCentimeterPerSecond: {"g·cm/s", 1e-5},
	PoundForceSecond:        {"lbf·s", 0.45359237 * 9.80665},
	PoundFootPerSecond:      {"lb·ft/s", 0.45359237 * 0.3048},
}
