# AGENTS.md

Guidance for AI coding agents (and humans) working in this repository. `CLAUDE.md` imports this file, so keep all shared rules here.

## What this is

`go-physics-units` (module `github.com/asaphin/go-physics-units`, package name `units`) is a small Go library for creating, converting, parsing and doing arithmetic on physical measurements (distance, area, mass, time, velocity) and user-defined measures (e.g. currency).

- **Published and pre-1.0** (tags `v0.0.0`–`v0.0.4`). Exported API changes affect real users: make them deliberately, call them out, and record them in `CHANGELOG.md`. `make api`, and the CI job "API changes", list every exported change since the latest tag; check that each one is in the CHANGELOG.
- **Go version:** `go.mod` declares `go 1.27.0`. CI tests that exact version and the latest stable release. Use modern Go freely: iterators and the `maps`/`slices` helpers, `strings.FieldsSeq`, `testing.B.Loop`, `sync.WaitGroup.Go`, and generic methods (Go 1.27). The `modernize` linter points out older forms. Raising the minimum again is its own deliberate change, with a CHANGELOG entry.
- **No third-party dependencies**, including in tests (stdlib `testing` only, no testify). Keep it that way unless asked.
- **Branches:** `dev` is the working branch, and it only accepts pull requests. Work on a short-lived branch and open a PR into `dev`. Releases are merged from `dev` into `main`.

## Design principles

The library rests on three principles. They take priority over convenience and over new features. Get the user's explicit agreement before making a change that weakens any of them, and say which one it weakens.

### 1. Performance

Measurements are created and combined in hot loops, so their cost is part of the API.

- **Calculations come first.** Conversions, arithmetic, comparisons and cross-measure operations are what users run in loops; keep them fast. Constructors and table building may get somewhat slower when a feature needs it: say so in the PR, but it isn't a blocker (the author's call, 2026-11-13).
- **Do the work once.** All conversion work happens at package init, or once per `NewMeasurementCreator`. A conversion at run time is one map lookup and one multiplication.
- **Hot paths** are `NewX`, `ConvertTo`, `Add`/`Sub`/`Mul`/`Div`, `Compare`/`Equal`, the cross-measure operations, `String` and `ParseString`. Keep them free of reflection, `fmt`, and any allocation beyond the result value.
- **Error paths too.** Build error messages by string concatenation, never with `fmt`, which formats named string types such as `MeasureType` through reflection. Return a constant where the message is fixed (`ErrZeroDivision`, the `errXConversion` constants): that costs no allocation. `*UnitError` builds its message in `Error()`, so returning one costs a single allocation. `TestAllocations` guards the error paths too.
- **Inlining.** Keep the `rates.Table` methods small enough to inline. Check with `go build -gcflags=-m ./internal/rates`.
- **Allocation guards.** `TestAllocations` (`alloc_test.go`) fixes the allocation count of each hot path, and CI fails if one rises. Unlike timings, allocation counts are deterministic. The race detector changes them, so the file is excluded from `-race` builds, and `make allocs` and CI run it separately. Lower a limit when you save an allocation; raise one only on purpose, and say why in the PR.
- **Measure.** For any change that touches a hot path, run `make bench-compare` (`BENCH_BASE=<ref>`, default the latest tag). It benchmarks both sides on the same Go toolchain and compares them with benchstat. Differences around ±10% are noise on this suite. Report every regression to the user; don't accept one silently.

### 2. Immutability after init (data security)

Shared state is built during package initialisation and never written afterwards. That makes the library safe for concurrent use without locks. It also means no importer can change how conversions behave for the rest of the program.

- **No exported package-level variables.** Any importer can reassign one, so exported names are constants or functions. `TestNoExportedPackageVariables` enforces this.
- **Never hand out a reference to shared state.** Return copies, as `ConversionFactors()` does, or read-only handles.
  - The predefined rate tables are unexported variables of the root package. Unit packages expose rates only as numbers, through `Rate(from, to)`.
  - No exported signature may use a type from an internal package. `TestNoInternalTypesInExportedAPI` enforces this.
  - `rates.Table` is a read-only handle passed by value: a struct whose only field points to unexported data. Never hand out a `*rates.Table`. A caller who can't name the type can still write `*p = other` through a pointer and overwrite a shared table.
