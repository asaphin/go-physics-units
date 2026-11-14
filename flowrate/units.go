// Package flowrate defines the volumetric flow rate units of the units package
// and their conversion factors to the base unit, the cubic metre per second.
//
// Use the units with units.NewFlowRate, from the root package
// github.com/asaphin/go-physics-units:
//
//	q := units.NewFlowRate(12, flowrate.LiterPerMinute)
//
// ParseUnit turns a symbol such as "L/min" into a Unit, and String turns a
// Unit into its symbol.
package flowrate

// Unit is a flow rate unit: one of the constants below. String returns its
// symbol and ParseUnit finds the unit of a symbol.
//
// The zero Unit is the base unit, CubicMeterPerSecond. The numbers behind the
// constants may change between versions: store and exchange units by their
// symbols.
type Unit uint8

// BaseUnit is the SI unit of flow rate that all factors are relative to.
const BaseUnit = CubicMeterPerSecond

// Flow rate units, with their symbols.
const (
	CubicMeterPerSecond Unit = iota // m³/s
	CubicMeterPerHour               // m³/h
	LiterPerSecond                  // L/s
	LiterPerMinute                  // L/min
	LiterPerHour                    // L/h
	MilliliterPerMinute             // mL/min
	MilliliterPerHour               // mL/h
	GallonPerMinute                 // gal/min: US liquid gallons per minute, gpm
	CubicFootPerMinute              // ft³/min: cfm
	CubicFootPerSecond              // ft³/s: cfs

	unitCount // the number of units; keep it last
)

// definitions holds the symbol of each unit and its factor: how many base
// units (m³/s) there are in one unit.
//
// Source: NIST SP 811, appendix B; the US liquid gallon is 231 in³.
var definitions = [unitCount]definition{
	CubicMeterPerSecond: {"m³/s", 1},
	CubicMeterPerHour:   {"m³/h", 1.0 / 3600},
	LiterPerSecond:      {"L/s", 1e-3},
	LiterPerMinute:      {"L/min", 1e-3 / 60},
	LiterPerHour:        {"L/h", 1e-3 / 3600},
	MilliliterPerMinute: {"mL/min", 1e-6 / 60},
	MilliliterPerHour:   {"mL/h", 1e-6 / 3600},
	GallonPerMinute:     {"gal/min", 231 * 0.0254 * 0.0254 * 0.0254 / 60},
	CubicFootPerMinute:  {"ft³/min", 0.3048 * 0.3048 * 0.3048 / 60},
	CubicFootPerSecond:  {"ft³/s", 0.3048 * 0.3048 * 0.3048},
}
