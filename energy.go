package units

import (
	"github.com/asaphin/go-physics-units/constant"
	"github.com/asaphin/go-physics-units/energy"
	"github.com/asaphin/go-physics-units/energydensity"
	"github.com/asaphin/go-physics-units/force"
	"github.com/asaphin/go-physics-units/frequency"
	"github.com/asaphin/go-physics-units/mass"
	"github.com/asaphin/go-physics-units/power"
	"github.com/asaphin/go-physics-units/specificenergy"
	"github.com/asaphin/go-physics-units/time"
	"github.com/asaphin/go-physics-units/voltage"
	"github.com/asaphin/go-physics-units/volume"
)

// Energy is an amount of energy: a value in an energy.Unit, such as 10 kWh. Create one with
// NewEnergy or StringToEnergy.
//
// An Energy is a small immutable value: every method returns a new value, so
// copy and share it freely. The zero Energy is 0 J.
type Energy struct {
	value float64
	unit  energy.Unit
}

// DivideByTime returns the power e/t (W) in power.BaseUnit. It fails with
// ErrZeroDivision if t is zero.
func (e Energy) DivideByTime(t Time) (Power, error) {
	if t.value == 0 {
		return Power{}, ErrZeroDivision
	}

	return Power{value: e.baseValue() / t.baseValue(), unit: power.BaseUnit}, nil
}

// DivideByDistance returns the force e/d (N) in force.BaseUnit. It fails with
// ErrZeroDivision if d is zero.
func (e Energy) DivideByDistance(d Distance) (Force, error) {
	if d.value == 0 {
		return Force{}, ErrZeroDivision
	}

	return Force{value: e.baseValue() / d.baseValue(), unit: force.BaseUnit}, nil
}

// DivideByCharge returns the voltage e/q (V) in voltage.BaseUnit. It fails with
// ErrZeroDivision if q is zero.
func (e Energy) DivideByCharge(q Charge) (Voltage, error) {
	if q.value == 0 {
		return Voltage{}, ErrZeroDivision
	}

	return Voltage{value: e.baseValue() / q.baseValue(), unit: voltage.BaseUnit}, nil
}

// DivideByPower returns the time e/p (s) in time.BaseUnit: how long the energy lasts at that power. It
// fails with ErrZeroDivision if p is zero.
func (e Energy) DivideByPower(p Power) (Time, error) {
	if p.value == 0 {
		return Time{}, ErrZeroDivision
	}

	return Time{value: e.baseValue() / p.baseValue(), unit: time.BaseUnit}, nil
}

// DivideByMass returns the specific energy e/m (J/kg) in specificenergy.BaseUnit. It fails with
// ErrZeroDivision if m is zero.
func (e Energy) DivideByMass(m Mass) (SpecificEnergy, error) {
	if m.value == 0 {
		return SpecificEnergy{}, ErrZeroDivision
	}

	return SpecificEnergy{value: e.baseValue() / m.baseValue(), unit: specificenergy.BaseUnit}, nil
}

// DivideBySpecificEnergy returns the mass e/se (kg) in mass.BaseUnit: the mass that stores e. It fails
// with ErrZeroDivision if se is zero.
func (e Energy) DivideBySpecificEnergy(se SpecificEnergy) (Mass, error) {
	if se.value == 0 {
		return Mass{}, ErrZeroDivision
	}

	return Mass{value: e.baseValue() / se.baseValue(), unit: mass.BaseUnit}, nil
}

// DivideByVolume returns the energy density e/v (J/m³) in energydensity.BaseUnit. It fails with
// ErrZeroDivision if v is zero.
func (e Energy) DivideByVolume(v Volume) (EnergyDensity, error) {
	if v.value == 0 {
		return EnergyDensity{}, ErrZeroDivision
	}

	return EnergyDensity{value: e.baseValue() / v.baseValue(), unit: energydensity.BaseUnit}, nil
}

// DivideByEnergyDensity returns the volume e/u (m³) in volume.BaseUnit: the volume that stores e. It
// fails with ErrZeroDivision if u is zero.
func (e Energy) DivideByEnergyDensity(u EnergyDensity) (Volume, error) {
	if u.value == 0 {
		return Volume{}, ErrZeroDivision
	}

	return Volume{value: e.baseValue() / u.baseValue(), unit: volume.BaseUnit}, nil
}

// MassEquivalent returns the mass e/c² (kg) in mass.BaseUnit equivalent to the
// energy e.
func (e Energy) MassEquivalent() Mass {
	return Mass{value: e.baseValue() / (constant.SpeedOfLight * constant.SpeedOfLight), unit: mass.BaseUnit}
}

// PhotonFrequency returns the frequency e/h (Hz) in frequency.BaseUnit of a
// photon of energy e.
func (e Energy) PhotonFrequency() Frequency {
	return Frequency{value: e.baseValue() / constant.PlanckConstant, unit: frequency.BaseUnit}
}
