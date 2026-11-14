package units

import (
	"github.com/asaphin/go-physics-units/charge"
	"github.com/asaphin/go-physics-units/current"
	"github.com/asaphin/go-physics-units/energy"
	"github.com/asaphin/go-physics-units/power"
	"github.com/asaphin/go-physics-units/resistance"
	"github.com/asaphin/go-physics-units/voltage"
)

// Voltage is a voltage: a value in a voltage.Unit, such as 230 V. Create one with
// NewVoltage or StringToVoltage.
//
// A Voltage is a small immutable value: every method returns a new value, so
// copy and share it freely. The zero Voltage is 0 V.
type Voltage struct {
	value float64
	unit  voltage.Unit
}

// MultiplyByCurrent returns the power v·i (W) in [power.BaseUnit].
func (v Voltage) MultiplyByCurrent(i Current) Power {
	return Power{value: v.baseValue() * i.baseValue(), unit: power.BaseUnit}
}

// DivideByCurrent returns the resistance v/i (Ω), by Ohm's law in [resistance.BaseUnit]. It fails with
// [ErrZeroDivision] if i is zero.
func (v Voltage) DivideByCurrent(i Current) (Resistance, error) {
	if i.value == 0 {
		return Resistance{}, ErrZeroDivision
	}

	return Resistance{value: v.baseValue() / i.baseValue(), unit: resistance.BaseUnit}, nil
}

// DivideByResistance returns the current v/r (A), by Ohm's law in [current.BaseUnit]. It fails with
// [ErrZeroDivision] if r is zero.
func (v Voltage) DivideByResistance(r Resistance) (Current, error) {
	if r.value == 0 {
		return Current{}, ErrZeroDivision
	}

	return Current{value: v.baseValue() / r.baseValue(), unit: current.BaseUnit}, nil
}

// MultiplyByCharge returns the energy v·q (J) in [energy.BaseUnit].
func (v Voltage) MultiplyByCharge(q Charge) Energy {
	return Energy{value: v.baseValue() * q.baseValue(), unit: energy.BaseUnit}
}

// MultiplyByCapacitance returns the charge v·c (C) in [charge.BaseUnit]: the charge a capacitor c holds.
func (v Voltage) MultiplyByCapacitance(c Capacitance) Charge {
	return Charge{value: v.baseValue() * c.baseValue(), unit: charge.BaseUnit}
}
