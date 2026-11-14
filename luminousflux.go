package units

import "github.com/asaphin/go-physics-units/luminousflux"

// LuminousFlux is a luminous flux: a value in a luminousflux.Unit, such as 800
// lm. Create one with NewLuminousFlux or StringToLuminousFlux.
//
// A LuminousFlux is a small immutable value: every method returns a new value,
// so copy and share it freely. The zero LuminousFlux is 0 lm.
type LuminousFlux struct {
	value float64
	unit  luminousflux.Unit
}
