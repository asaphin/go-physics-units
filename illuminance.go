package units

import (
	"github.com/asaphin/go-physics-units/illuminance"
	"github.com/asaphin/go-physics-units/luminousflux"
)

// Illuminance is an illuminance: a value in a illuminance.Unit, such as 500 lx.
// Create one with NewIlluminance or StringToIlluminance.
//
// An Illuminance is a small immutable value: every method returns a new value,
// so copy and share it freely. The zero Illuminance is 0 lx.
type Illuminance struct {
	value float64
	unit  illuminance.Unit
}

// MultiplyByArea returns the luminous flux e·a (lm) in luminousflux.BaseUnit: the luminous flux that
// lights a to e.
func (e Illuminance) MultiplyByArea(a Area) LuminousFlux {
	return LuminousFlux{value: e.baseValue() * a.baseValue(), unit: luminousflux.BaseUnit}
}
