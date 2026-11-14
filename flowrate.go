package units

import "github.com/asaphin/go-physics-units/flowrate"

// FlowRate is a flow rate: a value in a flowrate.Unit, such as 12 L/min. Create
// one with NewFlowRate or StringToFlowRate.
//
// A FlowRate is a small immutable value: every method returns a new value, so
// copy and share it freely. The zero FlowRate is 0 m³/s.
type FlowRate struct {
	value float64
	unit  flowrate.Unit
}
