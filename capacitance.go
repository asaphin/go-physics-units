package units

import (
	"github.com/asaphin/go-physics-units/capacitance"
	"github.com/asaphin/go-physics-units/charge"
	"github.com/asaphin/go-physics-units/time"
)

// Capacitance is a capacitance: a value in a capacitance.Unit, such as 100 μF.
// Create one with NewCapacitance or StringToCapacitance.
//
// A Capacitance is a small immutable value: every method returns a new value,
// so copy and share it freely. The zero Capacitance is 0 F.
type Capacitance struct {
	value float64
	unit  capacitance.Unit
}

// MultiplyByVoltage returns the charge c·v (C) in charge.BaseUnit: the charge a capacitor holds at v.
func (c Capacitance) MultiplyByVoltage(v Voltage) Charge {
	return Charge{value: c.baseValue() * v.baseValue(), unit: charge.BaseUnit}
}

// MultiplyByResistance returns the time c·r (s) in time.BaseUnit: the time constant of an RC circuit.
func (c Capacitance) MultiplyByResistance(r Resistance) Time {
	return Time{value: c.baseValue() * r.baseValue(), unit: time.BaseUnit}
}
