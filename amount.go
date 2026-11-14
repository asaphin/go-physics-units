package units

import (
	"github.com/asaphin/go-physics-units/amount"
	"github.com/asaphin/go-physics-units/constant"
	"github.com/asaphin/go-physics-units/mass"
	"github.com/asaphin/go-physics-units/molarity"
	"github.com/asaphin/go-physics-units/volume"
)

// Amount is an amount of substance: a value in a amount.Unit, such as 2.5 mmol.
// Create one with NewAmount or StringToAmount.
//
// An Amount is a small immutable value: every method returns a new value, so
// copy and share it freely. The zero Amount is 0 mol.
type Amount struct {
	value float64
	unit  amount.Unit
}

// MultiplyByMolarMass returns the mass n·mm (kg) in mass.BaseUnit.
func (n Amount) MultiplyByMolarMass(mm MolarMass) Mass {
	return Mass{value: n.baseValue() * mm.baseValue(), unit: mass.BaseUnit}
}

// DivideByVolume returns the molar concentration n/v (mol/m³) in molarity.BaseUnit: the concentration
// of n dissolved in v. It fails with ErrZeroDivision if v is zero.
func (n Amount) DivideByVolume(v Volume) (Molarity, error) {
	if v.value == 0 {
		return Molarity{}, ErrZeroDivision
	}

	return Molarity{value: n.baseValue() / v.baseValue(), unit: molarity.BaseUnit}, nil
}

// DivideByMolarity returns the volume n/c (m³) in volume.BaseUnit: the volume of a solution of c
// holding n. It fails with ErrZeroDivision if c is zero.
func (n Amount) DivideByMolarity(c Molarity) (Volume, error) {
	if c.value == 0 {
		return Volume{}, ErrZeroDivision
	}

	return Volume{value: n.baseValue() / c.baseValue(), unit: volume.BaseUnit}, nil
}

// Entities returns the number of entities (atoms, molecules, ions, …) in n:
// n·N_A, with the Avogadro constant.
func (n Amount) Entities() float64 {
	return n.baseValue() * constant.AvogadroConstant
}
