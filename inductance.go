package units

import (
	"github.com/asaphin/go-physics-units/inductance"
	"github.com/asaphin/go-physics-units/time"
)

// Inductance is an inductance: a value in a inductance.Unit, such as 10 mH.
// Create one with NewInductance or StringToInductance.
//
// An Inductance is a small immutable value: every method returns a new value,
// so copy and share it freely. The zero Inductance is 0 H.
type Inductance struct {
	value float64
	unit  inductance.Unit
}

// DivideByResistance returns the time l/r (s) in [time.BaseUnit]: the time constant of an RL circuit. It
// fails with [ErrZeroDivision] if r is zero.
func (l Inductance) DivideByResistance(r Resistance) (Time, error) {
	if r.value == 0 {
		return Time{}, ErrZeroDivision
	}

	return Time{value: l.baseValue() / r.baseValue(), unit: time.BaseUnit}, nil
}
