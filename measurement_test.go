package units_test

import (
	"testing"

	units "github.com/asaphin/go-physics-units"
	"github.com/asaphin/go-physics-units/distance"
	"github.com/asaphin/go-physics-units/mass"
)

func TestMeasurement_IsGenericUntilNarrowed(t *testing.T) {
	m, err := units.NewMeasurement(1, distance.Meter)
	if err != nil {
		t.Fatal(err)
	}

	if m.Type() != units.MeasureDistance {
		t.Errorf("Type() = %s, want %s", m.Type(), units.MeasureDistance)
	}

	if _, ok := m.(units.Distance); ok {
		t.Error("a generic Measurement shouldn't implement Distance; ToDistance must be used")
	}

	if _, err := units.ToDistance(m); err != nil {
		t.Errorf("ToDistance: %v", err)
	}
}

func TestNewMeasurement_UnknownUnit(t *testing.T) {
	m, err := units.NewMeasurement(1, "parsecs")
	if err == nil || m != nil {
		t.Errorf("NewMeasurement(1, parsecs) = %v, %v; want nil and an error", m, err)
	}
}

func TestDetectMeasureType(t *testing.T) {
	cases := map[string]units.MeasureType{
		"m":       units.MeasureDistance,
		"m²":      units.MeasureArea,
		"kg":      units.MeasureMass,
		"s":       units.MeasureTime,
		"m/s":     units.MeasureVelocity,
		"\u00b5s": units.MeasureTime,
		"":        units.MeasureCustom,
		"USD":     units.MeasureCustom,
	}

	for unit, want := range cases {
		if got := units.DetectMeasureType(unit); got != want {
			t.Errorf("DetectMeasureType(%q) = %s, want %s", unit, got, want)
		}
	}
}

func TestParseString(t *testing.T) {
	valid := []struct {
		in    string
		value float64
		unit  string
	}{
		{"10 km", 10, distance.Kilometer},
		{" 5   m ", 5, distance.Meter},
		{"1e3 m", 1000, distance.Meter},
		{"-2.5 kg", -2.5, mass.Kilogram},
		{"0 lb", 0, mass.Pound},
		{"1 \u00b5m", 1, distance.Micrometer},
	}

	for _, c := range valid {
		m, err := units.ParseString(c.in)
		if err != nil {
			t.Errorf("ParseString(%q): %v", c.in, err)

			continue
		}

		assertMeasurement(t, m, c.value, c.unit)
	}

	invalid := []string{
		"", "10", "10km", "km 10", "1 2 m", "x m", "NaN m", "+Inf m", "-Inf m", "1e400 m", "10 parsecs",
	}

	for _, in := range invalid {
		if m, err := units.ParseString(in); err == nil || m != nil {
			t.Errorf("ParseString(%q) = %v, %v; want nil and an error", in, m, err)
		}
	}
}

func TestMeasurement_Arithmetic(t *testing.T) {
	km, err := units.NewMeasurement(1, distance.Kilometer)
	if err != nil {
		t.Fatal(err)
	}

	m500, err := units.NewMeasurement(500, distance.Meter)
	if err != nil {
		t.Fatal(err)
	}

	sum, err := km.AddMeasurement(m500)
	if err != nil {
		t.Fatal(err)
	}

	assertMeasurement(t, sum, 1.5, distance.Kilometer)

	diff, err := km.SubMeasurement(m500)
	if err != nil {
		t.Fatal(err)
	}

	assertMeasurement(t, diff, 0.5, distance.Kilometer)

	assertMeasurement(t, km.MulMeasurement(3), 3, distance.Kilometer)

	quot, err := km.DivMeasurement(4)
	if err != nil {
		t.Fatal(err)
	}

	assertMeasurement(t, quot, 0.25, distance.Kilometer)

	if q, err := km.DivMeasurement(0); err == nil || q != nil {
		t.Errorf("DivMeasurement(0) = %v, %v; want nil and an error", q, err)
	}

	kg, err := units.NewMeasurement(1, mass.Kilogram)
	if err != nil {
		t.Fatal(err)
	}

	if s, err := km.AddMeasurement(kg); err == nil || s != nil {
		t.Errorf("adding mass to distance = %v, %v; want nil and an error", s, err)
	}

	if s, err := km.SubMeasurement(kg); err == nil || s != nil {
		t.Errorf("subtracting mass from distance = %v, %v; want nil and an error", s, err)
	}

	if s, err := km.AddMeasurement(nil); err == nil || s != nil {
		t.Errorf("adding nil = %v, %v; want nil and an error", s, err)
	}

	assertMeasurement(t, km, 1, distance.Kilometer)
}

func TestMeasurement_ConvertToMeasurement(t *testing.T) {
	km, err := units.NewMeasurement(1, distance.Kilometer)
	if err != nil {
		t.Fatal(err)
	}

	m, err := km.ConvertToMeasurement(distance.Meter)
	if err != nil {
		t.Fatal(err)
	}

	assertMeasurement(t, m, 1000, distance.Meter)

	if bad, err := km.ConvertToMeasurement(mass.Kilogram); err == nil || bad != nil {
		t.Errorf("ConvertToMeasurement(kg) = %v, %v; want nil and an error", bad, err)
	}
}
