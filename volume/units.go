// Package volume defines the volume units of the units package
// and their conversion factors to the base unit, the cubic metre.
package volume

// Unit is a volume unit: one of the constants below. String returns its
// symbol and ParseUnit finds the unit of a symbol.
//
// The zero Unit is the base unit, CubicMeter. The numbers behind the constants may
// change between versions: store and exchange units by their symbols.
type Unit uint8

// BaseUnit is the SI unit of volume that all volume factors are relative to.
const BaseUnit = CubicMeter

// Volume units, with their symbols.
const (
	CubicMeter      Unit = iota // m³
	CubicMillimeter             // mm³
	Microliter                  // μL, with U+03BC GREEK SMALL LETTER MU; "µL" (U+00B5 MICRO SIGN) is accepted too
	CubicCentimeter             // cm³
	Milliliter                  // mL
	CubicInch                   // in³
	Liter                       // L
	Pint                        // pt: US liquid pint, 1/8 US gallon
	Quart                       // qt: US liquid quart, 1/4 US gallon
	Gallon                      // gal: US liquid gallon, 231 in³
	CubicFoot                   // ft³
	Barrel                      // bbl: oil barrel, 42 US gallons
	CubicYard                   // yd³
	CubicKilometer              // km³

	unitCount // the number of units; keep it last
)

// definitions holds the symbol of each unit and its factor: how many base
// units (m³) there are in one unit.
//
// The factors are exact: defined values, written as constant expressions that Go
// evaluates exactly before rounding once.
//
// Source: NIST SP 811, appendix B; the inch, foot and yard are the
// international ones (exactly 0.0254, 0.3048 and 0.9144 m).
var definitions = [unitCount]definition{
	CubicMeter:      {"m³", 1},
	CubicMillimeter: {"mm³", 1e-9},
	Microliter:      {"μL", 1e-9},
	CubicCentimeter: {"cm³", 1e-6},
	Milliliter:      {"mL", 1e-6},
	CubicInch:       {"in³", 0.0254 * 0.0254 * 0.0254},
	Liter:           {"L", 1e-3},
	Pint:            {"pt", 231 * 0.0254 * 0.0254 * 0.0254 / 8},
	Quart:           {"qt", 231 * 0.0254 * 0.0254 * 0.0254 / 4},
	Gallon:          {"gal", 231 * 0.0254 * 0.0254 * 0.0254},
	CubicFoot:       {"ft³", 0.3048 * 0.3048 * 0.3048},
	Barrel:          {"bbl", 42 * 231 * 0.0254 * 0.0254 * 0.0254},
	CubicYard:       {"yd³", 0.9144 * 0.9144 * 0.9144},
	CubicKilometer:  {"km³", 1e9},
}
