package units_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"testing"

	units "github.com/asaphin/go-physics-units"
	"github.com/asaphin/go-physics-units/temperature"
)

func mustTemperature(t *testing.T, value float64, unit temperature.Unit) units.Temperature {
	t.Helper()

	tm, err := units.NewTemperature(value, unit)
	if err != nil {
		t.Fatal(err)
	}

	return tm
}

// TestTemperature_Conversions checks the claims of docs/rfc-temperature.md:
// conversions between the defined points are exact.
func TestTemperature_Conversions(t *testing.T) {
	cases := []struct {
		value    float64
		from, to temperature.Unit
		want     float64
		exact    bool
	}{
		{100, temperature.Celsius, temperature.Fahrenheit, 212, true},
		{0, temperature.Celsius, temperature.Fahrenheit, 32, true},
		{-40, temperature.Celsius, temperature.Fahrenheit, -40, true},
		{0, temperature.Celsius, temperature.Kelvin, 273.15, true},
		{0, temperature.Kelvin, temperature.Celsius, -273.15, true},
		{0, temperature.Kelvin, temperature.Fahrenheit, -459.67, true},
		{0, temperature.Rankine, temperature.Kelvin, 0, true},
		{491.67, temperature.Rankine, temperature.Kelvin, 273.15, false},
		{37, temperature.Celsius, temperature.Fahrenheit, 98.6, false},
		{212, temperature.Fahrenheit, temperature.Celsius, 100, false},
		{300, temperature.Kelvin, temperature.Rankine, 540, false},
		{20, temperature.Celsius, temperature.Celsius, 20, true},
	}

	for _, c := range cases {
		got := mustTemperature(t, c.value, c.from).ConvertTo(c.to)

		if got.Unit() != c.to || (c.exact && got.Value() != c.want) || !almostEqual(got.Value(), c.want, 1e-15) {
			t.Errorf("%v %v in %v = %v, want %v (exact: %v)", c.value, c.from, c.to, got, c.want, c.exact)
		}
	}
}

func TestTemperature_RoundTrips(t *testing.T) {
	for from := range temperature.Units() {
		for to := range temperature.Units() {
			for _, v := range []float64{0, 1, 36.6, 1e3, 1e6} {
				x := mustTemperature(t, v+from.AbsoluteZero()+500, from)

				if back := x.ConvertTo(to).ConvertTo(from); !back.EqualWithinRel(x, 1e-12) {
					t.Errorf("%v → %v → %v", x, x.ConvertTo(to), back)
				}
			}
		}
	}
}

func TestTemperature_AbsoluteZero(t *testing.T) {
	for u := range temperature.Units() {
		zero := mustTemperature(t, u.AbsoluteZero(), u)
		if zero.Kelvin() != 0 {
			t.Errorf("absolute zero in %v is %v K, want 0", u, zero.Kelvin())
		}

		// Converting absolute zero never rounds below it.
		for to := range temperature.Units() {
			if got := zero.ConvertTo(to); got.Value() < to.AbsoluteZero() {
				t.Errorf("%v in %v = %v, below absolute zero", zero, to, got)
			}
		}

		below := math.Nextafter(u.AbsoluteZero(), math.Inf(-1))
		for _, v := range []float64{below, u.AbsoluteZero() - 1, math.Inf(-1), math.NaN()} {
			if tm, err := units.NewTemperature(
				v,
				u,
			); !errors.Is(err, units.ErrOutOfRange) ||
				tm != (units.Temperature{}) {
				t.Errorf("NewTemperature(%v, %v) = %v, %v; want the zero value and ErrOutOfRange", v, u, tm, err)
			}
		}
	}

	if _, err := units.StringToTemperature("-300 °C"); !errors.Is(err, units.ErrOutOfRange) {
		t.Errorf("StringToTemperature(-300 °C): %v, want ErrOutOfRange", err)
	}
}

func TestTemperature_Arithmetic(t *testing.T) {
	room := mustTemperature(t, 20, temperature.Celsius)
	warm := mustTemperature(t, 50, temperature.Fahrenheit) // 10 °C

	diff := room.Sub(warm)
	assertMeasurement(t, diff, 10, "°C")

	if back := warm.Sub(room); !almostEqual(back.ConvertTo(temperature.Celsius).Value(), -10, 1e-12) {
		t.Errorf("50 °F - 20 °C = %v, want -18 °F", back)
	}

	up, err := room.Add(units.NewTemperatureDifference(18, temperature.Fahrenheit))
	if err != nil {
		t.Fatal(err)
	}

	assertMeasurement(t, up, 30, "°C")

	down, err := room.SubDifference(units.NewTemperatureDifference(20, temperature.Kelvin))
	if err != nil {
		t.Fatal(err)
	}

	assertMeasurement(t, down, 0, "°C")

	tooCold := units.NewTemperatureDifference(300, temperature.Kelvin)
	if got, err := room.SubDifference(tooCold); !errors.Is(err, units.ErrOutOfRange) || got != (units.Temperature{}) {
		t.Errorf("20 °C - 300 K = %v, %v; want ErrOutOfRange", got, err)
	}

	if got, err := room.Add(tooCold.Mul(-1)); !errors.Is(err, units.ErrOutOfRange) || got != (units.Temperature{}) {
		t.Errorf("20 °C + (-300 K) = %v, %v; want ErrOutOfRange", got, err)
	}

	// The ideal gas law needs absolute temperatures.
	if k := room.Kelvin(); k != 293.15 {
		t.Errorf("Kelvin() = %v, want 293.15", k)
	}
}

