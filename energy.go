package units

import (
	"github.com/asaphin/go-physics-units/energy"
	"github.com/asaphin/go-physics-units/force"
	"github.com/asaphin/go-physics-units/power"
	"github.com/asaphin/go-physics-units/voltage"
)

// Energy is an amount of energy: a value in an energy.Unit, such as 10 kWh. Create one with
// NewEnergy or StringToEnergy.
//
// An Energy is a small immutable value: every method returns a new value, so
// copy and share it freely. The zero Energy is 0 J.
type Energy struct {
	value float64
	unit  energy.Unit
}

// DivideByTime returns the power e/t (W) in power.BaseUnit. It fails with
// ErrZeroDivision if t is zero.
func (e Energy) DivideByTime(t Time) (Power, error) {
	if t.value == 0 {
		return Power{}, ErrZeroDivision
	}

	return Power{value: e.baseValue() / t.baseValue(), unit: power.BaseUnit}, nil
}

// DivideByDistance returns the force e/d (N) in force.BaseUnit. It fails with
// ErrZeroDivision if d is zero.
func (e Energy) DivideByDistance(d Distance) (Force, error) {
	if d.value == 0 {
		return Force{}, ErrZeroDivision
	}

	return Force{value: e.baseValue() / d.baseValue(), unit: force.BaseUnit}, nil
}

// DivideByCharge returns the voltage e/q (V) in voltage.BaseUnit. It fails with
// ErrZeroDivision if q is zero.
func (e Energy) DivideByCharge(q Charge) (Voltage, error) {
	if q.value == 0 {
		return Voltage{}, ErrZeroDivision
	}

	return Voltage{value: e.baseValue() / q.baseValue(), unit: voltage.BaseUnit}, nil
}
