package units_test

import (
	"github.com/asaphin/go-physics-units"
	"github.com/asaphin/go-physics-units/distance"
	"github.com/asaphin/go-physics-units/time"
	"github.com/asaphin/go-physics-units/velocity"
	"testing"
)

func TestDistanceConversionFactors_Safety(t *testing.T) {
	dcf := distance.ConversionFactors()

	dcf["km"] = 2000

	otherDcf := distance.ConversionFactors()

	if otherDcf["km"] == 2000 {
		t.Error("ConversionFactors function isn't safe")
	}
}

func TestDistance_ConvertTo(t *testing.T) {
	d1, err := units.NewDistance(1, distance.Meter)
	if err != nil {
		t.Error(err)
	}

	d2, err := d1.ConvertTo(distance.Kilometer)
	if err != nil {
		t.Error(err)
	}

	if !almostEqual(d2.Value(), 0.001, 0.001) {
		t.Errorf("value should be 0.001, but was %v", d2.Value())
	}
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
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			d, err := units.NewDistance(c.initialDistance, c.initialUnit)
			if err != nil {
				t.Error(err)
			}

			adder, err := units.NewDistance(c.adderDistance, c.adderUnit)
			if err != nil {
				t.Error(err)
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
		t.Error(err)
	}

	tm, err := units.NewTime(2, time.Hour)
	if err != nil {
		t.Error(err)
	}

	vel := dist.DivideByTime(tm)
	if vel.Unit() != velocity.BaseUnit {
		t.Errorf("result velocity should have base SI unit %s", velocity.BaseUnit)
	}

	var expectedValue = 13.89

	if !almostEqual(expectedValue, vel.Value(), 1e-3) {
		t.Errorf("actual velocity value %v, but expected %v", vel.Value(), expectedValue)
	}
}
