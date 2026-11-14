# Changelog

## Unreleased

## v0.2.0 - 2026-11-14

This release makes the library more scientific: 22 measures, units for science and everyday use in the existing measures, operations between measures, calculators, and two packages, `physics` (formulas) and `constant` (physical constants). Code written for v0.1.0 builds and runs unchanged; the one changed value is `mass.Dalton`, now the CODATA 2022 value.

### Added

- **22 measures:** `Torque`, `FlowRate`, `MassFlowRate`, `Momentum`, `Pace`, `FuelEconomy`, `FuelConsumption`, `Amount` (of substance), `MolarMass`, `Molarity`, `Radioactivity`, `AbsorbedDose`, `EquivalentDose`, `DoseRate`, `Capacitance`, `Inductance`, `MagneticField`, `Illuminance`, `LuminousFlux`, `SpecificEnergy`, `EnergyDensity` and `Ratio` (%, ‰, bp, ppm, ppb, ppt), each with its unit package, the typed API and a `MeasureType`. The rad of absorbed dose is `rd`, because `rad` is the radian.
- **Units in the existing measures:** TNT equivalents (`gTNT`, `tTNT`, `ktTNT`, `MtTNT`, 4.184 GJ per ton by convention), `keV` to `TeV`, `mJ` and `toe` (energy); `ha`, `ac` and `yd²` (area); `st`, `ct`, `gr` and `ozt` (mass); `km/s`; `kpc`, `Mpc`, `Gpc`, `mil` and `ftm` (distance); `bpm`; `mmH2O`, `cmH2O`, `inH2O` and `ksi` (pressure); `kip`, `ozf` and `tf` (force); `TR` and `L☉` (power); `EB` and `EiB`; `cL` and `dL`.
- **66 operations between measures,** in the result's SI base unit, among them `Distance.DivideByVelocity`, `Energy.DivideByPower`, momentum and impulse, torque and engine power (`Torque.MultiplyByAngularVelocity`, `Power.DivideByAngularVelocity`), flow rates, pace, fuel economy and consumption, amounts of substance and molarity, dose rates, `Q = C·V` and RC and L/R time constants, illuminance, specific energy, energy density and wavelengths. The README's methods table lists them all.
- **Calculators:** `Volume.CubeEdge` and `SphereRadius` (the equivalent cube and sphere), `Area.SquareSide` and `CircleRadius`, `Distance.CircleArea`, `SphereArea` and `SphereVolume`, `Mass.KineticEnergy`, `Weight` and `RestEnergy`, `Energy.MassEquivalent` and `PhotonFrequency`, `Frequency.PhotonEnergy`, `Velocity.LorentzFactor` and `Pace`, `Pace.Velocity`, `Time.Frequency`, `FuelEconomy.Consumption`, `FuelConsumption.Economy`, `Amount.Entities`, `Radioactivity.DecayAfter`, `AbsorbedDose.Equivalent` and `Ratio.Fraction`.
- **`Ratio(other)`** on every measure and on `Quantity`: the ratio of two values, exact when they share a unit.
- **Package `physics`:** formulas that combine measures: falling bodies, projectiles, pendulums and hydrostatic pressure; gravitation, orbits and black holes; light, sound, Mach numbers and time dilation; ideal gases and black-body radiation; capacitor and inductor energy and LC resonance.
- **Package `constant`:** the SI defining constants, exact derived constants (ħ, R, F, σ), the CODATA 2022 values of G, mₑ, mₚ, u, α, ε₀ and μ₀, and IAU astronomical constants.
- **Documentation:** every doc comment says how its function fails, with links to the sentinels and base units; the package documentation is a task guide (which function for which job, and the common mistakes); every unit package shows how its units are used.
- **USECASES.md:** four more scenarios: a car's power and a road trip, a lab solution, the radiation of a long-haul flight, and orbits.

### Changed

- **`mass.Dalton`** is the CODATA 2022 value, 1.66053906892e-27 kg (it was the CODATA 2018 value, 1.6605390666e-27 kg).
- **`ErrOutOfRange`** also reports a calculation without a real result, such as `Area.SquareSide` of a negative area.

### Performance

- The typed operations are unchanged: 2–5 ns, no allocations. `NewUnitSet` takes longer for a set built from a predefined measure that gained units (distance has 24 instead of 19).

## v0.1.0 - 2026-11-14

This release redesigns the API around value types and typed units, as decided in [docs/rfc-value-types.md](docs/rfc-value-types.md). Every program that uses the library needs changes; [MIGRATION.md](MIGRATION.md) lists them with before-and-after code. Most are mechanical, and the compiler reports each one.

The entries below describe the changes since v0.0.4.

### Changed (breaking)

