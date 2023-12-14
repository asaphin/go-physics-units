package units_test

import (
	units "github.com/asaphin/go-physics-units"
	"github.com/asaphin/go-physics-units/time"
	"github.com/asaphin/go-physics-units/velocity"
	"testing"
)

func TestVelocity_MultiplyByTime(t *testing.T) {
	v, err := units.NewVelocity(20, velocity.KilometerPerHour)
	if err != nil {
		t.Error(err)
	}

	tm, err := units.NewTime(30, time.Minute)
	if err != nil {
		t.Error(err)
	}

	d := v.MultiplyByTime(tm)

	if !almostEqual(d.Value(), 10000, 1e-3) {
		t.Errorf("should be 10000 m, but was %f", d.Value())
	}
}
