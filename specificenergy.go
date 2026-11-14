package units

import (
	"github.com/asaphin/go-physics-units/energy"
	"github.com/asaphin/go-physics-units/specificenergy"
)

// SpecificEnergy is a specific energy: a value in a specificenergy.Unit, such
// as 250 Wh/kg. Create one with NewSpecificEnergy or StringToSpecificEnergy.
//
// A SpecificEnergy is a small immutable value: every method returns a new
// value, so copy and share it freely. The zero SpecificEnergy is 0 J/kg.
type SpecificEnergy struct {
	value float64
	unit  specificenergy.Unit
}

// MultiplyByMass returns the energy e·m (J) in [energy.BaseUnit].
func (e SpecificEnergy) MultiplyByMass(m Mass) Energy {
	return Energy{value: e.baseValue() * m.baseValue(), unit: energy.BaseUnit}
}
