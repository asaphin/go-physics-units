// Package force defines the force units of the units package
// and their conversion factors to the base unit, the newton.
package force

// Unit is a force unit: one of the constants below. String returns its
// symbol and ParseUnit finds the unit of a symbol.
//
// The zero Unit is the base unit, Newton. The numbers behind the constants may
// change between versions: store and exchange units by their symbols.
type Unit uint8

// BaseUnit is the SI unit of force that all force factors are relative to.
const BaseUnit = Newton

// Force units, with their symbols.
const (
	Newton        Unit = iota // N
	Dyne                      // dyn
	Millinewton               // mN
	PoundForce                // lbf: the avoirdupois pound under standard gravity
	KilogramForce             // kgf: the kilogram under standard gravity
	Kilonewton                // kN: not to be confused with velocity.Knot, "kn"
	Meganewton                // MN
	Kip                       // kip: 1000 lbf
	OunceForce                // ozf: 1/16 lbf
	TonneForce                // tf: the tonne under standard gravity

	unitCount // the number of units; keep it last
)

// definitions holds the symbol of each unit and its factor: how many base
// units (N) there are in one unit.
//
// The factors are exact: defined values, written as constant expressions that Go
// evaluates exactly before rounding once.
//
// Source: NIST SP 811, appendix B.
var definitions = [unitCount]definition{
	Newton:        {"N", 1},
	Dyne:          {"dyn", 1e-5},
	Millinewton:   {"mN", 1e-3},
	PoundForce:    {"lbf", 0.45359237 * 9.80665},
	KilogramForce: {"kgf", 9.80665},
	Kilonewton:    {"kN", 1e3},
	Meganewton:    {"MN", 1e6},
	Kip:           {"kip", 1000 * 0.45359237 * 9.80665},
	OunceForce:    {"ozf", 0.45359237 * 9.80665 / 16},
	TonneForce:    {"tf", 1000 * 9.80665},
}