- **Value types.** `Distance`, `Mass`, `Time`, `Velocity` and the new measures are structs, not interfaces. Their methods have value receivers and return new values. The zero value is 0 in the base unit (0 m), not nil.
- **Typed units.** Each unit package declares `type Unit uint8`, and its unit constants are typed: `distance.Kilometer` is a `distance.Unit`, not the string `"km"`. A unit of another measure, or a string literal, no longer compiles. `String()` returns the symbol and `ParseUnit` parses one.
- **No errors where none can occur.** `NewX` and `ConvertTo` return just the value. `MustConvertTo` is deprecated: it is `ConvertTo`, and `go fix` rewrites calls to it.
- **`Measurement`** is implemented by every value, and `ParseString` and `NewMeasurement` return the measure's concrete type, so `m.(units.Distance)` works. Its `Unit()` method is now `Symbol()`, because the typed `Unit()` returns the typed unit. `CompareMeasurement` is new.
- **Custom measures:** `NewMeasurementCreator`, `MeasurementCreator` and `NewBaseMeasurement` are removed. A custom measure is a marker type whose `Units` method returns a `UnitSet`, and its values are `Quantity[M]` (see Added). Values of two unit sets combine only when they are the same set, even if the sets share a name.
- **Unit packages:** `ConversionRates()` is removed; `Rate(from, to Unit) float64` returns the multiplier between two units, and `ConversionFactors()` a copy of the factors.
- **`UnitConverter[T, U]`** takes the unit type as a second type parameter, and its `ConvertTo` no longer returns an error.
- **`Distance.DivideByTime`** returns `(Velocity, error)` and fails with `ErrZeroDivision` when the time is zero. It used to panic.
- **`mass.Quintal`** is `q`, the standard symbol, not `Qts`, and **`mass.Microgram`** is `μg` with U+03BC GREEK SMALL LETTER MU, like `distance.Micrometer` and `time.Microsecond`. Both micro signs are accepted as input.
- **Stricter input:** `ParseString` rejects NaN and infinite values. `NewUnitSet` rejects empty, zero, negative, NaN and infinite factors, empty symbols, symbols that differ only in the micro sign, and empty or predefined measure names (`distance`, …, `temperature`, `custom`).
- **Error messages** name the measure, for example `unknown distance unit kg`. Test for the new sentinel errors instead of matching messages.
- **`%#v`** prints `units.Distance{Value:5, Unit:"km"}`.
- **Go 1.27** is the minimum version.

### Added

