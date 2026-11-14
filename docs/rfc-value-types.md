# RFC: value types and typed units

ROADMAP step 2.1. **Status:** accepted on 2026-11-14 and implemented (ROADMAP 3.1–3.3A). Written 2026-11-13.

> **Decisions (2026-11-14):** Path A with option I (integer units); `Div` and the cross-measure divisions keep returning `ErrZeroDivision`; `NewMeasurementCreator` and `NewBaseMeasurement` are **removed** in the same release rather than deprecated.
>
> **As built**, some names differ from this proposal: the unit type of a custom measure is `units.Unit[M]` (here `U[M]`), the constructor is `NewQuantity` (here `New`), `NewUnitSet` takes a base unit like `NewMeasurementCreator` did, `Measure.Units` returns a `UnitSet` value handle rather than a `*UnitSet`, and the generic interface's `Unit()` is `Symbol()`. The prototype module, `docs/rfc-value-types/`, was deleted. The measurements below were taken with the prototype; the library matches them (`BenchmarkMeasures`, `BenchmarkSum`).

## Summary

Replace the interface-based measurements with **value types** and **typed integer units**, and give user-defined measures the same typed API through a **generic `Quantity[M]`**.

The prototype in [`docs/rfc-value-types/`](rfc-value-types/) does the hot operations in about 2 ns with no allocations, against 73–89 ns and one allocation today: 18 to 45 times faster. A loop that sums 1000 distances drops from 60 µs and 1001 allocations to 3.3 µs and none.

The cost is a breaking release. Constructors and `ConvertTo` lose their error results, and units can no longer be plain strings. Most call sites change mechanically, as the migration table below shows.

The recommendation is **Path A** of the roadmap, with option I for units. The decisions for the author are listed at the end.

## Why change

Every remaining large cost comes from the interface-based design:

- **Allocations (P1).** Every value is a heap-allocated interface. Every operation that returns a measurement allocates once and takes 65–200 ns.
- **Untyped units.** Units are strings, so `NewDistance(1, mass.Kilogram)` compiles and fails only at run time. As a result `NewX` and `ConvertTo` return errors that a type system could rule out.
- **Interface results.** Constructors return nil interfaces on error, and the generic operations lose the type: `ParseString("1 km")` returns a value that doesn't satisfy `units.Distance`.

## Measurements

The prototype measured three designs against the current library (`cd docs/rfc-value-types && go test -bench . -benchmem`). It ran on Go 1.27.0 on an Intel i5-6300HQ; the numbers are medians of 3 runs.

| Operation | current | S: string units | **I: integer units** | `Quantity[M]` (custom) |
|---|---|---|---|---|
| New | 88 ns, 1 alloc | 17 ns | **2.0 ns** | 2.0 ns |
| ConvertTo | 89 ns, 1 alloc | 31 ns | **2.0 ns** | 4.7 ns |
| Add (different units) | 73 ns, 1 alloc | 38 ns | **1.8 ns** | 5.3 ns |
| Compare (different units) | 39 ns | 40 ns | **4.4 ns** | 6.7 ns |
| Sum of 1000 mixed distances | 60 µs, 1001 allocs | — | **3.3 µs, 0 allocs** | — |
| Boxing a value into an interface | — | — | 26 ns, 1 alloc | — |

Option S looks units up in a map, so it still costs 17–40 ns. Option I reads two array elements. `Quantity[M]` reaches its unit set through the marker type's method; generic code calls that method through a dictionary, which adds about 3 ns.

## The five questions

### 1. Values

```go
type Distance struct {
	value float64
	unit  distance.Unit
}
```

- **Receivers.** Every method has a value receiver and returns a new value, so values are immutable (P2) and safe to copy and share.
- **Zero value.** `distance.Meter` is the zero `Unit`, so the zero `Distance` is a valid 0 m, not a nil interface.
- **Errors.** Constructors that can still fail, such as the parsers, return the zero value with the error.
- **Generic code.** `Distance` still implements `Measurement`, so code written against `Measurement` keeps working. Converting a value to the interface allocates once, about 26 ns, and only on that generic path. Values from `ParseString` can then be the concrete measure type, so `m.(units.Distance)` works. That removes the "XxxMeasurement methods lose the type" debt.

