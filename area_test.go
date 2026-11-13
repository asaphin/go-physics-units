package units_test

import (
	"testing"

	units "github.com/asaphin/go-physics-units"
	"github.com/asaphin/go-physics-units/area"
	"github.com/asaphin/go-physics-units/distance"
)

func TestArea_ConvertTo(t *testing.T) {
	a1, err := units.NewArea(1, area.SquareMeter)
	if err != nil {
		t.Fatal(err)
	}

	a2, err := a1.ConvertTo(area.SquareKilometer)
	if err != nil {
		t.Fatal(err)
	}

	assertMeasurement(t, a2, 0.000001, area.SquareKilometer)
}

func TestArea_DivideByDistance(t *testing.T) {
	a, err := units.NewArea(1, area.SquareKilometer)
	if err != nil {
		t.Fatal(err)
	}

	d, err := units.NewDistance(500, distance.Meter)
	if err != nil {
		t.Fatal(err)
	}

	assertMeasurement(t, a.DivideByDistance(d), 2000, distance.BaseUnit)

	zero, err := units.NewDistance(0, distance.Meter)
	if err != nil {
		t.Fatal(err)
	}

	assertPanics(t, func() { a.DivideByDistance(zero) })
}
