package units

import (
	"github.com/asaphin/go-physics-units/fuelconsumption"
	"github.com/asaphin/go-physics-units/fueleconomy"
	"github.com/asaphin/go-physics-units/volume"
)

// FuelConsumption is a fuel consumption: a value in a fuelconsumption.Unit,
// such as 6.5 L/100km. Create one with NewFuelConsumption or
// StringToFuelConsumption.
//
// A FuelConsumption is a small immutable value: every method returns a new
// value, so copy and share it freely. The zero FuelConsumption is 0 m³/m.
type FuelConsumption struct {
	value float64
	unit  fuelconsumption.Unit
}

// MultiplyByDistance returns the volume fc·d (m³) in volume.BaseUnit: the fuel needed to cover d.
func (fc FuelConsumption) MultiplyByDistance(d Distance) Volume {
	return Volume{value: fc.baseValue() * d.baseValue(), unit: volume.BaseUnit}
}

// Economy returns the fuel economy 1/fc (m/m³) in fueleconomy.BaseUnit: convert
// it to fueleconomy.MilePerGallon. It fails with ErrZeroDivision if fc is zero.
func (fc FuelConsumption) Economy() (FuelEconomy, error) {
	if fc.value == 0 {
		return FuelEconomy{}, ErrZeroDivision
	}

	return FuelEconomy{value: 1 / fc.baseValue(), unit: fueleconomy.BaseUnit}, nil
}
