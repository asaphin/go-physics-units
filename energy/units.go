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
	Joule                Unit = iota // J
	Electronvolt                     // eV: exact since the 2019 redefinition of the SI
	Erg                              // erg
	FootPound                        // ft·lbf: foot-pound force, with U+00B7 MIDDLE DOT
	Calorie                          // cal: thermochemical calorie
	Kilojoule                        // kJ
	BTU                              // Btu: International Table British thermal unit
	Kilocalorie                      // kcal: thermochemical kilocalorie, the dietary Calorie
	WattHour                         // Wh
	Megajoule                        // MJ
	KilowattHour                     // kWh
	Gigajoule                        // GJ
	MegawattHour                     // MWh
	GramOfTNT                        // gTNT: gram of TNT, 1 kcal by convention
	TonOfTNT                         // tTNT: ton of TNT, 1 Gcal by convention
	KilotonOfTNT                     // ktTNT: kiloton of TNT
	MegatonOfTNT                     // MtTNT: megaton of TNT
	Kiloelectronvolt                 // keV
	Megaelectronvolt                 // MeV
	Gigaelectronvolt                 // GeV
	Teraelectronvolt                 // TeV
	TonneOfOilEquivalent             // toe: tonne of oil equivalent

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
	Joule:                {"J", 1},
	Electronvolt:         {"eV", 1.602176634e-19},
	Erg:                  {"erg", 1e-7},
	FootPound:            {"ft·lbf", 0.3048 * 0.45359237 * 9.80665},
	Calorie:              {"cal", 4.184},
	Kilojoule:            {"kJ", 1e3},
	BTU:                  {"Btu", 1055.05585262},
	Kilocalorie:          {"kcal", 4184},
	WattHour:             {"Wh", 3600},
	Megajoule:            {"MJ", 1e6},
	KilowattHour:         {"kWh", 3.6e6},
	Gigajoule:            {"GJ", 1e9},
	MegawattHour:         {"MWh", 3.6e9},
	GramOfTNT:            {"gTNT", 4184},
	TonOfTNT:             {"tTNT", 4.184e9}, // TNT equivalent: a convention, not a measurement
	KilotonOfTNT:         {"ktTNT", 4.184e12},
	MegatonOfTNT:         {"MtTNT", 4.184e15},
	Kiloelectronvolt:     {"keV", 1.602176634e-16},
	Megaelectronvolt:     {"MeV", 1.602176634e-13},
	Gigaelectronvolt:     {"GeV", 1.602176634e-10},
	Teraelectronvolt:     {"TeV", 1.602176634e-7},
	TonneOfOilEquivalent: {"toe", 41.868e9}, // the IEA and OECD convention
}
