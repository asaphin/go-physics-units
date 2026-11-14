package units

import (
	"github.com/asaphin/go-physics-units/area"
	"github.com/asaphin/go-physics-units/distance"
	"github.com/asaphin/go-physics-units/velocity"
	"github.com/asaphin/go-physics-units/volume"
)

// Distance is a length: a value in a distance.Unit, such as 10 km. Create one with
// NewDistance or StringToDistance.
//
// A Distance is a small immutable value: every method returns a new value, so
// copy and share it freely. The zero Distance is 0 m.
type Distance struct {
	value float64
	unit  distance.Unit
}

// DivideByTime returns the velocity d/t (m/s) in velocity.BaseUnit. It fails with
// ErrZeroDivision if t is zero.
func (d Distance) DivideByTime(t Time) (Velocity, error) {
	if t.value == 0 {
		return Velocity{}, ErrZeroDivision
	}

	return Velocity{value: d.baseValue() / t.baseValue(), unit: velocity.BaseUnit}, nil
}

// Square returns the area d² in area.BaseUnit (m²).
func (d Distance) Square() Area {
	base := d.baseValue()

	return Area{value: base * base, unit: area.BaseUnit}
}

// Cube returns the volume d³ in volume.BaseUnit (m³).
func (d Distance) Cube() Volume {
	base := d.baseValue()

	return Volume{value: base * base * base, unit: volume.BaseUnit}
}
