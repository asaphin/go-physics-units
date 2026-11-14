// Package power defines the power units of the units package
// and their conversion factors to the base unit, the watt.
//
// Use the units with units.NewPower, from the root package
// github.com/asaphin/go-physics-units:
//
//	p := units.NewPower(100, power.Watt)
//
// ParseUnit turns a symbol such as "W" into a Unit, and String turns a
// Unit into its symbol.
package power

// Unit is a power unit: one of the constants below. String returns its
// symbol and ParseUnit finds the unit of a symbol.
//
// The zero Unit is the base unit, Watt. The numbers behind the constants may
// change between versions: store and exchange units by their symbols.
type Unit uint8

// BaseUnit is the SI unit of power that all power factors are relative to.
const BaseUnit = Watt

// Power units, with their symbols.
const (
	Watt               Unit = iota // W
	Milliwatt                      // mW
	BTUPerHour                     // Btu/h: International Table British thermal unit per hour
	MetricHorsepower               // PS: 75 kgf·m/s
	Horsepower                     // hp: mechanical horsepower, 550 ft·lbf/s
	Kilowatt                       // kW
	Megawatt                       // MW
	Gigawatt                       // GW
	TonOfRefrigeration             // TR: 12000 Btu/h, for air conditioning
	SolarLuminosity                // L☉: nominal, IAU 2015 Resolution B3

	unitCount // the number of units; keep it last
)

// definitions holds the symbol of each unit and its factor: how many base
// units (W) there are in one unit.
//
// The factors are exact: defined values, written as constant expressions that Go
// evaluates exactly before rounding once.
//
// Source: NIST SP 811, appendix B.
var definitions = [unitCount]definition{
	Watt:               {"W", 1},
	Milliwatt:          {"mW", 1e-3},
	BTUPerHour:         {"Btu/h", 1055.05585262 / 3600},
	MetricHorsepower:   {"PS", 75 * 9.80665},
	Horsepower:         {"hp", 550 * 0.3048 * 0.45359237 * 9.80665},
	Kilowatt:           {"kW", 1e3},
	Megawatt:           {"MW", 1e6},
	Gigawatt:           {"GW", 1e9},
	TonOfRefrigeration: {"TR", 12000 * 1055.05585262 / 3600},
	SolarLuminosity:    {"L☉", 3.828e26},
}