### 2. Units

Each unit package declares its own unit type. Three options:

| | S: `type Unit string` | **I: `type Unit uint8`, array index** | Unit values carrying a factor |
|---|---|---|---|
| Wrong measure caught at compile time | yes | yes | yes |
| Untyped literal `NewDistance(1, "km")` | compiles | **doesn't compile** | doesn't compile |
| Invalid unit possible | any string, checked at run time | only `distance.Unit(200)`, which panics | no |
| Conversion cost | 2 map lookups, 17–40 ns | 2 array loads, about 2 ns | 2 field loads |
| Errors from `NewX`/`ConvertTo` | still needed | **gone** | gone |
| P2 | map written only at init | arrays unexported, never written | predefined units would be **exported variables**: rejected |

**Recommendation: option I.** It alone meets the 5 ns criterion. In each unit package:

```go
type Unit uint8

const (
	Meter Unit = iota // the base unit is the zero value
	Kilometer
	// …
)

var (
	factors = [...]float64{1, 1000, /* … */} // unexported: no importer can change it
	symbols = [...]string{"m", "km", /* … */}
)

func (u Unit) String() string          // "km"
func (u Unit) Factor() float64         // metres in one u, inlined by the compiler
func ParseUnit(symbol string) (Unit, bool) // accepts both micro signs
func Units() iter.Seq[Unit]
func ConversionFactors() conversion.Factors // unchanged: a copy, keyed by symbol
```

What option I costs:

- **The unit set is closed.** A program can't add a unit to `distance` itself. That is already true today: extending a predefined measure means copying `ConversionFactors()` into a custom measure, and that keeps working.
- **Strings need parsing.** A unit stored as a string, for example in configuration, must go through `distance.ParseUnit` first. That is also the one place left for a wrong symbol to be reported.
- **Order is part of the ABI.** Adding a unit appends a constant, and code must never persist the integers. `String()` and `ParseUnit` give stable symbols for storage, and text and JSON encoding (step 4.1) would use them.

### 3. Custom measures (P3)

A generic value type, defined by a marker type that the user declares:

```go
type currency struct{}

var currencyUnits = units.MustNewUnitSet("currency", conversion.Factors{"USD": 1, "EUR": 1.1})

func (currency) Units() *units.UnitSet { return currencyUnits }

type Money = units.Quantity[currency]

usd := units.MustUnitOf[currency]("USD")
total := units.New(10, usd).Add(units.New(5, usd)) // a Money; Add, ConvertTo and Compare are typed
```

- **Typed API.** `Quantity[M]` has the full typed API: `ConvertTo`, `Add`, `Sub`, `Mul`, `Div`, `Compare`, `Equal`, `String` and `Format`. A custom measure needs these 5 lines, against a criterion of about 10. A runnable version is `Example_customMeasure` in the prototype.
- **Immutability (P2).** A `UnitSet` validates and copies the factors once, and its fields are unexported, so a `*UnitSet` is a read-only handle. There is no global registry. The user's `currencyUnits` variable lives in their own package and is theirs to keep unexported.
- **Units.** A unit of `M` is a `U[M]`, an index into M's set. Symbols are looked up once, with `UnitOf[M]`, which can fail. Arithmetic on the resulting values can't.
- **Go 1.27 generic methods.** Methods with their own type parameters can't implement interface methods, so they don't help `Quantity[M]` satisfy `Measurement`. Its plain methods already do. Generic methods could later express cross-measure operations between custom measures, such as `q.Mul[N, P]`, but that isn't needed for 2.1.
- **Existing entry points.** `NewMeasurementCreator` and `NewBaseMeasurement` become redundant. Deprecate them in favour of `Quantity`, with a `//go:fix inline` hint where a mechanical rewrite exists, and remove them at v1.0.0.

