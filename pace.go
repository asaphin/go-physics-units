package units

import "github.com/asaphin/go-physics-units/pace"

// Pace is a pace: a value in a pace.Unit, such as 5 min/km. Create one with
// NewPace or StringToPace.
//
// A Pace is a small immutable value: every method returns a new value, so copy
// and share it freely. The zero Pace is 0 s/m.
type Pace struct {
	value float64
	unit  pace.Unit
}
