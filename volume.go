package units

import (
	"github.com/asaphin/go-physics-units/area"
	"github.com/asaphin/go-physics-units/distance"
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

// DivideByArea returns the distance v/a (m) in distance.BaseUnit. It fails with
// ErrZeroDivision if a is zero.
func (v Volume) DivideByArea(a Area) (Distance, error) {
	if a.value == 0 {
		return Distance{}, ErrZeroDivision
	}

	return Distance{value: v.baseValue() / a.baseValue(), unit: distance.BaseUnit}, nil
}

// DivideByDistance returns the area v/d (m²) in area.BaseUnit. It fails with
// ErrZeroDivision if d is zero.
func (v Volume) DivideByDistance(d Distance) (Area, error) {
	if d.value == 0 {
		return Area{}, ErrZeroDivision
	}

	return Area{value: v.baseValue() / d.baseValue(), unit: area.BaseUnit}, nil
}
