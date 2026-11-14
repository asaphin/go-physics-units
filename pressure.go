package units

import (
	"github.com/asaphin/go-physics-units/force"
	"github.com/asaphin/go-physics-units/power"
	"github.com/asaphin/go-physics-units/pressure"
)

// Pressure is a pressure: a value in a pressure.Unit, such as 101.3 kPa. Create one with
// NewPressure or StringToPressure.
//
// A Pressure is a small immutable value: every method returns a new value, so
// copy and share it freely. The zero Pressure is 0 Pa.
type Pressure struct {
	value float64
	unit  pressure.Unit
}

// MultiplyByArea returns the force p·a (N) in force.BaseUnit.
func (p Pressure) MultiplyByArea(a Area) Force {
	return Force{value: p.baseValue() * a.baseValue(), unit: force.BaseUnit}
}

// MultiplyByFlowRate returns the power p·q (W) in power.BaseUnit: the hydraulic power.
func (p Pressure) MultiplyByFlowRate(q FlowRate) Power {
	return Power{value: p.baseValue() * q.baseValue(), unit: power.BaseUnit}
}
