package units

import (
	"errors"
	"fmt"
	"github.com/asaphin/go-physics-units/time"
)

var timeConversionFactors = newImmutableConversionFactors(time.ConversionFactors())

type Time struct {
	BaseMeasurement
}

func (t *Time) ConvertTo(unit string) (*Time, error) {
	m, err := t.convertTo(unit)
	if err != nil {
		return nil, err
	}

	return &Time{*m}, nil
}

func (t *Time) MustConvertTo(unit string) *Time {
	m, err := t.convertTo(unit)
	if err != nil {
		panic(err)
	}

	return &Time{*m}
}

func (t *Time) MustConvertToBaseUnits() *Time {
	m, err := t.convertTo(time.BaseUnit)
	if err != nil {
		panic(err)
	}

	return &Time{*m}
}

func NewTime(value float64, unit string) (*Time, error) {
	if _, ok := timeConversionFactors.HasFactor(unit); !ok {
		return nil, fmt.Errorf("unknown Time unit %s", unit)
	}

	m, err := newBaseMeasurement(value, unit, timeConversionFactors)
	if err != nil {
		return nil, err
	}

	return &Time{*m}, nil
}

var timeConversionErr = errors.New("not a time measure")

func ToTime(m Measurement) (*Time, error) {
	if m.Type() == MeasureTime {
		return m.(*Time), nil
	}

	return nil, timeConversionErr
}
