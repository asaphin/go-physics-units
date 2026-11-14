package units

import (
	"github.com/asaphin/go-physics-units/angle"
	"github.com/asaphin/go-physics-units/angularvelocity"
)

// Angle is an angle: a value in an angle.Unit, such as 90 °. Create one with
// NewAngle or StringToAngle.
//
// An Angle is a small immutable value: every method returns a new value, so
// copy and share it freely. The zero Angle is 0 rad.
type Angle struct {
	value float64
	unit  angle.Unit
}

// DivideByTime returns the angular velocity a/t (rad/s) in [angularvelocity.BaseUnit]. It fails with
// [ErrZeroDivision] if t is zero.
func (a Angle) DivideByTime(t Time) (AngularVelocity, error) {
	if t.value == 0 {
		return AngularVelocity{}, ErrZeroDivision
	}

	return AngularVelocity{value: a.baseValue() / t.baseValue(), unit: angularvelocity.BaseUnit}, nil
}
