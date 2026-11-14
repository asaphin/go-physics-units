# AGENTS.md

Guidance for AI coding agents (and humans) working in this repository. `CLAUDE.md` imports this file, so keep all shared rules here.

## What this is

`go-physics-units` (module `github.com/asaphin/go-physics-units`, package name `units`) is a small Go library for creating, converting, comparing, parsing, encoding and doing arithmetic on physical measurements (distance, area, volume, mass, time, velocity, acceleration, force, energy, power, pressure, frequency, angle, angular velocity, data size, density, data rate, temperature, electric current, voltage, electric charge, resistance) and user-defined measures (e.g. currency).

- **Published and pre-1.0** (tags `v0.0.0`–`v0.0.4`). The work since v0.0.4 redesigns the API (value types and typed units, ROADMAP Phase 3); `MIGRATION.md` is the upgrade guide. Exported API changes affect real users: make them deliberately, call them out, and record them in `CHANGELOG.md` and, if they break code, `MIGRATION.md`. `make api`, and the CI job "API changes", list every exported change since the latest tag; check that each one is in the CHANGELOG.
- **Go version:** `go.mod` declares `go 1.27.0`. CI tests that exact version and the latest stable release. Use modern Go freely: iterators and the `maps`/`slices` helpers, `strings.FieldsSeq`, `testing.B.Loop`, `sync.WaitGroup.Go`, `encoding.TextAppender`, and generic methods (Go 1.27). The `modernize` linter points out older forms. Raising the minimum again is its own deliberate change, with a CHANGELOG entry.
- **No third-party dependencies**, including in tests (stdlib `testing` only, no testify). Keep it that way unless asked.
- **Branches:** `dev` is the working branch, and it only accepts pull requests. Work on a short-lived branch and open a PR into `dev`. Releases are merged from `dev` into `main`.

## Design principles

The library rests on three principles. They take priority over convenience and over new features. Get the user's explicit agreement before making a change that weakens any of them, and say which one it weakens.

### 1. Performance

Measurements are created and combined in hot loops, so their cost is part of the API.

- **Calculations come first.** Conversions, arithmetic, comparisons and cross-measure operations are what users run in loops; keep them fast. Constructors of unit sets and parsers, and table building, may get somewhat slower when a feature needs it: say so in the PR, but it isn't a blocker (the author's call, 2026-11-13).
- **Do the work once.** All conversion data is built at package init, or once per `NewUnitSet` and `NewParser`. A conversion at run time is one array load and one multiplication.
- **Hot paths** are `NewX`, `ConvertTo`, `Add`/`Sub`/`Mul`/`Div`, `Compare`/`Equal`, the cross-measure operations, `String`, `StringToX`, `ParseString` and the `Quantity` methods. The typed ones must not allocate at all; `String` allocates its result, and the generic `Measurement` methods and `ParseString` allocate once, to box the result into the interface. Keep all of them free of reflection and `fmt`.
- **Error paths too.** Build error messages by string concatenation, never with `fmt`, which formats named string types such as `MeasureType` through reflection. Return a constant where the message is fixed (`ErrZeroDivision`, the `errXConversion` constants): that costs no allocation. `*UnitError` builds its message in `Error()`, so returning one costs a single allocation. `TestAllocations` guards the error paths too.
- **Inlining.** The unit packages' `Rate` and `Unit.Factor`, and the measures' `ConvertTo`, `Add`, `Sub` and `baseValue`, must stay small enough to inline. Check with `go build -gcflags=-m . ./distance`.
- **Allocation guards.** `TestAllocations` (`alloc_test.go`) fixes the allocation count of each hot path, and CI fails if one rises. Unlike timings, allocation counts are deterministic. The race detector changes them, so the file is excluded from `-race` builds, and `make allocs` and CI run it separately. Lower a limit when you save an allocation; raise one only on purpose, and say why in the PR. Box values into `Measurement` outside the measured function: the caller's boxing isn't the library's cost.
- **Measure.** For any change that touches a hot path, run `make bench-compare` (`BENCH_BASE=<ref>`, default the latest tag). It benchmarks both sides on the same Go toolchain and compares them with benchstat. Differences around ±10% are noise on this suite. Report every regression to the user; don't accept one silently.

### 2. Immutability after init (data security)

