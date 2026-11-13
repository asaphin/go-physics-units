# go-physics-units

[![Go Reference](https://pkg.go.dev/badge/github.com/asaphin/go-physics-units.svg)](https://pkg.go.dev/github.com/asaphin/go-physics-units)
[![CI](https://github.com/asaphin/go-physics-units/actions/workflows/ci.yml/badge.svg)](https://github.com/asaphin/go-physics-units/actions/workflows/ci.yml)

go-physics-units is a Go library for working with physical measurements such as distance, area, mass, time and velocity.
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
Operations that combine measures return their result in SI base units:

| Operation                  | Result                    |
|----------------------------|---------------------------|
| `Distance.DivideByTime`    | `Velocity` in m/s         |
| `Distance.Square`          | `Area` in m²              |
| `Area.DivideByDistance`    | `Distance` in m           |
| `Velocity.MultiplyByTime`  | `Distance` in m           |
| `ParseDuration`            | `Time` in s, from a `time.Duration` |

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

## Supported units

| Measure  | Package    | Base unit | Symbols |
|----------|------------|-----------|---------|
| Distance | `distance` | m         | lp, fm, pm, Å, nm, μm, mm, cm, dm, m, km, in, ft, yd, mi, nmi, AU, ly, pc |
| Area     | `area`     | m²        | mm², cm², in², ft², m², km², mi² |
| Mass     | `mass`     | kg        | Da, μg, mg, g, kg, oz, lb, q, t, M⊕, M☉ |
| Time     | `time`     | s         | tp, ns, μs, ms, s, min, h, d, wk, yr |
| Velocity | `velocity` | m/s       | km/h, ft/s, mph, kn, m/s, c |

- Every symbol has a named constant in its package, for example `distance.LightYear` or `mass.Pound`.
- Each unit package also provides `ConversionFactors()`, a copy of its factors, and `Rate(from, to)`, the multiplier between two units. For example, `distance.Rate(distance.Mile, distance.Kilometer)` returns 1.609344.
- Symbols are case-sensitive.
- The micro prefix can be written as `μ` (U+03BC GREEK SMALL LETTER MU) or `µ` (U+00B5 MICRO SIGN). Values always carry the `μ` form.
- `yr` is the Julian year of 365.25 days.
- Strings are parsed in the form `"<value> <unit>"`, with whitespace between the value and the unit. NaN and infinite values are rejected.

## Documentation

The full API reference is on [pkg.go.dev](https://pkg.go.dev/github.com/asaphin/go-physics-units).

## Development

```shell
make check   # formatting, go vet, golangci-lint and tests with the race detector
make cover   # coverage summary, plus an HTML report in coverage/
make bench   # benchmarks
make fuzz    # fuzz the string parser (FUZZTIME=30s by default)
```

`make check` needs golangci-lint v2.15.0, the version CI uses. With any other version it stops and prints the install command.

Guidelines for contributors and AI coding agents are in [AGENTS.md](AGENTS.md).

## License

MIT License
