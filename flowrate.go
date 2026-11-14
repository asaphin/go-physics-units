package units

import (
	"github.com/asaphin/go-physics-units/flowrate"
	"github.com/asaphin/go-physics-units/massflowrate"
	"github.com/asaphin/go-physics-units/power"
	"github.com/asaphin/go-physics-units/velocity"
	"github.com/asaphin/go-physics-units/volume"
)

// FlowRate is a flow rate: a value in a flowrate.Unit, such as 12 L/min. Create
// one with NewFlowRate or StringToFlowRate.
//
// A FlowRate is a small immutable value: every method returns a new value, so
// copy and share it freely. The zero FlowRate is 0 m³/s.
type FlowRate struct {
	value float64
	unit  flowrate.Unit
}

// MultiplyByTime returns the volume q·t (m³) in [volume.BaseUnit].
func (q FlowRate) MultiplyByTime(t Time) Volume {
	return Volume{value: q.baseValue() * t.baseValue(), unit: volume.BaseUnit}
}

// DivideByArea returns the velocity q/a (m/s) in [velocity.BaseUnit]: the mean speed of the flow through
// a cross-section a. It fails with [ErrZeroDivision] if a is zero.
func (q FlowRate) DivideByArea(a Area) (Velocity, error) {
	if a.value == 0 {
		return Velocity{}, ErrZeroDivision
	}

	return Velocity{value: q.baseValue() / a.baseValue(), unit: velocity.BaseUnit}, nil
}

// MultiplyByPressure returns the power q·p (W) in [power.BaseUnit]: the hydraulic power.
func (q FlowRate) MultiplyByPressure(p Pressure) Power {
	return Power{value: q.baseValue() * p.baseValue(), unit: power.BaseUnit}
}

// MultiplyByDensity returns the mass flow rate q·d (kg/s) in [massflowrate.BaseUnit].
func (q FlowRate) MultiplyByDensity(d Density) MassFlowRate {
	return MassFlowRate{value: q.baseValue() * d.baseValue(), unit: massflowrate.BaseUnit}
}
