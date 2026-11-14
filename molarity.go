package units

import (
	"github.com/asaphin/go-physics-units/amount"
	"github.com/asaphin/go-physics-units/molarity"
)

// Molarity is a molar concentration: a value in a molarity.Unit, such as 150
// mM. Create one with NewMolarity or StringToMolarity.
//
// A Molarity is a small immutable value: every method returns a new value, so
// copy and share it freely. The zero Molarity is 0 mol/m³.
type Molarity struct {
	value float64
	unit  molarity.Unit
}

// MultiplyByVolume returns the amount of substance c·v (mol) in [amount.BaseUnit].
func (c Molarity) MultiplyByVolume(v Volume) Amount {
	return Amount{value: c.baseValue() * v.baseValue(), unit: amount.BaseUnit}
}
