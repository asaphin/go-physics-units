package units_test

import (
	units "github.com/asaphin/go-physics-units"
	"github.com/asaphin/go-physics-units/area"
	"testing"
)

func TestArea_ConvertTo(t *testing.T) {
	a1, err := units.NewArea(1, area.SquareMeter)
	if err != nil {
		t.Error(err)
	}

	a2, err := a1.ConvertTo(area.SquareKilometer)
	if err != nil {
		t.Error(err)
	}

	if !almostEqual(a2.Value(), 0.000001, 0.001) {
		t.Errorf("value should be 0.000001, but was %v", a2.Value())
	}
}
