package units_test

import (
	"testing"

	units "github.com/asaphin/go-physics-units"
	"github.com/asaphin/go-physics-units/area"
	"github.com/asaphin/go-physics-units/distance"
	"github.com/asaphin/go-physics-units/time"
	"github.com/asaphin/go-physics-units/velocity"
)

func TestDistance_ConvertTo(t *testing.T) {
	d1, err := units.NewDistance(1, distance.Meter)
	if err != nil {
		t.Fatal(err)
	}

	d2, err := d1.ConvertTo(distance.Kilometer)
	if err != nil {
		t.Fatal(err)
	}

	assertMeasurement(t, d2, 0.001, distance.Kilometer)
}

func TestDistance_Add(t *testing.T) {
	epsilon := 1e-3

	cases := []struct {
		name            string
		initialDistance float64
		initialUnit     string
		adderDistance   float64
		adderUnit       string
		resultDistance  float64
	}{
		{
			name:            "adding of one decimeter to one meter",
			initialDistance: 1,
			initialUnit:     distance.Meter,
			adderDistance:   1,
			adderUnit:       distance.Decimeter,
			resultDistance:  1.1,
		},
		{
			name:            "adding of two kilometers to thousand meters",
			initialDistance: 1000,
			initialUnit:     distance.Meter,
			adderDistance:   2,
			adderUnit:       distance.Kilometer,
			resultDistance:  3000,
		},
		{
			name:            "adding of one km mile equivalent to one kilometer",
			initialDistance: 1,
			initialUnit:     distance.Kilometer,
			adderDistance:   0.621371,
			adderUnit:       distance.Mile,
			resultDistance:  2,
		},
		{
			name:            "adding of one millimeter to one meter",
			initialDistance: 1,
			initialUnit:     distance.Meter,
			adderDistance:   1,
			adderUnit:       distance.Millimeter,
			resultDistance:  1.001,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			d, err := units.NewDistance(c.initialDistance, c.initialUnit)
			if err != nil {
				t.Fatal(err)
			}

			adder, err := units.NewDistance(c.adderDistance, c.adderUnit)
			if err != nil {
				t.Fatal(err)
			}

			res := d.Add(adder)

			if res.Type() != units.MeasureDistance {
				t.Errorf("should be distance, but was %s", res.Type())
			}

			if res.Unit() != c.initialUnit {
				t.Errorf("should be %s, but was %s", c.initialUnit, res.Unit())
			}

			if !almostEqual(res.Value(), c.resultDistance, epsilon) {
				t.Errorf("value should be %f, but was %f", c.resultDistance, res.Value())
			}
		})
	}
}

func TestDistance_DivideByTime(t *testing.T) {
	dist, err := units.NewDistance(100, distance.Kilometer)
	if err != nil {
		t.Fatal(err)
	}

	tm, err := units.NewTime(2, time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	vel, err := dist.DivideByTime(tm)
	if err != nil {
		t.Fatal(err)
	}

	assertMeasurement(t, vel, 100000.0/7200, velocity.BaseUnit)
	assertMeasurement(t, vel.MustConvertTo(velocity.KilometerPerHour), 50, velocity.KilometerPerHour)

	zero, err := units.NewTime(0, time.Second)
	if err != nil {
		t.Fatal(err)
	}

	if v, err := dist.DivideByTime(zero); err == nil || v != nil {
		t.Errorf("DivideByTime(0 s) = %v, %v; want nil and an error", v, err)
	}
}

func TestDistance_Square(t *testing.T) {
	d, err := units.NewDistance(3, distance.Kilometer)
	if err != nil {
		t.Fatal(err)
	}

	a := d.Square()

	if a.Type() != units.MeasureArea {
		t.Errorf("Type() = %s, want %s", a.Type(), units.MeasureArea)
	}

	assertMeasurement(t, a, 9e6, area.BaseUnit)
	assertMeasurement(t, a.MustConvertTo(area.SquareKilometer), 9, area.SquareKilometer)
}
