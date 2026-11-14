package units

import (
	"github.com/asaphin/go-physics-units/angle"
)

// Angle is an angle: a value in a angle.Unit, such as 90 °. Create one with
// NewAngle or StringToAngle.
//
// A Angle is a small immutable value: every method returns a new value, so
// copy and share it freely. The zero Angle is 0 rad.
type Angle struct {
	value float64
	unit  angle.Unit
}
