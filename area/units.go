// Package area defines the area units of the units package and
// their conversion factors to the base unit, the square metre.
//
// Use the units with units.NewArea, from the root package
// github.com/asaphin/go-physics-units:
//
//	a := units.NewArea(10, area.SquareMeter)
//
// ParseUnit turns a symbol such as "m²" into a Unit, and String turns a
// Unit into its symbol.
package area

// Unit is an area unit: one of the constants below. String returns its
// symbol and ParseUnit finds the unit of a symbol.
//
// The zero Unit is the base unit, SquareMeter. The numbers behind the constants may
// change between versions: store and exchange units by their symbols.
type Unit uint8

// BaseUnit is the SI unit of area that all area factors are relative to.
const BaseUnit = SquareMeter

// Area units, with their symbols.
const (
	SquareMeter      Unit = iota // m²
	SquareMillimeter             // mm²
	SquareCentimeter             // cm²
	SquareInch                   // in²
	SquareFoot                   // ft²
	SquareKilometer              // km²
	SquareMile                   // mi²
	Hectare                      // ha
	Acre                         // ac: international acre, 43560 ft²
	SquareYard                   // yd²

	unitCount // the number of units; keep it last
)

// definitions holds the symbol of each unit and its factor: how many base
// units (m²) there are in one unit.
var definitions = [unitCount]definition{
	SquareMeter:      {"m²", 1},
	SquareMillimeter: {"mm²", 1e-6},
	SquareCentimeter: {"cm²", 1e-4},
	SquareInch:       {"in²", 0.00064516},
	SquareFoot:       {"ft²", 0.09290304},
	SquareKilometer:  {"km²", 1e6},
	SquareMile:       {"mi²", 2589988.110336},
	Hectare:          {"ha", 1e4},
	Acre:             {"ac", 43560 * 0.09290304},
	SquareYard:       {"yd²", 0.9144 * 0.9144},
}