- **Copy caller data on the way in.** `rates.NewTable` copies the factors, so later changes to the caller's map have no effect.
- **No writes after init.** That rules out lazy initialisation, caches and memoisation of shared state. `TestConcurrentUse` catches them under `-race`.
- **Immutable values.** Value fields are unexported, and every operation returns a new value. `Measurement` is sealed by an unexported method, so no outside type can pose as a measurement.
- **Exported errors are constants.** The sentinels are constants of the exported type `Error` (`errors.go`). Never add `var ErrX = errors.New(…)`: any importer could reassign it. Detail errors such as `*UnitError` are allocated on every call, so a caller who changes one changes only their own copy. Never return a shared pointer to an error.

### 3. Customization

User-defined measures are first-class: they convert, combine and parse like the predefined ones.

- **Two entry points.** `NewMeasurementCreator` builds the table once and names the measure. `NewBaseMeasurement` is for one-off values.
- **Custom symbols are independent of predefined ones.** A custom `"m"` is not a distance.
- **Predefined measures can be extended.** Copy their factors with `ConversionFactors()`, add units, and pass the result to `NewMeasurementCreator`. `TestMeasurementCreator_ExtendsPredefinedFactors` covers this.
- **Every restriction needs a correctness reason; don't add one without it.** The current restrictions and their reasons:
  - **Factors** must be positive and finite, because rates divide by them.
  - **Measure names** can't reuse a predefined name. Otherwise `ToDistance` would accept the value and convert it with the wrong table.
  - **`T`** must be satisfied by `Measurement` alone.
- **New features** for the predefined measures should, where it makes sense, work for custom measures too.

## Layout

| Path | Role |
|---|---|
| `doc.go` | Package documentation |
| `measurement.go` | `Measurement` interface (sealed: it has unexported methods), `baseMeasurement`, `NewMeasurement`, `ParseString`, `NewBaseMeasurement`, shared helpers |
| `measures.go` | `MeasureType` constants, the predefined rate tables (`distanceRates`, …), the registries `measureToConversionFactorsMapping` and `measureToConversionRatesMapping`, the `unitToMeasureType` index, `DetectMeasureType` |
| `common.go` | Generic `Arithmetics[T]` and `UnitConverter[T]` interfaces |
| `errors.go` | The sentinel errors (constants of type `Error`), `UnitError`, and the unexported error types and constructors the package uses |
| `distance.go`, `area.go`, `mass.go`, `time.go`, `velocity.go` | One file per measure: the public interface, the `xxxImplementation` struct embedding `baseMeasurement`, constructors and cross-measure operations. All five follow the same pattern |
| `distance/`, `area/`, `mass/`, `time/`, `velocity/` | One `units.go` per measure: unit symbol constants, `BaseUnit`, the `conversionFactors` table, `ConversionFactors()` (returns a copy) and `Rate(from, to)` |
| `conversion/` | Public `Factors` type alias (`map[string]float64`) and `CopyConversionFactors` |
| `measurement_creator.go` | Generic `NewMeasurementCreator[T]` for user-defined measures |
| `internal/rates` | `Table`: the immutable, validated N×N rate table of one measure, plus each unit's factor, a read-only handle passed by value. Also `CanonicalUnit` |
| `Makefile` | Task runner (see Commands) |
| `.github/workflows/ci.yml` | CI: `go mod tidy` check, vet, race tests and coverage (the `go.mod` Go version and stable), a short fuzz run, golangci-lint (pinned version), govulncheck |
| `CHANGELOG.md` | User-visible changes, under `## Unreleased` until the next tag |
| `ROADMAP.md` | Ordered plan of upcoming work. Take steps in order and tick them off when done |

## How it works

