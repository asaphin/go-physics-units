package units

import (
	"github.com/asaphin/go-physics-units/doserate"
	"github.com/asaphin/go-physics-units/equivalentdose"
)

// DoseRate is a dose rate: a value in a doserate.Unit, such as 0.12 μSv/h.
// Create one with NewDoseRate or StringToDoseRate.
//
// A DoseRate is a small immutable value: every method returns a new value, so
// copy and share it freely. The zero DoseRate is 0 Sv/s.
type DoseRate struct {
	value float64
	unit  doserate.Unit
}

// MultiplyByTime returns the equivalent dose r·t (Sv) in equivalentdose.BaseUnit.
func (r DoseRate) MultiplyByTime(t Time) EquivalentDose {
	return EquivalentDose{value: r.baseValue() * t.baseValue(), unit: equivalentdose.BaseUnit}
}
