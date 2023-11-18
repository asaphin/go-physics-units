package units_test

import (
	units "github.com/asaphin/go-physics-units"
	"testing"
)

func TestBaseMeasurement(t *testing.T) {
	ln, err := units.NewBaseMeasurement(1, "m", units.DistanceConversionFactors())
	if err != nil {
		t.Error(err)
	}

	ln2, err := ln.ConvertTo("km")
	if err != nil {
		t.Error(err)
	}

	ln2Value := ln2.Value()

	if ln2Value != 1000 {
		t.Errorf("value after konversion should be 1000, but was %v", ln2Value)
	}
}
