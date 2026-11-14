// Package time defines the time (duration) units of the units
// package and their conversion factors to the base unit, the second.
//
// It shares its name with the standard library package; import one of them
// under an alias when both are needed.
package time

// Unit is a time unit: one of the constants below. String returns its
// symbol and ParseUnit finds the unit of a symbol.
//
// The zero Unit is the base unit, Second. The numbers behind the constants may
// change between versions: store and exchange units by their symbols.
type Unit uint8

// BaseUnit is the SI unit of time that all time factors are relative to.
const BaseUnit = Second

// Time units, with their symbols.
const (
	Second      Unit = iota // s
	PlanckTime              // tp
	Nanosecond              // ns
	Microsecond             // μs, with U+03BC GREEK SMALL LETTER MU; "µs" (U+00B5 MICRO SIGN) is accepted too
	Millisecond             // ms
	Minute                  // min
	Hour                    // h
	Day                     // d
	Week                    // wk
	Year                    // yr: Julian year of 365.25 days, as used in astronomy

	unitCount // the number of units; keep it last
)

// PlankTime is a misspelled alias of PlanckTime.
//
// Deprecated: Use PlanckTime.
//
//go:fix inline
const PlankTime = PlanckTime

// definitions holds the symbol of each unit and its factor: how many base
// units (s) there are in one unit.
var definitions = [unitCount]definition{
	Second:      {"s", 1},
	PlanckTime:  {"tp", 5.391247e-44}, // CODATA 2018
	Nanosecond:  {"ns", 1e-9},
	Microsecond: {"μs", 1e-6},
	Millisecond: {"ms", 0.001},
	Minute:      {"min", 60},
	Hour:        {"h", 3600},
	Day:         {"d", 86400},
	Week:        {"wk", 604800},
	Year:        {"yr", 31557600},
}
