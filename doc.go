// Package units provides immutable physical measurements (distance, area, mass,
// time and velocity) with unit conversion, arithmetic and string parsing, plus
// user-defined measures such as currencies.
//
// Unit symbols are defined in the per-measure subpackages, for example
// distance.Kilometer or mass.Pound:
//
//	d, err := units.NewDistance(5, distance.Kilometer)
//	m, err := units.ParseString("10 km")      // generic Measurement
//	d, err = units.ToDistance(m)              // narrowed to Distance
//	nmi, err := d.ConvertTo(distance.NauticalMile)
//
// Measurements of the same measure can be added, subtracted and compared in any
// unit; the result of arithmetic keeps the receiver's unit. Operations that combine measures, such as
// Distance.DivideByTime, return results in the SI base unit of the result.
//
// Unit symbols are case-sensitive. U+00B5 MICRO SIGN and U+03BC GREEK SMALL
// LETTER MU are interchangeable in symbols ("µm" and "μm"); values carry the
// U+03BC form.
//
// User-defined measures are created with NewMeasurementCreator.
//
// Every error the package returns matches exactly one of the sentinel errors,
// such as ErrUnknownUnit; test for them with errors.Is. Unknown-unit errors are
// *UnitError values, which carry the measure and the unit.
//
// Measurements never change after they are created, and the package's
// conversion tables never change after initialisation. Values and functions are
// therefore safe for concurrent use by multiple goroutines without additional
// synchronisation.
package units
