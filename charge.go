package units

import (
	"github.com/asaphin/go-physics-units/capacitance"
	"github.com/asaphin/go-physics-units/charge"
	"github.com/asaphin/go-physics-units/current"
	"github.com/asaphin/go-physics-units/energy"
	"github.com/asaphin/go-physics-units/time"
	"github.com/asaphin/go-physics-units/voltage"
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

// DivideByTime returns the current q/t (A) in [current.BaseUnit]. It fails with
// [ErrZeroDivision] if t is zero.
func (q Charge) DivideByTime(t Time) (Current, error) {
	if t.value == 0 {
		return Current{}, ErrZeroDivision
	}

	return Current{value: q.baseValue() / t.baseValue(), unit: current.BaseUnit}, nil
}

// DivideByCurrent returns the time q/i (s), for example to drain a battery in [time.BaseUnit]. It fails with
// [ErrZeroDivision] if i is zero.
func (q Charge) DivideByCurrent(i Current) (Time, error) {
	if i.value == 0 {
		return Time{}, ErrZeroDivision
	}

	return Time{value: q.baseValue() / i.baseValue(), unit: time.BaseUnit}, nil
}

// MultiplyByVoltage returns the energy q·v (J) in [energy.BaseUnit].
func (q Charge) MultiplyByVoltage(v Voltage) Energy {
	return Energy{value: q.baseValue() * v.baseValue(), unit: energy.BaseUnit}
}

// DivideByVoltage returns the capacitance q/v (F) in [capacitance.BaseUnit]. It fails with
// [ErrZeroDivision] if v is zero.
func (q Charge) DivideByVoltage(v Voltage) (Capacitance, error) {
	if v.value == 0 {
		return Capacitance{}, ErrZeroDivision
	}

	return Capacitance{value: q.baseValue() / v.baseValue(), unit: capacitance.BaseUnit}, nil
}

// DivideByCapacitance returns the voltage q/c (V) in [voltage.BaseUnit]: the voltage across a capacitor
// c holding q. It fails with [ErrZeroDivision] if c is zero.
func (q Charge) DivideByCapacitance(c Capacitance) (Voltage, error) {
	if c.value == 0 {
		return Voltage{}, ErrZeroDivision
	}

	return Voltage{value: q.baseValue() / c.baseValue(), unit: voltage.BaseUnit}, nil
}
