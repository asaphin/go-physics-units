// Package angle defines the angle units of the units package
// and their conversion factors to the base unit, the radian.
//
// Use the units with units.NewAngle, from the root package
// github.com/asaphin/go-physics-units:
//
//	a := units.NewAngle(90, angle.Degree)
//
// ParseUnit turns a symbol such as "°" into a Unit, and String turns a
// Unit into its symbol.
package angle

import "math"

// Unit is an angle unit: one of the constants below. String returns its
// symbol and ParseUnit finds the unit of a symbol.
//
// The zero Unit is the base unit, Radian. The numbers behind the constants may
// change between versions: store and exchange units by their symbols.
type Unit uint8

// BaseUnit is the SI unit of angle that all angle factors are relative to.
const BaseUnit = Radian

// Angle units, with their symbols.
const (
	Radian      Unit = iota // rad
	Arcsecond               // arcsec: 1/3600 degree
	Arcminute               // arcmin: 1/60 degree
	Milliradian             // mrad
	Gradian                 // gon: 1/400 revolution
	Degree                  // °: U+00B0 DEGREE SIGN
	Revolution              // rev: one full turn

	unitCount // the number of units; keep it last
)

// definitions holds the symbol of each unit and its factor: how many base
// units (rad) there are in one unit.
//
// The factors are exact: defined values, written as constant expressions that Go
// evaluates exactly before rounding once.
//
// Source: the SI Brochure (9th edition), table 8. math.Pi is an untyped
// constant with far more precision than a float64, so each factor is rounded
// once, like the others.
var definitions = [unitCount]definition{
	Radian:      {"rad", 1},
	Arcsecond:   {"arcsec", math.Pi / 648000},
	Arcminute:   {"arcmin", math.Pi / 10800},
	Milliradian: {"mrad", 1e-3},
	Gradian:     {"gon", math.Pi / 200},
	Degree:      {"°", math.Pi / 180},
	Revolution:  {"rev", 2 * math.Pi},
}