Shared state is built during package initialisation and never written afterwards. That makes the library safe for concurrent use without locks. It also means no importer can change how conversions behave for the rest of the program.

- **No exported package-level variables.** Any importer can reassign one, so exported names are constants or functions. Units are typed constants for this reason: units that carry their own factor would have to be exported variables. `TestNoExportedPackageVariables` enforces this.
- **Never hand out a reference to shared state.** Return copies, as `ConversionFactors()` does, or read-only handles.
  - The unit packages' tables (`definitions`, `rateTable`, `bySymbol`) are unexported arrays and maps, written only by their initialisers.
  - No exported signature may use a type from an internal package. `TestNoInternalTypesInExportedAPI` enforces this.
  - `UnitSet` and `Parser` are read-only handles passed by value: structs whose only field points to unexported data. Never hand out a `*UnitSet`, a `*unitSet` or a `*Parser`. A caller who holds a pointer can write `*p = other` and overwrite shared data.
- **Copy caller data on the way in.** `NewUnitSet` and `WithAliases` copy their maps, so later changes to the caller's map have no effect.
- **No writes after init.** That rules out lazy initialisation, caches and memoisation of shared state. `TestConcurrentUse` catches them under `-race`.
- **Immutable values.** Values are small structs with unexported fields and value receivers; every operation returns a new value. The only pointer-receiver methods are `UnmarshalText`, which encoding packages require, and they change only the caller's own variable. `Measurement` is sealed by an unexported method, so no outside type can pose as a measurement.
- **Exported errors are constants.** The sentinels are constants of the exported type `Error` (`errors.go`). Never add `var ErrX = errors.New(…)`: any importer could reassign it. Detail errors such as `*UnitError` are allocated on every call, so a caller who changes one changes only their own copy. Never return a shared pointer to an error.

### 3. Customization

User-defined measures are first-class: they convert, combine, compare, parse and encode like the predefined ones.

- **One entry point.** A user-defined measure is a marker type `M` whose `Units()` returns a `UnitSet` built once with `NewUnitSet`; its values are `Quantity[M]`, with the typed API of the predefined measures. `Parser` learns it with `WithMeasure[M]`.
- **Custom symbols are independent of predefined ones.** A custom `"m"` is not a distance.
- **Predefined measures can be extended.** Copy their factors with `ConversionFactors()`, add units, and pass the result to `NewUnitSet`. `TestNewUnitSet_ExtendsPredefinedFactors` covers this.
- **Every restriction needs a correctness reason; don't add one without it.** The current restrictions and their reasons:
  - **Factors** must be positive and finite, because rates divide by them.
  - **Symbols** must be non-empty and unique once the micro sign is canonicalised, or a symbol would denote two units.
  - **Measure names** can't reuse a predefined name (including `temperature` and `custom`). `Type()` identifies a measure, so a custom measure named `distance` would pass checks such as `m.Type() == units.MeasureDistance` without being a `Distance`.
  - **The predefined unit sets are closed.** Integer units index fixed arrays (ROADMAP 2.1, option I); extending a measure means a custom measure, as above.
  - **Temperatures** can't be below absolute zero (the author's decision, 2026-11-14).
- **New features** for the predefined measures should, where it makes sense, work for `Quantity` too.

## Layout

