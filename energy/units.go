// Package energy defines the energy units of the units package
// and their conversion factors to the base unit, the joule.
package energy

// Unit is an energy unit: one of the constants below. String returns its
// symbol and ParseUnit finds the unit of a symbol.
//
// The zero Unit is the base unit, Joule. The numbers behind the constants may
// change between versions: store and exchange units by their symbols.
type Unit uint8

// BaseUnit is the SI unit of energy that all energy factors are relative to.
const BaseUnit = Joule

// Energy units, with their symbols.
const (
	Joule        Unit = iota // J
	Electronvolt             // eV: exact since the 2019 redefinition of the SI
	Erg                      // erg
	FootPound                // ft·lbf: foot-pound force, with U+00B7 MIDDLE DOT
	Calorie                  // cal: thermochemical calorie
	Kilojoule                // kJ
	BTU                      // Btu: International Table British thermal unit
	Kilocalorie              // kcal: thermochemical kilocalorie, the dietary Calorie
	WattHour                 // Wh
	Megajoule                // MJ
	KilowattHour             // kWh
	Gigajoule                // GJ
	MegawattHour             // MWh

	unitCount // the number of units; keep it last
)

// definitions holds the symbol of each unit and its factor: how many base
// units (J) there are in one unit.
//
// The factors are exact: defined values, written as constant expressions that Go
// evaluates exactly before rounding once.
//
// Source: NIST SP 811, appendix B; the electronvolt from the SI Brochure (9th
// edition).
var definitions = [unitCount]definition{
	Joule:        {"J", 1},
	Electronvolt: {"eV", 1.602176634e-19},
	Erg:          {"erg", 1e-7},
	FootPound:    {"ft·lbf", 0.3048 * 0.45359237 * 9.80665},
	Calorie:      {"cal", 4.184},
	Kilojoule:    {"kJ", 1e3},
	BTU:          {"Btu", 1055.05585262},
	Kilocalorie:  {"kcal", 4184},
	WattHour:     {"Wh", 3600},
	Megajoule:    {"MJ", 1e6},
	KilowattHour: {"kWh", 3.6e6},
	Gigajoule:    {"GJ", 1e9},
	MegawattHour: {"MWh", 3.6e9},
}
