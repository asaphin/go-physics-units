package units

import (
	"github.com/asaphin/go-physics-units/constant"
	"github.com/asaphin/go-physics-units/energy"
	"github.com/asaphin/go-physics-units/frequency"
	"github.com/asaphin/go-physics-units/time"
	"github.com/asaphin/go-physics-units/velocity"
)

// Frequency is a frequency: a value in a frequency.Unit, such as 50 Hz. Create one with
// NewFrequency or StringToFrequency.
//
// A Frequency is a small immutable value: every method returns a new value, so
// copy and share it freely. The zero Frequency is 0 Hz.
type Frequency struct {
	value float64
	unit  frequency.Unit
}

// Period returns the period 1/f in [time.BaseUnit] (s). It fails with
// [ErrZeroDivision] if f is zero.
func (f Frequency) Period() (Time, error) {
	if f.value == 0 {
		return Time{}, ErrZeroDivision
	}

	return Time{value: 1 / f.baseValue(), unit: time.BaseUnit}, nil
}

// MultiplyByDistance returns the velocity f·d (m/s) in [velocity.BaseUnit]: the speed of a wave of
// frequency f and wavelength d.
func (f Frequency) MultiplyByDistance(d Distance) Velocity {
	return Velocity{value: f.baseValue() * d.baseValue(), unit: velocity.BaseUnit}
}

// PhotonEnergy returns the energy h·f (J) in [energy.BaseUnit] of a photon of
// frequency f.
func (f Frequency) PhotonEnergy() Energy {
	return Energy{value: f.baseValue() * constant.PlanckConstant, unit: energy.BaseUnit}
}
