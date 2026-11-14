package units

import "github.com/asaphin/go-physics-units/fueleconomy"

// FuelEconomy is a fuel economy: a value in a fueleconomy.Unit, such as 15
// km/L. Create one with NewFuelEconomy or StringToFuelEconomy.
//
// A FuelEconomy is a small immutable value: every method returns a new value,
// so copy and share it freely. The zero FuelEconomy is 0 m/m³.
type FuelEconomy struct {
	value float64
	unit  fueleconomy.Unit
}