### 4. Errors

| Operation | today | with option I |
|---|---|---|
| `NewX(v, unit)` | `(X, error)` | `X` |
| `x.ConvertTo(unit)` | `(X, error)` | `X`; `MustConvertTo` is deprecated as redundant |
| `Add`, `Sub`, `Mul`, `Compare` | no error | no error |
| `Div(0)`, `DivideByTime(zero)` | `ErrZeroDivision` | unchanged, see the decisions below |
| `ParseString`, `StringToX`, `ParseUnit`, `UnitOf` | error | error |
| `AddMeasurement` (generic) | error for different measures | unchanged |

### 5. Migration

| Before (v0.0.4) | After (v0.1.0) | How |
|---|---|---|
| `d, err := units.NewDistance(5, distance.Kilometer)` | `d := units.NewDistance(5, distance.Kilometer)` | drop the error; the compiler finds every case |
| `units.NewDistance(5, "km")` | `units.NewDistance(5, distance.Kilometer)` | use the constant |
| `units.NewDistance(v, cfg.Unit)` (a string) | `u, ok := distance.ParseUnit(cfg.Unit)`, then `NewDistance(v, u)` | parse once |
| `km, err := d.ConvertTo(distance.Kilometer)` | `km := d.ConvertTo(distance.Kilometer)` | drop the error |
| `if d == nil` | gone: the zero `Distance` is 0 m | remove the check |
| `d.Unit() == "km"` | `d.Unit() == distance.Kilometer`, or `d.Unit().String() == "km"` | compare typed units |
| `NewMeasurementCreator[Currency](…)` | marker type plus `Quantity[M]` (above) | about 5 lines per measure |
| `func f(d units.Distance)` | unchanged | `Distance` is still the type's name |

Method names don't change: `ConvertTo`, `Add`, `Compare`, `DivideByTime` and the rest. So most code changes are mechanical, and the compiler reports every site.

## Acceptance criteria (from the roadmap)

| Criterion | Result |
|---|---|
| Hot paths at most 5 ns, no allocations | **met** for predefined measures (1.8–4.4 ns). `Quantity[M]` takes 4.7–6.7 ns: `Compare` is 1.7 ns over |
| No exported variables | **met**: units are constants, and tables are unexported arrays or unexported fields |
| A custom measure gets the typed API in about 10 lines | **met**: 5 lines |
| README examples need only mechanical changes | **met**: see the migration table |

## Plan if accepted

Phase 3 of the roadmap, adapted:

1. **3.1 (adapted)** Each unit package gets `type Unit uint8`, constants, the factor and symbol arrays, `String`, `Factor`, `ParseUnit` and `Units`, generated from one list per package so the arrays can't drift apart.
2. **3.2A** Value-type measures. The library now has fourteen, all with the same shape (the nine added in step 4.3 were generated from the area pattern), so generate them from one template plus hand-written cross-measure operations, rather than converting each by hand. Then lower the allocation guards to 0.
3. **3.3A** `UnitSet`, `U[M]` and `Quantity[M]`, plus the deprecation of `NewMeasurementCreator` and `NewBaseMeasurement`.
4. **3.4** `MIGRATION.md`, built from the table above, and the release.

The measures added in steps 4.3 (volume through data rate) convert the same way as the original five.

## Decisions for the author

1. **Path A (this RFC) or Path B** (keep the interfaces and only fix the typed results of the generic operations)?
2. **Units: option I or option S?** S keeps `NewDistance(5, "km")` compiling and is still 2–5 times faster than today, but `NewX` and `ConvertTo` would keep their error results.
3. **Division by zero:** keep `Div` returning `(X, error)`, or follow IEEE floats and return ±Inf or NaN without an error? This RFC keeps the error, as decided in 0.3; with value types it no longer costs an allocation.
4. **`NewMeasurementCreator` and `NewBaseMeasurement`:** deprecate them in v0.2.0 and remove them at v1.0.0 (proposed), or remove them in v0.2.0?

The prototype module was deleted after the decision.
