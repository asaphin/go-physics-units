package units

import (
	"math"

	"github.com/asaphin/go-physics-units/radioactivity"
)

// Radioactivity is a radioactivity: a value in a radioactivity.Unit, such as
// 3.7 GBq. Create one with NewRadioactivity or StringToRadioactivity.
//
// A Radioactivity is a small immutable value: every method returns a new value,
// so copy and share it freely. The zero Radioactivity is 0 Bq.
type Radioactivity struct {
	value float64
	unit  radioactivity.Unit
}

// DecayAfter returns the activity of a after elapsed, for a single radionuclide
// with the given half-life: a·2^(−elapsed/halfLife), in the receiver's unit. A
// negative elapsed time gives the activity before. It fails with
// ErrZeroDivision if halfLife is zero.
func (a Radioactivity) DecayAfter(elapsed, halfLife Time) (Radioactivity, error) {
	if halfLife.value == 0 {
		return Radioactivity{}, ErrZeroDivision
	}

	return Radioactivity{value: a.value * math.Exp2(-elapsed.baseValue()/halfLife.baseValue()), unit: a.unit}, nil
}
