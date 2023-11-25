package units

import (
	"errors"
	"fmt"
	"github.com/asaphin/go-physics-units/mass"
)

type Mass interface {
	Measurement
}

type massImplementation struct {
	baseMeasurement
}

func NewMass(value float64, unit string) (Mass, error) {
	rates := mass.ConversionRates()

	if _, ok := rates.Has(unit + unit); !ok {
		return nil, fmt.Errorf("unknown Mass unit %s", unit)
	}

	m, err := newBaseMeasurement(value, unit, rates)
	if err != nil {
		return nil, err
	}

	return &massImplementation{*m}, nil
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
