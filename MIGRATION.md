# Migrating from v0.0.4 to v0.1.0

v0.1.0 redesigns the API around value types and typed units (see [docs/rfc-value-types.md](docs/rfc-value-types.md)). Every program that uses the library needs changes, but most of them are mechanical, and the compiler reports each one. This guide lists them with before-and-after code.

After the redesign the typed operations take 2–5 ns and don't allocate, where they took about 30–90 ns and one allocation. Unit mix-ups become compile errors.

## The checklist

1. Set `go 1.27` or later in your `go.mod`.
2. Drop the error results of `NewX` and `ConvertTo`.
3. Replace string unit literals with the unit constants, and parse string variables with `ParseUnit`.
4. Replace nil checks of measurements with error checks or zero-value checks.
5. Replace `Unit()` with `Symbol()` where you hold a generic `Measurement`.
6. Rewrite custom measures with `NewUnitSet` and `Quantity`.
7. Replace `ConversionRates()` with `Rate(from, to)` or `ConversionFactors()`.
8. Check stored data for the `Qts` symbol, now `q`.
9. Replace error-message matching with `errors.Is`.

## Creating and converting values

Units are typed constants now, so constructors and conversions can't fail and return no error:

```go
// Before
d, err := units.NewDistance(5, distance.Kilometer)
if err != nil { … }
km, err := d.ConvertTo(distance.Kilometer)
if err != nil { … }

// After
d := units.NewDistance(5, distance.Kilometer)
km := d.ConvertTo(distance.Kilometer)
```

`MustConvertTo` still compiles but is deprecated: it is `ConvertTo`.

## Unit symbols

A unit of the wrong measure, or an untyped string, no longer compiles:

```go
units.NewDistance(5, "km")          // doesn't compile: use distance.Kilometer
units.NewDistance(5, mass.Kilogram) // doesn't compile: a mass unit
```

A unit that arrives as a string, for example from configuration, goes through its package's `ParseUnit` once:

```go
// Before
d, err := units.NewDistance(v, cfg.Unit)

// After
u, ok := distance.ParseUnit(cfg.Unit)
if !ok { … } // not a distance unit
d := units.NewDistance(v, u)
```

To parse a whole `"<value> <unit>"` string, `StringToDistance` and the other `StringToX` functions work as before.

Units are integers behind the constants, and the numbers may change between versions. Store and exchange units by their symbols: `u.String()` returns `"km"`, and `ParseUnit` goes back.

Comparing units:

```go
// Before
if d.Unit() == "km" { … }

// After
if d.Unit() == distance.Kilometer { … }
// or, by symbol:
if d.Unit().String() == "km" { … }
```

## Values instead of interfaces

`Distance`, `Mass`, `Time` and `Velocity` are structs now, not interfaces. Function signatures that take or return them stay the same:

```go
func trip(d units.Distance, t units.Time) (units.Velocity, error) { … } // unchanged
```

There is no nil value any more. The zero `Distance` is 0 m, and operations return the zero value together with an error:

```go
// Before
v, err := d.DivideByTime(t)
if v == nil { … }

// After
v, err := d.DivideByTime(t)
if err != nil { … } // ErrZeroDivision when t is zero
```

`Distance.DivideByTime` used to panic for a zero time; now it returns `ErrZeroDivision`.

## Generic code

Every value still implements `units.Measurement`, and `ParseString` still returns one. Two things changed:

- **`Unit()` is `Symbol()` on a `Measurement`**, because the typed `Unit()` returns the typed unit. On a `Distance`, `Unit()` returns a `distance.Unit` and `Symbol()` the string.
- **The concrete type is kept.** `ParseString("10 km")` returns a `Measurement` that holds a `units.Distance`, so a type assertion or type switch works. `ToDistance` still works too:

```go
m, err := units.ParseString("10 km")

switch v := m.(type) {
case units.Distance:
	fmt.Println(v.ConvertTo(distance.Mile))
case units.Mass:
	…
}
```

## Custom measures

`NewMeasurementCreator`, `MeasurementCreator` and `NewBaseMeasurement` are removed. A custom measure is now a marker type with a `UnitSet`, and its values are `Quantity[M]`, with the same typed API as the built-in measures:

```go
// Before
type Currency interface{ units.Measurement }

mc, err := units.NewMeasurementCreator[Currency]("currency", "USD", conversion.Factors{"USD": 1, "EUR": 1.1})
usd, err := mc.New(1, "USD")
eur, err := usd.ConvertToMeasurement("EUR")

// After
type currency struct{}

var currencyUnits = units.MustNewUnitSet("currency", "USD", conversion.Factors{"USD": 1, "EUR": 1.1})

func (currency) Units() units.UnitSet { return currencyUnits }

type Money = units.Quantity[currency]

var (
	usd = units.MustUnitOf[currency]("USD")
	eur = units.MustUnitOf[currency]("EUR")
)

price := units.NewQuantity(1, usd)
inEUR := price.ConvertTo(eur)                         // a Money, no error
parsed, err := units.ParseQuantity[currency]("2 EUR") // replaces mc.ParseString
```

- `NewUnitSet(measure, baseUnit, factors)` takes the same arguments as `NewMeasurementCreator` did, without the type parameter. Build the set once, in a package-level variable.
- One-off measures from `NewBaseMeasurement` become a `UnitSet` too.
- Values of two unit sets combine only if they are the same set. Two sets with the same name used to combine when their symbols matched; now `AddMeasurement` fails with `ErrIncompatibleMeasures`.
- A `Parser` learns a custom measure with `units.WithMeasure[currency]()`; parsed values are `Money`.

## Unit packages

- **`ConversionRates()`** is removed. For the multiplier between two units use `Rate(from, to)`, which takes typed units: `distance.Rate(distance.Mile, distance.Kilometer)` is 1.609344. For all factors use `ConversionFactors()`, a copy keyed by symbol.
- **New functions:** `ParseUnit(symbol)`, `Units()`, and the methods `String`, `Factor` and `Valid` on each unit.

## Behaviour that compiles unchanged

- **`mass.Quintal`** is `q`, not `Qts`. Stored strings such as `"3 Qts"` no longer parse.
- **`mass.Microgram`** is `μg` with U+03BC GREEK SMALL LETTER MU, not U+00B5 MICRO SIGN. Both forms are still accepted as input.
- **`ParseString` rejects** NaN and infinite values.
- **New measures** take over symbols: `L`, `N`, `J`, `W`, `Pa`, `Hz`, `K`, `°C` and others now belong to predefined measures, and their names (`volume`, `force`, `temperature`, …) can't be used for a custom measure.
- **Error messages** changed, for example to `unknown distance unit kg`. Test for the sentinel errors instead: `errors.Is(err, units.ErrUnknownUnit)`.
- **`%#v`** prints `units.Distance{Value:5, Unit:"km"}`.

## What you get

- Typed operations without allocations, in 2–5 ns instead of about 30–90 ns.
- Unit mix-ups as compile errors.
- New measures: area, volume, acceleration, force, energy, power, pressure, frequency, angle, data size, density, data rate and temperature.
- JSON and text encoding (`"12.5 km"`), `Compare` and `Equal`, a configurable `Parser`, `Time.ToDuration`, sentinel errors, and `fmt` verbs such as `%.2f`.

The [CHANGELOG](CHANGELOG.md) lists every change.
