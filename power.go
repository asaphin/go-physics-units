package units

import (
	"github.com/asaphin/go-physics-units/angularvelocity"
	"github.com/asaphin/go-physics-units/current"
	"github.com/asaphin/go-physics-units/energy"
	"github.com/asaphin/go-physics-units/power"
	"github.com/asaphin/go-physics-units/torque"
	"github.com/asaphin/go-physics-units/voltage"
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

// MultiplyByTime returns the energy p·t (J) in [energy.BaseUnit].
func (p Power) MultiplyByTime(t Time) Energy {
	return Energy{value: p.baseValue() * t.baseValue(), unit: energy.BaseUnit}
}

// DivideByVoltage returns the current p/v (A) in [current.BaseUnit]. It fails with
// [ErrZeroDivision] if v is zero.
func (p Power) DivideByVoltage(v Voltage) (Current, error) {
	if v.value == 0 {
		return Current{}, ErrZeroDivision
	}

	return Current{value: p.baseValue() / v.baseValue(), unit: current.BaseUnit}, nil
}

// DivideByCurrent returns the voltage p/i (V) in [voltage.BaseUnit]. It fails with
// [ErrZeroDivision] if i is zero.
func (p Power) DivideByCurrent(i Current) (Voltage, error) {
	if i.value == 0 {
		return Voltage{}, ErrZeroDivision
	}

	return Voltage{value: p.baseValue() / i.baseValue(), unit: voltage.BaseUnit}, nil
}

// DivideByAngularVelocity returns the torque p/w (N·m) in [torque.BaseUnit]: the torque of a shaft
// delivering p at w. It fails with [ErrZeroDivision] if w is zero.
func (p Power) DivideByAngularVelocity(w AngularVelocity) (Torque, error) {
	if w.value == 0 {
		return Torque{}, ErrZeroDivision
	}

	return Torque{value: p.baseValue() / w.baseValue(), unit: torque.BaseUnit}, nil
}

// DivideByTorque returns the angular velocity p/tq (rad/s) in [angularvelocity.BaseUnit]. It fails with
// [ErrZeroDivision] if tq is zero.
func (p Power) DivideByTorque(tq Torque) (AngularVelocity, error) {
	if tq.value == 0 {
		return AngularVelocity{}, ErrZeroDivision
	}

	return AngularVelocity{value: p.baseValue() / tq.baseValue(), unit: angularvelocity.BaseUnit}, nil
}
