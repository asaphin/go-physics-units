package units

import (
	"errors"
	"fmt"
	"github.com/asaphin/go-physics-units/velocity"
)

type Velocity interface {
	Measurement
}

type velocityImplenentation struct {
	baseMeasurement
}

func NewVelocity(value float64, unit string) (Velocity, error) {
	rates := velocity.ConversionRates()

	if _, ok := rates.Has(unit + unit); !ok {
		return nil, fmt.Errorf("unknown Velocity unit %s", unit)
	}

	m, err := newBaseMeasurement(value, unit, rates)
	if err != nil {
		return nil, err
	}

	return &velocityImplenentation{*m}, nil
}

var errVelocityConversion = errors.New("not a velocity measure")

func ToVelocity(m Measurement) (Velocity, error) {
	b, ok := m.(*baseMeasurement)
	if !ok {
		return nil, errBaseMeasurementConversion
	}

	if m.Type() == MeasureVelocity {
		return &velocityImplenentation{*b}, nil
	}

	return nil, errVelocityConversion
}
