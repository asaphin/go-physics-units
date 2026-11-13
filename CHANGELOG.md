# Changelog

## Unreleased

### Fixed

- Converting between metres and millimetres no longer gives a wrong result. The internal rate table keyed pairs by concatenating the two symbols, and `"m"+"mm"` equals `"mm"+"m"`, so one of the two directions was wrong in every process.
- Using measurements from several goroutines no longer causes a data race. The conversion tables were built lazily on first use, and are now built once at startup.
- `ToDistance`, `ToArea`, `ToMass`, `ToTime` and `ToVelocity` now accept a value that already is that measure, and return an error instead of panicking on `nil`.
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
  - `ParseString` is about twice as fast, because `DetectMeasureType` uses an index instead of a linear scan.
  - `String()` is about twice as fast and allocates once.
  - Creating a value with an unknown unit is almost twice as fast.
  - `NewBaseMeasurement` is two to four times as fast, depending on the number of units.
  - Most other operations allocate once instead of twice.
- The minimum Go version is now 1.22.

### Removed

- `ConversionRates()` in the unit packages. It returned `immutable.Float64Map`, a type from an internal package that callers couldn't name, keyed by concatenated symbols. Use `Rate(from, to)` instead.

### Added

- `Rate(from, to string) (float64, bool)` in each unit package: the multiplier that converts a value from one unit to another, for example `distance.Rate(distance.Mile, distance.Kilometer)` is 1.609344. Both forms of the micro sign are accepted.
- `distance.PlanckLength` and `time.PlanckTime`. The misspelled `PlankLength` and `PlankTime` remain as deprecated aliases.
