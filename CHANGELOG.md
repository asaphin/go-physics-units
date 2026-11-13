# Changelog

## Unreleased

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
- Error messages name the measure, for example `unknown distance unit kg`.
- Faster, with fewer allocations:
  - `ParseString` is about twice as fast and allocates once instead of three times. `DetectMeasureType` uses an index instead of a linear scan, and the parser no longer builds a slice of fields.
  - `String()` is about twice as fast and allocates once.
  - Creating a value with an unknown unit is almost twice as fast.
  - `NewBaseMeasurement` is two to four times as fast, depending on the number of units.
  - Most other operations allocate once instead of twice.
  - `StringToDistance` and the other `StringToX` functions allocate once instead of twice. They parse and construct directly instead of going through `ParseString`, so a unit of another measure is reported as unknown for that measure (`unknown distance unit kg`) instead of `not a distance measure`.
- `Distance.DivideByTime` returns `(Velocity, error)` and fails when the time is zero, as `Div` does. It used to panic.
- The minimum Go version is now 1.27.

### Removed

- `ConversionRates()` in the unit packages. It returned `immutable.Float64Map`, a type from an internal package that callers couldn't name, keyed by concatenated symbols. Use `Rate(from, to)` instead.

### Added

- Measurements implement `fmt.Formatter`. `%f`, `%e`, `%g` and the other float verbs format the value with the given flags, width and precision and append the unit: `fmt.Sprintf("%.2f", d)` gives `5.00 km`. `%v` and `%s` print `String()` with width, `%q` quotes it, and `%#v` prints `units.Measurement{Value:5, Unit:"km", Type:"distance"}`. Previously, verbs other than `%v` and `%s` printed the internal struct, including the whole conversion table.
- The area measure: the `area` package (mm², cm², in², ft², m², km², mi²), `Area`, `NewArea`, `ToArea`, `StringToArea` and `MeasureArea`, with the cross-measure operations `Distance.Square` and `Area.DivideByDistance`.
- `Rate(from, to string) (float64, bool)` in each unit package: the multiplier that converts a value from one unit to another, for example `distance.Rate(distance.Mile, distance.Kilometer)` is 1.609344. Both forms of the micro sign are accepted.
- `distance.PlanckLength` and `time.PlanckTime`. The misspelled `PlankLength` and `PlankTime` remain as deprecated aliases.
