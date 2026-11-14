package units

import (
	"github.com/asaphin/go-physics-units/distance"
	"github.com/asaphin/go-physics-units/fuelconsumption"
	"github.com/asaphin/go-physics-units/fueleconomy"
)

// FuelEconomy is a fuel economy: a value in a fueleconomy.Unit, such as 15
// km/L. Create one with NewFuelEconomy or StringToFuelEconomy.
//
// A FuelEconomy is a small immutable value: every method returns a new value,
// so copy and share it freely. The zero FuelEconomy is 0 m/m³.
type FuelEconomy struct {
	value float64
	unit  fueleconomy.Unit
}

// MultiplyByVolume returns the distance fe·v (m) in [distance.BaseUnit]: the range on v.
func (fe FuelEconomy) MultiplyByVolume(v Volume) Distance {
	return Distance{value: fe.baseValue() * v.baseValue(), unit: distance.BaseUnit}
}

// Consumption returns the fuel consumption 1/fe (m³/m) in [fuelconsumption.BaseUnit]:
// convert it to fuelconsumption.LiterPer100Kilometers. It fails with
// [ErrZeroDivision] if fe is zero.
func (fe FuelEconomy) Consumption() (FuelConsumption, error) {
	if fe.value == 0 {
		return FuelConsumption{}, ErrZeroDivision
	}

	return FuelConsumption{value: 1 / fe.baseValue(), unit: fuelconsumption.BaseUnit}, nil
}