1. **Factors.** Each unit subpackage defines how many base units are in one unit: `distance.Kilometer: 1000` means 1 km = 1000 m. Base units are coherent SI units (m, m², kg, s, m/s).
2. **Rate table.** `rates.NewTable(measure, factors)` validates the factors and expands them into a `map[pair{from, to}]float64` with value `f[from] / f[to]`. Validation requires every factor to be positive and finite and every symbol to be non-empty. The root package builds the predefined tables at package init in `measures.go` with `rates.MustNewTable`, so a bad factor panics at startup. A `Table` is a one-pointer handle to unexported data: copying it is cheap, and a copy can't be used to change the original (see Design principles).
3. **Values.** A measurement is `baseMeasurement{value, unit, conversionRates}`. Measure-specific types embed it (`distanceImplementation{baseMeasurement}`) and are returned as interfaces. Values are immutable: every operation allocates and returns a new value.
4. **Measure identity.** A measurement's measure is the name of its table: `Type()` returns `MeasureType(table.Measure())`.
   - Predefined tables are named after their `MeasureType` constant, and `TestMeasureRegistriesAreConsistent` enforces this.
   - `NewMeasurementCreator` names the table after the measure type you pass; predefined names are rejected.
   - `NewBaseMeasurement` names it `custom`.
   - `DetectMeasureType(unit)` maps a predefined symbol to its measure through an index. It is only used where no table is known yet (`NewMeasurement`, `ParseString`). Because of it, **predefined symbols must be unique across all measures**; `TestMeasuresSelfCheck` enforces this.
5. **Micro sign.** `rates.CanonicalUnit` maps U+00B5 MICRO SIGN to U+03BC GREEK SMALL LETTER MU. Every entry point that receives a unit string canonicalises it (`newBaseMeasurement`, `convertTo`, `DetectMeasureType`), and tables store canonical keys. Internal paths that only handle units already stored in a value skip this step.
6. **Combining values.**
   - `AddMeasurement` and `SubMeasurement` check that the types match *and* that a rate between the two units exists.
   - Typed `Add`/`Sub` need neither check: two values of one Go type share the same table.
   - Cross-measure operations (`Distance.DivideByTime`, `Distance.Square`, `Area.DivideByDistance`, `Velocity.MultiplyByTime`) convert to base units and return results in the target measure's base unit. Dividing by a zero measurement returns an error, as `Div` does for a zero divisor.
   - **Comparisons** (`Compare`, `Equal`, `EqualWithinRel`, `CompareMeasurement`) compare two values in the same unit exactly. Otherwise they multiply each value by its unit's factor (`Table.Factor`), which gives the value in the base unit. The result doesn't depend on the order of the operands, so `a.Equal(b) == b.Equal(a)`; converting one value into the other's unit would break that at the last bit. `TestCompare_Symmetric` enforces it. `CompareMeasurement` falls back to the receiver's unit only for two tables that share a measure name.
7. **Parsing.** `ParseString` accepts exactly `"<float> <unit>"`, two whitespace-separated fields. It rejects NaN and ±Inf.

Constructor conventions per measure:

- `NewX(value, unit) (X, error)` validates the unit via `newBaseMeasurement`.
- Unexported `newX(value, unit) X` skips validation. Use it only internally, with valid, canonical units.
- `ToX(Measurement) (X, error)` narrows any measurement of that measure (generic or typed) through the unexported `base()` method.
- `StringToX(s)` is `ParseString` followed by `ToX`.

On error, always return a nil interface (`return nil, err`, or `var zero T` in generic code), never an interface that wraps a nil pointer.

## Commands

```sh
make check      # what CI runs: go mod tidy check, format check, go vet, golangci-lint, go test -race, allocation guards
make test       # quick tests without the race detector
make cover      # coverage summary; HTML report in coverage/ (gitignored)
make bench      # benchmarks (root package only)
make bench-compare  # benchstat: BENCH_BASE (default the latest tag) vs the working tree
make api        # gorelease: exported API changes since the latest tag, and the suggested version
make fuzz       # FuzzParseString; FUZZTIME=30s by default
make fmt        # golangci-lint fmt: gofumpt, goimports, golines (120 columns)
```

