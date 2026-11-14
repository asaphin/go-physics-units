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
- **Type-safe.** `units.NewDistance(5, mass.Kilogram)` doesn't compile. The typed operations can't fail, except division by zero and the few calculations without a real result.
- **Immutable and safe for concurrent use.** Values never change, and the library's tables are fixed once the package is initialised. No importer can alter them, and no locks are needed.
- **Extensible.** Your own measures, such as currencies, get the same typed API as the built-in ones, and can extend them.
- **Broad.** 43 measures with 356 units, from distance and time to torque, molarity and radiation dose. More than 100 operations between measures (force is mass × acceleration, power is torque × angular velocity, …) and calculators (the equivalent cube of a volume, rest energy, the Lorentz factor, …), a `physics` package of formulas, and the CODATA physical constants. Plus `fmt` verbs, JSON and text encoding, a configurable parser, and temperatures with their offsets.
- **Dependency-free and well tested.** Only the standard library, 100% test coverage, fuzzed parsers, and allocation counts checked in CI.

> [!NOTE]
> The API can still change before v1.0.0; see [Versioning](#versioning).
> Upgrading from v0.0.4? [MIGRATION.md](MIGRATION.md) lists every change, with before-and-after code.

## Contents

- [Installation](#installation)
- [Quick start](#quick-start), and [use cases](USECASES.md) from real-world domains
- [Usage](#usage): [values and units](#values-and-units), [comparing](#comparing-values), [combining measures](#combining-measures), [calculators](#calculators), [physics formulas](#physics-formulas), [constants](#constants), [formatting](#formatting), [parsing](#parsing), [JSON](#json-and-text-encoding), [temperature](#temperature), [custom measures](#custom-measures), [errors](#errors)
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

The `time` subpackage shadows the standard library's `time`. Import the standard one under another name, such as `builtintime "time"` as the examples do, when you need both.

For complete, real-world scenarios, from a fitness tracker and IoT telemetry to electronics and aviation fuel, see [USECASES.md](USECASES.md).

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

Methods combine measures and return their result in the SI base unit of the result (m/s, N, J, W, …):

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

| Value             | Methods and their results |
|-------------------|--------------------------|
| `Distance`        | `DivideByTime(Time)` → `Velocity`, `Square()` → `Area`, `Cube()` → `Volume`, `DivideByVelocity(Velocity)` → `Time`, `MultiplyByPace(Pace)` → `Time`, `DivideByVolume(Volume)` → `FuelEconomy`, `DivideByFuelEconomy(FuelEconomy)` → `Volume`, `MultiplyByFuelConsumption(FuelConsumption)` → `Volume`, `CircleArea()` → `Area`, `SphereArea()` → `Area`, `SphereVolume()` → `Volume` |
| `Area`            | `DivideByDistance(Distance)` → `Distance`, `MultiplyByDistance(Distance)` → `Volume`, `SquareSide()` → `Distance`, `CircleRadius()` → `Distance` |
| `Volume`          | `DivideByArea(Area)` → `Distance`, `DivideByDistance(Distance)` → `Area`, `DivideByTime(Time)` → `FlowRate`, `DivideByFlowRate(FlowRate)` → `Time`, `DivideByFuelConsumption(FuelConsumption)` → `Distance`, `CubeEdge()` → `Distance`, `SphereRadius()` → `Distance` |
| `Angle`           | `DivideByTime(Time)` → `AngularVelocity` |
| `Time`            | `ToDuration()` → `time.Duration`, `DivideByDistance(Distance)` → `Pace`, `Frequency()` → `Frequency` |
| `Frequency`       | `Period()` → `Time`, `MultiplyByDistance(Distance)` → `Velocity`, `PhotonEnergy()` → `Energy` |
| `Velocity`        | `MultiplyByTime(Time)` → `Distance`, `DivideByTime(Time)` → `Acceleration`, `MultiplyByMass(Mass)` → `Momentum`, `MultiplyByArea(Area)` → `FlowRate`, `DivideByFrequency(Frequency)` → `Distance`, `DivideByDistance(Distance)` → `Frequency`, `Pace()` → `Pace`, `LorentzFactor()` → `float64` |
| `Acceleration`    | `MultiplyByTime(Time)` → `Velocity`, `MultiplyByMass(Mass)` → `Force` |
| `AngularVelocity` | `MultiplyByTime(Time)` → `Angle`, `MultiplyByTorque(Torque)` → `Power` |
| `Pace`            | `MultiplyByDistance(Distance)` → `Time`, `Velocity()` → `Velocity` |
| `Mass`            | `MultiplyByAcceleration(Acceleration)` → `Force`, `DivideByVolume(Volume)` → `Density`, `DivideByDensity(Density)` → `Volume`, `MultiplyByVelocity(Velocity)` → `Momentum`, `DivideByTime(Time)` → `MassFlowRate`, `DivideByMassFlowRate(MassFlowRate)` → `Time`, `DivideByMolarMass(MolarMass)` → `Amount`, `DivideByAmount(Amount)` → `MolarMass`, `MultiplyBySpecificEnergy(SpecificEnergy)` → `Energy`, `KineticEnergy(Velocity)` → `Energy`, `Weight()` → `Force`, `RestEnergy()` → `Energy` |
| `Density`         | `MultiplyByVolume(Volume)` → `Mass`, `MultiplyByFlowRate(FlowRate)` → `MassFlowRate` |
| `Force`           | `DivideByMass(Mass)` → `Acceleration`, `DivideByAcceleration(Acceleration)` → `Mass`, `DivideByArea(Area)` → `Pressure`, `MultiplyByDistance(Distance)` → `Energy`, `MultiplyByTime(Time)` → `Momentum`, `TorqueAt(Distance)` → `Torque` |
| `Pressure`        | `MultiplyByArea(Area)` → `Force`, `MultiplyByFlowRate(FlowRate)` → `Power` |
| `Torque`          | `DivideByDistance(Distance)` → `Force`, `MultiplyByAngularVelocity(AngularVelocity)` → `Power`, `MultiplyByAngle(Angle)` → `Energy` |
| `Momentum`        | `DivideByMass(Mass)` → `Velocity`, `DivideByVelocity(Velocity)` → `Mass`, `DivideByTime(Time)` → `Force` |
| `Energy`          | `DivideByTime(Time)` → `Power`, `DivideByDistance(Distance)` → `Force`, `DivideByCharge(Charge)` → `Voltage`, `DivideByPower(Power)` → `Time`, `DivideByMass(Mass)` → `SpecificEnergy`, `DivideBySpecificEnergy(SpecificEnergy)` → `Mass`, `DivideByVolume(Volume)` → `EnergyDensity`, `DivideByEnergyDensity(EnergyDensity)` → `Volume`, `MassEquivalent()` → `Mass`, `PhotonFrequency()` → `Frequency` |
| `Power`           | `MultiplyByTime(Time)` → `Energy`, `DivideByVoltage(Voltage)` → `Current`, `DivideByCurrent(Current)` → `Voltage`, `DivideByAngularVelocity(AngularVelocity)` → `Torque`, `DivideByTorque(Torque)` → `AngularVelocity` |
| `SpecificEnergy`  | `MultiplyByMass(Mass)` → `Energy` |
| `EnergyDensity`   | `MultiplyByVolume(Volume)` → `Energy` |
| `FlowRate`        | `MultiplyByTime(Time)` → `Volume`, `DivideByArea(Area)` → `Velocity`, `MultiplyByPressure(Pressure)` → `Power`, `MultiplyByDensity(Density)` → `MassFlowRate` |
| `MassFlowRate`    | `DivideByDensity(Density)` → `FlowRate`, `MultiplyByTime(Time)` → `Mass` |
| `FuelEconomy`     | `MultiplyByVolume(Volume)` → `Distance`, `Consumption()` → `FuelConsumption` |
| `FuelConsumption` | `MultiplyByDistance(Distance)` → `Volume`, `Economy()` → `FuelEconomy` |
| `Current`         | `MultiplyByVoltage(Voltage)` → `Power`, `MultiplyByResistance(Resistance)` → `Voltage`, `MultiplyByTime(Time)` → `Charge` |
| `Voltage`         | `MultiplyByCurrent(Current)` → `Power`, `DivideByCurrent(Current)` → `Resistance`, `DivideByResistance(Resistance)` → `Current`, `MultiplyByCharge(Charge)` → `Energy`, `MultiplyByCapacitance(Capacitance)` → `Charge` |
| `Charge`          | `DivideByTime(Time)` → `Current`, `DivideByCurrent(Current)` → `Time`, `MultiplyByVoltage(Voltage)` → `Energy`, `DivideByVoltage(Voltage)` → `Capacitance`, `DivideByCapacitance(Capacitance)` → `Voltage` |
| `Resistance`      | `MultiplyByCurrent(Current)` → `Voltage`, `MultiplyByCapacitance(Capacitance)` → `Time` |
| `Capacitance`     | `MultiplyByVoltage(Voltage)` → `Charge`, `MultiplyByResistance(Resistance)` → `Time` |
| `Inductance`      | `DivideByResistance(Resistance)` → `Time` |
| `Illuminance`     | `MultiplyByArea(Area)` → `LuminousFlux` |
| `LuminousFlux`    | `DivideByArea(Area)` → `Illuminance` |
| `Amount`          | `MultiplyByMolarMass(MolarMass)` → `Mass`, `DivideByVolume(Volume)` → `Molarity`, `DivideByMolarity(Molarity)` → `Volume`, `Entities()` → `float64` |
| `MolarMass`       | `MultiplyByAmount(Amount)` → `Mass` |
| `Molarity`        | `MultiplyByVolume(Volume)` → `Amount` |
| `Radioactivity`   | `DecayAfter(Time, Time)` → `Radioactivity` |
| `AbsorbedDose`    | `Equivalent(float64)` → `EquivalentDose` |
| `EquivalentDose`  | `DivideByTime(Time)` → `DoseRate`, `DivideByDoseRate(DoseRate)` → `Time` |
| `DoseRate`        | `MultiplyByTime(Time)` → `EquivalentDose` |
| `DataSize`        | `DivideByTime(Time)` → `DataRate`, `DivideByDataRate(DataRate)` → `Time` |
| `DataRate`        | `MultiplyByTime(Time)` → `DataSize` |
| `Ratio`           | `Fraction()` → `float64` |

Divisions and reciprocals also return an error: `ErrZeroDivision` for a zero divisor, as `Div(0)` does. `SquareSide`, `CircleRadius` and `LorentzFactor` fail with `ErrOutOfRange` when the result isn't real.
`ParseDuration` turns a `time.Duration` into a `Time`, and the `Ratio` method of every measure gives the ratio of two values, such as a share in percent.

### Calculators

Some methods are formulas of a single value: the equivalent cube or sphere of a volume, the kinetic, rest and photon energies, the Lorentz factor, radioactive decay:

```go
// All the gold ever mined, about 216 000 t, would fit in one cube.
gold, err := units.NewMass(216000, mass.Tonne).
	DivideByDensity(units.NewDensity(19.32, density.GramPerCubicCentimeter))
if err != nil {
	panic(err) // the density is zero
}

fmt.Printf("%.1f\n", gold.CubeEdge()) // 22.4 m

// The energy equivalent of one gram of matter, as TNT.
gram := units.NewMass(1, mass.Gram)

fmt.Printf("%.1f\n", gram.RestEnergy().ConvertTo(energy.KilotonOfTNT)) // 21.5 ktTNT
```

### Physics formulas

The `physics` package has the formulas with several inputs, an approximation or a measured constant. They take and return the same typed values:

| Topic                     | Functions                                                                                                          |
|---------------------------|--------------------------------------------------------------------------------------------------------------------|
| Falling bodies            | `FreeFallTime`, `FreeFallSpeed`, `FallDistance`, `PotentialEnergy`, `ProjectileRange`, `PendulumPeriod`, `HydrostaticPressure` |
| Gravitation and orbits    | `GravitationalForce`, `SurfaceGravity`, `EscapeVelocity`, `OrbitalVelocity`, `OrbitalPeriod`, `SchwarzschildRadius` |
| Light, sound, relativity  | `LightWavelength`, `LightFrequency`, `PhotonEnergy`, `SpeedOfSoundInAir`, `MachNumber`, `TimeDilation`             |
| Thermodynamics            | `IdealGasPressure`, `IdealGasVolume`, `IdealGasAmount`, `BlackBodyPower`, `WienPeakWavelength`                     |
| Circuits                  | `CapacitorEnergy`, `InductorEnergy`, `ResonantFrequency`                                                           |

```go
earth := units.NewMass(1, mass.EarthMass)

v, err := physics.EscapeVelocity(earth, units.NewDistance(6371, distance.Kilometer))
if err != nil {
	panic(err)
}

fmt.Printf("%.1f\n", v.ConvertTo(velocity.KilometerPerSecond)) // 11.2 km/s

// A geostationary orbit goes round once a sidereal day.
t, err := physics.OrbitalPeriod(earth, units.NewDistance(42164, distance.Kilometer))
if err != nil {
	panic(err)
}

fmt.Printf("%.2f\n", t.ConvertTo(time.Hour)) // 23.93 h
```

They fail with `ErrZeroDivision` for a zero divisor and with `ErrOutOfRange` when the result isn't real, such as the time to fall a height against gravity.

### Constants

The `constant` package has physical constants in SI units, as untyped Go constants:

- the defining constants of the SI, which are exact: c, h, e, k and N_A;
- constants derived from them: ħ, R, F and σ;
- the CODATA 2022 values of measured constants: G, mₑ, mₚ, u, α, ε₀ and μ₀;
- astronomical constants of the IAU: the astronomical unit, the light-year, the parsec and nominal solar and terrestrial values.

```go
c := units.NewVelocity(constant.SpeedOfLight, velocity.MeterPerSecond)

fmt.Println(c.ConvertTo(velocity.KilometerPerSecond)) // 299792.458 km/s
```

They are constants, not variables, so no importer can change them.

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
| `ErrOutOfRange`           | a temperature below absolute zero, `Time.ToDuration` beyond ±292 years, marshalling NaN, or a result that isn't real, such as `Area.SquareSide` of a negative area |

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

| Measure               | Package           | Base unit | Symbols                                                                                                             |
|-----------------------|-------------------|-----------|---------------------------------------------------------------------------------------------------------------------|
| Distance              | `distance`        | m         | lp, fm, pm, Å, nm, μm, mil, mm, cm, in, dm, ft, yd, m, ftm, km, mi, nmi, AU, ly, pc, kpc, Mpc, Gpc                  |
| Area                  | `area`            | m²        | mm², cm², in², ft², yd², m², ac, ha, km², mi²                                                                       |
| Volume                | `volume`          | m³        | mm³, μL, cm³, mL, cL, in³, dL, pt, qt, L, gal, ft³, bbl, yd³, m³, km³                                               |
| Angle                 | `angle`           | rad       | arcsec, arcmin, mrad, gon, °, rad, rev                                                                              |
| Time                  | `time`            | s         | tp, ns, μs, ms, s, min, h, d, wk, yr                                                                                |
| Frequency             | `frequency`       | Hz        | rpm, bpm, Hz, kHz, MHz, GHz, THz                                                                                    |
| Velocity              | `velocity`        | m/s       | km/h, ft/s, mph, kn, m/s, km/s, c                                                                                   |
| Acceleration          | `acceleration`    | m/s²      | cm/s², ft/s², m/s², gₙ                                                                                              |
| Angular velocity      | `angularvelocity` | rad/s     | °/s, rev/min, rad/s, rev/s                                                                                          |
| Pace                  | `pace`            | s/m       | s/mi, s/km, s/500m, s/100m, min/mi, min/km, s/m                                                                     |
| Mass                  | `mass`            | kg        | Da, μg, mg, gr, ct, g, oz, ozt, lb, kg, st, q, t, M⊕, M☉                                                            |
| Density               | `density`         | kg/m³     | kg/m³, g/L, lb/ft³, lb/gal, kg/L, g/mL, g/cm³                                                                       |
| Force                 | `force`           | N         | dyn, mN, ozf, N, lbf, kgf, kN, kip, tf, MN                                                                          |
| Pressure              | `pressure`        | Pa        | Pa, mmH2O, cmH2O, hPa, mbar, Torr, mmHg, inH2O, kPa, inHg, psi, kgf/cm², bar, atm, MPa, ksi, GPa                    |
| Torque                | `torque`          | N·m       | mN·m, ozf·in, N·cm, kgf·cm, lbf·in, N·m, lbf·ft, kgf·m, kN·m                                                        |
| Momentum              | `momentum`        | kg·m/s    | g·cm/s, lb·ft/s, kg·m/s, N·s, lbf·s                                                                                 |
| Energy                | `energy`          | J         | eV, keV, MeV, GeV, erg, TeV, mJ, J, ft·lbf, cal, kJ, Btu, Wh, kcal, gTNT, MJ, kWh, GJ, MWh, tTNT, toe, ktTNT, MtTNT |
| Power                 | `power`           | W         | mW, Btu/h, W, PS, hp, kW, TR, MW, GW, L☉                                                                            |
| Specific energy       | `specificenergy`  | J/kg      | J/kg, kJ/kg, Btu/lb, Wh/kg, kcal/kg, MJ/kg, kWh/kg                                                                  |
| Energy density        | `energydensity`   | J/m³      | J/m³, kJ/m³, MJ/m³, Wh/L, kWh/m³, MJ/L, kWh/L                                                                       |
| Flow rate             | `flowrate`        | m³/s      | mL/h, mL/min, L/h, L/min, gal/min, m³/h, ft³/min, L/s, ft³/s, m³/s                                                  |
| Mass flow rate        | `massflowrate`    | kg/s      | lb/h, kg/h, g/s, lb/min, kg/min, t/h, lb/s, kg/s                                                                    |
| Fuel economy          | `fueleconomy`     | m/m³      | m/m³, mpg(imp), mpg, km/L                                                                                           |
| Fuel consumption      | `fuelconsumption` | m³/m      | L/100km, gal/100mi, L/km, m³/m                                                                                      |
| Temperature           | `temperature`     | K         | K, °C, °F, °R                                                                                                       |
| Electric current      | `current`         | A         | μA, mA, A, kA                                                                                                       |
| Voltage               | `voltage`         | V         | μV, mV, V, kV, MV                                                                                                   |
| Electric charge       | `charge`          | C         | μC, mC, C, mAh, Ah                                                                                                  |
| Resistance            | `resistance`      | Ω         | mΩ, Ω, kΩ, MΩ, GΩ                                                                                                   |
| Capacitance           | `capacitance`     | F         | pF, nF, μF, mF, F                                                                                                   |
| Inductance            | `inductance`      | H         | nH, μH, mH, H                                                                                                       |
| Magnetic flux density | `magneticfield`   | T         | nT, mG, μT, G, mT, T                                                                                                |
| Illuminance           | `illuminance`     | lx        | lx, fc, klx, ph                                                                                                     |
| Luminous flux         | `luminousflux`    | lm        | lm, klm                                                                                                             |
| Amount of substance   | `amount`          | mol       | pmol, nmol, μmol, mmol, mol, kmol                                                                                   |
| Molar mass            | `molarmass`       | kg/mol    | g/mol, kg/mol                                                                                                       |
| Molar concentration   | `molarity`        | mol/m³    | pM, nmol/L, nM, μmol/L, μM, mol/m³, mmol/L, mM, mol/L, M                                                            |
| Radioactivity         | `radioactivity`   | Bq        | Bq, nCi, kBq, μCi, MBq, mCi, GBq, Ci, TBq                                                                           |
| Absorbed dose         | `absorbeddose`    | Gy        | μGy, mGy, cGy, rd, Gy                                                                                               |
| Equivalent dose       | `equivalentdose`  | Sv        | nSv, μSv, mrem, mSv, rem, Sv                                                                                        |
| Dose rate             | `doserate`        | Sv/s      | nSv/h, mSv/yr, μSv/h, mrem/h, mSv/h, Sv/h, Sv/s                                                                     |
| Data size             | `datasize`        | B         | bit, B, kbit, kB, KiB, Mbit, MB, MiB, Gbit, GB, GiB, TB, TiB, PB, PiB, EB, EiB                                      |
| Data rate             | `datarate`        | B/s       | bit/s, B/s, kbit/s, kB/s, KiB/s, Mbit/s, MB/s, MiB/s, Gbit/s, GB/s, GiB/s                                           |
| Ratio                 | `ratio`           | 1         | ppt, ppb, ppm, bp, ‰, %, 1                                                                                          |

- Every unit is a typed constant in its package, for example `distance.LightYear` or `mass.Pound`; `String()` returns its symbol.
- Each unit package also provides `ParseUnit(symbol)`, `Units()`, `ConversionFactors()`, a copy of its factors, and `Rate(from, to)`, the multiplier between two units. For example, `distance.Rate(distance.Mile, distance.Kilometer)` returns 1.609344. The `temperature` package has `Convert` and `DifferenceRate` instead.
- Store and exchange units by their symbols, never as numbers: the numbers behind the constants may change between versions.
- Symbols are case-sensitive.
- The micro prefix can be written as `μ` (U+03BC GREEK SMALL LETTER MU) or `µ` (U+00B5 MICRO SIGN). Values always carry the `μ` form.
- `yr` is the Julian year of 365.25 days.
- `pt`, `qt` and `gal` are US liquid measures, `bbl` the 42-gallon oil barrel, `cal` the thermochemical calorie, `Btu` the International Table one, `hp` mechanical horsepower and `PS` metric horsepower. `mmHg` and `inHg` are the conventional values, and `gₙ` (with U+2099) is standard gravity, 9.80665 m/s².
- Mind the case: `kN` is a kilonewton and `kn` a knot.
- `kgf/cm²` is the kilogram-force per square centimetre, also called the technical atmosphere and often written kg/cm². A `Parser` accepts that spelling with `WithAliases(map[string]string{"kg/cm²": "kgf/cm²"})`.
- `℃` (U+2103) and `℉` (U+2109) are accepted for `°C` and `°F`, `Ω` (U+2126 OHM SIGN) for `Ω` (U+03A9), `Å` (U+212B ANGSTROM SIGN) for `Å` (U+00C5), and `K` (U+212A KELVIN SIGN) for `K`.
- Angular velocity's `rev/min` counts radians; frequency's `rpm` counts cycles. Convert between them through `Angle.DivideByTime` or the factor 2π.
- Data sizes have both decimal (`kB` = 1000 B) and binary (`KiB` = 1024 B) prefixes, and bits (`bit`, `Mbit`, …).
- Symbols are listed in order of size.
- `tTNT`, `ktTNT` and `MtTNT` are TNT equivalents, 4.184 GJ per ton by convention, and `toe` is the tonne of oil equivalent, 41.868 GJ by the IEA convention.
- `mpg` is miles per US gallon and `mpg(imp)` per imperial gallon. Fuel economy (`km/L`, `mpg`) and fuel consumption (`L/100km`) are reciprocals: `Consumption` and `Economy` convert between them.
- `M`, `mM`, `μM`, … are the chemists' molar units of mol/L. The rad of absorbed dose is `rd`, because `rad` is the radian.
- `mmH2O`, `cmH2O` and `inH2O` are conventional water columns, `bpm` and `rpm` both count per minute, and `L☉` is the IAU's nominal solar luminosity.
- `ppm`, `ppb` and `ppt` are on the short scale (1e-6, 1e-9 and 1e-12), `bp` is a basis point, 0.01 %, and `1` is the unit of a plain fraction.
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
[CHANGELOG.md](CHANGELOG.md) records every user-visible change, and [MIGRATION.md](MIGRATION.md) covers the upgrade from v0.0.4; v0.2.0 needs no changes after v0.1.0.
Deprecated names keep working until v1.0.0, and `go fix` rewrites your code to their replacements, such as `MustConvertTo` to `ConvertTo`.

The full API reference, with runnable examples for every measure, is on [pkg.go.dev](https://pkg.go.dev/github.com/asaphin/go-physics-units).

## Contributing

Issues and pull requests are welcome. Open pull requests against the `dev` branch; `main` holds the releases.

```shell
make tools    # once per clone: installs the pinned golangci-lint into ./bin
make check    # what CI runs: formatting, go vet, golangci-lint, tests with the race detector, allocation guards
make cover    # coverage summary, plus an HTML report in coverage/
make bench    # benchmarks
make fuzz     # fuzz the string parser (FUZZTIME=30s by default)
make fix      # apply every automatic fix: go fix, golangci-lint --fix, formatting
make generate # regenerate the code every measure shares, after changing internal/gen
make help     # every target
```

A change is done when `make check` is clean and coverage stays at 100%.
[AGENTS.md](AGENTS.md) describes the design rules, the layout of the code, and step-by-step recipes for adding a unit or a measure. It is written for human contributors and AI coding agents alike.

## License

go-physics-units is released under the [MIT License](LICENSE).
