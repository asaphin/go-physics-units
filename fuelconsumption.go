package units

import "github.com/asaphin/go-physics-units/fuelconsumption"

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
