package units

import (
	"github.com/asaphin/go-physics-units/energy"
	"github.com/asaphin/go-physics-units/energydensity"
)

// EnergyDensity is an energy density: a value in a energydensity.Unit, such as
// 700 Wh/L. Create one with NewEnergyDensity or StringToEnergyDensity.
//
// An EnergyDensity is a small immutable value: every method returns a new
// value, so copy and share it freely. The zero EnergyDensity is 0 J/m³.
type EnergyDensity struct {
	value float64
	unit  energydensity.Unit
}

// MultiplyByVolume returns the energy u·v (J) in [energy.BaseUnit].
func (u EnergyDensity) MultiplyByVolume(v Volume) Energy {
	return Energy{value: u.baseValue() * v.baseValue(), unit: energy.BaseUnit}
}
