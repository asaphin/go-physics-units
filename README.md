# go-physics-units

[![Go Reference](https://pkg.go.dev/badge/github.com/asaphin/go-physics-units.svg)](https://pkg.go.dev/github.com/asaphin/go-physics-units)
[![CI](https://github.com/asaphin/go-physics-units/actions/workflows/ci.yml/badge.svg)](https://github.com/asaphin/go-physics-units/actions/workflows/ci.yml)

go-physics-units is a Go library for working with physical measurements: distance, area, volume, mass, time, velocity, acceleration, force, energy, power, pressure, frequency, angle, data size, density, data rate and temperature.
You can create, parse, convert, compare, combine and encode measurements, and define measures of your own, such as currencies.
It has no dependencies.

- **Fast:** measurements are plain values, and units are typed constants that index tables built at startup. Converting, adding or comparing takes 2–5 ns and never allocates.
- **Safe:** a unit of the wrong measure doesn't compile, and the typed operations can't fail, except division by zero.
- **Immutable:** values and the library's tables never change after they are created. No importer can alter them, and they are safe to use from many goroutines.
- **Customizable:** your own measures get the same typed API as the built-in ones, and can extend them.

> [!WARNING]
> The library is under active development and its API is not stable yet (v0).
> Upgrading from v0.0.4? See [MIGRATION.md](MIGRATION.md).

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

Each measure is a small value type, such as `units.Distance`. Its zero value is 0 in the base unit (0 m), and every method returns a new value.

Adding or subtracting values of the same measure works across units, and the result keeps the receiver's unit.
`Compare`, `Equal` and `EqualWithinRel` compare values across units too:

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
| `Mass.DivideByVolume`         | `Density` in kg/m³                                 |
| `Mass.DivideByDensity`        | `Volume` in m³                                     |
| `Density.MultiplyByVolume`    | `Mass` in kg                                       |
| `DataSize.DivideByTime`       | `DataRate` in B/s                                  |
| `DataSize.DivideByDataRate`   | `Time` in s                                        |
| `DataRate.MultiplyByTime`     | `DataSize` in B                                    |
| `ParseDuration`               | `Time` in s, from a `time.Duration`                |
| `Time.ToDuration`             | `time.Duration`, rounded to the nearest nanosecond |

The divisions fail with `ErrZeroDivision` when the divisor is zero, as `Div(0)` does.

## Generic code and parsing

Every value implements the `units.Measurement` interface, for code that handles any measure.
`units.ParseString` detects the measure from the unit symbol and returns a `Measurement` holding the measure's concrete type:

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

A unit symbol that arrives as a string, for example from configuration, converts to a typed unit with its package's `ParseUnit`: `u, ok := distance.ParseUnit("km")`.

## Parsing with options

`ParseString` is strict: exactly `"<value> <unit>"` with a predefined, case-sensitive symbol.
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

## JSON and text

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

## Temperature

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

## Errors

Every error matches exactly one sentinel, which you can test for with `errors.Is`:

| Sentinel                  | Returned when                                                                             |
|---------------------------|-------------------------------------------------------------------------------------------|
| `ErrUnknownUnit`          | a unit symbol doesn't belong to the measure (a `*UnitError`)                              |
| `ErrIncompatibleMeasures` | the measurements belong to different measures, or one is nil                              |
| `ErrZeroDivision`         | the divisor is zero                                                                       |
| `ErrInvalidFormat`        | the string isn't `"<value> <unit>"` with a finite value                                   |
| `ErrInvalidFactors`       | custom conversion factors are invalid or lack the base unit                               |
| `ErrInvalidMeasure`       | `NewUnitSet` gets an empty or predefined measure name, or a measure has no units          |
| `ErrAmbiguousUnit`        | a case-insensitive `Parser` gets a unit that matches several symbols, such as `KN`        |
| `ErrInvalidOption`        | `NewParser` gets an invalid alias, or two measures that share a symbol                    |
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

## Custom measures

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

### Extending a predefined measure

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
| Density      | `density`      | kg/m³     | kg/m³, g/L, lb/ft³, lb/gal, kg/L, g/mL, g/cm³                             |
| Data rate    | `datarate`     | B/s       | bit/s, B/s, kbit/s, kB/s, KiB/s, Mbit/s, MB/s, MiB/s, Gbit/s, GB/s, GiB/s |
| Temperature  | `temperature`  | K         | K, °C, °F, °R                                                             |

- Every unit is a typed constant in its package, for example `distance.LightYear` or `mass.Pound`; `String()` returns its symbol.
- Each unit package also provides `ParseUnit(symbol)`, `Units()`, `ConversionFactors()`, a copy of its factors, and `Rate(from, to)`, the multiplier between two units. For example, `distance.Rate(distance.Mile, distance.Kilometer)` returns 1.609344. The `temperature` package has `Convert` and `DifferenceRate` instead.
- Store and exchange units by their symbols, never as numbers: the numbers behind the constants may change between versions.
- Symbols are case-sensitive.
- The micro prefix can be written as `μ` (U+03BC GREEK SMALL LETTER MU) or `µ` (U+00B5 MICRO SIGN). Values always carry the `μ` form.
- `yr` is the Julian year of 365.25 days.
- `pt`, `qt` and `gal` are US liquid measures, `bbl` the 42-gallon oil barrel, `cal` the thermochemical calorie, `Btu` the International Table one, `hp` mechanical horsepower and `PS` metric horsepower. `mmHg` and `inHg` are the conventional values, and `gₙ` (with U+2099) is standard gravity, 9.80665 m/s².
- Mind the case: `kN` is a kilonewton and `kn` a knot.
- `℃` (U+2103) and `℉` (U+2109) are accepted for `°C` and `°F`.
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
make generate # regenerate the code every measure shares, after changing internal/gen
```

`make check` needs golangci-lint v2.15.0, the version CI uses; `make tools` installs it into `./bin`, which `make` prefers over the one on your `PATH`. `make help` lists all targets.

Guidelines for contributors and AI coding agents are in [AGENTS.md](AGENTS.md).

## License

MIT License
