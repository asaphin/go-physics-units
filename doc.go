// Package units provides physical measurements, with unit conversion,
// arithmetic, comparison, parsing and text encoding: 43 predefined measures,
// from distance, time and mass to torque, molar concentration, radiation dose
// and luminous flux, plus temperature and user-defined measures such as
// currencies.
//
//	dist := units.NewDistance(10, distance.Kilometer)
//	vel, err := dist.DivideByTime(units.NewTime(10, time.Minute))
//	fmt.Println(vel.ConvertTo(velocity.KilometerPerHour)) // 60 km/h
//
// # Concepts
//
//   - A measure is a kind of quantity, such as distance or mass. Each has a
//     value type, such as [Distance], and a package of units, such as
//     [github.com/asaphin/go-physics-units/distance].
//   - A unit is a typed constant of that package, such as distance.Kilometer.
//     A unit of another measure doesn't compile where a distance unit is
//     expected.
//   - A symbol is a unit's text form, such as "km". String returns it, and the
//     unit package's ParseUnit turns it back into the unit.
//   - The base unit of a measure is its coherent SI unit (m, kg, s, m/s, …),
//     and a unit's factor is how many base units there are in one of it.
//   - A value is a float64 and a unit. Values are immutable: every method
//     returns a new value. The zero value of a measure is 0 in its base unit.
//
// # Which function do I need?
//
//   - A value in a known unit: [NewDistance] and its siblings, one per measure.
//   - A string in a known measure, such as "10 km": [StringToDistance] and its
//     siblings.
//   - A string in any measure: [ParseString], then a type switch or
//     [ToDistance] and its siblings.
//   - Text that people type ("10KM", "5 metres"): a [Parser] from [NewParser].
//   - A unit symbol from configuration: the unit package's ParseUnit.
//   - Another unit: ConvertTo, or ConvertToBaseUnits for the SI base unit.
//   - Display: the fmt verbs, such as %.2f, or String.
//   - Storage and JSON: the values implement encoding.TextMarshaler and
//     encoding.TextUnmarshaler, in the form "12.5 km".
//   - Two values of one measure: Add, Sub, Compare, Equal, EqualWithinRel and
//     Ratio.
//   - Two measures: the methods that combine them, such as
//     [Distance.DivideByTime] (a velocity) or [Torque.MultiplyByAngularVelocity]
//     (a power). The README lists them all.
//   - A formula of one value: [Volume.CubeEdge], [Mass.RestEnergy],
//     [Velocity.LorentzFactor], [Radioactivity.DecayAfter] and others.
//   - A formula with several inputs, such as an escape velocity or the
//     pressure of an ideal gas: [github.com/asaphin/go-physics-units/physics].
//   - Physical constants: [github.com/asaphin/go-physics-units/constant].
//   - Units of your own, such as currencies: [NewUnitSet] and [Quantity].
//   - A time.Duration: [ParseDuration] and [Time.ToDuration].
//   - Temperatures, which have offsets: [Temperature] and
//     [TemperatureDifference].
//
// # Values and units
//
// Units are typed constants, so a unit of the wrong measure doesn't compile:
//
//	d := units.NewDistance(5, distance.Kilometer)
//	nmi := d.ConvertTo(distance.NauticalMile)
//	total := d.Add(units.NewDistance(300, distance.Meter)) // 5.3 km
//
// The typed operations don't allocate, and can't fail except for division by
// zero and the few calculations without a real result. Values of the same
// measure can be added, subtracted and compared in any unit; arithmetic keeps
// the receiver's unit.
//
// # Combining measures, calculators and formulas
//
// Methods that combine measures return their result in the SI base unit of the
// result: [Distance.DivideByTime] gives a [Velocity] in m/s, and
// [Torque.MultiplyByAngularVelocity] a [Power] in W. Convert the result for
// display. Other methods are formulas of one value, such as [Volume.CubeEdge],
// the edge of the cube with that volume, and the Ratio method of every measure
// divides two values of one measure.
//
// The physics package has the formulas with several inputs, an approximation
// or a measured constant, and the constant package the physical constants of
// the SI, CODATA and the IAU.
//
// # Parsing and generic code
//
// [ParseString] detects the measure from the symbol and returns a
// [Measurement], the interface that every value implements, holding the
// measure's concrete type:
//
//	m, err := units.ParseString("10 km")
//	d, ok := m.(units.Distance)
//
// ParseString is strict: exactly "<value> <symbol>", with a space and a
// case-sensitive predefined symbol. [NewParser] builds a [Parser] that accepts
// aliases, any case, no space, and user-defined measures.
//
// U+00B5 MICRO SIGN and U+03BC GREEK SMALL LETTER MU are interchangeable in
// symbols ("µm" and "μm"), as are the other look-alike characters that Unicode
// NFKC maps together; values carry the canonical form.
//
// # User-defined measures
//
// A user-defined measure is a marker type whose Units method returns a
// [UnitSet] built once with [NewUnitSet]. [Quantity] gives it the typed API of
// the predefined measures; see [Measure].
//
// # Temperature
//
// Temperature scales have offsets, so a [Temperature] is a point on a scale and
// a [TemperatureDifference] a change. Temperatures can't be below absolute zero
// and aren't Measurements: parse them with [StringToTemperature].
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
// such as [ErrUnknownUnit]; test for them with errors.Is. Each function's
// documentation names the sentinels it fails with. Unknown-unit errors are
// [*UnitError] values, which carry the measure and the unit.
//
// # Common mistakes
//
//   - Comparing converted values with ==. Conversions round, so 12 in and 1 ft
//     differ in the last bit: use Equal, or EqualWithinRel after a conversion.
//   - Storing units as numbers. The numbers behind the unit constants may change
//     between versions: store and exchange symbols (String, ParseUnit) or the
//     text form of values.
//   - Using the result of a combining method in its base unit for display: 10
//     km in 10 min is 16.666… m/s until you convert it to km/h.
//   - Confusing look-alike symbols. Symbols are case-sensitive: kN is a
//     kilonewton and kn a knot, mM millimolar and mm a millimetre, G a gauss
//     and g a gram. rpm is a frequency (cycles per minute), rev/min an angular
//     velocity.
//   - Converting a fuel economy (km/L, mpg) into a fuel consumption (L/100km)
//     with ConvertTo. They are reciprocal measures: use
//     [FuelEconomy.Consumption] and [FuelConsumption.Economy].
//   - Adding two temperatures. Add a [TemperatureDifference] to a [Temperature]
//     instead.
//   - Importing the standard library's time package as time. This module's time
//     package has the same name: import the standard one under another name,
//     such as builtintime "time".
//
// # Concurrency
//
// Values never change after they are created, and the package's unit tables
// never change after initialisation. Values and functions are therefore safe
// for concurrent use by multiple goroutines without additional
// synchronisation.
package units
