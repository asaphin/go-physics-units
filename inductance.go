package units

import "github.com/asaphin/go-physics-units/inductance"

// Inductance is an inductance: a value in a inductance.Unit, such as 10 mH.
// Create one with NewInductance or StringToInductance.
//
// An Inductance is a small immutable value: every method returns a new value,
// so copy and share it freely. The zero Inductance is 0 H.
type Inductance struct {
	value float64
	unit  inductance.Unit
}
