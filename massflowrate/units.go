// Package massflowrate defines the mass flow rate units of the units package
// and their conversion factors to the base unit, the kilogram per second.
package massflowrate

// Unit is a mass flow rate unit: one of the constants below. String returns its
// symbol and ParseUnit finds the unit of a symbol.
//
// The zero Unit is the base unit, KilogramPerSecond. The numbers behind the
// constants may change between versions: store and exchange units by their
// symbols.
type Unit uint8

// BaseUnit is the SI unit of mass flow rate that all factors are relative to.
const BaseUnit = KilogramPerSecond

// Mass flow rate units, with their symbols.
const (
	KilogramPerSecond Unit = iota // kg/s
	GramPerSecond                 // g/s
	KilogramPerMinute             // kg/min
	KilogramPerHour               // kg/h
	TonnePerHour                  // t/h
	PoundPerSecond                // lb/s
	PoundPerMinute                // lb/min
	PoundPerHour                  // lb/h

	unitCount // the number of units; keep it last
)

// definitions holds the symbol of each unit and its factor: how many base
// units (kg/s) there are in one unit.
//
// Source: NIST SP 811, appendix B.
var definitions = [unitCount]definition{
	KilogramPerSecond: {"kg/s", 1},
	GramPerSecond:     {"g/s", 1e-3},
	KilogramPerMinute: {"kg/min", 1.0 / 60},
	KilogramPerHour:   {"kg/h", 1.0 / 3600},
	TonnePerHour:      {"t/h", 1000.0 / 3600},
	PoundPerSecond:    {"lb/s", 0.45359237},
	PoundPerMinute:    {"lb/min", 0.45359237 / 60},
	PoundPerHour:      {"lb/h", 0.45359237 / 3600},
}
