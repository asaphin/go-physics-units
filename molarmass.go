package units

import (
	"github.com/asaphin/go-physics-units/mass"
	"github.com/asaphin/go-physics-units/molarmass"
)

// MolarMass is a molar mass: a value in a molarmass.Unit, such as 18.015 g/mol.
// Create one with NewMolarMass or StringToMolarMass.
//
// A MolarMass is a small immutable value: every method returns a new value, so
// copy and share it freely. The zero MolarMass is 0 kg/mol.
type MolarMass struct {
	value float64
	unit  molarmass.Unit
}

// MultiplyByAmount returns the mass mm·n (kg) in mass.BaseUnit.
func (mm MolarMass) MultiplyByAmount(n Amount) Mass {
	return Mass{value: mm.baseValue() * n.baseValue(), unit: mass.BaseUnit}
}
