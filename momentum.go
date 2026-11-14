package units

import (
	"github.com/asaphin/go-physics-units/force"
	"github.com/asaphin/go-physics-units/mass"
	"github.com/asaphin/go-physics-units/momentum"
	"github.com/asaphin/go-physics-units/velocity"
)

// Momentum is a momentum: a value in a momentum.Unit, such as 10 N·s. Create
// one with NewMomentum or StringToMomentum.
//
// A Momentum is a small immutable value: every method returns a new value, so
// copy and share it freely. The zero Momentum is 0 kg·m/s.
type Momentum struct {
	value float64
	unit  momentum.Unit
}

// DivideByMass returns the velocity p/m (m/s) in [velocity.BaseUnit]. It fails with [ErrZeroDivision] if m
// is zero.
func (p Momentum) DivideByMass(m Mass) (Velocity, error) {
	if m.value == 0 {
		return Velocity{}, ErrZeroDivision
	}

	return Velocity{value: p.baseValue() / m.baseValue(), unit: velocity.BaseUnit}, nil
}

// DivideByVelocity returns the mass p/v (kg) in [mass.BaseUnit]. It fails with [ErrZeroDivision] if v is
// zero.
func (p Momentum) DivideByVelocity(v Velocity) (Mass, error) {
	if v.value == 0 {
		return Mass{}, ErrZeroDivision
	}

	return Mass{value: p.baseValue() / v.baseValue(), unit: mass.BaseUnit}, nil
}

// DivideByTime returns the force p/t (N) in [force.BaseUnit]: the mean force. It fails with
// [ErrZeroDivision] if t is zero.
func (p Momentum) DivideByTime(t Time) (Force, error) {
	if t.value == 0 {
		return Force{}, ErrZeroDivision
	}

	return Force{value: p.baseValue() / t.baseValue(), unit: force.BaseUnit}, nil
}
