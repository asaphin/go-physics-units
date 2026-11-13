# Changelog

## Unreleased

This release breaks the API in a few places; the next section lists what to change. Everything else is fixes and additions.

### Upgrading from v0.0.4

Compile-time breaks:

- **Go 1.27** is the minimum version.
- **`ConversionRates()`** is gone from the unit packages. For the multiplier between two units use `Rate(from, to)`, for example `distance.Rate(distance.Mile, distance.Kilometer)`; for the whole table use `ConversionFactors()`.
- **`Distance.DivideByTime`** returns `(Velocity, error)`. Handle the error, which is `ErrZeroDivision` when the time is zero; the method used to panic.

Changes in behaviour that compile unchanged:

- **`mass.Quintal`** is `"q"`, not `"Qts"`. Stored strings such as `"3 Qts"` no longer parse; replace the symbol.
- **`mass.Microgram`** is `"μg"` with U+03BC GREEK SMALL LETTER MU, not U+00B5 MICRO SIGN. Both forms are still accepted as input, but values report `"μg"`, so compare units with the constant rather than a literal.
- **Custom measures:** `Type()` of a value from `NewMeasurementCreator` returns the measure name you gave, not `custom`. Code that checked for `MeasureCustom` should check for that name. `NewBaseMeasurement` values are still `custom`.
- **Stricter input.** `NewMeasurementCreator` and `NewBaseMeasurement` reject empty, zero, negative, NaN and infinite factors, empty unit symbols, and, for the creator, empty or predefined measure names. `ParseString` rejects NaN and infinite values.
- **`AddMeasurement` and `SubMeasurement`** return an error when the units can't be converted into each other. They used to count the other operand as zero.
- **New measures** reserve their names: a custom measure named `volume`, `force`, `energy`, `power`, `pressure`, `frequency`, `angle`, `acceleration` or `data size` is now rejected by `NewMeasurementCreator`; rename it. Symbols such as `L`, `N`, `J`, `W`, `Pa` and `Hz` now parse as those measures.
- **Error messages** changed, for example to `unknown distance unit kg`. Instead of matching messages, test for the new sentinel errors with `errors.Is`, such as `errors.Is(err, units.ErrUnknownUnit)`.

### Fixed

- Converting between metres and millimetres no longer gives a wrong result. The internal rate table keyed pairs by concatenating the two symbols, and `"m"+"mm"` equals `"mm"+"m"`, so one of the two directions was wrong in every process.
- Using measurements from several goroutines no longer causes a data race. The conversion tables were built lazily on first use, and are now built once at startup.
- `ToDistance`, `ToMass`, `ToTime` and `ToVelocity` now accept a value that already is that measure, and return an error instead of panicking on `nil`.
- `NewMeasurementCreator[T]` returns an error instead of panicking later when `T` is a predefined measure interface such as `units.Distance`.
- On error, `MeasurementCreator.New`, `MeasurementCreator.ParseString` and `ConvertToMeasurement` return a nil interface, not an interface that wraps a nil pointer.
- `Type()` of a custom measurement returns the measure type given to `NewMeasurementCreator`, not `custom`.
- A custom unit that shares a symbol with a predefined unit (for example `"m"`) is no longer treated as that predefined measure.
- `AddMeasurement` and `SubMeasurement` return an error for measurements whose units can't be converted into each other. Previously the other operand was silently counted as zero.
- More precise conversion factors: km/h, knot, ounce, pound, light-year, parsec and Planck time. For example, 100 km / 1 h now gives exactly 100 km/h.

### Changed

- The micro prefix is accepted in both forms, `μ` (U+03BC) and `µ` (U+00B5), and stored as `μ`. `mass.Microgram` is now `"μg"` (U+03BC), consistent with `distance.Micrometer` and `time.Microsecond`.
- `mass.Quintal` is now `"q"`, the standard symbol. It was `"Qts"`.
- `ParseString` and `MeasurementCreator.ParseString` reject NaN and infinite values.
- `NewMeasurementCreator` and `NewBaseMeasurement` reject:
  - conversion factors that are zero, negative, NaN or infinite;
  - empty unit symbols;
  - measure types that are empty or the same as a predefined one.
- `NewBaseMeasurement` always returns a measurement of type `custom`.
- The names of the new measures (`volume`, `acceleration`, `force`, `energy`, `power`, `pressure`, `frequency`, `angle` and `data size`) are now reserved, so `NewMeasurementCreator` rejects them for custom measures. Their predefined symbols, such as `L` and `N`, now belong to those measures in `ParseString` and `DetectMeasureType`.
- Error messages name the measure, for example `unknown distance unit kg`.
- Faster, with fewer allocations:
  - `ParseString` and `NewMeasurement` are three to four times as fast and allocate once instead of three times. One index takes a symbol straight to its unit and table, and the parser no longer builds a slice of fields.
  - `String()` is about twice as fast and allocates once.
  - Errors are cheaper. Creating a value with an unknown unit is about three times as fast and allocates once instead of four times. Parsing a malformed string is about twice as fast and allocates twice instead of four times.
  - `NewBaseMeasurement` is two to four times as fast, depending on the number of units.
  - Most other operations allocate once instead of twice.
  - `StringToDistance` and the other `StringToX` functions allocate once instead of twice. They parse and construct directly instead of going through `ParseString`, so a unit of another measure is reported as unknown for that measure (`unknown distance unit kg`) instead of `not a distance measure`.
