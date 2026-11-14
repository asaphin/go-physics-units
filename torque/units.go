// Package torque defines the torque units of the units package and their
// conversion factors to the base unit, the newton metre.
package torque

// Unit is a torque unit: one of the constants below. String returns its symbol
// and ParseUnit finds the unit of a symbol.
//
// The zero Unit is the base unit, NewtonMeter. The numbers behind the constants
// may change between versions: store and exchange units by their symbols.
type Unit uint8

// BaseUnit is the SI unit of torque that all factors are relative to.
const BaseUnit = NewtonMeter

// Torque units, with their symbols.
const (
	NewtonMeter             Unit = iota // N·m
	NewtonCentimeter                    // N·cm
	MillinewtonMeter                    // mN·m
	KilonewtonMeter                     // kN·m
	KilogramForceMeter                  // kgf·m
	KilogramForceCentimeter             // kgf·cm
	PoundForceFoot                      // lbf·ft: not to be confused with energy.FootPound, "ft·lbf"
	PoundForceInch                      // lbf·in
	OunceForceInch                      // ozf·in

	unitCount // the number of units; keep it last
)

// definitions holds the symbol of each unit and its factor: how many base
// units (N·m) there are in one unit.
//
// Source: NIST SP 811, appendix B; the pound-force and the kilogram-force under
// standard gravity.
var definitions = [unitCount]definition{
	NewtonMeter:             {"N·m", 1},
	NewtonCentimeter:        {"N·cm", 0.01},
	MillinewtonMeter:        {"mN·m", 1e-3},
	KilonewtonMeter:         {"kN·m", 1e3},
	KilogramForceMeter:      {"kgf·m", 9.80665},
	KilogramForceCentimeter: {"kgf·cm", 9.80665 / 100},
	PoundForceFoot:          {"lbf·ft", 0.45359237 * 9.80665 * 0.3048},
	PoundForceInch:          {"lbf·in", 0.45359237 * 9.80665 * 0.0254},
	OunceForceInch:          {"ozf·in", 0.45359237 * 9.80665 / 16 * 0.0254},
}
