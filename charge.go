package units

import (
	"github.com/asaphin/go-physics-units/charge"
	"github.com/asaphin/go-physics-units/current"
	"github.com/asaphin/go-physics-units/energy"
	"github.com/asaphin/go-physics-units/time"
)

// Charge is an electric charge: a value in a charge.Unit, such as 3000 mAh. Create one with
// NewCharge or StringToCharge.
//
// A Charge is a small immutable value: every method returns a new value, so
// copy and share it freely. The zero Charge is 0 C.
type Charge struct {
	value float64
	unit  charge.Unit
}

// DivideByTime returns the current q/t (A) in current.BaseUnit. It fails with
// ErrZeroDivision if t is zero.
func (q Charge) DivideByTime(t Time) (Current, error) {
	if t.value == 0 {
		return Current{}, ErrZeroDivision
	}

	return Current{value: q.baseValue() / t.baseValue(), unit: current.BaseUnit}, nil
}

// DivideByCurrent returns the time q/i (s), for example to drain a battery in time.BaseUnit. It fails with
// ErrZeroDivision if i is zero.
func (q Charge) DivideByCurrent(i Current) (Time, error) {
	if i.value == 0 {
		return Time{}, ErrZeroDivision
	}

	return Time{value: q.baseValue() / i.baseValue(), unit: time.BaseUnit}, nil
}

// MultiplyByVoltage returns the energy q·v (J) in energy.BaseUnit.
func (q Charge) MultiplyByVoltage(v Voltage) Energy {
	return Energy{value: q.baseValue() * v.baseValue(), unit: energy.BaseUnit}
}