func TestTemperature_Compare(t *testing.T) {
	freezing := mustTemperature(t, 0, temperature.Celsius)
	f32 := mustTemperature(t, 32, temperature.Fahrenheit)
	body := mustTemperature(t, 98.6, temperature.Fahrenheit)

	if freezing.Compare(body) != -1 || body.Compare(freezing) != 1 || body.Compare(body.ConvertTo(body.Unit())) != 0 {
		t.Error("Compare orders temperatures wrongly")
	}

	if !freezing.EqualWithinRel(f32, 1e-12) || !f32.EqualWithinRel(freezing, 1e-12) {
		t.Errorf("0 °C and 32 °F aren't equal within 1e-12: %v K and %v K", freezing.Kelvin(), f32.Kelvin())
	}

	if freezing.Equal(f32) != f32.Equal(freezing) {
		t.Error("Equal isn't symmetric")
	}

	if !freezing.Equal(freezing.ConvertTo(temperature.Celsius)) {
		t.Error("a temperature doesn't equal itself")
	}
}

func TestTemperatureDifference(t *testing.T) {
	rise := units.NewTemperatureDifference(10, temperature.Celsius)

	assertMeasurement(t, rise.ConvertTo(temperature.Fahrenheit), 18, "°F")
	assertMeasurement(t, rise.ConvertTo(temperature.Kelvin), 10, "K")
	assertMeasurement(t, rise.ConvertToBaseUnits(), 10, "K")
	assertMeasurement(t, rise.Add(units.NewTemperatureDifference(9, temperature.Fahrenheit)), 15, "°C")
	assertMeasurement(t, rise.Sub(units.NewTemperatureDifference(5, temperature.Kelvin)), 5, "°C")
	assertMeasurement(t, rise.Mul(3), 30, "°C")

	half, err := rise.Div(2)
	if err != nil {
		t.Fatal(err)
	}

	assertMeasurement(t, half, 5, "°C")

	if bad, err := rise.Div(0); !errors.Is(err, units.ErrZeroDivision) || bad != (units.TemperatureDifference{}) {
		t.Errorf("Div(0) = %v, %v; want ErrZeroDivision", bad, err)
	}

	f18 := units.NewTemperatureDifference(18, temperature.Fahrenheit)
	if !rise.EqualWithinRel(f18, 1e-12) || rise.Compare(rise.Mul(2)) != -1 || !rise.Equal(rise.Mul(1)) {
		t.Errorf("comparisons: %v vs %v", rise, f18)
	}

	if rise.Compare(f18) != -f18.Compare(rise) {
		t.Error("Compare isn't antisymmetric")
	}

	var zero units.TemperatureDifference
	if zero.Value() != 0 || zero.Unit() != temperature.Kelvin || zero.Symbol() != "K" {
		t.Errorf("the zero TemperatureDifference is %v", zero)
	}
}

func TestTemperature_Parsing(t *testing.T) {
	for s, want := range map[string]units.Temperature{
		"20 °C":     mustTemperature(t, 20, temperature.Celsius),
		"20 \u2103": mustTemperature(t, 20, temperature.Celsius),
		"68 \u2109": mustTemperature(t, 68, temperature.Fahrenheit),
		"300 K":     mustTemperature(t, 300, temperature.Kelvin),
		"540 °R":    mustTemperature(t, 540, temperature.Rankine),
	} {
		if got, err := units.StringToTemperature(s); err != nil || got != want {
			t.Errorf("StringToTemperature(%q) = %v, %v; want %v", s, got, err, want)
		}
	}

	for s, want := range map[string]error{
		"20 C":  units.ErrUnknownUnit,
		"20 °K": units.ErrUnknownUnit,
		"20°C":  units.ErrInvalidFormat,
		"NaN K": units.ErrInvalidFormat,
	} {
		if _, err := units.StringToTemperature(s); !errors.Is(err, want) {
			t.Errorf("StringToTemperature(%q): %v, want %v", s, err, want)
		}

		if _, err := units.StringToTemperatureDifference(s); !errors.Is(err, want) {
			t.Errorf("StringToTemperatureDifference(%q): %v, want %v", s, err, want)
		}
	}

	d, err := units.StringToTemperatureDifference("-5 °F")
	if err != nil || d != units.NewTemperatureDifference(-5, temperature.Fahrenheit) {
		t.Errorf("StringToTemperatureDifference(-5 °F) = %v, %v", d, err)
	}

	// Temperatures aren't Measurements: the generic parsers say so.
	if _, err := units.ParseString("20 °C"); !errors.Is(err, units.ErrIncompatibleMeasures) {
		t.Errorf("ParseString(20 °C): %v, want ErrIncompatibleMeasures", err)
	}

	if _, err := mustParser(t).Parse("20 K"); !errors.Is(err, units.ErrIncompatibleMeasures) {
		t.Errorf("Parser.Parse(20 K): %v, want ErrIncompatibleMeasures", err)
	}

	if mt := units.DetectMeasureType("°F"); mt != units.MeasureTemperature {
		t.Errorf("DetectMeasureType(°F) = %s, want temperature", mt)
	}

	if _, err := units.NewUnitSet(
		"temperature",
		"K",
		map[string]float64{"K": 1},
	); !errors.Is(
		err,
		units.ErrInvalidMeasure,
	) {
		t.Errorf("a custom measure named temperature: %v, want ErrInvalidMeasure", err)
	}
}

