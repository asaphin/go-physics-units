package units

import (
	"math"

	"github.com/asaphin/go-physics-units/area"
	"github.com/asaphin/go-physics-units/distance"
	"github.com/asaphin/go-physics-units/flowrate"
	"github.com/asaphin/go-physics-units/time"
	"github.com/asaphin/go-physics-units/volume"
)

// Volume is a volume: a value in a volume.Unit, such as 10 L. Create one with
// NewVolume or StringToVolume.
//
// A Volume is a small immutable value: every method returns a new value, so
// copy and share it freely. The zero Volume is 0 m³.
type Volume struct {
	value float64
	unit  volume.Unit
}

// DivideByArea returns the distance v/a (m) in [distance.BaseUnit]. It fails with
// [ErrZeroDivision] if a is zero.
func (v Volume) DivideByArea(a Area) (Distance, error) {
	if a.value == 0 {
		return Distance{}, ErrZeroDivision
	}

	return Distance{value: v.baseValue() / a.baseValue(), unit: distance.BaseUnit}, nil
}

// DivideByDistance returns the area v/d (m²) in [area.BaseUnit]. It fails with
// [ErrZeroDivision] if d is zero.
func (v Volume) DivideByDistance(d Distance) (Area, error) {
	if d.value == 0 {
		return Area{}, ErrZeroDivision
	}

	return Area{value: v.baseValue() / d.baseValue(), unit: area.BaseUnit}, nil
}

// DivideByTime returns the flow rate v/t (m³/s) in [flowrate.BaseUnit]. It fails with [ErrZeroDivision] if
// t is zero.
func (v Volume) DivideByTime(t Time) (FlowRate, error) {
	if t.value == 0 {
		return FlowRate{}, ErrZeroDivision
	}

	return FlowRate{value: v.baseValue() / t.baseValue(), unit: flowrate.BaseUnit}, nil
}

// DivideByFlowRate returns the time v/q (s) in [time.BaseUnit]: the time to move v at q. It fails with
// [ErrZeroDivision] if q is zero.
func (v Volume) DivideByFlowRate(q FlowRate) (Time, error) {
	if q.value == 0 {
		return Time{}, ErrZeroDivision
	}

	return Time{value: v.baseValue() / q.baseValue(), unit: time.BaseUnit}, nil
}

// DivideByFuelConsumption returns the distance v/fc (m) in [distance.BaseUnit]: the range on v. It fails
// with [ErrZeroDivision] if fc is zero.
func (v Volume) DivideByFuelConsumption(fc FuelConsumption) (Distance, error) {
	if fc.value == 0 {
		return Distance{}, ErrZeroDivision
	}

	return Distance{value: v.baseValue() / fc.baseValue(), unit: distance.BaseUnit}, nil
}

// CubeEdge returns the edge ∛v (m) in [distance.BaseUnit] of a cube of volume v:
// the edge of the equivalent cube.
func (v Volume) CubeEdge() Distance {
	return Distance{value: math.Cbrt(v.baseValue()), unit: distance.BaseUnit}
}

// SphereRadius returns the radius ∛(3v/4π) (m) in [distance.BaseUnit] of a sphere
// of volume v: the radius of the volume-equivalent sphere. Double it for the
// equivalent spherical diameter of particle sizing.
func (v Volume) SphereRadius() Distance {
	return Distance{value: math.Cbrt(3 * v.baseValue() / (4 * math.Pi)), unit: distance.BaseUnit}
}
