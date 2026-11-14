// Package angularvelocity defines the angular velocity units of the units package and their
// conversion factors to the base unit, the radian per second.
package angularvelocity

import "math"

// Unit is an angular velocity unit: one of the constants below. String returns its
// symbol and ParseUnit finds the unit of a symbol.
//
// The zero Unit is the base unit, RadianPerSecond. The numbers behind the constants may
// change between versions: store and exchange units by their symbols.
type Unit uint8

// BaseUnit is the SI unit of angular velocity that all factors are relative to.
const BaseUnit = RadianPerSecond

// Angular velocity units, with their symbols.
const (
	RadianPerSecond     Unit = iota // rad/s
	DegreePerSecond                 // °/s: with U+00B0 DEGREE SIGN
	RevolutionPerMinute             // rev/min: not frequency.RevolutionPerMinute, "rpm", which counts cycles, not radians
	RevolutionPerSecond             // rev/s

	unitCount // the number of units; keep it last
)

// definitions holds the symbol of each unit and its factor: how many base
// units (rad/s) there are in one unit.
//
// The degree is π/180 rad and the revolution 2π rad exactly (SI Brochure, table 8).
var definitions = [unitCount]definition{
	RadianPerSecond:     {"rad/s", 1},
	DegreePerSecond:     {"°/s", math.Pi / 180},
	RevolutionPerMinute: {"rev/min", 2 * math.Pi / 60},
	RevolutionPerSecond: {"rev/s", 2 * math.Pi},
}