| Path | Role |
|---|---|
| `doc.go` | Package documentation, including the `fmt` verbs |
| `measurement.go` | The `Measurement` interface (sealed by the unexported `generic` method), `NewMeasurement`, `ParseString`, string parsing and formatting helpers |
| `measures.go` | `MeasureType` constants, `predefinedUnits` (every predefined symbol), `DetectMeasureType`, the `go:generate` directive |
| `measures_gen.go` | **Generated** (`internal/gen`): the shared API of every predefined measure (`NewX`, `ConvertTo`, `Add`, …, `ToX`, `StringToX`, text encoding) and their unit sets. Don't edit it by hand |
| `distance.go`, `area.go`, … (one per measure) | The measure's struct type with its documentation, and its cross-measure operations. Hand-written |
| `time.go` | `Time` plus `ToDuration` and `ParseDuration` |
| `temperature.go` | `Temperature` and `TemperatureDifference`, hand-written: affine conversions, outside `Measurement` |
| `unitset.go` | `UnitSet` (the exported read-only handle), `unitSet` (the units of one measure, for the generic code), `NewUnitSet` |
| `quantity.go` | `Measure`, `Unit[M]`, `UnitOf`, `Quantity[M]` and its methods |
| `common.go` | The generic interfaces `Arithmetics[T]`, `UnitConverter[T, U]` and `Comparer[T]` |
| `errors.go` | The sentinel errors (constants of type `Error`), `UnitError`, and the unexported error types and constructors the package uses |
| `parser.go` | `Parser`, `NewParser` and its options |
| `distance/`, `area/`, … (one per measure) | `units.go`: the `Unit` type, its constants, `BaseUnit` and the `definitions` array (symbol and factor of each unit). `unit_gen.go`, **generated**: `String`, `Factor`, `Valid`, `ParseUnit`, `Units`, `ConversionFactors`, `Rate` |
| `temperature/` | The temperature units, hand-written: scales with offsets, `Convert`, `DifferenceRate`, `AbsoluteZero` |
| `conversion/` | Public `Factors` type alias (`map[string]float64`) and `CopyConversionFactors` |
| `internal/rates` | `Validate` (unit definitions), `Matrix` (the rates between every pair of units), `CanonicalUnit` |
| `internal/gen` | The code generator, and `TestGenerated`, which fails when a generated file is stale |
| `Makefile` | Task runner (see Commands) |
| `.github/workflows/ci.yml` | CI: `go mod tidy` check, vet, race tests and coverage (the `go.mod` Go version and stable), allocation guards, short fuzz runs, golangci-lint (pinned version), the exported-API report, govulncheck |
| `.github/dependabot.yml`, `.github/pull_request_template.md` | Weekly bumps of the CI actions (into `dev`); the PR checklist |
| `docs/` | Design documents: `docs/rfc-value-types.md` (ROADMAP 2.1, accepted) and `docs/rfc-temperature.md` (ROADMAP 4.4, accepted) |
| `CHANGELOG.md`, `MIGRATION.md` | User-visible changes, under `## Unreleased` until the next tag; the upgrade guide from v0.0.4 |
| `ROADMAP.md` | Ordered plan of upcoming work. Take steps in order and tick them off when done |

## How it works

