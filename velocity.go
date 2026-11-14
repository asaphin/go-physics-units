package units

import (
	"math"

	"github.com/asaphin/go-physics-units/acceleration"
	"github.com/asaphin/go-physics-units/constant"
	"github.com/asaphin/go-physics-units/distance"
	"github.com/asaphin/go-physics-units/flowrate"
	"github.com/asaphin/go-physics-units/frequency"
	"github.com/asaphin/go-physics-units/momentum"
	"github.com/asaphin/go-physics-units/pace"
	"github.com/asaphin/go-physics-units/velocity"
)

// Velocity is a speed: a value in a velocity.Unit, such as 10 km/h. Create one with
// NewVelocity or StringToVelocity.
//
// A Velocity is a small immutable value: every method returns a new value, so
// copy and share it freely. The zero Velocity is 0 m/s.
type Velocity struct {
	value float64
	unit  velocity.Unit
}

// MultiplyByTime returns the distance v·t (m) in [distance.BaseUnit].
func (v Velocity) MultiplyByTime(t Time) Distance {
	return Distance{value: v.baseValue() * t.baseValue(), unit: distance.BaseUnit}
}

// DivideByTime returns the acceleration v/t (m/s²) in [acceleration.BaseUnit]. It fails with
// [ErrZeroDivision] if t is zero.
func (v Velocity) DivideByTime(t Time) (Acceleration, error) {
	if t.value == 0 {
		return Acceleration{}, ErrZeroDivision
	}

	return Acceleration{value: v.baseValue() / t.baseValue(), unit: acceleration.BaseUnit}, nil
}

// MultiplyByMass returns the momentum v·m (kg·m/s) in [momentum.BaseUnit].
func (v Velocity) MultiplyByMass(m Mass) Momentum {
	return Momentum{value: v.baseValue() * m.baseValue(), unit: momentum.BaseUnit}
}

// MultiplyByArea returns the flow rate v·a (m³/s) in [flowrate.BaseUnit]: the flow through a cross-
// section a.
func (v Velocity) MultiplyByArea(a Area) FlowRate {
	return FlowRate{value: v.baseValue() * a.baseValue(), unit: flowrate.BaseUnit}
}

// DivideByFrequency returns the distance v/f (m) in [distance.BaseUnit]: the wavelength of a wave of
// frequency f travelling at v. It fails with [ErrZeroDivision] if f is zero.
func (v Velocity) DivideByFrequency(f Frequency) (Distance, error) {
	if f.value == 0 {
		return Distance{}, ErrZeroDivision
	}

	return Distance{value: v.baseValue() / f.baseValue(), unit: distance.BaseUnit}, nil
}

// DivideByDistance returns the frequency v/d (Hz) in [frequency.BaseUnit]: the frequency of a wave of
// wavelength d travelling at v. It fails with [ErrZeroDivision] if d is zero.
func (v Velocity) DivideByDistance(d Distance) (Frequency, error) {
	if d.value == 0 {
		return Frequency{}, ErrZeroDivision
	}

	return Frequency{value: v.baseValue() / d.baseValue(), unit: frequency.BaseUnit}, nil
}

// Pace returns the pace 1/v (s/m) in [pace.BaseUnit]: the time per distance at
// the speed v. It fails with [ErrZeroDivision] if v is zero.
func (v Velocity) Pace() (Pace, error) {
	if v.value == 0 {
		return Pace{}, ErrZeroDivision
	}

	return Pace{value: 1 / v.baseValue(), unit: pace.BaseUnit}, nil
}

// LorentzFactor returns γ = 1/√(1 − v²/c²), the factor by which time dilates and
// lengths contract at the speed v in special relativity. It fails with
// [ErrOutOfRange] unless the speed is below the speed of light.
func (v Velocity) LorentzFactor() (float64, error) {
	beta := v.baseValue() / constant.SpeedOfLight
	if math.IsNaN(beta) || math.Abs(beta) >= 1 {
		return 0, errNotBelowLight
	}

	return 1 / math.Sqrt(1-beta*beta), nil
}
