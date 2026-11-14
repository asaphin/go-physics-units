// Package specificenergy defines the specific energy (energy per mass) units of
// the units package and their conversion factors to the base unit, the joule
// per kilogram.
//
// Use the units with units.NewSpecificEnergy, from the root package
// github.com/asaphin/go-physics-units:
//
//	e := units.NewSpecificEnergy(250, specificenergy.WattHourPerKilogram)
//
// ParseUnit turns a symbol such as "Wh/kg" into a Unit, and String turns a
// Unit into its symbol.
package specificenergy

// Unit is a specific energy unit: one of the constants below. String returns
// its symbol and ParseUnit finds the unit of a symbol.
//
// The zero Unit is the base unit, JoulePerKilogram. The numbers behind the
// constants may change between versions: store and exchange units by their
// symbols.
type Unit uint8

// BaseUnit is the SI unit of specific energy that all factors are relative to.
const BaseUnit = JoulePerKilogram

// Specific energy units, with their symbols.
const (
	JoulePerKilogram        Unit = iota // J/kg
	KilojoulePerKilogram                // kJ/kg
	MegajoulePerKilogram                // MJ/kg
	WattHourPerKilogram                 // Wh/kg
	KilowattHourPerKilogram             // kWh/kg
	KilocaloriePerKilogram              // kcal/kg: thermochemical
	BTUPerPound                         // Btu/lb: International Table, 2326 J/kg exactly

	unitCount // the number of units; keep it last
)

// definitions holds the symbol of each unit and its factor: how many base
// units (J/kg) there are in one unit.
//
// Source: NIST SP 811, appendix B.
var definitions = [unitCount]definition{
	JoulePerKilogram:        {"J/kg", 1},
	KilojoulePerKilogram:    {"kJ/kg", 1e3},
	MegajoulePerKilogram:    {"MJ/kg", 1e6},
	WattHourPerKilogram:     {"Wh/kg", 3600},
	KilowattHourPerKilogram: {"kWh/kg", 3.6e6},
	KilocaloriePerKilogram:  {"kcal/kg", 4184},
	BTUPerPound:             {"Btu/lb", 2326},
}
