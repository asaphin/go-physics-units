package units

import (
	"github.com/asaphin/go-physics-units/acceleration"
	"github.com/asaphin/go-physics-units/distance"
	"github.com/asaphin/go-physics-units/velocity"
)

// Velocity is a speed: a value in a velocity.Unit, such as 10 km/h. Create one with
// NewVelocity or StringToVelocity.
//
// A Velocity is a small immutable value: every method returns a new value, so
// copy and share it freely. The zero Velocity is 0 m/s.
type Velocity struct {
	value float64
	unit  velocity.Unit
}

// MultiplyByTime returns the distance v·t (m) in distance.BaseUnit.
func (v Velocity) MultiplyByTime(t Time) Distance {
	return Distance{value: v.baseValue() * t.baseValue(), unit: distance.BaseUnit}
}

// DivideByTime returns the acceleration v/t (m/s²) in acceleration.BaseUnit. It fails with
// ErrZeroDivision if t is zero.
func (v Velocity) DivideByTime(t Time) (Acceleration, error) {
	if t.value == 0 {
		return Acceleration{}, ErrZeroDivision
	}

	return Acceleration{value: v.baseValue() / t.baseValue(), unit: acceleration.BaseUnit}, nil
}
