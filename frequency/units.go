// Package frequency defines the frequency units of the units package
// and their conversion factors to the base unit, the hertz.
//
// Use the units with units.NewFrequency, from the root package
// github.com/asaphin/go-physics-units:
//
//	f := units.NewFrequency(50, frequency.Hertz)
//
// ParseUnit turns a symbol such as "Hz" into a Unit, and String turns a
// Unit into its symbol.
package frequency

// Unit is a frequency unit: one of the constants below. String returns its
// symbol and ParseUnit finds the unit of a symbol.
//
// The zero Unit is the base unit, Hertz. The numbers behind the constants may
// change between versions: store and exchange units by their symbols.
type Unit uint8

// BaseUnit is the SI unit of frequency that all frequency factors are relative to.
const BaseUnit = Hertz

// Frequency units, with their symbols.
const (
	Hertz               Unit = iota // Hz
	RevolutionPerMinute             // rpm: one revolution per minute, as a rotational frequency
	Kilohertz                       // kHz
	Megahertz                       // MHz
	Gigahertz                       // GHz
	Terahertz                       // THz
	BeatPerMinute                   // bpm: one beat per minute: heart rate, tempo

	unitCount // the number of units; keep it last
)

// definitions holds the symbol of each unit and its factor: how many base
// units (Hz) there are in one unit.
//
// The factors are exact: defined values, written as constant expressions that Go
// evaluates exactly before rounding once.
//
// Source: the SI Brochure (9th edition).
var definitions = [unitCount]definition{
	Hertz:               {"Hz", 1},
	RevolutionPerMinute: {"rpm", 1.0 / 60},
	Kilohertz:           {"kHz", 1e3},
	Megahertz:           {"MHz", 1e6},
	Gigahertz:           {"GHz", 1e9},
	Terahertz:           {"THz", 1e12},
	BeatPerMinute:       {"bpm", 1.0 / 60},
}