1. **Units.** Each unit package declares `type Unit uint8` and typed constants, the base unit first so that it is the zero value. The `definitions` array, indexed by `Unit`, holds each unit's symbol and factor: how many base units are in one unit (`Kilometer: {"km", 1000}`). Base units are coherent SI units (m, m², kg, s, m/s, …). The generated `build` validates the definitions with `rates.MustValidate` at package init, so a bad factor panics at startup, and computes `rateTable[from][to] = factor[from] / factor[to]` and the symbol index.
2. **Values.** A value is a struct `{value float64; unit X.Unit}` with value receivers. `ConvertTo(u)` is `value * X.Rate(from, u)`, one array load and one multiplication, inlined. An invalid unit, such as `distance.Unit(200)`, can only come from a conversion of an arbitrary integer; it panics on the first conversion, and `Valid` reports it.
3. **Measure identity.** Each predefined measure has a `unitSet` (`distanceSet`, …, in `measures_gen.go`) for the generic code: the symbols, factors and rates indexed by the same numbers as the constants, and a `box` function that makes a `Measurement` of the concrete type. `Type()` is the set's name. A `Quantity[M]` uses `M.Units()`'s set. Two values belong to the same measure when they share the set, not merely its name.
4. **Symbols.** `predefinedUnits` (`measures.go`) indexes every predefined symbol by its set and unit. `NewMeasurement`, `ParseString` and `DetectMeasureType` use it, and every `Parser` starts from a copy of it. Because of it, **predefined symbols must be unique across all measures, temperature included**; `TestMeasuresSelfCheck` enforces this.
5. **Look-alike characters.** `rates.CanonicalUnit` maps U+00B5 MICRO SIGN to U+03BC GREEK SMALL LETTER MU, U+2126 OHM SIGN to U+03A9 GREEK CAPITAL LETTER OMEGA, U+212B ANGSTROM SIGN to U+00C5 and U+212A KELVIN SIGN to K, as Unicode NFKC does. Every entry point that receives a symbol canonicalises it (`ParseUnit`, `UnitOf`, `NewMeasurement`, `ConvertToMeasurement`, `DetectMeasureType`, `Parser`), and definitions and sets store canonical symbols.
6. **Combining values.**
   - Typed `Add`/`Sub` need no check: two values of one Go type share their units. The result keeps the receiver's unit.
   - The generic `AddMeasurement`, `SubMeasurement` and `CompareMeasurement` check that the operand belongs to the same unit set (`unitSet.operand`).
   - Cross-measure operations (`Distance.DivideByTime`, `Force.DivideByArea`, `Frequency.Period` and the others in the README's operations table) convert to base units with `baseValue()` and return results in the target measure's base unit. Dividing by a zero value returns `ErrZeroDivision`, as `Div` does for a zero divisor.
   - **Comparisons** (`Compare`, `Equal`, `EqualWithinRel`, `CompareMeasurement`) compare two values in the same unit exactly. Otherwise they multiply each value by its unit's factor, which gives the value in the base unit. The result doesn't depend on the order of the operands, so `a.Equal(b) == b.Equal(a)`; converting one value into the other's unit would break that at the last bit. `TestCompare_Symmetric` enforces it.
7. **Custom measures.** `NewUnitSet` validates and copies the factors, puts the base unit at index 0 and sorts the other symbols, so the indexes don't depend on map iteration. `Quantity[M]` reaches the set through `M.Units()` on every operation; generic code calls that method through a dictionary, which costs about 3 ns.
8. **Parsing.** `ParseString` and `StringToX` accept exactly `"<float> <unit>"`, two whitespace-separated fields. They reject NaN and ±Inf.
   - `Parser` (`parser.go`) is the configurable variant, a value handle: `NewParser` builds an immutable index of symbols and aliases (plus a lower-case index for `CaseInsensitive`), and `Parse` is one or two map lookups. Aliases and custom units live in the parser, never in the global tables (P2).
   - **Its contract:** anything `ParseString` accepts, every `Parser` accepts with the same result; `FuzzParser` checks this, and the String round trip.
9. **Text encoding.** `AppendText`, `MarshalText` and `UnmarshalText` use the `String` form. Marshalling fails for values that aren't finite and for invalid units, so whatever is encoded parses back; `FuzzParseString` checks the round trip.
10. **Temperature.** Temperature scales have offsets, so they don't fit the factor model. The `temperature` package defines each scale by an exact rational offset and scale, and computes the coefficients `y = x·a + b` of every pair with `math/big` at init, rounding each once: 100 °C converts to exactly 212 °F. `Convert` rounds the product explicitly, so no platform fuses it with the addition. `Temperature` rejects values below absolute zero (checked in the value's own unit, so the defined zero points are exact), and `ConvertTo` clamps rounding at absolute zero. `TemperatureDifference` converts by `a` alone. Neither implements `Measurement`: its `AddMeasurement` and `MulMeasurement` would be meaningless for temperatures.

Constructor conventions per measure:

- `NewX(value, unit) X` can't fail. Temperature is the exception: `NewTemperature` returns an error below absolute zero.
- `ToX(Measurement) (X, error)` is a type assertion, with an `errXConversion` constant as the error.
- `StringToX(s)` parses the value and the symbol, then the unit with `X.ParseUnit`, without boxing.

On error, return the zero value of a concrete type (`return Distance{}, err`) or a nil interface (`return nil, err`), never a boxed zero value as a `Measurement`.

## Commands

```sh
make tools      # once per clone: the pinned golangci-lint into ./bin (gitignored), built with the go.mod Go
make check      # what CI runs: go mod tidy check, format check, go vet, golangci-lint, go test -race, allocation guards
make help       # every documented target
make generate   # regenerate measures_gen.go and the unit packages' unit_gen.go (go generate ./...)
make test       # quick tests without the race detector
make cover      # coverage summary; HTML report in coverage/ (gitignored)
make bench      # benchmarks (root package only)
make bench-compare  # benchstat: BENCH_BASE (default the latest tag) vs the working tree
make api        # gorelease: exported API changes since the latest tag, and the suggested version
make fuzz       # FuzzParseString; FUZZTIME=30s by default
make fmt        # golangci-lint fmt: gofumpt, goimports, golines (120 columns)
```

