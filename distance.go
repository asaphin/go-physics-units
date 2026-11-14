package units

import (
	"math"

	"github.com/asaphin/go-physics-units/area"
	"github.com/asaphin/go-physics-units/distance"
	"github.com/asaphin/go-physics-units/fueleconomy"
	"github.com/asaphin/go-physics-units/time"
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

// DivideByVelocity returns the time d/v (s) in time.BaseUnit: the time it takes to cover d at v. It
// fails with ErrZeroDivision if v is zero.
func (d Distance) DivideByVelocity(v Velocity) (Time, error) {
	if v.value == 0 {
		return Time{}, ErrZeroDivision
	}

	return Time{value: d.baseValue() / v.baseValue(), unit: time.BaseUnit}, nil
}

// MultiplyByPace returns the time d·p (s) in time.BaseUnit: the time to cover d at the pace.
func (d Distance) MultiplyByPace(p Pace) Time {
	return Time{value: d.baseValue() * p.baseValue(), unit: time.BaseUnit}
}

// DivideByVolume returns the fuel economy d/v (m/m³) in fueleconomy.BaseUnit: the fuel economy of
// covering d on v. It fails with ErrZeroDivision if v is zero.
func (d Distance) DivideByVolume(v Volume) (FuelEconomy, error) {
	if v.value == 0 {
		return FuelEconomy{}, ErrZeroDivision
	}

	return FuelEconomy{value: d.baseValue() / v.baseValue(), unit: fueleconomy.BaseUnit}, nil
}

// DivideByFuelEconomy returns the volume d/fe (m³) in volume.BaseUnit: the fuel needed to cover d. It
// fails with ErrZeroDivision if fe is zero.
func (d Distance) DivideByFuelEconomy(fe FuelEconomy) (Volume, error) {
	if fe.value == 0 {
		return Volume{}, ErrZeroDivision
	}

	return Volume{value: d.baseValue() / fe.baseValue(), unit: volume.BaseUnit}, nil
}

// MultiplyByFuelConsumption returns the volume d·fc (m³) in volume.BaseUnit: the fuel needed to cover
// d.
func (d Distance) MultiplyByFuelConsumption(fc FuelConsumption) Volume {
	return Volume{value: d.baseValue() * fc.baseValue(), unit: volume.BaseUnit}
}

// CircleArea returns the area π·d² (m²) in area.BaseUnit of a circle whose
// radius is d.
func (d Distance) CircleArea() Area {
	r := d.baseValue()

	return Area{value: math.Pi * r * r, unit: area.BaseUnit}
}

// SphereArea returns the surface area 4π·d² (m²) in area.BaseUnit of a sphere
// whose radius is d.
func (d Distance) SphereArea() Area {
	r := d.baseValue()

	return Area{value: 4 * math.Pi * r * r, unit: area.BaseUnit}
}

// SphereVolume returns the volume 4π·d³/3 (m³) in volume.BaseUnit of a sphere
// whose radius is d.
func (d Distance) SphereVolume() Volume {
	r := d.baseValue()

	return Volume{value: 4 * math.Pi / 3 * r * r * r, unit: volume.BaseUnit}
}
