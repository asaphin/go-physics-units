package units

import (
	"errors"

	"github.com/asaphin/go-physics-units/distance"
	"github.com/asaphin/go-physics-units/velocity"
)

// Velocity is a velocity (speed) measurement. Its unit symbols are defined in the
// velocity package.
type Velocity interface {
	Measurement
	UnitConverter[Velocity]
	Arithmetics[Velocity]

	// MultiplyByTime returns the distance v·t in distance.BaseUnit (m).
	MultiplyByTime(t Time) Distance
}

// velocityImplementation is the concrete implementation of Velocity.
type velocityImplementation struct {
	baseMeasurement
}

// ConvertTo returns the velocity expressed in unit.
func (v *velocityImplementation) ConvertTo(unit string) (Velocity, error) {
	msm, err := v.convertTo(unit)
	if err != nil {
		return nil, err
	}

	return &velocityImplementation{msm}, nil
}

// MustConvertTo is like ConvertTo but panics if unit isn't a velocity unit.
func (v *velocityImplementation) MustConvertTo(unit string) Velocity {
	msm, err := v.convertTo(unit)
	if err != nil {
		panic(err)
	}

	return &velocityImplementation{msm}
}

// ConvertToBaseUnits returns the velocity expressed in metres per second.
func (v *velocityImplementation) ConvertToBaseUnits() Velocity {
	return &velocityImplementation{v.unsafeConvertTo(velocity.BaseUnit)}
}

func (v *velocityImplementation) valueInBaseUnits() float64 {
	return v.unsafeGetValueIn(velocity.BaseUnit)
}

func (v *velocityImplementation) valueInUnits(unit string) float64 {
	return v.unsafeGetValueIn(unit)
}

// Add returns the sum of the velocities, expressed in the receiver's unit.
func (v *velocityImplementation) Add(other Velocity) Velocity {
	return &velocityImplementation{v.withValue(v.value + other.valueInUnits(v.unit))}
}

// Sub returns the difference of the velocities, expressed in the receiver's unit.
func (v *velocityImplementation) Sub(other Velocity) Velocity {
	return &velocityImplementation{v.withValue(v.value - other.valueInUnits(v.unit))}
}

// Mul returns the velocity multiplied by multiplier.
func (v *velocityImplementation) Mul(multiplier float64) Velocity {
	return &velocityImplementation{v.withValue(v.value * multiplier)}
}

// Div returns the velocity divided by divisor. It fails if divisor is zero.
func (v *velocityImplementation) Div(divisor float64) (Velocity, error) {
	if divisor == 0 {
		return nil, errZeroDivision
	}

	return &velocityImplementation{v.withValue(v.value / divisor)}, nil
}

// MultiplyByTime returns the distance v·t in m.
func (v *velocityImplementation) MultiplyByTime(t Time) Distance {
	return newDistance(v.valueInBaseUnits()*t.valueInBaseUnits(), distance.BaseUnit)
}

// NewVelocity creates a velocity of value in unit, one of the symbols defined in
// the velocity package.
func NewVelocity(value float64, unit string) (Velocity, error) {
	msm, err := newBaseMeasurement(value, unit, velocity.ConversionRates())
	if err != nil {
		return nil, err
	}

	return &velocityImplementation{msm}, nil
}

// newVelocity is an unchecked constructor for internal use with valid,
// canonical units.
func newVelocity(value float64, unit string) Velocity {
	return &velocityImplementation{baseMeasurement{value: value, unit: unit, conversionRates: velocity.ConversionRates()}}
}

var errVelocityConversion = errors.New("not a velocity measure")

// ToVelocity converts m, for example a result of ParseString, into a Velocity.
// It fails if m isn't a velocity.
func ToVelocity(m Measurement) (Velocity, error) {
	if m == nil || m.Type() != MeasureVelocity {
		return nil, errVelocityConversion
	}

	return &velocityImplementation{*m.base()}, nil
}

// StringToVelocity parses a velocity from a string such as "10 km/h".
func StringToVelocity(s string) (Velocity, error) {
	m, err := ParseString(s)
	if err != nil {
		return nil, err
	}

	return ToVelocity(m)
}
