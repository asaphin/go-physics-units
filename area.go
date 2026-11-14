package units

import (
	"math"

	"github.com/asaphin/go-physics-units/area"
	"github.com/asaphin/go-physics-units/distance"
	"github.com/asaphin/go-physics-units/volume"
)

// Area is an area: a value in an area.Unit, such as 10 m². Create one with
// NewArea or StringToArea.
//
// An Area is a small immutable value: every method returns a new value, so
// copy and share it freely. The zero Area is 0 m².
type Area struct {
	value float64
	unit  area.Unit
}

// DivideByDistance returns the distance a/d (m) in [distance.BaseUnit]. It fails with
// [ErrZeroDivision] if d is zero.
func (a Area) DivideByDistance(d Distance) (Distance, error) {
	if d.value == 0 {
		return Distance{}, ErrZeroDivision
	}

	return Distance{value: a.baseValue() / d.baseValue(), unit: distance.BaseUnit}, nil
}

// MultiplyByDistance returns the volume a·d (m³) in [volume.BaseUnit].
func (a Area) MultiplyByDistance(d Distance) Volume {
	return Volume{value: a.baseValue() * d.baseValue(), unit: volume.BaseUnit}
}

// SquareSide returns the side √a (m) in [distance.BaseUnit] of a square of area
// a: the side of the equivalent square. It fails with [ErrOutOfRange] if a is
// negative.
func (a Area) SquareSide() (Distance, error) {
	if a.value < 0 {
		return Distance{}, errNegativeArea
	}

	return Distance{value: math.Sqrt(a.baseValue()), unit: distance.BaseUnit}, nil
}

// CircleRadius returns the radius √(a/π) (m) in [distance.BaseUnit] of a circle of
// area a: the radius of the equivalent circle. Double it for the area-equivalent
// diameter. It fails with [ErrOutOfRange] if a is negative.
func (a Area) CircleRadius() (Distance, error) {
	if a.value < 0 {
		return Distance{}, errNegativeArea
	}

	return Distance{value: math.Sqrt(a.baseValue() / math.Pi), unit: distance.BaseUnit}, nil
}