- **`Quantity[M]`**, the typed API for custom measures: `NewUnitSet` and `MustNewUnitSet` build a `UnitSet`; a marker type implements `Measure` by returning it; `UnitOf` and `MustUnitOf` look up a `Unit[M]`; `NewQuantity`, `ParseQuantity` and `ToQuantity` create values, which have `ConvertTo`, `ConvertToBaseUnits`, `Add`, `Sub`, `Mul`, `Div`, `Compare`, `Equal`, `EqualWithinRel`, text encoding and the `Measurement` methods. A predefined measure is extended by copying its `ConversionFactors()` into a set of your own.
- **Eighteen measures:** area, volume, acceleration, force, energy, power, pressure, frequency, angle, angular velocity, data size, density, data rate, electric current, voltage, electric charge and resistance, each with its unit package (`area`, …, `angularvelocity`, `current`, `voltage`, `charge`, `resistance`), `NewX`, `ToX`, `StringToX` and `MeasureX`, and temperature (below). The factors are exact definitions from NIST SP 811, the SI Brochure and IEC 80000-13.
- **Cross-measure operations** that link the measures, in SI base units: `Distance.Square` and `Cube`, `Area.DivideByDistance` and `MultiplyByDistance`, `Volume.DivideByArea` and `DivideByDistance`, `Velocity.DivideByTime`, `Acceleration.MultiplyByTime` and `MultiplyByMass`, `Mass.MultiplyByAcceleration`, `DivideByVolume` and `DivideByDensity`, `Force.DivideByMass`, `DivideByAcceleration`, `DivideByArea` and `MultiplyByDistance`, `Energy.DivideByTime` and `DivideByDistance`, `Power.MultiplyByTime`, `Pressure.MultiplyByArea`, `Frequency.Period`, `Density.MultiplyByVolume`, `DataSize.DivideByTime` and `DivideByDataRate`, and `DataRate.MultiplyByTime`. For the electric measures: `Voltage.MultiplyByCurrent`, `DivideByCurrent` (Ohm's law), `DivideByResistance` and `MultiplyByCharge`, `Current.MultiplyByVoltage`, `MultiplyByResistance` and `MultiplyByTime`, `Resistance.MultiplyByCurrent`, `Charge.DivideByTime`, `DivideByCurrent` (how long a battery lasts) and `MultiplyByVoltage`, `Power.DivideByVoltage` and `DivideByCurrent`, and `Energy.DivideByCharge`. For rotation: `Angle.DivideByTime` and `AngularVelocity.MultiplyByTime`. The divisions fail with `ErrZeroDivision` for a zero divisor.
- **Temperature:** the `temperature` package (K, °C, °F, °R; `℃` and `℉` are accepted as input), with `Temperature`, a point on a scale, and `TemperatureDifference`, a change. A temperature can't be below absolute zero: `NewTemperature`, `StringToTemperature`, `Temperature.Add` and `SubDifference` fail with `ErrOutOfRange` instead. `Temperature.Sub` gives a difference, and `Kelvin` the absolute value for physical formulas. Each pair of scales converts with coefficients computed exactly and rounded once, so 100 °C is exactly 212 °F. Temperatures don't implement `Measurement`; `ParseString` rejects their units with an error that points to `StringToTemperature`.
- **Text and JSON encoding:** every value implements `encoding.TextMarshaler`, `encoding.TextAppender` and `encoding.TextUnmarshaler` in the `"<value> <unit>"` form, so values encode in JSON as strings such as `"12.5 km"`. Unmarshalling checks the unit; values that aren't finite can't be marshalled.
- **`Parser`**, a configurable parser built with `NewParser` and options: `WithAliases` (alternative unit names such as `metres`), `CaseInsensitive` (exact matches win, and an ambiguous match such as `KN` fails with `ErrAmbiguousUnit`), `AllowNoSpace` (`10km`) and `WithMeasure[M]` (the units of a custom measure). It is immutable and safe for concurrent use, and the zero `Parser` behaves like `ParseString`. `Parser.Unit` resolves an alias to its symbol.
- **Comparisons:** `Compare`, `Equal` and `EqualWithinRel` on every value, from the new `Comparer[T]` interface. `Compare` follows `cmp.Compare`, so `slices.SortFunc(ds, units.Distance.Compare)` sorts distances in any mix of units. Values in different units are compared in the base unit, so the result doesn't depend on the order of the operands. `EqualWithinRel` allows a relative tolerance, because conversions round: 12 in and 1 ft differ in the last bit.
- **Sentinel errors** for `errors.Is`: `ErrUnknownUnit`, `ErrIncompatibleMeasures`, `ErrZeroDivision`, `ErrInvalidFormat`, `ErrInvalidFactors`, `ErrInvalidMeasure`, `ErrOutOfRange`, `ErrAmbiguousUnit` and `ErrInvalidOption`. Every error the package returns matches exactly one of them and keeps its detailed message. They are constants of the type `Error`, so no importer can reassign them. Unknown-unit errors are `*UnitError` values, which carry the measure and the unit, and an unparsable value wraps its `*strconv.NumError`.
- **`fmt.Formatter`** on every value: `%f`, `%e`, `%g` and the other float verbs format the value with the given flags, width and precision and append the unit, so `fmt.Sprintf("%.2f", d)` gives `5.00 km`. `%v` and `%s` print `String()` with width, and `%q` quotes it.
- **`Time.ToDuration`**, the reverse of `ParseDuration`: it rounds to the nearest nanosecond and fails with `ErrOutOfRange` for NaN or beyond about ±292 years.
- **Unit packages:** `ParseUnit`, `Units`, and the `Unit` methods `String`, `Factor` and `Valid`.
- **Look-alike characters** in symbols are accepted as Unicode NFKC maps them: `µ` (MICRO SIGN) for `μ`, `Ω` (OHM SIGN) for `Ω`, `Å` (ANGSTROM SIGN) for `Å` and `K` (KELVIN SIGN) for `K`.
- **`distance.PlanckLength`** and **`time.PlanckTime`**. The misspelled `PlankLength` and `PlankTime` remain as deprecated aliases, which `go fix` replaces.
- Runnable examples on pkg.go.dev for every measure, every cross-measure operation, custom measures, JSON, parsing and formatting.
- **[USECASES.md](USECASES.md):** real-world scenarios (a fitness tracker, IoT telemetry, service configuration, importing mixed data, electronics, home energy, aviation fuel, pressure and recipe scaling), each a verified example. The README is reorganised, with a quick start, a guide to every feature and performance figures.

### Fixed

- Converting between metres and millimetres no longer gives a wrong result. The rate table keyed pairs by concatenating the two symbols, and `"m"+"mm"` equals `"mm"+"m"`, so one of the two directions was wrong in every process.
- Using measurements from several goroutines no longer causes a data race. The conversion tables were built lazily on first use; they are now built once at startup.
- `ToDistance`, `ToMass`, `ToTime` and `ToVelocity` accept a value that already is that measure, and return an error instead of panicking on `nil`.
- A custom unit that shares a symbol with a predefined unit, such as `"m"`, is no longer treated as that predefined measure.
- `AddMeasurement` and `SubMeasurement` return an error for measurements of different measures. They used to count the other operand as zero.
- More precise conversion factors: km/h, knot, ounce, pound, light-year, parsec and Planck time. For example, 100 km / 1 h now gives exactly 100 km/h.

### Performance

- The typed operations don't allocate. Creating, converting, adding and multiplying a value takes about 2–3 ns, and comparing values in different units about 5 ns; the interface design took 30–90 ns and one allocation. Summing 1000 distances in mixed units takes about 3.3 µs.
- `StringToDistance` and the other `StringToX` functions don't allocate. `ParseString` and `Parser.Parse` allocate once, 16 bytes, to return the `Measurement` interface.
- Error paths are cheap too: an unknown unit costs one allocation, a malformed string two, and constant errors such as `ErrZeroDivision` none.
- `TestAllocations` fixes these counts, and CI fails if one rises.
