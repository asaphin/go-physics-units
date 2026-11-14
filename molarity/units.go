// Package molarity defines the molar concentration (molarity) units of the
// units package and their conversion factors to the base unit, the mole per
// cubic metre.
//
// Use the units with units.NewMolarity, from the root package
// github.com/asaphin/go-physics-units:
//
//	c := units.NewMolarity(150, molarity.Millimolar)
//
// ParseUnit turns a symbol such as "mM" into a Unit, and String turns a
// Unit into its symbol.
package molarity

// Unit is a molar concentration unit: one of the constants below. String
// returns its symbol and ParseUnit finds the unit of a symbol.
//
// The zero Unit is the base unit, MolePerCubicMeter. The numbers behind the
// constants may change between versions: store and exchange units by their
// symbols.
type Unit uint8

// BaseUnit is the SI unit of molar concentration that all factors are relative
// to.
const BaseUnit = MolePerCubicMeter

// Molar concentration units, with their symbols.
const (
	MolePerCubicMeter Unit = iota // mol/m³
	MolePerLiter                  // mol/L
	MillimolePerLiter             // mmol/L
	MicromolePerLiter             // μmol/L: with U+03BC GREEK SMALL LETTER MU
	NanomolePerLiter              // nmol/L
	Molar                         // M: 1 mol/L
	Millimolar                    // mM
	Micromolar                    // μM: with U+03BC GREEK SMALL LETTER MU
	Nanomolar                     // nM
	Picomolar                     // pM

	unitCount // the number of units; keep it last
)

// definitions holds the symbol of each unit and its factor: how many base
// units (mol/m³) there are in one unit.
//
// Source: the SI Brochure (9th edition); M is the chemists' symbol for mol/L.
var definitions = [unitCount]definition{
	MolePerCubicMeter: {"mol/m³", 1},
	MolePerLiter:      {"mol/L", 1e3},
	MillimolePerLiter: {"mmol/L", 1},
	MicromolePerLiter: {"μmol/L", 1e-3},
	NanomolePerLiter:  {"nmol/L", 1e-6},
	Molar:             {"M", 1e3},
	Millimolar:        {"mM", 1},
	Micromolar:        {"μM", 1e-3},
	Nanomolar:         {"nM", 1e-6},
	Picomolar:         {"pM", 1e-9},
}
