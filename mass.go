package units

import (
	"github.com/asaphin/go-physics-units/density"
	"github.com/asaphin/go-physics-units/force"
	"github.com/asaphin/go-physics-units/mass"
	"github.com/asaphin/go-physics-units/volume"
)

// Mass is a mass: a value in a mass.Unit, such as 10 kg. Create one with
// NewMass or StringToMass.
//
// A Mass is a small immutable value: every method returns a new value, so
// copy and share it freely. The zero Mass is 0 kg.
type Mass struct {
	value float64
	unit  mass.Unit
}

// MultiplyByAcceleration returns the force m·a (N) in force.BaseUnit.
func (m Mass) MultiplyByAcceleration(a Acceleration) Force {
	return Force{value: m.baseValue() * a.baseValue(), unit: force.BaseUnit}
}

// DivideByVolume returns the density m/v (kg/m³) in density.BaseUnit. It fails with
// ErrZeroDivision if v is zero.
func (m Mass) DivideByVolume(v Volume) (Density, error) {
	if v.value == 0 {
		return Density{}, ErrZeroDivision
	}

	return Density{value: m.baseValue() / v.baseValue(), unit: density.BaseUnit}, nil
}

// DivideByDensity returns the volume m/d (m³) in volume.BaseUnit. It fails with
// ErrZeroDivision if d is zero.
func (m Mass) DivideByDensity(d Density) (Volume, error) {
	if d.value == 0 {
		return Volume{}, ErrZeroDivision
	}

	return Volume{value: m.baseValue() / d.baseValue(), unit: volume.BaseUnit}, nil
}
