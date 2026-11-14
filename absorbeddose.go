package units

import "github.com/asaphin/go-physics-units/absorbeddose"

// AbsorbedDose is an absorbed dose: a value in a absorbeddose.Unit, such as 2
// Gy. Create one with NewAbsorbedDose or StringToAbsorbedDose.
//
// An AbsorbedDose is a small immutable value: every method returns a new value,
// so copy and share it freely. The zero AbsorbedDose is 0 Gy.
type AbsorbedDose struct {
	value float64
	unit  absorbeddose.Unit
}
