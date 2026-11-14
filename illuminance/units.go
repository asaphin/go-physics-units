// Package illuminance defines the illuminance units of the units package and
// their conversion factors to the base unit, the lux.
//
// Use the units with units.NewIlluminance, from the root package
// github.com/asaphin/go-physics-units:
//
//	e := units.NewIlluminance(500, illuminance.Lux)
//
// ParseUnit turns a symbol such as "lx" into a Unit, and String turns a
// Unit into its symbol.
package illuminance

// Unit is an illuminance unit: one of the constants below. String returns its
// symbol and ParseUnit finds the unit of a symbol.
//
// The zero Unit is the base unit, Lux. The numbers behind the constants may
// change between versions: store and exchange units by their symbols.
type Unit uint8

// BaseUnit is the SI unit of illuminance that all factors are relative to.
const BaseUnit = Lux

// Illuminance units, with their symbols.
const (
	Lux        Unit = iota // lx
	Kilolux                // klx
	FootCandle             // fc: one lumen per square foot
	Phot                   // ph

	unitCount // the number of units; keep it last
)

// definitions holds the symbol of each unit and its factor: how many base
// units (lx) there are in one unit.
//
// Source: the SI Brochure (9th edition); the foot-candle and the phot from NIST
// SP 811.
var definitions = [unitCount]definition{
	Lux:        {"lx", 1},
	Kilolux:    {"klx", 1e3},
	FootCandle: {"fc", 1 / (0.3048 * 0.3048)},
	Phot:       {"ph", 1e4},
}