- `Distance.DivideByTime` returns `(Velocity, error)` and fails when the time is zero, as `Div` does. It used to panic.
- The minimum Go version is now 1.27.

### Removed

- `ConversionRates()` in the unit packages. It returned `immutable.Float64Map`, a type from an internal package that callers couldn't name, keyed by concatenated symbols. Use `Rate(from, to)` instead.

### Added

- `Parser`, a configurable parser built with `NewParser` and options: `WithAliases` (alternative unit names such as `metres`), `CaseInsensitive` (exact matches win, and an ambiguous match such as `KN` fails with the new `ErrAmbiguousUnit`), `AllowNoSpace` (`10km`) and `WithCreator` (the units of a custom measure). A `Parser` is immutable and safe for concurrent use, and the zero `Parser` behaves like `ParseString`. `Parser.Unit` resolves an alias to its symbol. Invalid options fail with the new `ErrInvalidOption`.
- Nine measures: volume, acceleration, force, energy, power, pressure, frequency, angle and data size, each with its unit package (`volume`, `acceleration`, …, `angle`, `datasize`), `NewX`, `ToX`, `StringToX` and `MeasureX`. The factors are exact definitions from NIST SP 811 and the SI Brochure. They come with cross-measure operations that link them to each other and to the existing measures: `Distance.Cube`, `Area.MultiplyByDistance`, `Volume.DivideByArea` and `DivideByDistance`, `Velocity.DivideByTime`, `Acceleration.MultiplyByTime` and `MultiplyByMass`, `Mass.MultiplyByAcceleration`, `Force.DivideByMass`, `DivideByAcceleration`, `DivideByArea` and `MultiplyByDistance`, `Energy.DivideByTime` and `DivideByDistance`, `Power.MultiplyByTime`, `Pressure.MultiplyByArea` and `Frequency.Period`.
- Runnable examples on pkg.go.dev for every measure, the cross-measure operations, extending a predefined measure, `NewBaseMeasurement`, formatting and `distance.Rate`.
- `Time.ToDuration`, the reverse of `ParseDuration`: it converts a `Time` into a `time.Duration`, rounded to the nearest nanosecond, and fails with `ErrOutOfRange` for NaN or beyond about ±292 years.
- `Compare`, `Equal` and `EqualWithinRel` on every measure, from the new `Comparer[T]` interface, and `CompareMeasurement` on every measurement, including custom ones. `Compare` follows `cmp.Compare`, so `slices.SortFunc(ds, units.Distance.Compare)` sorts distances in any mix of units. Values in different units are compared in the base unit, so the result doesn't depend on the order of the operands. `EqualWithinRel` allows a relative tolerance, because conversions round: 12 in and 1 ft differ in the last bit.
- Sentinel errors for `errors.Is`: `ErrUnknownUnit`, `ErrIncompatibleMeasures`, `ErrZeroDivision`, `ErrInvalidFormat`, `ErrInvalidFactors`, `ErrInvalidMeasure` and `ErrOutOfRange`. Every error the package returns matches exactly one of them and keeps its detailed message. They are constants of the new type `Error`, so no importer can reassign them, and `errors.AsType[units.Error](err)` finds the kind. Unknown-unit errors are `*UnitError` values, which carry the measure and the unit, and an unparsable value wraps its `*strconv.NumError`.
- Measurements implement `fmt.Formatter`. `%f`, `%e`, `%g` and the other float verbs format the value with the given flags, width and precision and append the unit: `fmt.Sprintf("%.2f", d)` gives `5.00 km`. `%v` and `%s` print `String()` with width, `%q` quotes it, and `%#v` prints `units.Measurement{Value:5, Unit:"km", Type:"distance"}`. Previously, verbs other than `%v` and `%s` printed the internal struct, including the whole conversion table.
- The area measure: the `area` package (mm², cm², in², ft², m², km², mi²), `Area`, `NewArea`, `ToArea`, `StringToArea` and `MeasureArea`, with the cross-measure operations `Distance.Square` and `Area.DivideByDistance`.
- `Rate(from, to string) (float64, bool)` in each unit package: the multiplier that converts a value from one unit to another, for example `distance.Rate(distance.Mile, distance.Kilometer)` is 1.609344. Both forms of the micro sign are accepted.
- `distance.PlanckLength` and `time.PlanckTime`. The misspelled `PlankLength` and `PlankTime` remain as deprecated aliases.
