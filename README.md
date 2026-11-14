# go-physics-units

**Type-safe physical units for Go.** Create, convert, combine, compare, parse and encode measurements in a few nanoseconds, without allocations or dependencies.

[![Go Reference](https://pkg.go.dev/badge/github.com/asaphin/go-physics-units.svg)](https://pkg.go.dev/github.com/asaphin/go-physics-units)
[![CI](https://github.com/asaphin/go-physics-units/actions/workflows/ci.yml/badge.svg)](https://github.com/asaphin/go-physics-units/actions/workflows/ci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/asaphin/go-physics-units)](https://goreportcard.com/report/github.com/asaphin/go-physics-units)
[![Release](https://img.shields.io/github/v/release/asaphin/go-physics-units?sort=semver)](https://github.com/asaphin/go-physics-units/releases)
[![Go version](https://img.shields.io/github/go-mod/go-version/asaphin/go-physics-units)](go.mod)
[![License: MIT](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

```go
dist := units.NewDistance(10, distance.Kilometer)
vel, _ := dist.DivideByTime(units.NewTime(10, time.Minute))

fmt.Println(vel.ConvertTo(velocity.KilometerPerHour)) // 60 km/h
```

## Highlights

- **Fast.** Values are plain structs, and units are typed constants that index tables built at startup. Converting, adding or comparing takes 2–5 ns and never allocates.
- **Type-safe.** `units.NewDistance(5, mass.Kilogram)` doesn't compile. The typed operations can't fail, except division by zero.
- **Immutable and safe for concurrent use.** Values never change, and the library's tables are fixed once the package is initialised. No importer can alter them, and no locks are needed.
- **Extensible.** Your own measures, such as currencies, get the same typed API as the built-in ones, and can extend them.
- **Broad.** 22 measures with 184 units, more than 40 operations between measures (force is mass × acceleration, power is voltage × current, …), `fmt` verbs, JSON and text encoding, a configurable parser, and temperatures with their offsets.
- **Dependency-free and well tested.** Only the standard library, 100% test coverage, fuzzed parsers, and allocation counts checked in CI.

> [!NOTE]
> The API can still change before v1.0.0; see [Versioning](#versioning).
> Upgrading from v0.0.4? [MIGRATION.md](MIGRATION.md) lists every change, with before-and-after code.

## Contents

- [Installation](#installation)
- [Quick start](#quick-start)
- [Usage](#usage): [values and units](#values-and-units), [comparing](#comparing-values), [combining measures](#combining-measures), [formatting](#formatting), [parsing](#parsing), [JSON](#json-and-text-encoding), [temperature](#temperature), [custom measures](#custom-measures), [errors](#errors)
- [Supported units](#supported-units)
- [Performance](#performance)
- [Design](#design)
- [Versioning](#versioning)
- [Contributing](#contributing)
- [License](#license)

## Installation

```shell
go get github.com/asaphin/go-physics-units@latest
```

Requires Go 1.27 or later. The root package is `units`; each measure's units live in a subpackage of the same name, such as `distance` or `mass`.

## Quick start

```go
package main

import (
	"fmt"

	units "github.com/asaphin/go-physics-units"
	"github.com/asaphin/go-physics-units/distance"
	"github.com/asaphin/go-physics-units/time"
	"github.com/asaphin/go-physics-units/velocity"
)

func main() {
	// Creating a distance of 5 kilometers. Units are typed constants, so
	// a unit of another measure doesn't compile.
	dist := units.NewDistance(5, distance.Kilometer)

	fmt.Println(dist) // 5 km

	// Parsing a distance from a string.
	dist, err := units.StringToDistance("10 km")
	if err != nil {
		panic(err)
	}

	// Converting the distance to nautical miles.
	nmi := dist.ConvertTo(distance.NauticalMile)

	fmt.Println(nmi) // 5.399568034557235 nmi

	// fmt verbs format the value and keep the unit.
	fmt.Printf("%.2f\n", nmi) // 5.40 nmi

	// Accessing the value and unit.
	fmt.Println(dist.Value(), dist.Unit()) // 10 km

	// Dividing a distance by a time gives a velocity in m/s.
	tm := units.NewTime(10, time.Minute)

	vel, err := dist.DivideByTime(tm)
	if err != nil {
		panic(err) // the time is zero
	}

	fmt.Println(vel)                                      // 16.666666666666668 m/s
	fmt.Println(vel.ConvertTo(velocity.KilometerPerHour)) // 60 km/h
}
```

The `time` subpackage shadows the standard library's `time`. Import the standard one under another name, such as `stdtime "time"`, when you need both.

## Usage

### Values and units

Each measure is a small value type, such as `units.Distance`, holding a `float64` and a unit.
Its zero value is 0 in the base unit (0 m), and every method returns a new value.

Adding or subtracting values of the same measure works across units, and the result keeps the receiver's unit:

```go
km := units.NewDistance(1, distance.Kilometer)
m := units.NewDistance(250, distance.Meter)

fmt.Println(km.Add(m)) // 1.25 km
fmt.Println(m.Add(km)) // 1250 m
```

`Mul` and `Div` scale a value by a number; `Div(0)` returns `ErrZeroDivision`.
A unit symbol that arrives as a string, for example from configuration, converts to a typed unit with its package's `ParseUnit`: `u, ok := distance.ParseUnit("km")`.

### Comparing values

`Compare`, `Equal` and `EqualWithinRel` compare values across units, and give the same answer whichever operand comes first:

```go
mile := units.NewDistance(1, distance.Mile)
foot := units.NewDistance(1, distance.Feet)
inches := units.NewDistance(12, distance.Inch)

fmt.Println(mile.Compare(foot)) // 1

// Conversions round: in metres, 12 in and 1 ft differ in the last bit.
fmt.Println(foot.Equal(inches), foot.EqualWithinRel(inches, 1e-12)) // false true

ds := []units.Distance{mile, foot, inches}
slices.SortFunc(ds, units.Distance.Compare)
fmt.Println(ds) // [12 in 1 ft 1 mi]
```

Values are `float64`, so conversions round to the nearest representable number. Use `EqualWithinRel` for values that went through a conversion.

### Combining measures

Operations that combine measures return their result in the SI base unit of the result (m/s, N, J, W, …):

```go
climber := units.NewMass(75, mass.Kilogram)
g := units.NewAcceleration(1, acceleration.StandardGravity)
height := units.NewDistance(3, distance.Meter)
climb := units.NewTime(4, time.Second)

weight := climber.MultiplyByAcceleration(g)
work := weight.MultiplyByDistance(height)

p, err := work.DivideByTime(climb)
if err != nil {
	panic(err) // the time is zero
}

fmt.Printf("%.1f, %.1f, %.1f = %.2f\n", weight, work, p, p.ConvertTo(power.Horsepower))
// 735.5 N, 2206.5 J, 551.6 W = 0.74 hp
```

| Value             | Operations and their results                                                                                                    |
|-------------------|---------------------------------------------------------------------------------------------------------------------------------|
| `Distance`        | `Square` → `Area`, `Cube` → `Volume`, `DivideByTime` → `Velocity`                                                               |
| `Area`            | `MultiplyByDistance` → `Volume`, `DivideByDistance` → `Distance`                                                                |
| `Volume`          | `DivideByArea` → `Distance`, `DivideByDistance` → `Area`                                                                        |
| `Angle`           | `DivideByTime` → `AngularVelocity`                                                                                              |
| `Time`            | `ToDuration` → `time.Duration`, rounded to the nearest nanosecond; `ParseDuration` converts back                                |
| `Frequency`       | `Period` → `Time`                                                                                                               |
| `Velocity`        | `MultiplyByTime` → `Distance`, `DivideByTime` → `Acceleration`                                                                  |
| `Acceleration`    | `MultiplyByTime` → `Velocity`, `MultiplyByMass` → `Force`                                                                       |
| `AngularVelocity` | `MultiplyByTime` → `Angle`                                                                                                      |
| `Mass`            | `MultiplyByAcceleration` → `Force`, `DivideByVolume` → `Density`, `DivideByDensity` → `Volume`                                  |
| `Density`         | `MultiplyByVolume` → `Mass`                                                                                                     |
| `Force`           | `DivideByMass` → `Acceleration`, `DivideByAcceleration` → `Mass`, `DivideByArea` → `Pressure`, `MultiplyByDistance` → `Energy`  |
| `Pressure`        | `MultiplyByArea` → `Force`                                                                                                      |
| `Energy`          | `DivideByTime` → `Power`, `DivideByDistance` → `Force`, `DivideByCharge` → `Voltage`                                            |
| `Power`           | `MultiplyByTime` → `Energy`, `DivideByVoltage` → `Current`, `DivideByCurrent` → `Voltage`                                       |
| `Current`         | `MultiplyByVoltage` → `Power`, `MultiplyByResistance` → `Voltage`, `MultiplyByTime` → `Charge`                                  |
| `Voltage`         | `MultiplyByCurrent` → `Power`, `DivideByCurrent` → `Resistance` (Ohm's law), `DivideByResistance` → `Current`, `MultiplyByCharge` → `Energy` |
| `Resistance`      | `MultiplyByCurrent` → `Voltage`                                                                                                 |
| `Charge`          | `DivideByTime` → `Current`, `DivideByCurrent` → `Time` (how long a battery lasts), `MultiplyByVoltage` → `Energy`               |
| `DataSize`        | `DivideByTime` → `DataRate`, `DivideByDataRate` → `Time`                                                                        |
| `DataRate`        | `MultiplyByTime` → `DataSize`                                                                                                   |

The divisions fail with `ErrZeroDivision` when the divisor is zero, as `Div(0)` does.

### Formatting

`String` gives `"<value> <unit>"`, and every value implements `fmt.Formatter`, so the usual verbs format the number and keep the unit:

| Verb            | `units.NewDistance(1234.5678, distance.Meter)` |
|-----------------|------------------------------------------------|
| `%v`, `%s`      | `1234.5678 m`                                  |
| `%.1f`          | `1234.6 m`                                     |
| `[%12.2f]`      | `[     1234.57 m]`                             |
| `%.3e`          | `1.235e+03 m`                                  |
| `%q`            | `"1234.5678 m"`                                |
| `%#v`           | `units.Distance{Value:1234.5678, Unit:"m"}`    |

### Parsing

`StringToDistance` and its siblings parse `"<value> <unit>"` into a typed value without allocating.
When the measure isn't known in advance, `units.ParseString` detects it from the unit symbol. It returns a `units.Measurement`, the interface every value implements, holding the measure's concrete type:

```go
m, err := units.ParseString("1500 m")
if err != nil {
	panic(err)
}

fmt.Println(m.Type()) // distance

if d, ok := m.(units.Distance); ok {
	fmt.Println(d.ConvertTo(distance.Kilometer)) // 1.5 km
}
```

`ParseString` is strict: exactly `"<value> <unit>"` with a predefined, case-sensitive symbol, and a finite value.
A `Parser` accepts more, with the options you choose. It is immutable once built and safe for concurrent use:

```go
p, err := units.NewParser(
	units.WithAliases(map[string]string{"metres": "m", "kph": "km/h"}),
	units.CaseInsensitive(), // "10 KM" is 10 km, but "1 KN" is ambiguous: kN or kn
	units.AllowNoSpace(),    // "100kph"
)
if err != nil {
	panic(err)
}

m, err := p.Parse("5 metres") // 5 m, a units.Distance
```

`WithMeasure[M]()` adds the units of a custom measure, and `p.Unit("metres")` resolves an alias to its symbol.

### JSON and text encoding

Values implement `encoding.TextMarshaler` and `encoding.TextUnmarshaler` in the `"<value> <unit>"` form, so they work in JSON, YAML, `flag.TextVar` and configuration loaders:

```go
type Route struct {
	Name   string         `json:"name"`
	Length units.Distance `json:"length"`
}

b, _ := json.Marshal(Route{"loop", units.NewDistance(12.5, distance.Kilometer)})
fmt.Println(string(b)) // {"name":"loop","length":"12.5 km"}
```

Unmarshalling checks the unit: `"3.2 kg"` is an error for a `Distance`. Values that aren't finite can't be marshalled.

### Temperature

Temperature scales have offsets (0 °C is 273.15 K), so temperatures have two types:

- `Temperature` is a point on a scale, such as 20 °C. It can't be below absolute zero: constructors and operations that would go below it fail with `ErrOutOfRange`.
- `TemperatureDifference` is a change, such as a rise of 10 °C, which is a rise of 18 °F.

```go
room, err := units.NewTemperature(20, temperature.Celsius)
if err != nil {
	panic(err) // below absolute zero
}

fmt.Println(room.ConvertTo(temperature.Fahrenheit)) // 68 °F
fmt.Println(room.Kelvin())                           // 293.15

noon, _ := units.StringToTemperature("77 °F")
rise := noon.Sub(room) // a TemperatureDifference: 9 °F
warmer, err := room.Add(units.NewTemperatureDifference(5, temperature.Kelvin)) // 25 °C
```

Conversions between the defined points are exact: 100 °C is exactly 212 °F.
Adding two temperatures, or multiplying one, isn't offered, since it has no physical meaning.
Temperatures don't implement `Measurement`; `ParseString` rejects their units and points to `StringToTemperature`.

### Custom measures

A measure of your own takes four declarations: a marker type, its units, the method that connects them, and a name for the value type.
The factors say how many base units there are in one unit of each symbol:

```go
package main

import (
	"fmt"

	units "github.com/asaphin/go-physics-units"
	"github.com/asaphin/go-physics-units/conversion"
)

type currency struct{}

var currencyUnits = units.MustNewUnitSet("currency", "USD", conversion.Factors{
	"USD": 1,
	"EUR": 1.1,
	"UAH": 0.027,
})

func (currency) Units() units.UnitSet { return currencyUnits }

type Money = units.Quantity[currency]

func main() {
	usd := units.MustUnitOf[currency]("USD")
	uah := units.MustUnitOf[currency]("UAH")

	price := units.NewQuantity(1, usd)
	tip, err := units.ParseQuantity[currency]("0.5 EUR")
	if err != nil {
		panic(err)
	}

	var total Money = price.Add(tip)

	fmt.Println(total)                // 1.55 USD
	fmt.Println(total.ConvertTo(uah)) // 57.40740740740741 UAH
	fmt.Println(total.Type())         // currency
}
```

`Quantity[M]` has the typed API of the built-in measures: `ConvertTo`, `Add`, `Sub`, `Mul`, `Div`, `Compare`, `Equal`, text encoding, and the `Measurement` interface.
Look units up by symbol once, with `UnitOf`, and reuse them.
A custom measure's symbols are its own: a custom `"m"` is not a distance.

#### Extending a predefined measure

`ConversionFactors()` returns a copy of a measure's factors. Add units to it and define a measure of your own; the predefined measure stays unchanged:

```go
type surveying struct{}

var surveyingUnits = func() units.UnitSet {
	factors := distance.ConversionFactors() // a copy: changing it is safe
	factors["fur"] = 201.168                // furlong, 660 international feet
	factors["ch"] = 20.1168                 // chain, 66 international feet

	return units.MustNewUnitSet("surveying", distance.BaseUnit.String(), factors)
}()

func (surveying) Units() units.UnitSet { return surveyingUnits }

furlong, err := units.ParseQuantity[surveying]("1 fur")
if err != nil {
	panic(err)
}

fmt.Println(furlong.ConvertTo(units.MustUnitOf[surveying]("ch"))) // 10 ch
fmt.Println(furlong.ConvertToBaseUnits())                         // 201.168 m

_, err = units.StringToDistance("1 fur")
fmt.Println(err) // unknown distance unit fur
```

### Errors

Every error matches exactly one sentinel, which you can test for with `errors.Is`:

| Sentinel                  | Returned when                                                                              |
|---------------------------|--------------------------------------------------------------------------------------------|
| `ErrUnknownUnit`          | a unit symbol doesn't belong to the measure (a `*UnitError`)                               |
| `ErrIncompatibleMeasures` | the measurements belong to different measures, or one is nil                               |
| `ErrZeroDivision`         | the divisor is zero                                                                        |
| `ErrInvalidFormat`        | the string isn't `"<value> <unit>"` with a finite value                                    |
| `ErrInvalidFactors`       | custom conversion factors are invalid or lack the base unit                                |
| `ErrInvalidMeasure`       | `NewUnitSet` gets an empty or predefined measure name, or a measure has no units           |
| `ErrAmbiguousUnit`        | a case-insensitive `Parser` gets a unit that matches several symbols, such as `KN`         |
| `ErrInvalidOption`        | `NewParser` gets an invalid alias, or two measures that share a symbol                     |
| `ErrOutOfRange`           | a temperature below absolute zero, `Time.ToDuration` beyond ±292 years, or marshalling NaN |

```go
_, err := units.StringToDistance("3 kg")

fmt.Println(errors.Is(err, units.ErrUnknownUnit)) // true

// Unknown-unit errors carry the measure and the unit.
if ue, ok := errors.AsType[*units.UnitError](err); ok {
	fmt.Println(ue.Measure, ue.Unit) // distance kg
}

fmt.Println(err) // unknown distance unit kg
```

The sentinels are constants of type `units.Error`, so no other package can reassign them.

## Supported units

| Measure          | Package           | Base unit | Symbols                                                                   |
|------------------|-------------------|-----------|---------------------------------------------------------------------------|
| Distance         | `distance`        | m         | lp, fm, pm, Å, nm, μm, mm, cm, dm, m, km, in, ft, yd, mi, nmi, AU, ly, pc |
| Area             | `area`            | m²        | mm², cm², in², ft², m², km², mi²                                          |
| Volume           | `volume`          | m³        | mm³, μL, cm³, mL, in³, L, pt, qt, gal, ft³, bbl, yd³, m³, km³             |
| Angle            | `angle`           | rad       | arcsec, arcmin, mrad, gon, °, rad, rev                                    |
| Time             | `time`            | s         | tp, ns, μs, ms, s, min, h, d, wk, yr                                      |
| Frequency        | `frequency`       | Hz        | rpm, Hz, kHz, MHz, GHz, THz                                               |
| Velocity         | `velocity`        | m/s       | km/h, ft/s, mph, kn, m/s, c                                               |
| Acceleration     | `acceleration`    | m/s²      | cm/s², ft/s², m/s², gₙ                                                    |
| Angular velocity | `angularvelocity` | rad/s     | rad/s, °/s, rev/min, rev/s                                                |
| Mass             | `mass`            | kg        | Da, μg, mg, g, kg, oz, lb, q, t, M⊕, M☉                                   |
| Density          | `density`         | kg/m³     | kg/m³, g/L, lb/ft³, lb/gal, kg/L, g/mL, g/cm³                             |
| Force            | `force`           | N         | dyn, mN, N, lbf, kgf, kN, MN                                              |
| Pressure         | `pressure`        | Pa        | Pa, hPa, mbar, kPa, Torr, mmHg, inHg, psi, bar, atm, MPa, GPa             |
| Energy           | `energy`          | J         | eV, erg, J, ft·lbf, cal, kJ, Btu, kcal, Wh, MJ, kWh, GJ, MWh              |
| Power            | `power`           | W         | mW, W, Btu/h, PS, hp, kW, MW, GW                                          |
| Temperature      | `temperature`     | K         | K, °C, °F, °R                                                             |
| Electric current | `current`         | A         | A, μA, mA, kA                                                             |
| Voltage          | `voltage`         | V         | V, μV, mV, kV, MV                                                         |
| Electric charge  | `charge`          | C         | C, μC, mC, mAh, Ah                                                        |
| Resistance       | `resistance`      | Ω         | Ω, mΩ, kΩ, MΩ, GΩ                                                         |
| Data size        | `datasize`        | B         | bit, B, kbit, kB, KiB, Mbit, MB, MiB, Gbit, GB, GiB, TB, TiB, PB, PiB     |
| Data rate        | `datarate`        | B/s       | bit/s, B/s, kbit/s, kB/s, KiB/s, Mbit/s, MB/s, MiB/s, Gbit/s, GB/s, GiB/s |

- Every unit is a typed constant in its package, for example `distance.LightYear` or `mass.Pound`; `String()` returns its symbol.
- Each unit package also provides `ParseUnit(symbol)`, `Units()`, `ConversionFactors()`, a copy of its factors, and `Rate(from, to)`, the multiplier between two units. For example, `distance.Rate(distance.Mile, distance.Kilometer)` returns 1.609344. The `temperature` package has `Convert` and `DifferenceRate` instead.
- Store and exchange units by their symbols, never as numbers: the numbers behind the constants may change between versions.
- Symbols are case-sensitive.
- The micro prefix can be written as `μ` (U+03BC GREEK SMALL LETTER MU) or `µ` (U+00B5 MICRO SIGN). Values always carry the `μ` form.
- `yr` is the Julian year of 365.25 days.
- `pt`, `qt` and `gal` are US liquid measures, `bbl` the 42-gallon oil barrel, `cal` the thermochemical calorie, `Btu` the International Table one, `hp` mechanical horsepower and `PS` metric horsepower. `mmHg` and `inHg` are the conventional values, and `gₙ` (with U+2099) is standard gravity, 9.80665 m/s².
- Mind the case: `kN` is a kilonewton and `kn` a knot.
- `℃` (U+2103) and `℉` (U+2109) are accepted for `°C` and `°F`, `Ω` (U+2126 OHM SIGN) for `Ω` (U+03A9), `Å` (U+212B ANGSTROM SIGN) for `Å` (U+00C5), and `K` (U+212A KELVIN SIGN) for `K`.
- Angular velocity's `rev/min` counts radians; frequency's `rpm` counts cycles. Convert between them through `Angle.DivideByTime` or the factor 2π.
- Data sizes have both decimal (`kB` = 1000 B) and binary (`KiB` = 1024 B) prefixes, and bits (`bit`, `Mbit`, …).
- Strings are parsed in the form `"<value> <unit>"`, with whitespace between the value and the unit. NaN and infinite values are rejected.

## Performance

The typed operations are a few array loads and floating-point operations, inlined by the compiler.
Typical timings on an Intel Core i5-6300HQ (2.3 GHz, linux/amd64, Go 1.27):

| Operation                                          | Time    | Allocations           |
|----------------------------------------------------|--------:|-----------------------|
| `NewDistance`                                      | 2.5 ns  | 0                     |
| `ConvertTo`                                        | 3.2 ns  | 0                     |
| `Add`                                              | 3.1 ns  | 0                     |
| `Mul`                                              | 2.5 ns  | 0                     |
| `Compare`                                          | 4.9 ns  | 0                     |
| `distance.Rate`                                    | 2 ns    | 0                     |
| Summing 1,000 distances in mixed units             | 3.3 μs  | 0                     |
| `StringToDistance("1234.5 mi")`                    | 105 ns  | 0                     |
| `String`                                           | 100 ns  | 1, the result         |
| `ParseString("10 km")`                             | 150 ns  | 1, the interface value |
| `Quantity`: `NewQuantity`, `ConvertTo`, `Add`      | 2–5 ns  | 0                     |

Reproduce them with `make bench`, or compare a change against the latest release with `make bench-compare`.
CI fails if a change adds an allocation to any of these paths. Building a `UnitSet` or a `Parser` costs more (about 1 μs for a set of three units), which is why you build them once.

## Design

Three principles guide the library, and they take priority over convenience:

1. **Performance.** All conversion data is built once, at package initialisation or by `NewUnitSet` and `NewParser`. A conversion at run time is one array load and one multiplication.
2. **Immutability after initialisation.** The package exports no variables: units and sentinel errors are constants, and the conversion tables are unexported and never written after `init`. No importer can change how conversions behave for the rest of the program, and no locks are needed.
3. **Customization.** User-defined measures are first-class: they convert, combine, compare, parse and encode like the built-in ones.

The design decisions are recorded in [docs/rfc-value-types.md](docs/rfc-value-types.md) and [docs/rfc-temperature.md](docs/rfc-temperature.md), and the plan ahead in [ROADMAP.md](ROADMAP.md).

## Versioning

The module follows [semantic versioning](https://semver.org).
Before v1.0.0, a minor release (v0.x.0) may break the API, and a patch release doesn't.
[CHANGELOG.md](CHANGELOG.md) records every user-visible change, and [MIGRATION.md](MIGRATION.md) covers the upgrade from v0.0.4.
Deprecated names keep working until v1.0.0; `go fix` rewrites calls to the ones marked with `//go:fix inline`, such as `MustConvertTo`.

The full API reference, with runnable examples for every measure, is on [pkg.go.dev](https://pkg.go.dev/github.com/asaphin/go-physics-units).

## Contributing

Issues and pull requests are welcome. Open pull requests against the `dev` branch; `main` holds the releases.

```shell
make tools    # once per clone: installs the pinned golangci-lint into ./bin
make check    # what CI runs: formatting, go vet, golangci-lint, tests with the race detector, allocation guards
make cover    # coverage summary, plus an HTML report in coverage/
make bench    # benchmarks
make fuzz     # fuzz the string parser (FUZZTIME=30s by default)
make generate # regenerate the code every measure shares, after changing internal/gen
make help     # every target
```

A change is done when `make check` is clean and coverage stays at 100%.
[AGENTS.md](AGENTS.md) describes the design rules, the layout of the code, and step-by-step recipes for adding a unit or a measure. It is written for human contributors and AI coding agents alike.

## License

go-physics-units is released under the [MIT License](LICENSE).
