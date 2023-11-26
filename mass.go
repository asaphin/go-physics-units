package units

import (
	"errors"
	"fmt"
	"github.com/asaphin/go-physics-units/mass"
)

type Mass interface {
	Measurement
	UnitConverter[Mass]
	Arithmetics[Mass]
}

type massImplementation struct {
	baseMeasurement
}

func (m *massImplementation) ConvertTo(unit string) (Mass, error) {
	msm, err := m.convertTo(unit)
	if err != nil {
		return nil, err
	}

	return &massImplementation{*msm}, nil
}

func (m *massImplementation) MustConvertTo(unit string) Mass {
	msm, err := m.convertTo(unit)
	if err != nil {
		panic(err)
	}

	return &massImplementation{*msm}
}

func (m *massImplementation) ConvertToBaseUnits() Mass {
	msm := m.unsafeConvertTo(mass.BaseUnit)

	return &massImplementation{*msm}
}

func (m *massImplementation) valueInBaseUnits() float64 {
	return m.unsafeGetValueIn(mass.BaseUnit)
}

func (m *massImplementation) valueInUnits(unit string) float64 {
	return m.unsafeGetValueIn(unit)
}

func (m *massImplementation) Add(mass Mass) Mass {
	m1 := m.value
	m2 := mass.valueInUnits(m.unit)

	return newMass(m1+m2, m.unit)
}

func (m *massImplementation) Sub(mass Mass) Mass {
	m1 := m.value
	m2 := mass.valueInUnits(m.unit)

	return newMass(m1-m2, m.unit)
}

func (m *massImplementation) Mul(multiplier float64) Mass {
	msm := newUnsafeBaseMeasurement(m.value*multiplier, m.unit, m.conversionRates)

	return &massImplementation{*msm}
}

func (m *massImplementation) Div(divisor float64) (Mass, error) {
	if divisor == 0 {
		return nil, errZeroDivision
	}

	msm := newUnsafeBaseMeasurement(m.value/divisor, m.unit, m.conversionRates)

	return &massImplementation{*msm}, nil
}

func NewMass(value float64, unit string) (Mass, error) {
	rates := mass.ConversionRates()

	if _, ok := rates.Has(unit + unit); !ok {
		return nil, fmt.Errorf("unknown Mass unit %s", unit)
	}

	m := newUnsafeBaseMeasurement(value, unit, rates)

	return &massImplementation{*m}, nil
}

// newMass unsafe mass constructor without checks for internal usage.
func newMass(value float64, unit string) Mass {
	m := newUnsafeBaseMeasurement(value, unit, mass.ConversionRates())

	return &massImplementation{*m}
}

var errMassConversion = errors.New("not a mass measure")

func ToMass(m Measurement) (Mass, error) {
	b, ok := m.(*baseMeasurement)
	if !ok {
		return nil, errBaseMeasurementConversion
	}

	if m.Type() == MeasureMass {
		return &massImplementation{*b}, nil
	}

	return nil, errMassConversion
}
