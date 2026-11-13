package units_test

import (
	"testing"

	units "github.com/asaphin/go-physics-units"
	"github.com/asaphin/go-physics-units/distance"
	"github.com/asaphin/go-physics-units/time"
	"github.com/asaphin/go-physics-units/velocity"
)

func TestVelocity_MultiplyByTime(t *testing.T) {
	v, err := units.NewVelocity(20, velocity.KilometerPerHour)
	if err != nil {
		t.Fatal(err)
	}

	tm, err := units.NewTime(30, time.Minute)
	if err != nil {
		t.Fatal(err)
	}

	assertMeasurement(t, v.MultiplyByTime(tm), 10000, distance.BaseUnit)
}
