package units

import "github.com/asaphin/go-physics-units/equivalentdose"

// EquivalentDose is an equivalent dose: a value in a equivalentdose.Unit, such
// as 2.4 mSv. Create one with NewEquivalentDose or StringToEquivalentDose.
//
// An EquivalentDose is a small immutable value: every method returns a new
// value, so copy and share it freely. The zero EquivalentDose is 0 Sv.
type EquivalentDose struct {
	value float64
	unit  equivalentdose.Unit
}