- **golangci-lint:** `.golangci.yaml` (v2 format) lists every enabled linter explicitly, grouped by purpose, and gives a reason for each linter it leaves out. It is adapted from the golden config by Marat Reimers. golangci-lint also runs staticcheck with all checks, so there is no separate staticcheck step.
- **Pinned version:** golangci-lint `v2.15.0`, set in the Makefile (`GOLANGCI_LINT_VERSION`) and in CI. `make lint` refuses to run any other version and points to `make tools`, which builds that version into `./bin` with `GOTOOLCHAIN` set to the `go.mod` Go (golangci-lint's own `go.mod` would otherwise select an older Go, too old for this module).
- **Lint Go version:** CI runs golangci-lint with the `go.mod` Go version. A golangci-lint binary can't read the export data of Go releases newer than it supports: v2.13.2 fails on Go 1.27.4, for example. Linting with the latest Go would therefore break whenever Go is released.
- **Upgrading golangci-lint:** upgrade on purpose, in its own PR. Bump both pins, read the changelog for renamed options, then fix the new findings or disable a check with a reason. Enabling a new linter is a separate, reviewed change.
- **`nolint`:** add a `//nolint:<linter> // reason` directive only with a reason. `nolintlint` flags unused directives.
- **Generated code:** after changing `internal/gen` or the measures list in it, run `make generate` and commit the result. `TestGenerated` fails when a generated file differs from what the generator produces.
- **Done means:** `make check` is clean, coverage of the library is still 100% (`make cover`; only `internal/gen`'s `main` is uncovered), and benchmarks show no hot-path regression.

## Recipes

### Add a unit to an existing measure

1. Append a constant to the unit block in `<measure>/units.go`, before `unitCount`, with its symbol as the line comment, and add its entry to `definitions`. Use exact defined values (NIST SP 811, the BIPM SI Brochure, CODATA) and note the source in a comment. Write derived factors as constant expressions, not rounded decimals: `1000.0 / 3600`, not `0.277778`.
2. Append, don't insert: users mustn't persist the numbers, but appending keeps them stable anyway.
3. Check that the symbol is unique across **all** measures (`go test -run TestMeasuresSelfCheck .`).
4. For the micro prefix, use U+03BC (`μ`). U+00B5 input is canonicalised to it.
5. Add an independent reference value to `TestConversion_ReferenceValues` in `conversion_test.go`. `testUnitPackage` checks every unit's factor, rate and symbol automatically.
6. Update the units table in `README.md` and add a `CHANGELOG.md` entry.

### Add a new measure type

1. **`<measure>/units.go`:** the package documentation, `type Unit uint8`, the constants (base unit first) ending with `unitCount`, `BaseUnit`, and `definitions`. Mirror `area/units.go`.
2. **`internal/gen/main.go`:** add the measure to `measures`, then `make generate`. That writes `<measure>/unit_gen.go` and the measure's shared API in `measures_gen.go`.
3. **`<measure>.go` in the root package:** the struct type `X {value float64; unit <measure>.Unit}` with its documentation, and any cross-measure operations. Mirror `area.go`.
4. **`measures.go`:** add a `MeasureType` constant.
5. **Tests:**
   - an `xCase` value in `measure_types_test.go`, run by `TestTypedMeasures` (the whole shared API, plus `testUnitPackage`) and `BenchmarkMeasures`;
   - an entry in `predefinedFactors` and in `TestConversionFactors_ReturnsCopy` (`conversion_test.go`);
   - reference values;
   - rows in `TestCrossMeasureOperations` (and `_ZeroDivision`) for any cross-measure operations.
6. **Docs:** an `ExampleNewX` in `example_measures_test.go` (and one per cross-measure operation), `README.md` (units table and operations table) and `CHANGELOG.md`.

The factor model is multiplicative only. Temperature, with its offsets, has its own hand-written types (`temperature.go`). Don't shoehorn other non-linear scales (dB, pH) into either; raise them as a design question instead.

## Code conventions

- Format with `make fmt`. It runs gofumpt and goimports, with stdlib imports first, a blank line, then module imports, and golines wraps lines longer than 120 columns.
- Error strings are lower-case with no trailing punctuation and include the offending unit.
- **Every returned error matches exactly one sentinel** (`TestErrors` checks each one). Use `unknownUnitError(mt, unit)` for unknown units, `ErrZeroDivision` itself, an `incompatibleError` constant for a fixed message, or `newError(kind, msg, cause)` for anything else; `cause` is the underlying error, or nil. A new kind of error needs a new sentinel, added to `sentinels` in `errors_test.go` and to the README table.
- Return errors rather than panicking. The only intended panics are `MustNewUnitSet`, `MustUnitOf`, `rates.MustValidate` on constant data, and conversions in an invalid unit.
- Exported identifiers need doc comments that start with the identifier's name, as do package clauses (`revive` enforces both).
- **Look-alike characters:** U+00B5 and U+03BC look identical, as do ℃ (U+2103) and °C. Wherever the difference matters in code (comparisons, tests), write them as Go escape sequences, as `internal/rates.CanonicalUnit` and the tests do. Comments and the `definitions` arrays keep the readable characters; `TestMeasuresSelfCheck` checks that definitions are canonical. Some editing tools turn escapes into the literal characters when they write a file; check with `grep` afterwards.
- **Shadowed `time` package:** the `time` subpackage shadows the stdlib one. Import the stdlib package as `builtintime "time"`, as `time.go` and `time_test.go` do.
- **Renaming exported names** breaks users. Add the new name and keep the old one as an alias with a `// Deprecated: Use X.` paragraph, as `distance.PlankLength` does. Unexported names can change freely.
- **Generated files** (`*_gen.go`) start with `// Code generated … DO NOT EDIT.`; change the generator instead.

## Testing conventions

- Black-box tests live in package `units_test`. The only internal test is `measures_test.go` (`package units`, `//nolint:testpackage`), because it inspects the unexported unit sets. `internal/rates` and `internal/gen` have their own tests.
- Compare floats with `almostEqual(a, b, epsilon)` or `assertMeasurement(t, m, value, symbol)` from `helpers_test.go`, which check relative error. Never use `==` on converted values; values that weren't converted compare with `==` as structs.
- Prefer table-driven tests with `t.Run`. Use `t.Fatal` when a setup step fails, and `assertPanics` for intended panics.
- On error paths, assert that the result is the zero value (or nil, for a `Measurement`) as well as that the error matches its sentinel.
- `measure_types_test.go` runs the whole shared API against every measure and against `Quantity` (`testTypedMeasure`, `testGenericMethods`, `testText`, `testUnitPackage`). A new shared method needs a check there.
- `Example*` functions in `example_test.go` mirror `README.md` and have verified `// Output:` blocks. Keep the README snippets in sync with them.
- `example_measures_test.go` has an example for each measure, each cross-measure operation, custom measures, temperature and formatting, and the unit packages have their own (`distance/example_test.go`). They are what pkg.go.dev shows; a new measure or cross-measure operation needs one.
- If a test fails because of a library bug, fix the library or skip the test with `t.Skip("known defect: …")`. Never change expected values to match wrong output.
- Benchmarks are in `measures_benchmark_test.go` (every measure, custom measures, a sum of mixed distances, error paths, the parser, `Rate`) and `distance_benchmark_test.go` (random data). Keep benchmark names stable, because `make bench-compare` matches them across versions. They loop with `b.Loop()` (`for i := 0; b.Loop(); i++` when an index is needed), which keeps the compiler from optimizing the loop body away, and use the `generateRandom*` helpers, which draw from `math/rand/v2`.
- Package-level sinks for benchmarks and allocation tests are in `helpers_test.go`, not in `alloc_test.go`, which `-race` builds exclude.
- `immutability_test.go` enforces the immutability principle (no exported variables, no internal types in the exported API, race-free concurrent use). If one of its tests fails, fix the code, not the test.

## Known defects and design debt

None known. Verified on 2026-11-14, after the value-type redesign removed the three earlier entries (interface-based values, typed results lost by the generic operations, and the costly `NewBaseMeasurement`). Add an entry here for each defect that needs an exported API change, and delete it once it is fixed.

## Commits

Match the existing history:

- **Prefixes:** `feat:`, `fix:`, `test:`, `ref:` (refactor), `perf:`, `docs:`, `chore:` (tooling).
- **Summary:** short and lower-case, e.g. `feat: added area measure`.
- **API breaks:** note any break to the exported API in the message body, and in `CHANGELOG.md`.
- **No commits or pushes** unless asked.
