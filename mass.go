package units

import (
	"errors"
	"fmt"
	"github.com/asaphin/go-physics-units/mass"
)

var massConversionFactors = newImmutableConversionFactors(mass.ConversionFactors())

type Mass struct {
	BaseMeasurement
}

func NewMass(value float64, unit string) (*Mass, error) {
	if _, ok := massConversionFactors.HasFactor(unit); !ok {
		return nil, fmt.Errorf("unknown Mass unit %s", unit)
	}

	m, err := newBaseMeasurement(value, unit, massConversionFactors)
	if err != nil {
		return nil, err
	}

	return &Mass{*m}, nil
}

var errMassConversion = errors.New("not a mass measure")

func ToMass(m Measurement) (*Mass, error) {
	if m.Type() == MeasureMass {
		ms, ok := m.(*Mass)
		if !ok {
			return nil, errMassConversion
		}

		return ms, nil
	}

	return nil, errMassConversion
}
