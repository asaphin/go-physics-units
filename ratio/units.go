// Package ratio defines the dimensionless ratio units of the units package and
// their conversion factors to the base unit, the one, written 1.
package ratio

// Unit is a ratio unit: one of the constants below. String returns its symbol
// and ParseUnit finds the unit of a symbol.
//
// The zero Unit is the base unit, One. The numbers behind the constants may
// change between versions: store and exchange units by their symbols.
type Unit uint8

// BaseUnit is the SI unit of ratio that all factors are relative to.
const BaseUnit = One

// Ratio units, with their symbols.
const (
	One              Unit = iota // 1: a plain fraction: 0.25 1 is 25 %
	Percent                      // %
	PerMille                     // ‰: with U+2030 PER MILLE SIGN
	BasisPoint                   // bp: one hundredth of a percent
	PartsPerMillion              // ppm
	PartsPerBillion              // ppb: the short scale: 1e-9
	PartsPerTrillion             // ppt: the short scale: 1e-12

	unitCount // the number of units; keep it last
)

// definitions holds the symbol of each unit and its factor: how many base
// units (1) there are in one unit.
//
// Source: ISO 80000-1; ppm, ppb and ppt on the short scale.
var definitions = [unitCount]definition{
	One:              {"1", 1},
	Percent:          {"%", 1e-2},
	PerMille:         {"‰", 1e-3},
	BasisPoint:       {"bp", 1e-4},
	PartsPerMillion:  {"ppm", 1e-6},
	PartsPerBillion:  {"ppb", 1e-9},
	PartsPerTrillion: {"ppt", 1e-12},
}
