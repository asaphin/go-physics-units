package units

import (
	"github.com/asaphin/go-physics-units/flowrate"
	"github.com/asaphin/go-physics-units/mass"
	"github.com/asaphin/go-physics-units/massflowrate"
)

// MassFlowRate is a mass flow rate: a value in a massflowrate.Unit, such as 250
// kg/h. Create one with NewMassFlowRate or StringToMassFlowRate.
//
// A MassFlowRate is a small immutable value: every method returns a new value,
// so copy and share it freely. The zero MassFlowRate is 0 kg/s.
type MassFlowRate struct {
	value float64
	unit  massflowrate.Unit
}

// DivideByDensity returns the flow rate mr/d (m³/s) in flowrate.BaseUnit. It fails with
// ErrZeroDivision if d is zero.
func (mr MassFlowRate) DivideByDensity(d Density) (FlowRate, error) {
	if d.value == 0 {
		return FlowRate{}, ErrZeroDivision
	}

	return FlowRate{value: mr.baseValue() / d.baseValue(), unit: flowrate.BaseUnit}, nil
}

// MultiplyByTime returns the mass mr·t (kg) in mass.BaseUnit.
func (mr MassFlowRate) MultiplyByTime(t Time) Mass {
	return Mass{value: mr.baseValue() * t.baseValue(), unit: mass.BaseUnit}
}
