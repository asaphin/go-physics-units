package units

import (
	"github.com/asaphin/go-physics-units/resistance"
	"github.com/asaphin/go-physics-units/time"
	"github.com/asaphin/go-physics-units/voltage"
)

// Resistance is an electrical resistance: a value in a resistance.Unit, such as 4.7 kΩ. Create one with
// NewResistance or StringToResistance.
//
// A Resistance is a small immutable value: every method returns a new value, so
// copy and share it freely. The zero Resistance is 0 Ω.
type Resistance struct {
	value float64
	unit  resistance.Unit
}

// MultiplyByCurrent returns the voltage r·i (V), by Ohm's law in [voltage.BaseUnit].
func (r Resistance) MultiplyByCurrent(i Current) Voltage {
	return Voltage{value: r.baseValue() * i.baseValue(), unit: voltage.BaseUnit}
}

// MultiplyByCapacitance returns the time r·c (s) in [time.BaseUnit]: the time constant of an RC circuit.
func (r Resistance) MultiplyByCapacitance(c Capacitance) Time {
	return Time{value: r.baseValue() * c.baseValue(), unit: time.BaseUnit}
}
