package units

import (
	"github.com/asaphin/go-physics-units/amount"
	"github.com/asaphin/go-physics-units/constant"
	"github.com/asaphin/go-physics-units/density"
	"github.com/asaphin/go-physics-units/energy"
	"github.com/asaphin/go-physics-units/force"
	"github.com/asaphin/go-physics-units/mass"
	"github.com/asaphin/go-physics-units/massflowrate"
	"github.com/asaphin/go-physics-units/molarmass"
	"github.com/asaphin/go-physics-units/momentum"
	"github.com/asaphin/go-physics-units/time"
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

// MultiplyByAcceleration returns the force m·a (N) in [force.BaseUnit].
func (m Mass) MultiplyByAcceleration(a Acceleration) Force {
	return Force{value: m.baseValue() * a.baseValue(), unit: force.BaseUnit}
}

// DivideByVolume returns the density m/v (kg/m³) in [density.BaseUnit]. It fails with
// [ErrZeroDivision] if v is zero.
func (m Mass) DivideByVolume(v Volume) (Density, error) {
	if v.value == 0 {
		return Density{}, ErrZeroDivision
	}

	return Density{value: m.baseValue() / v.baseValue(), unit: density.BaseUnit}, nil
}

// DivideByDensity returns the volume m/d (m³) in [volume.BaseUnit]. It fails with
// [ErrZeroDivision] if d is zero.
func (m Mass) DivideByDensity(d Density) (Volume, error) {
	if d.value == 0 {
		return Volume{}, ErrZeroDivision
	}

	return Volume{value: m.baseValue() / d.baseValue(), unit: volume.BaseUnit}, nil
}

// MultiplyByVelocity returns the momentum m·v (kg·m/s) in [momentum.BaseUnit].
func (m Mass) MultiplyByVelocity(v Velocity) Momentum {
	return Momentum{value: m.baseValue() * v.baseValue(), unit: momentum.BaseUnit}
}

// DivideByTime returns the mass flow rate m/t (kg/s) in [massflowrate.BaseUnit]. It fails with
// [ErrZeroDivision] if t is zero.
func (m Mass) DivideByTime(t Time) (MassFlowRate, error) {
	if t.value == 0 {
		return MassFlowRate{}, ErrZeroDivision
	}

	return MassFlowRate{value: m.baseValue() / t.baseValue(), unit: massflowrate.BaseUnit}, nil
}

// DivideByMassFlowRate returns the time m/mr (s) in [time.BaseUnit]: how long m lasts at mr. It fails
// with [ErrZeroDivision] if mr is zero.
func (m Mass) DivideByMassFlowRate(mr MassFlowRate) (Time, error) {
	if mr.value == 0 {
		return Time{}, ErrZeroDivision
	}

	return Time{value: m.baseValue() / mr.baseValue(), unit: time.BaseUnit}, nil
}

// DivideByMolarMass returns the amount of substance m/mm (mol) in [amount.BaseUnit]. It fails with
// [ErrZeroDivision] if mm is zero.
func (m Mass) DivideByMolarMass(mm MolarMass) (Amount, error) {
	if mm.value == 0 {
		return Amount{}, ErrZeroDivision
	}

	return Amount{value: m.baseValue() / mm.baseValue(), unit: amount.BaseUnit}, nil
}

// DivideByAmount returns the molar mass m/n (kg/mol) in [molarmass.BaseUnit]. It fails with
// [ErrZeroDivision] if n is zero.
func (m Mass) DivideByAmount(n Amount) (MolarMass, error) {
	if n.value == 0 {
		return MolarMass{}, ErrZeroDivision
	}

	return MolarMass{value: m.baseValue() / n.baseValue(), unit: molarmass.BaseUnit}, nil
}

// MultiplyBySpecificEnergy returns the energy m·e (J) in [energy.BaseUnit].
func (m Mass) MultiplyBySpecificEnergy(e SpecificEnergy) Energy {
	return Energy{value: m.baseValue() * e.baseValue(), unit: energy.BaseUnit}
}

// KineticEnergy returns the kinetic energy m·v²/2 (J) in [energy.BaseUnit] of the
// mass m moving at v. It is the classical formula: use it far below the speed
// of light.
func (m Mass) KineticEnergy(v Velocity) Energy {
	s := v.baseValue()

	return Energy{value: m.baseValue() * s * s / 2, unit: energy.BaseUnit}
}

// Weight returns the force m·gₙ (N) in [force.BaseUnit] that gravity exerts on m
// at standard gravity, 9.80665 m/s²: the weight of 1 kg is 1 kgf.
func (m Mass) Weight() Force {
	return Force{value: m.baseValue() * constant.StandardGravity, unit: force.BaseUnit}
}

// RestEnergy returns the energy m·c² (J) in [energy.BaseUnit] equivalent to the
// mass m.
func (m Mass) RestEnergy() Energy {
	return Energy{value: m.baseValue() * (constant.SpeedOfLight * constant.SpeedOfLight), unit: energy.BaseUnit}
}
