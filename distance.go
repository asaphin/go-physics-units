package units

import (
	"errors"
	"fmt"
	"github.com/asaphin/go-physics-units/distance"
)

var distanceConversionFactors = newImmutableConversionFactors(distance.ConversionFactors())

type Distance struct {
	BaseMeasurement
}

func (d *Distance) ConvertTo(unit string) (*Distance, error) {
	m, err := d.convertTo(unit)
	if err != nil {
		return nil, err
	}

	return &Distance{*m}, nil
}

func NewDistance(value float64, unit string) (*Distance, error) {
	if _, ok := distanceConversionFactors.HasFactor(unit); !ok {
		return nil, fmt.Errorf("unknown Distance unit %s", unit)
	}

	m, err := newBaseMeasurement(value, unit, distanceConversionFactors)
	if err != nil {
		return nil, err
	}

	return &Distance{*m}, nil
}

var distanceConversionErr = errors.New("not a distance measure")

func ToDistance(m Measurement) (*Distance, error) {
	if m.Type() == MeasureDistance {
		return &Distance{*(m.(*BaseMeasurement))}, nil
	}

	return nil, distanceConversionErr
}
