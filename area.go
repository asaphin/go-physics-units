package units

import (
	"github.com/asaphin/go-physics-units/area"
	"github.com/asaphin/go-physics-units/distance"
	"github.com/asaphin/go-physics-units/volume"
)

// Area is an area: a value in a area.Unit, such as 10 m². Create one with
// NewArea or StringToArea.
//
// A Area is a small immutable value: every method returns a new value, so
// copy and share it freely. The zero Area is 0 m².
type Area struct {
	value float64
	unit  area.Unit
}

// DivideByDistance returns the distance a/d (m) in distance.BaseUnit. It fails with
// ErrZeroDivision if d is zero.
func (a Area) DivideByDistance(d Distance) (Distance, error) {
	if d.value == 0 {
		return Distance{}, ErrZeroDivision
	}

	return Distance{value: a.baseValue() / d.baseValue(), unit: distance.BaseUnit}, nil
}

// MultiplyByDistance returns the volume a·d (m³) in volume.BaseUnit.
func (a Area) MultiplyByDistance(d Distance) Volume {
	return Volume{value: a.baseValue() * d.baseValue(), unit: volume.BaseUnit}
}
