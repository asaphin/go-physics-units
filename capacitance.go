package units

import "github.com/asaphin/go-physics-units/capacitance"

// Capacitance is a capacitance: a value in a capacitance.Unit, such as 100 μF.
// Create one with NewCapacitance or StringToCapacitance.
//
// A Capacitance is a small immutable value: every method returns a new value,
// so copy and share it freely. The zero Capacitance is 0 F.
type Capacitance struct {
	value float64
	unit  capacitance.Unit
}
