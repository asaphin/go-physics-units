# Roadmap

This is the ordered plan for improving go-physics-units. Take the steps from top to bottom. It was last reviewed on 2026-11-13.

## How to use this plan

- **One step, one PR.** Each step is one PR into `dev`. It keeps `make check` green and leaves the library releasable.
- **Principles first.** Every step keeps the three design principles in [AGENTS.md](AGENTS.md#design-principles): **P1** performance, **P2** immutability after init, **P3** customization. Each step names the principles it serves.
- **Decisions.** Steps marked **Decision** change the exported API or the design. Agree on them before starting.
- **Finishing a step.** Tick it in the table below, add its `CHANGELOG.md` entry, and delete any entry it fixes from "Known defects and design debt" in AGENTS.md.

Two rules keep the steps small:

- **Added methods break no one.** The measure interfaces are sealed by an unexported method, so no code outside the package implements them. Adding a method is therefore not a breaking change.
- **Breaking changes are batched.** Users have to migrate when the API breaks, so breaks go into as few releases as possible. These are v0.1.0, which already breaks the API, and v0.2.0, the redesign.

## Overview

Size is relative effort: S < M < L.

| Done | Step | What | Serves | Breaking | Size | Release |
|:-:|---|---|---|---|:-:|---|
| ☑ | [0.1](#01-commit-the-current-work-and-get-ci-green) | Commit the current work, first CI run | — | — | S | — |
| ☑ | [0.2](#02-decision-replace-conversionrates-with-rate) | Replace `ConversionRates()` with `Rate()` | P2 | yes, **Decision** | S | v0.1.0 |
| ☑ | [0.3](#03-decision-zero-division-returns-an-error-everywhere) | Zero division returns an error everywhere | correctness | yes, **Decision** | S | v0.1.0 |
| ☐ | [0.4](#04-release-v010) | Release v0.1.0 | — | — | S | v0.1.0 |
| ☑ | [1.1](#11-allocation-guards-and-wider-benchmarks) | Allocation guards and wider benchmarks | P1 | no | S | — |
| ☑ | [1.2](#12-api-change-report-in-ci) | API-change report in CI | all | no | S | — |
| ☑ | [1.3](#13-formatting-with-fmt-verbs) | `fmt.Formatter` (`%.2f` → `5.00 km`) | P1, P2 | no | S | v0.1.x |
| ☑ | [1.4](#14-immutable-exported-errors) | Immutable exported errors for `errors.Is` | P2 | no | M | v0.1.x |
| ☑ | [1.5](#15-compare-and-equal) | `Compare` and `Equal` | P3 | no | S | v0.1.x |
| ☑ | [1.6](#16-timetoduration) | `Time.ToDuration` | P3 | no | S | v0.1.x |
| ☑ | [1.7](#17-examples-for-every-measure) | Examples for every measure | docs | no | S | v0.1.x |
| ☑ | [2.1](#21-decision-gate-design-rfc-and-prototype) | RFC and prototype: value types, typed units | P1–P3 | **Decision gate** | M | — |
| ☐ | [3.1](#31-typed-unit-symbols) | Typed unit symbols | P1, safety | yes | M | v0.2.0 |
| ☐ | [3.2A](#32a-value-type-measures) | Value-type measures (Path A) | P1, P2 | yes | L | v0.2.0 |
| ☐ | [3.3A](#33a-typed-custom-measures) | Typed custom measures (Path A) | P3 | yes | M | v0.2.0 |
| ☐ | [3.2B](#32b-typed-results-from-generic-operations) | Typed results from generic operations (Path B) | P3 | no | M | v0.2.0 |
| ☐ | [3.4](#34-migration-guide-and-release-v020) | Migration guide, release v0.2.0 | — | — | S | v0.2.0 |
| ☐ | [4.1](#41-text-and-json-encoding) | Text and JSON encoding | P3 | no | M | v0.2.x |
| ☐ | [4.2](#42-configurable-parser) | Configurable parser with aliases | P3, P2 | no | M | v0.2.x |
| ☐ | [4.3](#43-new-measures-one-per-step) | New measures, one per step | P3 | no | S each | v0.2.x |
| ☐ | [4.4](#44-decision-temperature) | Temperature (design first) | P3 | **Decision** | M | v0.2.x |
| ☐ | [5.1](#51-stabilise-and-release-v100) | Stabilise and release v1.0.0 | — | removals | M | v1.0.0 |

## Phase 0: finish and ship v0.1.0

The work since v0.0.4 already breaks the API: it changes the `mass.Quintal` symbol, the `ConversionRates()` type, the names of custom measures and the minimum Go version. Steps 0.2 and 0.3 add the two remaining small breaks to the same release, so users migrate once.

### 0.1 Commit the current work and get CI green

- Split the uncommitted work into reviewable commits on `dev`:
  - `fix:` the rate table, the data race, `ToX`, the creator and precision;
  - `perf:`;
  - `test:`;
  - `chore:` the Makefile, golangci-lint, CI and the Go version;
  - `docs:`.
- Push `dev`. This is the first run of `.github/workflows/ci.yml`. Check that every pinned action version resolves and that all jobs pass: tests on the `go.mod` version and stable, lint and govulncheck.
- **Done when** CI is green on `dev`.

### 0.2 Decision: replace `ConversionRates()` with `Rate()`

**Why:** v0.1.0 already breaks `ConversionRates()`. It used to return an `immutable.Float64Map` whose concatenated keys caused the m↔mm bug, and it now returns `rates.Table`, a type from an internal package (a known defect). Finishing the change now fixes that defect within the same migration.

- Add `Rate(from, to string) (float64, bool)` to each unit package. Compute it from the factors, accept both micro signs, and don't allocate.
- Remove `ConversionRates()`. Build the predefined tables in the root package (`measures.go`) from `ConversionFactors()`, so no exported signature mentions an internal type.
- Update `immutability_test.go` to match.
- **Alternative:** keep `ConversionRates()` returning `rates.Table` and live with the internal type until v0.2.0.

### 0.3 Decision: zero division returns an error everywhere

**Why:** `Distance.DivideByTime` and `Area.DivideByDistance` panic on a zero divisor, while `Div` returns an error (a known defect).

- Change both to return `(X, error)` with the same zero-division error as `Div`. Call sites gain an error check.
- It goes into v0.1.0 for the same reason as 0.2: one migration.
- **Alternative:** keep the panic, document it on both methods, and fix it in v0.2.0.

### 0.4 Release v0.1.0

- **Version:** use a minor bump, not v0.0.5. The release breaks the API, and in v0 the minor number signals that.
- **CHANGELOG:** rename `## Unreleased` to `## v0.1.0 - <date>` and add a short "Upgrading from v0.0.4" list covering every break.
- **Tag:**
  - merge `dev` into `main` through a PR, as for earlier releases;
  - tag `v0.1.0` on `main` and push the tag;
  - create a GitHub release from the CHANGELOG section.
- **pkg.go.dev:** request the version from the proxy so the docs update: `GOPROXY=https://proxy.golang.org go list -m github.com/asaphin/go-physics-units@v0.1.0`.

## Phase 1: non-breaking improvements (v0.1.x)

Everything in this phase stays useful whatever Phase 2 decides. Tag a v0.1.x release after each user-visible step, or after a group of them.

### 1.1 Allocation guards and wider benchmarks

Serves **P1**. Timing benchmarks on shared CI runners are too noisy to fail a build on, but allocation counts are deterministic.

- **Allocation tests.** Add `testing.AllocsPerRun` tests for the hot paths, with today's counts as limits: 1 for `NewDistance`, `ConvertTo`, `Add`, `Mul`, `String`, `ParseString` and the creator's `New`, and 2 for the unknown-unit error. Put them in a file with `//go:build !race`, because the race detector can change allocation counts.
- **Wider benchmarks.** Cover every measure, not only distance, plus custom measures (`New` on a creator, `NewBaseMeasurement`), formatting and error paths.
- **`make bench-compare`.** It benchmarks the latest tag and the working tree, then compares them with benchstat. Run benchstat with `go run golang.org/x/perf/cmd/benchstat@<pinned version>`: it's a tool, not a module dependency.
- **Done when** undoing any allocation saving makes a test fail.

### 1.2 API-change report in CI

- Add a CI job that runs `gorelease` (`golang.org/x/exp/cmd/gorelease`, pinned) against the latest tag. It prints the incompatible API changes and suggests the next version.
- The job only reports until v1.0.0; step 5.1 makes it fail the build.
- It shows in every PR whether the exported API changed, which AGENTS.md requires to be deliberate.

### 1.3 Formatting with fmt verbs

Serves **P1** and **P2**. Today, `fmt.Printf("%.2f", d)` prints the internal struct, including the whole distance rate table, instead of `5.00 km`. Only `%v` and `%s` use `String()`.

- **Implement `fmt.Formatter`** on every measurement:
  - `%v` and `%s` stay as they are;
  - `%f`, `%e` and `%g`, with width and precision, format the value and append the unit;
  - `%#v` never shows internals;
  - other verbs produce fmt's usual `%!verb(…)` output.
- Keep `String()` on its current fast path: 1 allocation, no `fmt`.
- **Tests:** a table of verbs. The fuzz test's round trip through `%v` must keep passing.

### 1.4 Immutable exported errors

Serves **P2**. Today, callers can only tell errors apart by their message.

- **Exported errors** for `errors.Is`: unknown unit, incompatible measures, division by zero, invalid format and invalid factors.
- **Constants, not variables.** Make them constants of an exported error type, such as `const ErrUnknownUnit = Error("unknown unit")`. Any importer could reassign `var ErrX = errors.New(…)`.
- **Details.** A detail type carries the measure and the unit for `errors.As`, and keeps today's messages unchanged.
- **Performance (P1).** Error paths stay at one allocation, without `fmt`.

### 1.5 Compare and Equal

Serves **P3**.

- Add `Compare(other X) int` (semantics of `cmp.Compare`) and `Equal(other X) bool` to each typed measure.
- Conversions round, so values that should be equal can differ in the last bit. Document this, and add a variant with a relative tolerance.
- **As built:** values in different units are compared in the base unit, not in the receiver's unit, so the comparison is symmetric. `CompareMeasurement` gives custom measures the same comparison.

### 1.6 Time.ToDuration

Serves **P3**.

- `ParseDuration` turns a `time.Duration` into a `Time`. Add the reverse, `Time.ToDuration() (time.Duration, error)`.
- It rounds to the nearest nanosecond and returns an error outside the range of `time.Duration`, about ±292 years.

### 1.7 Examples for every measure

- Add `Example` functions for each measure, for the cross-measure operations and for custom measures, so that pkg.go.dev shows runnable examples.
- Keep the README snippets in sync with them.

## Phase 2: decide the core design

### 2.1 Decision gate: design RFC and prototype

**Why:** the remaining large costs all come from the interface-based design.

- **Allocations (P1).** Every value is a heap-allocated interface, at about 65–200 ns and one allocation per operation. The earlier value-type prototype took about 1.4 ns with no allocations.
- **Untyped units.** Units are plain strings, so `NewDistance(1, mass.Kilogram)` compiles and fails only at run time.
- **Interface results.** Constructors return nil interfaces on error, and generic operations return values whose type is lost. For example, `ParseString("1 km")` doesn't satisfy `units.Distance`.

Write `docs/rfc-value-types.md` and answer these questions in it:

1. **Values.** Switch to `type Distance struct{…}` with value receivers instead of an interface. What does the zero value mean, and what do constructors return on error?
2. **Units.** This question is where the principles pull hardest against each other:
   - **Typed string constants** (`type Unit string`) are simple and catch a wrong measure at compile time. Conversions still need a map lookup.
   - **Typed integer constants indexing an immutable array** are the fastest option, but the set of units is closed.
   - **Unit values that carry their own factor** let users add units. The predefined units would then have to be exported variables, which P2 forbids.
3. **Custom measures (P3).** How does a user-defined measure get the same typed API as the predefined ones? One option is a generic `Quantity[M]` with a marker type supplied by the user. Go 1.27's generic methods widen the options, though they can't implement interface methods. It must not need a global registry that changes after init (P2).
4. **Errors.** Which operations can stop returning errors once units are typed?
5. **Migration.** Keep the method names, so that most code compiles unchanged, and list what breaks, such as nil checks and string variables used as units.

Prototype on a branch and benchmark it against v0.1.x.

- **Accept when:**
  - hot paths take 5 ns or less with no allocations;
  - there are still no exported variables;
  - a custom measure gets the typed API in about 10 lines;
  - the README examples need only mechanical changes.
- **Outcome:** Path A (value types) or Path B (keep the interfaces) in Phase 3.
- **Status:** the RFC and the prototype are written: [docs/rfc-value-types.md](docs/rfc-value-types.md). It recommends Path A with integer units; the decision is the author's.

## Phase 3: the breaking batch (v0.2.0)

Each step is its own PR on `dev`, but nothing is tagged until the batch is complete, so users migrate once.

### 3.1 Typed unit symbols

This step is part of both paths. If the RFC chooses integer units, adapt the step.

- **Typed constants.** Each unit package gets `type Unit string` and typed constants, for example `const Kilometer Unit = "km"`.
- **Typed parameters.** Constructors and `ConvertTo` take the measure's `Unit` type, so `NewDistance(1, mass.Kilogram)` no longer compiles.
- **What breaks.** Untyped literals such as `NewDistance(1, "km")` still compile. String variables need a conversion: `distance.Unit(s)`.
- **What stays.** The generic `Measurement` methods keep string units. Deprecated aliases such as `PlankLength` become typed too.

### Path A: value types (recommended if the prototype meets the criteria)

#### 3.2A Value-type measures

Serves **P1** and **P2**.

- **The change.** Replace each measure interface with a struct type, one measure per PR.
- **Order.** Start with mass, which has no cross-measure operations, then time, velocity, area and distance. The code compiles at every point in between.
- **Generic code.** Each type keeps implementing `Measurement`, so generic code still works. Boxing allocates only on that generic path.
- **Afterwards.** This removes the "interface-based design" and "XxxMeasurement methods lose the type" debt. Lower the allocation limits from step 1.1 to 0.

#### 3.3A Typed custom measures

Serves **P3**.

- Implement the custom-measure design from the RFC, so user measures get the typed API (`Add`, `ConvertTo`, `Compare`, …) instead of the generic `Measurement`.
- Decide what happens to `NewBaseMeasurement`, which builds an N×N table on every call. Either keep it as the one-off path or deprecate it in favour of the creator.
- `TestMeasurementCreator_ExtendsPredefinedFactors` keeps passing, in its new form.

### Path B: keep the interfaces

#### 3.2B Typed results from generic operations

Serves **P3**. This change doesn't break the API: static types stay `Measurement`.

- **The change.** `ParseString`, `NewMeasurement`, `ConvertToMeasurement` and the `XxxMeasurement` methods return values whose dynamic type is the measure's interface. `m.(units.Distance)` and type switches then work, which they don't today.
- **How.** Add one helper that wraps a `baseMeasurement` by its measure, using an immutable map built at init (P2).
- **`NewBaseMeasurement`.** Decide its fate as in 3.3A.

### 3.4 Migration guide and release v0.2.0

- Write `MIGRATION.md` with before-and-after code for every break.
- Remove anything deprecated during v0.1.x.
- Add the CHANGELOG section and release as in step 0.4.

## Phase 4: features (v0.2.x)

Each step is non-breaking and ships on its own.

### 4.1 Text and JSON encoding

Serves **P3**.

- **The change.** Add `MarshalText` and `UnmarshalText` in the `"<value> <unit>"` form. That one change gives support for JSON, YAML, `flag.TextVar` and configuration loaders.
- **Decide in the step** whether to add a `{"value": …, "unit": …}` JSON form as well.
- **Why after Phase 3.** Unmarshalling needs a concrete type and can't target an interface.

### 4.2 Configurable parser

Serves **P3** and **P2**.

- **The parser.** `units.NewParser(options…)` returns a parser that is immutable once built.
- **Options:**
  - aliases, such as `metre`, `meters` and `kph`;
  - case-insensitive matching;
  - the form without a space (`10km`);
  - extra user units.
- **Defaults.** The package-level `ParseString` stays strict.
- **Immutability (P2).** Aliases live in the parser, not in the global tables.

### 4.3 New measures, one per step

Serves **P3**.

- **How.** Follow the "Add a new measure type" recipe in AGENTS.md.
- **Candidates**, in rough order of demand: volume, acceleration, force, energy, power, pressure, frequency, angle and data size. Each comes with its cross-measure operations, such as area × distance → volume and velocity ÷ time → acceleration.
- **Symbols.** Every symbol must stay unique across all measures (`TestMeasuresSelfCheck`). Check new symbols for clashes before starting a measure.
- **Status:** all nine are done, ahead of Phase 3 at the author's request: volume, acceleration, force, energy, power, pressure, frequency, angle and data size. Further candidates: angular velocity, data rate, density, electric units.

### 4.4 Decision: temperature

- The factor model is multiplicative only. Temperature needs an offset (°C, °F).
- It also needs to tell absolute temperatures apart from differences: 20 °C plus a 10 K difference is 30 °C, but adding two absolute temperatures is meaningless.
- Design the model for affine units, and the operations they allow, before writing any code.

### Not planned unless there is demand

- **Full dimensional analysis** (arbitrary unit products such as kg·m/s²) costs performance (P1) and simplicity.
- **Non-linear scales** such as dB and pH don't fit a conversion model of factors and offsets.

## Phase 5: v1.0.0

### 5.1 Stabilise and release v1.0.0

- **Clean up.** Remove the deprecated names: `PlankLength`, `PlankTime` and anything deprecated since.
- **Review the API.** Read the whole exported API with `go doc -all`, and make the step 1.2 `gorelease` check fail the build.
- **Then release.** After v1.0.0, a breaking change needs a new `/v2` module path. Complete Phase 3 before this step.

## Ongoing

- **Go version.** The minimum is Go 1.27. Raise it only when a feature needs it, and give each raise its own CHANGELOG entry. CI's stable test job shows early whether a new Go release breaks anything.
- **Tools.** Upgrade the pinned tools on purpose, one at a time, and handle new findings: golangci-lint, gorelease, benchstat and the CI actions.
- **Agent guidance.** Keep AGENTS.md in step with the code. Whenever a step changes how the library works, update the "How it works" and "Known defects" sections.
