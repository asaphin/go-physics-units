package units

import (
	"github.com/asaphin/go-physics-units/density"
	"github.com/asaphin/go-physics-units/mass"
)

// Density is a density: a value in a density.Unit, such as 1.2 kg/m³. Create one with
// NewDensity or StringToDensity.
//
// A Density is a small immutable value: every method returns a new value, so
// copy and share it freely. The zero Density is 0 kg/m³.
type Density struct {
	value float64
	unit  density.Unit
}

// MultiplyByVolume returns the mass d·v (kg) in mass.BaseUnit.
func (d Density) MultiplyByVolume(v Volume) Mass {
	return Mass{value: d.baseValue() * v.baseValue(), unit: mass.BaseUnit}
}
