package units

import (
	"github.com/asaphin/go-physics-units/acceleration"
	"github.com/asaphin/go-physics-units/force"
	"github.com/asaphin/go-physics-units/velocity"
)

// Acceleration is an acceleration: a value in an acceleration.Unit, such as 9.8 m/s². Create one with
// NewAcceleration or StringToAcceleration.
//
// An Acceleration is a small immutable value: every method returns a new value, so
// copy and share it freely. The zero Acceleration is 0 m/s².
type Acceleration struct {
	value float64
	unit  acceleration.Unit
}

// MultiplyByTime returns the velocity a·t (m/s) in [velocity.BaseUnit].
func (a Acceleration) MultiplyByTime(t Time) Velocity {
	return Velocity{value: a.baseValue() * t.baseValue(), unit: velocity.BaseUnit}
}

// MultiplyByMass returns the force m·a (N) in [force.BaseUnit].
func (a Acceleration) MultiplyByMass(m Mass) Force {
	return Force{value: a.baseValue() * m.baseValue(), unit: force.BaseUnit}
}
