package units

import (
	"github.com/asaphin/go-physics-units/energy"
	"github.com/asaphin/go-physics-units/force"
	"github.com/asaphin/go-physics-units/power"
	"github.com/asaphin/go-physics-units/torque"
)

// Torque is a torque: a value in a torque.Unit, such as 40 N·m. Create one with
// NewTorque or StringToTorque.
//
// A Torque is a small immutable value: every method returns a new value, so
// copy and share it freely. The zero Torque is 0 N·m.
type Torque struct {
	value float64
	unit  torque.Unit
}

// DivideByDistance returns the force tq/d (N) in [force.BaseUnit]: the force at a lever arm d. It fails
// with [ErrZeroDivision] if d is zero.
func (tq Torque) DivideByDistance(d Distance) (Force, error) {
	if d.value == 0 {
		return Force{}, ErrZeroDivision
	}

	return Force{value: tq.baseValue() / d.baseValue(), unit: force.BaseUnit}, nil
}

// MultiplyByAngularVelocity returns the power tq·w (W) in [power.BaseUnit]: the power of a shaft turning
// at w.
func (tq Torque) MultiplyByAngularVelocity(w AngularVelocity) Power {
	return Power{value: tq.baseValue() * w.baseValue(), unit: power.BaseUnit}
}

// MultiplyByAngle returns the energy tq·a (J) in [energy.BaseUnit]: the work done turning through a.
func (tq Torque) MultiplyByAngle(a Angle) Energy {
	return Energy{value: tq.baseValue() * a.baseValue(), unit: energy.BaseUnit}
}
