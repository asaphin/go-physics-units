package units

import (
	"github.com/asaphin/go-physics-units/pace"
	"github.com/asaphin/go-physics-units/time"
	"github.com/asaphin/go-physics-units/velocity"
)

// Pace is a pace: a value in a pace.Unit, such as 5 min/km. Create one with
// NewPace or StringToPace.
//
// A Pace is a small immutable value: every method returns a new value, so copy
// and share it freely. The zero Pace is 0 s/m.
type Pace struct {
	value float64
	unit  pace.Unit
}

// MultiplyByDistance returns the time p·d (s) in [time.BaseUnit]: the time to cover d at the pace.
func (p Pace) MultiplyByDistance(d Distance) Time {
	return Time{value: p.baseValue() * d.baseValue(), unit: time.BaseUnit}
}

// Velocity returns the speed 1/p (m/s) in [velocity.BaseUnit] at the pace p. It
// fails with [ErrZeroDivision] if p is zero.
func (p Pace) Velocity() (Velocity, error) {
	if p.value == 0 {
		return Velocity{}, ErrZeroDivision
	}

	return Velocity{value: 1 / p.baseValue(), unit: velocity.BaseUnit}, nil
}
