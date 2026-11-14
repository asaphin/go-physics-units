// Package doserate defines the equivalent dose rate units of the units package
// and their conversion factors to the base unit, the sievert per second.
package doserate

// Unit is a dose rate unit: one of the constants below. String returns its
// symbol and ParseUnit finds the unit of a symbol.
//
// The zero Unit is the base unit, SievertPerSecond. The numbers behind the
// constants may change between versions: store and exchange units by their
// symbols.
type Unit uint8

// BaseUnit is the SI unit of dose rate that all factors are relative to.
const BaseUnit = SievertPerSecond

// Dose rate units, with their symbols.
const (
	SievertPerSecond    Unit = iota // Sv/s
	SievertPerHour                  // Sv/h
	MillisievertPerHour             // mSv/h
	MicrosievertPerHour             // μSv/h: with U+03BC GREEK SMALL LETTER MU
	NanosievertPerHour              // nSv/h
	MillisievertPerYear             // mSv/yr: per Julian year, as time.Year
	MilliremPerHour                 // mrem/h

	unitCount // the number of units; keep it last
)

// definitions holds the symbol of each unit and its factor: how many base
// units (Sv/s) there are in one unit.
//
// Source: the SI Brochure (9th edition); the rem from NIST SP 811.
var definitions = [unitCount]definition{
	SievertPerSecond:    {"Sv/s", 1},
	SievertPerHour:      {"Sv/h", 1.0 / 3600},
	MillisievertPerHour: {"mSv/h", 1e-3 / 3600},
	MicrosievertPerHour: {"μSv/h", 1e-6 / 3600},
	NanosievertPerHour:  {"nSv/h", 1e-9 / 3600},
	MillisievertPerYear: {"mSv/yr", 1e-3 / 31557600},
	MilliremPerHour:     {"mrem/h", 1e-5 / 3600},
}
