package units

import (
	"github.com/asaphin/go-physics-units/angle"
	"github.com/asaphin/go-physics-units/angularvelocity"
)

// AngularVelocity is an angular velocity: a value in an angularvelocity.Unit, such as 33.3 rev/min. Create one with
// NewAngularVelocity or StringToAngularVelocity.
//
// An AngularVelocity is a small immutable value: every method returns a new value, so
// copy and share it freely. The zero AngularVelocity is 0 rad/s.
type AngularVelocity struct {
	value float64
	unit  angularvelocity.Unit
}

// MultiplyByTime returns the angle ω·t (rad) in angle.BaseUnit.
func (w AngularVelocity) MultiplyByTime(t Time) Angle {
	return Angle{value: w.baseValue() * t.baseValue(), unit: angle.BaseUnit}
}
