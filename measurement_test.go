package units_test

import (
	units "github.com/asaphin/go-physics-units"
	"testing"
)

func TestMeasurement_ImplicitConversion(t *testing.T) {
	ln, err := units.NewMeasurement(1, "m")
	if err != nil {
		t.Error(err)
	}

	if ln.Type() != units.MeasureDistance {
		t.Error("not distance measure type")
	}

	_, ok := ln.(*units.Distance)
	if ok {
		t.Error("implicit conversion should be successful")
	}
}
