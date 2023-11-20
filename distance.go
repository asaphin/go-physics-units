package units

import (
	"errors"
	"fmt"
	"github.com/asaphin/go-physics-units/distance"
	"github.com/asaphin/go-physics-units/velocity"
)

var distanceConversionFactors = newImmutableConversionFactors(distance.ConversionFactors())

// Distance interface represents a distance measurement.
type Distance interface {
	Measurement
	ConvertTo(unit string) (Distance, error)
	MustConvertTo(unit string) Distance
	ConvertToBaseUnits() Distance
	DivideByTime(t Time) Velocity
}

// distanceImplementation is a concrete implementation of the Distance interface.
type distanceImplementation struct {
	baseMeasurement
}

// NewDistance creates a new Distance instance.
func NewDistance(value float64, unit string) (Distance, error) {
	if _, ok := distanceConversionFactors.HasFactor(unit); !ok {
		return nil, fmt.Errorf("unknown Distance unit %s", unit)
	}

	m, err := newBaseMeasurement(value, unit, distanceConversionFactors)
	if err != nil {
		return nil, err
	}

	return &distanceImplementation{*m}, nil
}

// ConvertTo implements the ConvertTo method of the Distance interface.
func (d *distanceImplementation) ConvertTo(unit string) (Distance, error) {
	m, err := d.convertTo(unit)
	if err != nil {
		return nil, err
	}

	return &distanceImplementation{*m}, nil
}

// MustConvertTo implements the MustConvertTo method of the Distance interface.
func (d *distanceImplementation) MustConvertTo(unit string) Distance {
	m, err := d.convertTo(unit)
	if err != nil {
		panic(err)
	}

	return &distanceImplementation{*m}
}

// ConvertToBaseUnits implements the MustConvertToBaseUnits method of the Distance interface.
func (d *distanceImplementation) ConvertToBaseUnits() Distance {
	m, err := d.convertTo(distance.BaseUnit)
	if err != nil {
		panic(err)
	}

	return &distanceImplementation{*m}
}

// DivideByTime implements the DivideByTime method of the Distance interface.
func (d *distanceImplementation) DivideByTime(t Time) Velocity {
	if t.Value() == 0 {
		panic("division by zero")
	}

	baseD := d.ConvertToBaseUnits()
	baseT := t.ConvertToBaseUnits()

	v, _ := NewVelocity(baseD.Value()/baseT.Value(), velocity.BaseUnit)

	return v
}

var distanceConversionErr = errors.New("not a distance measure")

func ToDistance(m Measurement) (Distance, error) {
	if m.Type() == MeasureDistance {
		return &distanceImplementation{*m.(*baseMeasurement)}, nil
	}

	return nil, distanceConversionErr
}
