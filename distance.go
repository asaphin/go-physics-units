package units

import (
	"errors"

	"github.com/asaphin/go-physics-units/area"
	"github.com/asaphin/go-physics-units/distance"
	"github.com/asaphin/go-physics-units/velocity"
)

// Distance is a length measurement. Its unit symbols are defined in the
// distance package.
type Distance interface {
	Measurement
	UnitConverter[Distance]
	Arithmetics[Distance]

	// DivideByTime returns the velocity d/t in velocity.BaseUnit (m/s).
	// It panics if t is zero.
	DivideByTime(t Time) Velocity
	// Square returns the area d² in area.BaseUnit (m²).
	Square() Area
}

// distanceImplementation is the concrete implementation of Distance.
type distanceImplementation struct {
	baseMeasurement
}

// ConvertTo returns the distance expressed in unit.
func (d *distanceImplementation) ConvertTo(unit string) (Distance, error) {
	msm, err := d.convertTo(unit)
	if err != nil {
		return nil, err
	}

	return &distanceImplementation{msm}, nil
}

// MustConvertTo is like ConvertTo but panics if unit isn't a distance unit.
func (d *distanceImplementation) MustConvertTo(unit string) Distance {
	msm, err := d.convertTo(unit)
	if err != nil {
		panic(err)
	}

	return &distanceImplementation{msm}
}

// ConvertToBaseUnits returns the distance expressed in metres.
func (d *distanceImplementation) ConvertToBaseUnits() Distance {
	return &distanceImplementation{d.unsafeConvertTo(distance.BaseUnit)}
}

func (d *distanceImplementation) valueInBaseUnits() float64 {
	return d.unsafeGetValueIn(distance.BaseUnit)
}

func (d *distanceImplementation) valueInUnits(unit string) float64 {
	return d.unsafeGetValueIn(unit)
}

// Add returns the sum of the distances, expressed in the receiver's unit.
func (d *distanceImplementation) Add(other Distance) Distance {
	return &distanceImplementation{d.withValue(d.value + other.valueInUnits(d.unit))}
}

// Sub returns the difference of the distances, expressed in the receiver's unit.
func (d *distanceImplementation) Sub(other Distance) Distance {
	return &distanceImplementation{d.withValue(d.value - other.valueInUnits(d.unit))}
}

// Mul returns the distance multiplied by multiplier.
func (d *distanceImplementation) Mul(multiplier float64) Distance {
	return &distanceImplementation{d.withValue(d.value * multiplier)}
}

// Div returns the distance divided by divisor. It fails if divisor is zero.
func (d *distanceImplementation) Div(divisor float64) (Distance, error) {
	if divisor == 0 {
		return nil, errZeroDivision
	}

	return &distanceImplementation{d.withValue(d.value / divisor)}, nil
}

// DivideByTime returns the velocity d/t in m/s. It panics if t is zero.
func (d *distanceImplementation) DivideByTime(t Time) Velocity {
	if t.Value() == 0 {
		panic(errZeroDivision)
	}

	return newVelocity(d.valueInBaseUnits()/t.valueInBaseUnits(), velocity.BaseUnit)
}

// Square returns the area d² in m².
func (d *distanceImplementation) Square() Area {
	baseD := d.valueInBaseUnits()

	return newArea(baseD*baseD, area.BaseUnit)
}

// NewDistance creates a distance of value in unit, one of the symbols defined in
// the distance package.
func NewDistance(value float64, unit string) (Distance, error) {
	msm, err := newBaseMeasurement(value, unit, distance.ConversionRates())
	if err != nil {
		return nil, err
	}

	return &distanceImplementation{msm}, nil
}

// newDistance is an unchecked constructor for internal use with valid,
// canonical units.
func newDistance(value float64, unit string) Distance {
	return &distanceImplementation{baseMeasurement{value: value, unit: unit, conversionRates: distance.ConversionRates()}}
}

var errDistanceConversion = errors.New("not a distance measure")

// ToDistance converts m, for example a result of ParseString, into a Distance.
// It fails if m isn't a distance.
func ToDistance(m Measurement) (Distance, error) {
	if m == nil || m.Type() != MeasureDistance {
		return nil, errDistanceConversion
	}

	return &distanceImplementation{*m.base()}, nil
}

// StringToDistance parses a distance from a string such as "10 km".
func StringToDistance(s string) (Distance, error) {
	m, err := ParseString(s)
	if err != nil {
		return nil, err
	}

	return ToDistance(m)
}