- **golangci-lint:** `.golangci.yaml` (v2 format) lists every enabled linter explicitly, grouped by purpose, and gives a reason for each linter it leaves out. It is adapted from the golden config by Marat Reimers. golangci-lint also runs staticcheck with all checks, so there is no separate staticcheck step.
- **Pinned version:** golangci-lint `v2.15.0`, set in the Makefile (`GOLANGCI_LINT_VERSION`) and in CI. `make lint` refuses to run any other version and prints the install command.
- **Lint Go version:** CI runs golangci-lint with the `go.mod` Go version. A golangci-lint binary can't read the export data of Go releases newer than it supports: v2.13.2 fails on Go 1.27.4, for example. Linting with the latest Go would therefore break whenever Go is released.
- **Upgrading golangci-lint:** upgrade on purpose, in its own PR. Bump both pins, read the changelog for renamed options, then fix the new findings or disable a check with a reason. Enabling a new linter is a separate, reviewed change.
- **`nolint`:** add a `//nolint:<linter> // reason` directive only with a reason. `nolintlint` flags unused directives.
- **Done means:** `make check` is clean, coverage is still around 99% (`make cover`), and benchmarks show no hot-path regression.

## Recipes

### Add a unit to an existing measure

1. Add the symbol constant and its factor to `<measure>/units.go`. Use exact defined values (NIST SP 811, the BIPM SI Brochure, CODATA) and note the source in a comment. Write derived factors as constant expressions, not rounded decimals: `1000.0 / 3600`, not `0.277778`.
2. Check that the symbol is unique across **all** measures (`go test -run TestMeasuresSelfCheck .`).
3. For the micro prefix, use U+03BC (`μ`). U+00B5 input is canonicalised to it.
4. Add an independent reference value to `TestConversion_ReferenceValues` in `conversion_test.go`. `TestConversion_AllUnitPairs` picks up the new unit automatically.
5. Update the units table in `README.md` and add a `CHANGELOG.md` entry.

### Add a new measure type

Mirror an existing measure, e.g. `area.go` and `area/units.go`, and update all of these:

1. **`<measure>/units.go`:** constants, `BaseUnit`, `conversionFactors` (with canonical symbols), `ConversionFactors()` and `Rate()`.
2. **`<measure>.go` in the root package:**
   - a public interface embedding `Measurement`, `UnitConverter[X]`, `Arithmetics[X]` and `Comparer[X]`, plus an `xImplementation` struct;
   - the methods `ConvertTo`, `MustConvertTo`, `ConvertToBaseUnits`, `valueInBaseUnits`, `valueInUnits`, `Add`, `Sub`, `Mul`, `Div`, `Compare`, `Equal`, `EqualWithinRel`;
   - the functions `NewX`, `ToX`, `StringToX` (and `newX` if a cross-measure operation returns this measure) and an `errXConversion` constant of type `incompatibleError`.
3. **`measures.go`:** add a `MeasureType` constant, build the table with `xRates = rates.MustNewTable(string(MeasureX), x.ConversionFactors())`, and register the measure in **both** mapping maps.
4. **Tests:**
   - an `xCase` value in `measure_types_test.go`, used by `TestTypedMeasures` (the whole shared API) and `BenchmarkMeasures`;
   - an entry in `predefinedFactors` and in `TestConversionFactors_ReturnsCopy` (`conversion_test.go`);
   - reference values;
   - tests for any cross-measure operations.
5. **Docs:** update `README.md` (units table and operations table) and `CHANGELOG.md`.

The factor model is multiplicative only. It can't express measures with an offset (°C, °F) or non-linear scales (dB, pH). Don't shoehorn them in; raise it as a design question instead.

## Code conventions

