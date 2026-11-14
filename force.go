package units

import (
	"github.com/asaphin/go-physics-units/acceleration"
	"github.com/asaphin/go-physics-units/energy"
	"github.com/asaphin/go-physics-units/force"
	"github.com/asaphin/go-physics-units/mass"
	"github.com/asaphin/go-physics-units/momentum"
	"github.com/asaphin/go-physics-units/pressure"
	"github.com/asaphin/go-physics-units/torque"
)

// Force is a force: a value in a force.Unit, such as 10 N. Create one with
// NewForce or StringToForce.
//
// A Force is a small immutable value: every method returns a new value, so
// copy and share it freely. The zero Force is 0 N.
type Force struct {
	value float64
	unit  force.Unit
}

// DivideByMass returns the acceleration f/m (m/s²) in acceleration.BaseUnit. It fails with
// ErrZeroDivision if m is zero.
func (f Force) DivideByMass(m Mass) (Acceleration, error) {
	if m.value == 0 {
		return Acceleration{}, ErrZeroDivision
	}

	return Acceleration{value: f.baseValue() / m.baseValue(), unit: acceleration.BaseUnit}, nil
}

// DivideByAcceleration returns the mass f/a (kg) in mass.BaseUnit. It fails with
// ErrZeroDivision if a is zero.
func (f Force) DivideByAcceleration(a Acceleration) (Mass, error) {
	if a.value == 0 {
		return Mass{}, ErrZeroDivision
	}

	return Mass{value: f.baseValue() / a.baseValue(), unit: mass.BaseUnit}, nil
}

// DivideByArea returns the pressure f/a (Pa) in pressure.BaseUnit. It fails with
// ErrZeroDivision if a is zero.
func (f Force) DivideByArea(a Area) (Pressure, error) {
	if a.value == 0 {
		return Pressure{}, ErrZeroDivision
	}

	return Pressure{value: f.baseValue() / a.baseValue(), unit: pressure.BaseUnit}, nil
}

// MultiplyByDistance returns the energy (work) f·d (J) in energy.BaseUnit.
func (f Force) MultiplyByDistance(d Distance) Energy {
	return Energy{value: f.baseValue() * d.baseValue(), unit: energy.BaseUnit}
}

// MultiplyByTime returns the momentum f·t (kg·m/s) in momentum.BaseUnit: the impulse, the change of
// momentum.
func (f Force) MultiplyByTime(t Time) Momentum {
	return Momentum{value: f.baseValue() * t.baseValue(), unit: momentum.BaseUnit}
}

// TorqueAt returns the torque f·arm (N·m) in torque.BaseUnit: the torque of f at a lever arm,
// perpendicular to it.
func (f Force) TorqueAt(arm Distance) Torque {
	return Torque{value: f.baseValue() * arm.baseValue(), unit: torque.BaseUnit}
}
