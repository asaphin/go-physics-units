package units

import (
	"errors"
	"fmt"
	"github.com/asaphin/go-physics-units/distance"
	"github.com/asaphin/go-physics-units/velocity"
)

// Distance interface represents a distance measurement.
type Distance interface {
	Measurement
	UnitConverter[Distance]
	Arithmetics[Distance]

	DivideByTime(t Time) Velocity
}

// distanceImplementation is a concrete implementation of the Distance interface.
type distanceImplementation struct {
	baseMeasurement
}

// ConvertTo implements the ConvertTo method of the Distance interface.
func (d *distanceImplementation) ConvertTo(unit string) (Distance, error) {
	msm, err := d.convertTo(unit)
	if err != nil {
		return nil, err
	}

	return &distanceImplementation{*msm}, nil
}

// MustConvertTo implements the MustConvertTo method of the Distance interface.
func (d *distanceImplementation) MustConvertTo(unit string) Distance {
	msm, err := d.convertTo(unit)
	if err != nil {
		panic(err)
	}

	return &distanceImplementation{*msm}
}

// ConvertToBaseUnits implements the MustConvertToBaseUnits method of the Distance interface.
func (d *distanceImplementation) ConvertToBaseUnits() Distance {
	msm := d.unsafeConvertTo(distance.BaseUnit)

	return &distanceImplementation{*msm}
}

func (d *distanceImplementation) valueInBaseUnits() float64 {
	return d.unsafeGetValueIn(distance.BaseUnit)
}

func (d *distanceImplementation) valueInUnits(unit string) float64 {
	return d.unsafeGetValueIn(unit)
}

func (d *distanceImplementation) Add(dist Distance) Distance {
	d1 := d.value
	d2 := dist.valueInUnits(d.unit)

	return newDistance(d1+d2, d.unit)
}

func (d *distanceImplementation) Sub(dist Distance) Distance {
	d1 := d.value
	d2 := dist.valueInUnits(d.unit)

	return newDistance(d1-d2, d.unit)
}

func (d *distanceImplementation) Mul(multiplier float64) Distance {
	msm := newUnsafeBaseMeasurement(d.value*multiplier, d.unit, d.conversionRates)

	return &distanceImplementation{*msm}
}

func (d *distanceImplementation) Div(divisor float64) (Distance, error) {
	if divisor == 0 {
		return nil, errZeroDivision
	}

	msm := newUnsafeBaseMeasurement(d.value/divisor, d.unit, d.conversionRates)

	return &distanceImplementation{*msm}, nil
}

// DivideByTime implements the DivideByTime method of the Distance interface.
func (d *distanceImplementation) DivideByTime(t Time) Velocity {
	if t.Value() == 0 {
		panic("division by zero")
	}

	baseD := d.valueInBaseUnits()
	baseT := t.valueInBaseUnits()

	return newVelocity(baseD/baseT, velocity.BaseUnit)
}

// NewDistance creates a new Distance instance.
func NewDistance(value float64, unit string) (Distance, error) {
	rates := distance.ConversionRates()

	if _, ok := rates.Has(unit + unit); !ok {
		return nil, fmt.Errorf("unknown Distance unit %s", unit)
	}

	msm := newUnsafeBaseMeasurement(value, unit, rates)

	return &distanceImplementation{*msm}, nil
}

// newDistance unsafe distance constructor without checks for internal usage.
func newDistance(value float64, unit string) Distance {
	msm := newUnsafeBaseMeasurement(value, unit, distance.ConversionRates())

	return &distanceImplementation{*msm}
}

var errDistanceConversion = errors.New("not a distance measure")

func ToDistance(m Measurement) (Distance, error) {
	if m.Type() == MeasureDistance {
		b, ok := m.(*baseMeasurement)
		if !ok {
			return nil, errBaseMeasurementConversion
		}

		return &distanceImplementation{*b}, nil
	}

	return nil, errDistanceConversion
}

func StringToDistance(s string) (Distance, error) {
	m, err := ParseString(s)
	if err != nil {
		return nil, err
	}

	return ToDistance(m)
}
