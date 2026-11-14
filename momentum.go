package units

import "github.com/asaphin/go-physics-units/momentum"

// Momentum is a momentum: a value in a momentum.Unit, such as 10 N·s. Create
// one with NewMomentum or StringToMomentum.
//
// A Momentum is a small immutable value: every method returns a new value, so
// copy and share it freely. The zero Momentum is 0 kg·m/s.
type Momentum struct {
	value float64
	unit  momentum.Unit
}
