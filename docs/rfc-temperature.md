# RFC: temperature

ROADMAP step 4.4. **Status:** accepted on 2026-11-14 and implemented (`temperature.go`, `temperature/`). Written 2026-11-13.

> **Decisions (2026-11-14):** the proposal as written, except decision 2: temperatures below absolute zero are **rejected** with `ErrOutOfRange`, by the constructors, the parsers, `Temperature.Add` and `SubDifference`. NaN is rejected too. `ConvertTo` clamps rounding at absolute zero, so a valid temperature never converts below it. `ParseString` rejects temperature symbols with `ErrIncompatibleMeasures` and a message that points to `StringToTemperature`.

## Why temperature needs its own design

Every measure so far converts by multiplication: one km is 1000 m. Temperature breaks that model in two ways.

1. **Offsets.** Celsius and Fahrenheit are affine scales: 0 °C is 273.15 K, and 0 °F is 459.67 × 5/9 K. A single factor can't express this. The table needs a scale and an offset per unit.
2. **Two kinds of quantity.** A temperature (20 °C) and a temperature *difference* (a rise of 10 °C) share their symbols but behave differently:
   - 20 °C + 10 °C (a difference) = 30 °C, but 20 °C + 30 °C (two temperatures) is meaningless;
   - 20 °C is 68 °F, but a 20 °C *rise* is a 36 °F rise, with no offset;
   - 40 °C isn't "twice as hot" as 20 °C, but 40 K is twice 20 K.

The roadmap therefore says to settle the design before writing any code.

## Proposal

### Two types

| | `Temperature` (a point) | `TemperatureDifference` (an interval) |
|---|---|---|
| Units | K, °C, °F, °R | K, °C, °F, °R (the same symbols) |
| Conversion | affine: `x·a + b` | linear: `x·a` |
| 20 °C in °F | 68 °F | 36 °F |
| Base unit | K | K |

The operations follow the rules of an affine space. Each one is typed, so the meaningless operations don't compile:

| Operation | Result |
|---|---|
| `Temperature.Sub(Temperature)` | `TemperatureDifference` |
| `Temperature.Add(TemperatureDifference)`, `Temperature.SubDifference(TemperatureDifference)` | `Temperature` |
| `TemperatureDifference.Add/Sub(TemperatureDifference)`, `.Mul(float64)`, `.Div(float64)` | `TemperatureDifference` |
| `Temperature.Compare/Equal/EqualWithinRel` | in kelvin |
| `Temperature.Kelvin() float64` | for physics that needs ratios of absolute temperatures, such as the ideal gas law |
| `Temperature + Temperature`, `Temperature × 2` | **not offered** |

### Units and precision

| Unit | Symbol | Scale to K | Offset to K |
|---|---|---|---|
| kelvin | `K` | 1 | 0 |
| degree Celsius | `°C` | 1 | 273.15 |
| degree Fahrenheit | `°F` | 5/9 | 459.67 × 5/9 |
| degree Rankine | `°R` | 5/9 | 0 |

**Pair coefficients.** Going through kelvin at run time rounds several times: 100 °C → 373.15 K → 211.99999999999991 °F, and 0 °C gives 31.999999999999936 °F. Instead, compute the direct coefficients for each pair (`°C→°F: ×9/5 + 32`) at init, with `math/big.Rat`, from the exact rational definitions above. Each coefficient then rounds once: 100 °C gives exactly 212 °F, and 0 °C exactly 32 °F. Values whose decimal isn't exact in binary still round in the last bit (37 °C gives 98.600000000000009 °F), as with every measure. These numbers were checked with Go 1.27. This makes table construction slower, which is acceptable (P1: calculations come first). A conversion stays one multiplication and one addition.

**Symbols.** `°C` and `°F` don't clash with the angle's `°`, because symbols are whole strings. U+2103 `℃` and U+2109 `℉` should be accepted and canonicalised to `°C` and `°F`, the way the micro sign is.

**Generic parsing.** `ParseString("20 °C")` yields a `Temperature`, the common case. Differences come only from `NewTemperatureDifference`, `StringToTemperatureDifference` and `Temperature.Sub`. The generic index can't map one symbol to two measures (`TestMeasuresSelfCheck`).

## Fitting the core design (ROADMAP 2.1)

- **Path A (value types).** Two plain structs, `Temperature` and `TemperatureDifference`, each with a typed `temperature.Unit`, and an affine table of unexported arrays. They are allocation-free like the other measures. This is the natural fit, so **implement temperature after Phase 3**.
- **Path B (interfaces).** It works too, but `Temperature` would have to implement `Measurement`. Then `AddMeasurement` of two temperatures must fail at run time with a new sentinel, such as `ErrUndefinedOperation`, which is the kind of error typed APIs exist to avoid.

Either way, temperature needs its own small table type (`internal/rates.AffineTable`). It must not stretch `rates.Table`: every other measure relies on its rates being plain ratios, including `Compare` (through `Factor`) and the cross-measure operations.

## Decisions for the author

1. **Two types** (proposed) or a single `Temperature` whose operations take the absolute or difference reading from context?
2. **Below absolute zero:** accept any value, like the other measures, so calculations never fail (proposed)? Or reject it in constructors with `ErrOutOfRange`, which `Temperature.SubDifference` would then also have to return?
3. **`Measurement` membership:** should `Temperature` implement the generic interface, and have the undefined operations fail? Or stay outside it, implementing only `Stringer` and `Formatter` (proposed under Path A)?
4. **When:** after Phase 3 (proposed) or now, in the interface design?
5. **`℃` and `℉` canonicalisation** (proposed: yes).
