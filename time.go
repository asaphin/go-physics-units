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
