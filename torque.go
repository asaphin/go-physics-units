package units

import "github.com/asaphin/go-physics-units/torque"

// Torque is a torque: a value in a torque.Unit, such as 40 N·m. Create one with
// NewTorque or StringToTorque.
//
// A Torque is a small immutable value: every method returns a new value, so
// copy and share it freely. The zero Torque is 0 N·m.
type Torque struct {
	value float64
	unit  torque.Unit
}
