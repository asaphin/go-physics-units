# go-physics-units

[![Go Reference](https://pkg.go.dev/badge/github.com/asaphin/go-physics-units.svg)](https://pkg.go.dev/github.com/asaphin/go-physics-units)
[![CI](https://github.com/asaphin/go-physics-units/actions/workflows/ci.yml/badge.svg)](https://github.com/asaphin/go-physics-units/actions/workflows/ci.yml)

go-physics-units is a Go library for working with physical measurements: distance, area, volume, mass, time, velocity, acceleration, force, energy, power, pressure, frequency, angle and data size.
You can create, parse, convert and combine measurements, and define measures of your own, such as currencies.
It has no dependencies.

- **Fast:** conversion tables are built once, at startup. A conversion is one lookup and one multiplication.
- **Immutable:** measurements and the library's tables never change after they are created. No importer can alter them, and they are safe to use from many goroutines.
- **Customizable:** your own measures and units work like the built-in ones, and can extend them.

> [!WARNING]
> The library is under active development and its API is not stable yet (v0).

## Installation

```shell
go get github.com/asaphin/go-physics-units@latest
```

Requires Go 1.27 or later.

## Usage

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
	// Creating a distance of 5 kilometers.
	dist, err := units.NewDistance(5, distance.Kilometer)
	if err != nil {
		panic(err)
	}

	fmt.Println(dist) // 5 km

	// Parsing a distance from a string.
	dist, err = units.StringToDistance("10 km")
	if err != nil {
		panic(err)
	}

	// Converting the distance to nautical miles.
	nmi, err := dist.ConvertTo(distance.NauticalMile)
	if err != nil {
		panic(err)
	}

	fmt.Println(nmi) // 5.399568034557235 nmi

	// fmt verbs format the value and keep the unit.
	fmt.Printf("%.2f\n", nmi) // 5.40 nmi

	// Accessing the value and unit.
	fmt.Println(dist.Value(), dist.Unit()) // 10 km

	// Dividing a distance by a time gives a velocity in m/s.
	tm, err := units.NewTime(10, time.Minute)
	if err != nil {
		panic(err)
	}

	vel, err := dist.DivideByTime(tm)
	if err != nil {
		panic(err)
	}

	fmt.Println(vel)                                          // 16.666666666666668 m/s
	fmt.Println(vel.MustConvertTo(velocity.KilometerPerHour)) // 60 km/h
}
```

`units.ParseString` parses a string into a generic `Measurement` when you don't know the measure in advance.
`ToDistance`, `ToMass` and the other `ToX` functions narrow it to a specific measure.

Adding or subtracting measurements of the same measure works across units, and the result keeps the receiver's unit.
`Compare`, `Equal` and `EqualWithinRel` compare measurements across units too:

```go
mile, _ := units.NewDistance(1, distance.Mile)
foot, _ := units.NewDistance(1, distance.Feet)
inches, _ := units.NewDistance(12, distance.Inch)

fmt.Println(mile.Compare(foot)) // 1

// Conversions round: in metres, 12 in and 1 ft differ in the last bit.
fmt.Println(foot.Equal(inches), foot.EqualWithinRel(inches, 1e-12)) // false true

ds := []units.Distance{mile, foot, inches}
slices.SortFunc(ds, units.Distance.Compare)
fmt.Println(ds) // [12 in 1 ft 1 mi]
```

Operations that combine measures return their result in SI base units:

| Operation                     | Result                                             |
|-------------------------------|----------------------------------------------------|
| `Distance.DivideByTime`       | `Velocity` in m/s                                  |
| `Distance.Square`             | `Area` in m²                                       |
| `Area.DivideByDistance`       | `Distance` in m                                    |
| `Velocity.MultiplyByTime`     | `Distance` in m                                    |
| `Velocity.DivideByTime`       | `Acceleration` in m/s²                             |
| `Distance.Cube`               | `Volume` in m³                                     |
| `Area.MultiplyByDistance`     | `Volume` in m³                                     |
| `Volume.DivideByArea`         | `Distance` in m                                    |
| `Volume.DivideByDistance`     | `Area` in m²                                       |
| `Acceleration.MultiplyByTime` | `Velocity` in m/s                                  |
| `Mass.MultiplyByAcceleration` | `Force` in N                                       |
| `Acceleration.MultiplyByMass` | `Force` in N                                       |
| `Force.DivideByMass`          | `Acceleration` in m/s²                             |
| `Force.DivideByAcceleration`  | `Mass` in kg                                       |
| `Force.DivideByArea`          | `Pressure` in Pa                                   |
| `Force.MultiplyByDistance`    | `Energy` in J                                      |
| `Energy.DivideByTime`         | `Power` in W                                       |
| `Energy.DivideByDistance`     | `Force` in N                                       |
| `Power.MultiplyByTime`        | `Energy` in J                                      |
| `Pressure.MultiplyByArea`     | `Force` in N                                       |
| `Frequency.Period`            | `Time` in s                                        |
| `ParseDuration`               | `Time` in s, from a `time.Duration`                |
| `Time.ToDuration`             | `time.Duration`, rounded to the nearest nanosecond |

## Parsing with options

`ParseString` is strict: exactly `"<value> <unit>"` with a predefined, case-sensitive symbol.
A `Parser` accepts more, with the options you choose. It is immutable once built and safe for concurrent use:

```go
p, err := units.NewParser(
	units.WithAliases(map[string]string{"metres": distance.Meter, "kph": velocity.KilometerPerHour}),
	units.CaseInsensitive(), // "10 KM" is 10 km, but "1 KN" is ambiguous: kN or kn
	units.AllowNoSpace(),    // "100kph"
)
if err != nil {
	panic(err)
}

m, err := p.Parse("5 metres") // 5 m, a distance
```

`WithCreator(mc)` adds the units of a custom measure, and `p.Unit("metres")` resolves an alias to its symbol for the typed constructors.

## Errors

Every error matches exactly one sentinel, which you can test for with `errors.Is`:

| Sentinel                  | Returned when                                                                            |
|---------------------------|------------------------------------------------------------------------------------------|
| `ErrUnknownUnit`          | the unit doesn't belong to the measure (a `*UnitError`)                                  |
| `ErrIncompatibleMeasures` | the measurements belong to different measures, or one is nil                             |
| `ErrZeroDivision`         | the divisor is zero                                                                      |
| `ErrInvalidFormat`        | the string isn't `"<value> <unit>"` with a finite value                                  |
| `ErrInvalidFactors`       | custom conversion factors are invalid or lack the base unit                              |
| `ErrInvalidMeasure`       | `NewMeasurementCreator` gets an empty or predefined measure name, or an unsupported type |
| `ErrAmbiguousUnit`        | a case-insensitive `Parser` gets a unit that matches several symbols, such as `KN`       |
| `ErrInvalidOption`        | `NewParser` gets an invalid alias, or two measures that share a symbol                   |
| `ErrOutOfRange`           | `Time.ToDuration` gets NaN or more than about 292 years                                  |

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

## Custom measures

You can also create a factory for a measure of your own.
The factors say how many base units there are in one unit of each symbol:

```go
package main

import (
	"fmt"

	units "github.com/asaphin/go-physics-units"
	"github.com/asaphin/go-physics-units/conversion"
)

func main() {
	type Currency interface {
		units.Measurement
	}

	currencyFactors := conversion.Factors{
		"EUR": 1.1,
		"USD": 1,
		"UAH": 0.027,
	}

	mc, err := units.NewMeasurementCreator[Currency]("currency", "USD", currencyFactors)
	if err != nil {
		panic(err)
	}

	usd, err := mc.New(1, "USD")
	if err != nil {
		panic(err)
	}

	uah, err := usd.ConvertToMeasurement("UAH")
	if err != nil {
		panic(err)
	}

	fmt.Println(uah)        // 37.03703703703704 UAH
	fmt.Println(uah.Type()) // currency
}
```

### Extending a predefined measure

`ConversionFactors()` returns a copy of a measure's factors. Add units to it and create a measure of your own; the predefined measure stays unchanged:

```go
type Surveying interface{ units.Measurement }

factors := distance.ConversionFactors() // a copy: changing it is safe
factors["fur"] = 201.168                // furlong, 660 international feet
factors["ch"] = 20.1168                 // chain, 66 international feet

mc, err := units.NewMeasurementCreator[Surveying]("surveying", distance.BaseUnit, factors)
if err != nil {
	panic(err)
}

furlong, err := mc.New(1, "fur")
if err != nil {
	panic(err)
}

chains, _ := furlong.ConvertToMeasurement("ch")
metres, _ := furlong.ConvertToMeasurement(distance.Meter)

fmt.Println(chains, metres, furlong.Type()) // 10 ch 201.168 m surveying

_, err = units.NewDistance(1, "fur")
fmt.Println(err) // unknown distance unit fur
```

## Supported units

| Measure      | Package        | Base unit | Symbols                                                                   |
|--------------|----------------|-----------|---------------------------------------------------------------------------|
| Distance     | `distance`     | m         | lp, fm, pm, Å, nm, μm, mm, cm, dm, m, km, in, ft, yd, mi, nmi, AU, ly, pc |
| Area         | `area`         | m²        | mm², cm², in², ft², m², km², mi²                                          |
| Mass         | `mass`         | kg        | Da, μg, mg, g, kg, oz, lb, q, t, M⊕, M☉                                   |
| Time         | `time`         | s         | tp, ns, μs, ms, s, min, h, d, wk, yr                                      |
| Velocity     | `velocity`     | m/s       | km/h, ft/s, mph, kn, m/s, c                                               |
| Volume       | `volume`       | m³        | mm³, μL, cm³, mL, in³, L, pt, qt, gal, ft³, bbl, yd³, m³, km³             |
| Acceleration | `acceleration` | m/s²      | cm/s², ft/s², m/s², gₙ                                                    |
| Force        | `force`        | N         | dyn, mN, N, lbf, kgf, kN, MN                                              |
| Energy       | `energy`       | J         | eV, erg, J, ft·lbf, cal, kJ, Btu, kcal, Wh, MJ, kWh, GJ, MWh              |
| Power        | `power`        | W         | mW, W, Btu/h, PS, hp, kW, MW, GW                                          |
| Pressure     | `pressure`     | Pa        | Pa, hPa, mbar, kPa, Torr, mmHg, inHg, psi, bar, atm, MPa, GPa             |
| Frequency    | `frequency`    | Hz        | rpm, Hz, kHz, MHz, GHz, THz                                               |
| Angle        | `angle`        | rad       | arcsec, arcmin, mrad, gon, °, rad, rev                                    |
| Data size    | `datasize`     | B         | bit, B, kbit, kB, KiB, Mbit, MB, MiB, Gbit, GB, GiB, TB, TiB, PB, PiB     |

- Every symbol has a named constant in its package, for example `distance.LightYear` or `mass.Pound`.
- Each unit package also provides `ConversionFactors()`, a copy of its factors, and `Rate(from, to)`, the multiplier between two units. For example, `distance.Rate(distance.Mile, distance.Kilometer)` returns 1.609344.
- Symbols are case-sensitive.
- The micro prefix can be written as `μ` (U+03BC GREEK SMALL LETTER MU) or `µ` (U+00B5 MICRO SIGN). Values always carry the `μ` form.
- `yr` is the Julian year of 365.25 days.
- `pt`, `qt` and `gal` are US liquid measures, `bbl` the 42-gallon oil barrel, `cal` the thermochemical calorie, `Btu` the International Table one, `hp` mechanical horsepower and `PS` metric horsepower. `mmHg` and `inHg` are the conventional values, and `gₙ` (with U+2099) is standard gravity, 9.80665 m/s².
- Mind the case: `kN` is a kilonewton and `kn` a knot.
- Data sizes have both decimal (`kB` = 1000 B) and binary (`KiB` = 1024 B) prefixes, and bits (`bit`, `Mbit`, …).
- Strings are parsed in the form `"<value> <unit>"`, with whitespace between the value and the unit. NaN and infinite values are rejected.

## Documentation

The full API reference is on [pkg.go.dev](https://pkg.go.dev/github.com/asaphin/go-physics-units).

## Development

```shell
make tools   # once: installs the pinned golangci-lint into ./bin
make check   # formatting, go vet, golangci-lint and tests with the race detector
make cover   # coverage summary, plus an HTML report in coverage/
make bench   # benchmarks
make fuzz    # fuzz the string parser (FUZZTIME=30s by default)
```

`make check` needs golangci-lint v2.15.0, the version CI uses; `make tools` installs it into `./bin`, which `make` prefers over the one on your `PATH`. `make help` lists all targets.

Guidelines for contributors and AI coding agents are in [AGENTS.md](AGENTS.md).

## License

MIT License
