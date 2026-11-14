// Package acceleration defines the acceleration units of the units package
// and their conversion factors to the base unit, the metre per second squared.
package acceleration

// Unit is an acceleration unit: one of the constants below. String returns its
// symbol and ParseUnit finds the unit of a symbol.
//
// The zero Unit is the base unit, MeterPerSecondSquared. The numbers behind the constants may
// change between versions: store and exchange units by their symbols.
type Unit uint8

// BaseUnit is the SI unit of acceleration that all acceleration factors are relative to.
const BaseUnit = MeterPerSecondSquared

// Acceleration units, with their symbols.
const (
	MeterPerSecondSquared      Unit = iota // m/s²
	CentimeterPerSecondSquared             // cm/s²
	FootPerSecondSquared                   // ft/s²
	StandardGravity                        // gₙ: standard acceleration of gravity, gₙ with U+2099 LATIN SUBSCRIPT SMALL LETTER N; exact by definition (3rd CGPM, 1901)

	unitCount // the number of units; keep it last
)

// definitions holds the symbol of each unit and its factor: how many base
// units (m/s²) there are in one unit.
//
// The factors are exact: defined values, written as constant expressions that Go
// evaluates exactly before rounding once.
//
// Source: NIST SP 811, appendix B.
var definitions = [unitCount]definition{
	MeterPerSecondSquared:      {"m/s²", 1},
	CentimeterPerSecondSquared: {"cm/s²", 0.01},
	FootPerSecondSquared:       {"ft/s²", 0.3048},
	StandardGravity:            {"gₙ", 9.80665},
}
