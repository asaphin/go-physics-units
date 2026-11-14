package units

import "github.com/asaphin/go-physics-units/massflowrate"

// MassFlowRate is a mass flow rate: a value in a massflowrate.Unit, such as 250
// kg/h. Create one with NewMassFlowRate or StringToMassFlowRate.
//
// A MassFlowRate is a small immutable value: every method returns a new value,
// so copy and share it freely. The zero MassFlowRate is 0 kg/s.
type MassFlowRate struct {
	value float64
	unit  massflowrate.Unit
}
