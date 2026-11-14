package units

import (
	"github.com/asaphin/go-physics-units/energy"
	"github.com/asaphin/go-physics-units/power"
)

// Power is a power: a value in a power.Unit, such as 100 W. Create one with
// NewPower or StringToPower.
//
// A Power is a small immutable value: every method returns a new value, so
// copy and share it freely. The zero Power is 0 W.
type Power struct {
	value float64
	unit  power.Unit
}

// MultiplyByTime returns the energy p·t (J) in energy.BaseUnit.
func (p Power) MultiplyByTime(t Time) Energy {
	return Energy{value: p.baseValue() * t.baseValue(), unit: energy.BaseUnit}
}
