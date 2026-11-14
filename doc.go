// Package units provides physical measurements (distance, area, volume, mass,
// time, velocity, acceleration, force, energy, power, pressure, frequency,
// angle, angular velocity, data size, density, data rate, temperature,
// electric current, voltage, electric charge and resistance) with unit
// conversion, arithmetic, comparison, parsing and text encoding, plus
// user-defined measures such as currencies.
//
// # Values and units
//
// Each measure is a small value type, such as Distance, holding a float64 and
// a unit. Units are typed constants of the per-measure subpackages, such as
// distance.Kilometer or mass.Pound, so a unit of the wrong measure doesn't
// compile:
//
//	d := units.NewDistance(5, distance.Kilometer)
//	nmi := d.ConvertTo(distance.NauticalMile)
//	total := d.Add(units.NewDistance(300, distance.Meter)) // 5.3 km
//
// The typed operations can't fail, except division by zero, and don't
// allocate. Values are immutable: every method returns a new value. The zero
// value of each measure is 0 in its base unit (0 m for Distance).
//
// Values of the same measure can be added, subtracted and compared in any
// unit; arithmetic keeps the receiver's unit. Operations that combine
// measures, such as Distance.DivideByTime, return results in the SI base unit
// of the result.
//
// # Parsing and generic code
//
// StringToDistance and its siblings parse "<value> <symbol>" strings into a
// typed value. ParseString detects the measure from the symbol and returns a
// Measurement, the interface that every value implements, holding the
// measure's concrete type:
//
//	m, err := units.ParseString("10 km")
//	d, ok := m.(units.Distance)
//
// A unit symbol from configuration converts to a typed unit with the unit
// package's ParseUnit. NewParser builds a parser with aliases, case-insensitive
// symbols and user-defined measures.
//
// Unit symbols are case-sensitive. U+00B5 MICRO SIGN and U+03BC GREEK SMALL
// LETTER MU are interchangeable in symbols ("µm" and "μm"); values carry the
// U+03BC form.
//
// # User-defined measures
//
// A user-defined measure is a marker type whose Units method returns a
// UnitSet. Quantity gives it the typed API of the predefined measures; see
// Measure.
//
// # Formatting and encoding
//
// Every value implements fmt.Formatter:
//
//   - %v and %s print String(), "<value> <symbol>", honoring width and the '-'
//     flag: "5 km".
//   - %b, %e, %E, %f, %F, %g, %G, %x and %X format the value with the given
//     flags, width and precision, then append the symbol: %.2f gives "5.00 km".
//   - %q prints String() quoted: "\"5 km\"".
//   - %#v prints a Go-like form: units.Distance{Value:5, Unit:"km"}.
//   - Other verbs print fmt's usual error form: %!d(units.Distance=5 km).
//
// Values implement encoding.TextMarshaler and encoding.TextUnmarshaler in the
// String form, so they encode in JSON as strings such as "5 km".
//
// # Errors
//
// Every error the package returns matches exactly one of the sentinel errors,
// such as ErrUnknownUnit; test for them with errors.Is. Unknown-unit errors are
// *UnitError values, which carry the measure and the unit.
//
// # Concurrency
//
// Values never change after they are created, and the package's unit tables
// never change after initialisation. Values and functions are therefore safe
// for concurrent use by multiple goroutines without additional
// synchronisation.
package units
