package units

import "github.com/asaphin/go-physics-units/illuminance"

// Illuminance is an illuminance: a value in a illuminance.Unit, such as 500 lx.
// Create one with NewIlluminance or StringToIlluminance.
//
// An Illuminance is a small immutable value: every method returns a new value,
// so copy and share it freely. The zero Illuminance is 0 lx.
type Illuminance struct {
	value float64
	unit  illuminance.Unit
}
