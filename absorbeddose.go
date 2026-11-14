package units

import (
	"github.com/asaphin/go-physics-units/absorbeddose"
	"github.com/asaphin/go-physics-units/equivalentdose"
)

// AbsorbedDose is an absorbed dose: a value in a absorbeddose.Unit, such as 2
// Gy. Create one with NewAbsorbedDose or StringToAbsorbedDose.
//
// An AbsorbedDose is a small immutable value: every method returns a new value,
// so copy and share it freely. The zero AbsorbedDose is 0 Gy.
type AbsorbedDose struct {
	value float64
	unit  absorbeddose.Unit
}

// Equivalent returns the equivalent dose d·w_R (Sv) in [equivalentdose.BaseUnit]
// of the absorbed dose d of radiation with the weighting factor w_R: 1 for
// photons and electrons, 20 for alpha particles (ICRP Publication 103).
func (d AbsorbedDose) Equivalent(weightingFactor float64) EquivalentDose {
	return EquivalentDose{value: d.baseValue() * weightingFactor, unit: equivalentdose.BaseUnit}
}
