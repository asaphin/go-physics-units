package units

import (
	"github.com/asaphin/go-physics-units/charge"
	"github.com/asaphin/go-physics-units/current"
	"github.com/asaphin/go-physics-units/power"
	"github.com/asaphin/go-physics-units/voltage"
)

// Current is an electric current: a value in a current.Unit, such as 2.5 mA. Create one with
// NewCurrent or StringToCurrent.
//
// A Current is a small immutable value: every method returns a new value, so
// copy and share it freely. The zero Current is 0 A.
type Current struct {
	value float64
	unit  current.Unit
}

// MultiplyByVoltage returns the power i·v (W) in power.BaseUnit.
func (i Current) MultiplyByVoltage(v Voltage) Power {
	return Power{value: i.baseValue() * v.baseValue(), unit: power.BaseUnit}
}

// MultiplyByResistance returns the voltage i·r (V), by Ohm's law in voltage.BaseUnit.
func (i Current) MultiplyByResistance(r Resistance) Voltage {
	return Voltage{value: i.baseValue() * r.baseValue(), unit: voltage.BaseUnit}
}

// MultiplyByTime returns the charge i·t (C) in charge.BaseUnit.
func (i Current) MultiplyByTime(t Time) Charge {
	return Charge{value: i.baseValue() * t.baseValue(), unit: charge.BaseUnit}
}
