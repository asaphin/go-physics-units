// Package pace defines the pace (time per distance) units of the units package
// and their conversion factors to the base unit, the second per metre.
//
// Use the units with units.NewPace, from the root package
// github.com/asaphin/go-physics-units:
//
//	p := units.NewPace(5, pace.MinutePerKilometer)
//
// ParseUnit turns a symbol such as "min/km" into a Unit, and String turns a
// Unit into its symbol.
package pace

// Unit is a pace unit: one of the constants below. String returns its symbol
// and ParseUnit finds the unit of a symbol.
//
// The zero Unit is the base unit, SecondPerMeter. The numbers behind the
// constants may change between versions: store and exchange units by their
// symbols.
type Unit uint8

// BaseUnit is the SI unit of pace that all factors are relative to.
const BaseUnit = SecondPerMeter

// Pace units, with their symbols.
const (
	SecondPerMeter     Unit = iota // s/m
	SecondPerKilometer             // s/km
	MinutePerKilometer             // min/km
	SecondPerMile                  // s/mi
	MinutePerMile                  // min/mi
	SecondPer100Meters             // s/100m: swimming
	SecondPer500Meters             // s/500m: rowing

	unitCount // the number of units; keep it last
)

// definitions holds the symbol of each unit and its factor: how many base
// units (s/m) there are in one unit.
//
// Source: the minute (60 s), the kilometre and the international mile (1609.344
// m).
var definitions = [unitCount]definition{
	SecondPerMeter:     {"s/m", 1},
	SecondPerKilometer: {"s/km", 1.0 / 1000},
	MinutePerKilometer: {"min/km", 60.0 / 1000},
	SecondPerMile:      {"s/mi", 1 / 1609.344},
	MinutePerMile:      {"min/mi", 60 / 1609.344},
	SecondPer100Meters: {"s/100m", 1.0 / 100},
	SecondPer500Meters: {"s/500m", 1.0 / 500},
}
