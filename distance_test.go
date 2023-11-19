package units_test

import (
	"github.com/asaphin/go-physics-units"
	"github.com/asaphin/go-physics-units/distance"
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

	if d2.Value() != 0.001 {
		t.Errorf("value should be 0.001, but was %v", d2.Value())
	}
}
