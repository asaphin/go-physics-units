package units_test

import (
	"errors"
	"testing"

	units "github.com/asaphin/go-physics-units"
	"github.com/asaphin/go-physics-units/distance"
	"github.com/asaphin/go-physics-units/mass"
)

func TestNewMeasurement_HasTheConcreteType(t *testing.T) {
	m, err := units.NewMeasurement(1, "km")
	if err != nil {
		t.Fatal(err)
	}

	if m.Type() != units.MeasureDistance {
		t.Errorf("Type() = %s, want %s", m.Type(), units.MeasureDistance)
	}

	if d, ok := m.(units.Distance); !ok || d != units.NewDistance(1, distance.Kilometer) {
		t.Errorf("NewMeasurement(1, km) holds %T %v, want the Distance 1 km", m, m)
	}

	// The generic and typed APIs agree.
	if sum, err := m.AddMeasurement(units.NewDistance(500, distance.Meter)); err != nil ||
		sum != units.NewDistance(1.5, distance.Kilometer) {
		t.Errorf("1 km + 500 m = %v, %v; want the Distance 1.5 km", sum, err)
	}
}

func TestNewMeasurement_UnknownUnit(t *testing.T) {
	m, err := units.NewMeasurement(1, "parsecs")
	if !errors.Is(err, units.ErrUnknownUnit) || m != nil {
		t.Errorf("NewMeasurement(1, parsecs) = %v, %v; want nil and ErrUnknownUnit", m, err)
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
		in   string
		want units.Measurement
	}{
		{"10 km", units.NewDistance(10, distance.Kilometer)},
		{" 5   m ", units.NewDistance(5, distance.Meter)},
		{"1e3 m", units.NewDistance(1000, distance.Meter)},
		{"-2.5 kg", units.NewMass(-2.5, mass.Kilogram)},
		{"0 lb", units.NewMass(0, mass.Pound)},
		{"1 \u00b5m", units.NewDistance(1, distance.Micrometer)},
	}

	for _, c := range valid {
		m, err := units.ParseString(c.in)
		if err != nil || m != c.want {
			t.Errorf("ParseString(%q) = %v, %v; want %v", c.in, m, err, c.want)
		}
	}

	invalid := map[string]error{
		"":           units.ErrInvalidFormat,
		"10":         units.ErrInvalidFormat,
		"10km":       units.ErrInvalidFormat,
		"km 10":      units.ErrInvalidFormat,
		"1 2 m":      units.ErrInvalidFormat,
		"x m":        units.ErrInvalidFormat,
		"NaN m":      units.ErrInvalidFormat,
		"+Inf m":     units.ErrInvalidFormat,
		"-Inf m":     units.ErrInvalidFormat,
		"1e400 m":    units.ErrInvalidFormat,
		"10 parsecs": units.ErrUnknownUnit,
	}

	for in, want := range invalid {
		if m, err := units.ParseString(in); !errors.Is(err, want) || m != nil {
			t.Errorf("ParseString(%q) = %v, %v; want nil and %v", in, m, err, want)
		}
	}
}
