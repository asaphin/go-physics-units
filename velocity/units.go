// Package velocity defines the velocity (speed) units of the units package and
// their conversion factors to the base unit, the metre per second.
package velocity

// Unit is a velocity unit: one of the constants below. String returns its
// symbol and ParseUnit finds the unit of a symbol.
//
// The zero Unit is the base unit, MeterPerSecond. The numbers behind the constants may
// change between versions: store and exchange units by their symbols.
type Unit uint8

// BaseUnit is the SI unit of velocity that all velocity factors are relative to.
const BaseUnit = MeterPerSecond

// Velocity units, with their symbols.
const (
	MeterPerSecond     Unit = iota // m/s
	KilometerPerHour               // km/h
	FeetPerSecond                  // ft/s
	MilePerHour                    // mph
	Knot                           // kn
	SpeedOfLight                   // c
	KilometerPerSecond             // km/s

	unitCount // the number of units; keep it last
)

// definitions holds the symbol of each unit and its factor: how many base
// units (m/s) there are in one unit.
var definitions = [unitCount]definition{
	MeterPerSecond:     {"m/s", 1},
	KilometerPerHour:   {"km/h", 1000.0 / 3600},
	FeetPerSecond:      {"ft/s", 0.3048},
	MilePerHour:        {"mph", 0.44704},
	Knot:               {"kn", 1852.0 / 3600}, // one nautical mile per hour
	SpeedOfLight:       {"c", 299792458},      // exact by definition
	KilometerPerSecond: {"km/s", 1000},
}