- Format with `make fmt`. It runs gofumpt and goimports, with stdlib imports first, a blank line, then module imports, and golines wraps lines longer than 120 columns.
- Error strings are lower-case with no trailing punctuation and include the offending unit.
- **Every returned error matches exactly one sentinel** (`TestErrors` checks each one). Use `unknownUnitError(mt, unit)` for unknown units, `ErrZeroDivision` itself, an `incompatibleError` constant for a fixed message, or `newError(kind, msg, cause)` for anything else; `cause` is the underlying error, or nil. A new kind of error needs a new sentinel, added to `sentinels` in `errors_test.go` and to the README table.
- Return errors rather than panicking. The only intended panics are `Must*` methods and `rates.MustNewTable` on constant data.
- Exported identifiers need doc comments that start with the identifier's name, as do package clauses (`revive` enforces both).
- **Look-alike characters:** U+00B5 and U+03BC look identical. Wherever the difference matters in code (comparisons, tests), write them as Go escape sequences, as `internal/rates.CanonicalUnit` and the tests do. Unit constants keep the readable literal symbol.
- **Shadowed `time` package:** the `time` subpackage shadows the stdlib one. Import the stdlib package as `builtintime "time"`, as `time.go` and `time_test.go` do.
- **Renaming exported names** breaks users. Add the new name and keep the old one as an alias with a `// Deprecated: Use X.` paragraph, as `distance.PlankLength` does. Unexported names can change freely.

## Testing conventions

- Black-box tests live in package `units_test`. The only internal test is `measures_test.go` (`package units`, `//nolint:testpackage`), because it inspects the unexported registries. `internal/rates` has its own tests.
- Compare floats with `almostEqual(a, b, epsilon)` or `assertMeasurement(t, m, value, unit)` from `helpers_test.go`, which check relative error. Never use `==` on converted values.
- Prefer table-driven tests with `t.Run`. Use `t.Fatal` when a setup step fails, and `assertPanics` for intended panics.
- On error paths, assert that the result is nil as well as that the error is non-nil.
- `Example*` functions in `example_test.go` mirror `README.md` and have verified `// Output:` blocks. Keep the README snippets in sync with them.
- If a test fails because of a library bug, fix the library or skip the test with `t.Skip("known defect: …")`. Never change expected values to match wrong output.
- Benchmarks are in `measures_benchmark_test.go` (every measure, custom measures, error paths, `Rate`) and `distance_benchmark_test.go` (random data). Keep benchmark names stable, because `make bench-compare` matches them across versions. They loop with `b.Loop()` (`for i := 0; b.Loop(); i++` when an index is needed) and use the `generateRandom*` helpers, which draw from `math/rand/v2`.
- `immutability_test.go` enforces the immutability principle (no exported variables, no internal types in the exported API, race-free concurrent use). If one of its tests fails, fix the code, not the test.

## Known defects and design debt

Verified on 2026-11-13. Delete each entry once it is fixed. Each of these needs an exported API change, so agree on the fix before starting. `ROADMAP.md` schedules the fixes.

- **`XxxMeasurement` methods lose the type.** `AddMeasurement` and the other `XxxMeasurement` methods on a typed value (e.g. a `Distance`) return a generic `Measurement`, not a `Distance`.
- **Costly `NewBaseMeasurement`.** It builds a full N×N table on every call; `NewMeasurementCreator` is the efficient path.
- **Interface-based design.** Every value is a heap-allocated interface (about 65–200 ns per operation, one allocation). Units are plain strings, so `NewDistance(1, mass.Kilogram)` compiles. A value-type redesign with typed units was prototyped at about 1.4 ns with 0 allocations, but it hasn't been decided.

## Commits

Match the existing history:

- **Prefixes:** `feat:`, `fix:`, `test:`, `ref:` (refactor), `perf:`, `docs:`, `chore:` (tooling).
- **Summary:** short and lower-case, e.g. `feat: added area measure`.
- **API breaks:** note any break to the exported API in the message body, and in `CHANGELOG.md`.
- **No commits or pushes** unless asked.
