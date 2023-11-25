package units

import (
	"errors"
	"fmt"
	"github.com/asaphin/go-physics-units/time"
)

var timeConversionFactors = newImmutableConversionFactors(time.ConversionFactors())

type Time interface {
	Measurement
	ConvertTo(unit string) (Time, error)
	MustConvertTo(unit string) Time
	ConvertToBaseUnits() Time
}

type timeImplementation struct {
	baseMeasurement
}

func (t *timeImplementation) ConvertTo(unit string) (Time, error) {
	m, err := t.convertTo(unit)
	if err != nil {
		return nil, err
	}

	return &timeImplementation{*m}, nil
}

func (t *timeImplementation) MustConvertTo(unit string) Time {
	m, err := t.convertTo(unit)
	if err != nil {
		panic(err)
	}

	return &timeImplementation{*m}
}

func (t *timeImplementation) ConvertToBaseUnits() Time {
	m, err := t.convertTo(time.BaseUnit)
	if err != nil {
		panic(err)
	}

	return &timeImplementation{*m}
}

func NewTime(value float64, unit string) (Time, error) {
	if _, ok := timeConversionFactors.HasFactor(unit); !ok {
		return nil, fmt.Errorf("unknown Time unit %s", unit)
	}

	m, err := newBaseMeasurement(value, unit, timeConversionFactors)
	if err != nil {
		return nil, err
	}

	return &timeImplementation{*m}, nil
}

var errTimeConversion = errors.New("not a time measure")

func ToTime(m Measurement) (Time, error) {
	b, ok := m.(*baseMeasurement)
	if !ok {
		return nil, errBaseMeasurementConversion
	}

	if m.Type() == MeasureTime {
		return &timeImplementation{*b}, nil
	}

	return nil, errTimeConversion
}
