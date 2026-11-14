package units

import (
	"github.com/asaphin/go-physics-units/doserate"
	"github.com/asaphin/go-physics-units/equivalentdose"
	"github.com/asaphin/go-physics-units/time"
)

// EquivalentDose is an equivalent dose: a value in a equivalentdose.Unit, such
// as 2.4 mSv. Create one with NewEquivalentDose or StringToEquivalentDose.
//
// An EquivalentDose is a small immutable value: every method returns a new
// value, so copy and share it freely. The zero EquivalentDose is 0 Sv.
type EquivalentDose struct {
	value float64
	unit  equivalentdose.Unit
}

// DivideByTime returns the dose rate h/t (Sv/s) in [doserate.BaseUnit]. It fails with [ErrZeroDivision] if
// t is zero.
func (h EquivalentDose) DivideByTime(t Time) (DoseRate, error) {
	if t.value == 0 {
		return DoseRate{}, ErrZeroDivision
	}

	return DoseRate{value: h.baseValue() / t.baseValue(), unit: doserate.BaseUnit}, nil
}

// DivideByDoseRate returns the time h/r (s) in [time.BaseUnit]: how long it takes to receive h at r. It
// fails with [ErrZeroDivision] if r is zero.
func (h EquivalentDose) DivideByDoseRate(r DoseRate) (Time, error) {
	if r.value == 0 {
		return Time{}, ErrZeroDivision
	}

	return Time{value: h.baseValue() / r.baseValue(), unit: time.BaseUnit}, nil
}