func TestTemperature_Text(t *testing.T) {
	type reading struct {
		Air  units.Temperature           `json:"air"`
		Rise units.TemperatureDifference `json:"rise"`
	}

	in := reading{mustTemperature(t, 21.5, temperature.Celsius), units.NewTemperatureDifference(2, temperature.Kelvin)}

	b, err := json.Marshal(in)
	if err != nil {
		t.Fatal(err)
	}

	if string(b) != `{"air":"21.5 °C","rise":"2 K"}` || in.Air.String() != "21.5 °C" || in.Rise.String() != "2 K" {
		t.Errorf("json = %s", b)
	}

	var out reading
	if err := json.Unmarshal(b, &out); err != nil || out != in {
		t.Errorf("round trip = %v, %v; want %v", out, err, in)
	}

	if err := json.Unmarshal([]byte(`{"air":"-500 °C"}`), &out); !errors.Is(err, units.ErrOutOfRange) {
		t.Errorf("unmarshalling -500 °C: %v, want ErrOutOfRange", err)
	}

	if err := json.Unmarshal([]byte(`{"rise":"2 furlong"}`), &out); !errors.Is(err, units.ErrUnknownUnit) {
		t.Errorf("unmarshalling 2 furlong: %v, want ErrUnknownUnit", err)
	}

	invalid := units.NewTemperatureDifference(1, temperature.Unit(9))
	if _, err := invalid.MarshalText(); !errors.Is(err, units.ErrUnknownUnit) {
		t.Errorf("MarshalText in an invalid unit: %v, want ErrUnknownUnit", err)
	}

	if tm, err := units.NewTemperature(1, temperature.Unit(9)); !errors.Is(err, units.ErrUnknownUnit) {
		t.Errorf("NewTemperature in an invalid unit = %v, %v; want ErrUnknownUnit", tm, err)
	}

	hot := mustTemperature(t, math.Inf(1), temperature.Kelvin)
	if _, err := hot.MarshalText(); !errors.Is(err, units.ErrOutOfRange) {
		t.Errorf("MarshalText of +Inf K: %v, want ErrOutOfRange", err)
	}

	if got := fmt.Sprintf(
		"%.1f|%v|%#v",
		in.Air,
		in.Rise,
		in.Air,
	); got != `21.5 °C|2 K|units.Temperature{Value:21.5, Unit:"°C"}` {
		t.Errorf("Sprintf = %q", got)
	}

	if got := fmt.Sprintf("%#v", in.Rise); got != `units.TemperatureDifference{Value:2, Unit:"K"}` {
		t.Errorf("Sprintf(%%#v) = %q", got)
	}
}

func TestTemperature_UnitPackage(t *testing.T) {
	n := 0

	for u := range temperature.Units() {
		if temperature.Unit(n) != u || !u.Valid() {
			t.Errorf("unit %d is %v", n, u)
		}

		n++

		if parsed, ok := temperature.ParseUnit(u.String()); !ok || parsed != u {
			t.Errorf("ParseUnit(%q) = %v, %v", u.String(), parsed, ok)
		}

		if r := temperature.DifferenceRate(u, u); r != 1 {
			t.Errorf("DifferenceRate(%v, %v) = %v, want 1", u, u, r)
		}
	}

	if n != 4 {
		t.Errorf("%d units, want 4", n)
	}

	for range temperature.Units() {
		break // the iterator stops when asked to
	}

	invalid := temperature.Unit(9)
	if invalid.Valid() || invalid.String() != "temperature.Unit(9)" {
		t.Errorf("invalid unit: Valid %v, String %q", invalid.Valid(), invalid.String())
	}

	if _, ok := temperature.ParseUnit("C"); ok {
		t.Error("ParseUnit(C) found a unit")
	}

	if r := temperature.DifferenceRate(temperature.Celsius, temperature.Fahrenheit); r != 1.8 {
		t.Errorf("DifferenceRate(°C, °F) = %v, want 1.8", r)
	}

	if got := temperature.Convert(100, temperature.Celsius, temperature.Fahrenheit); got != 212 {
		t.Errorf("Convert(100 °C, °F) = %v, want exactly 212", got)
	}
}
