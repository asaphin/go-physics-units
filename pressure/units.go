// Package pressure defines the pressure units of the units package
// and their conversion factors to the base unit, the pascal.
//
// Use the units with units.NewPressure, from the root package
// github.com/asaphin/go-physics-units:
//
//	p := units.NewPressure(101.3, pressure.Kilopascal)
//
// ParseUnit turns a symbol such as "kPa" into a Unit, and String turns a
// Unit into its symbol.
package pressure

// Unit is a pressure unit: one of the constants below. String returns its
// symbol and ParseUnit finds the unit of a symbol.
//
// The zero Unit is the base unit, Pascal. The numbers behind the constants may
// change between versions: store and exchange units by their symbols.
type Unit uint8

// BaseUnit is the SI unit of pressure that all pressure factors are relative to.
const BaseUnit = Pascal

// Pressure units, with their symbols.
const (
	Pascal                           Unit = iota // Pa
	Hectopascal                                  // hPa
	Millibar                                     // mbar
	Kilopascal                                   // kPa
	Torr                                         // Torr: 1/760 standard atmosphere
	MillimeterOfMercury                          // mmHg: conventional: 1 mm of mercury of density 13595.1 kg/m³ under standard gravity
	InchOfMercury                                // inHg: conventional, as for MillimeterOfMercury
	PoundPerSquareInch                           // psi: pound-force per square inch
	Bar                                          // bar
	Atmosphere                                   // atm: standard atmosphere
	Megapascal                                   // MPa
	Gigapascal                                   // GPa
	KilogramForcePerSquareCentimeter             // kgf/cm²: the technical atmosphere, 1 kgf on 1 cm², often written kg/cm²
	MillimeterOfWater                            // mmH2O: conventional: 1 mm of water of density 1000 kg/m³ under standard gravity
	CentimeterOfWater                            // cmH2O: conventional, as for MillimeterOfWater
	InchOfWater                                  // inH2O: conventional, as for MillimeterOfWater
	KilopoundPerSquareInch                       // ksi: 1000 psi

	unitCount // the number of units; keep it last
)

// definitions holds the symbol of each unit and its factor: how many base
// units (Pa) there are in one unit.
//
// The factors are exact: defined values, written as constant expressions that Go
// evaluates exactly before rounding once.
//
// Source: NIST SP 811, appendix B.
var definitions = [unitCount]definition{
	Pascal:                           {"Pa", 1},
	Hectopascal:                      {"hPa", 100},
	Millibar:                         {"mbar", 100},
	Kilopascal:                       {"kPa", 1e3},
	Torr:                             {"Torr", 101325.0 / 760},
	MillimeterOfMercury:              {"mmHg", 13595.1 * 9.80665 / 1000},
	InchOfMercury:                    {"inHg", 13595.1 * 9.80665 * 0.0254},
	PoundPerSquareInch:               {"psi", 0.45359237 * 9.80665 / (0.0254 * 0.0254)},
	Bar:                              {"bar", 1e5},
	Atmosphere:                       {"atm", 101325},
	Megapascal:                       {"MPa", 1e6},
	Gigapascal:                       {"GPa", 1e9},
	KilogramForcePerSquareCentimeter: {"kgf/cm²", 9.80665 / (0.01 * 0.01)},
	MillimeterOfWater:                {"mmH2O", 1000 * 9.80665 / 1000},
	CentimeterOfWater:                {"cmH2O", 1000 * 9.80665 / 100},
	InchOfWater:                      {"inH2O", 1000 * 9.80665 * 0.0254},
	KilopoundPerSquareInch:           {"ksi", 1000 * 0.45359237 * 9.80665 / (0.0254 * 0.0254)},
}
