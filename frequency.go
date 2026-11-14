package units

import (
	"github.com/asaphin/go-physics-units/frequency"
	"github.com/asaphin/go-physics-units/time"
)

// Frequency is a frequency: a value in a frequency.Unit, such as 50 Hz. Create one with
// NewFrequency or StringToFrequency.
//
// A Frequency is a small immutable value: every method returns a new value, so
// copy and share it freely. The zero Frequency is 0 Hz.
type Frequency struct {
	value float64
	unit  frequency.Unit
}

// Period returns the period 1/f in time.BaseUnit (s). It fails with
// ErrZeroDivision if f is zero.
func (f Frequency) Period() (Time, error) {
	if f.value == 0 {
		return Time{}, ErrZeroDivision
	}

	return Time{value: 1 / f.baseValue(), unit: time.BaseUnit}, nil
}
