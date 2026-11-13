package units

import (
	"cmp"

	"github.com/asaphin/go-physics-units/acceleration"
	"github.com/asaphin/go-physics-units/force"
	"github.com/asaphin/go-physics-units/velocity"
)

// Acceleration is an acceleration measurement. Its unit symbols are defined in the acceleration
// package.
type Acceleration interface {
	Measurement
	UnitConverter[Acceleration]
	Arithmetics[Acceleration]
	Comparer[Acceleration]

	// MultiplyByTime returns the velocity a·t in velocity.BaseUnit (m/s).
	MultiplyByTime(t Time) Velocity
	// MultiplyByMass returns the force m·a in force.BaseUnit (N).
	MultiplyByMass(m Mass) Force
}

// accelerationImplementation is the concrete implementation of Acceleration.
type accelerationImplementation struct {
	baseMeasurement
}

// ConvertTo returns the acceleration expressed in unit.
func (a *accelerationImplementation) ConvertTo(unit string) (Acceleration, error) {
	msm, err := a.convertTo(unit)
	if err != nil {
		return nil, err
	}

	return &accelerationImplementation{msm}, nil
}

// MustConvertTo is like ConvertTo but panics if unit isn't an acceleration unit.
func (a *accelerationImplementation) MustConvertTo(unit string) Acceleration {
	msm, err := a.convertTo(unit)
	if err != nil {
		panic(err)
	}

	return &accelerationImplementation{msm}
}

// ConvertToBaseUnits returns the acceleration expressed in metres per second squared.
func (a *accelerationImplementation) ConvertToBaseUnits() Acceleration {
	return &accelerationImplementation{a.unsafeConvertTo(acceleration.BaseUnit)}
}

func (a *accelerationImplementation) valueInBaseUnits() float64 {
	return a.unsafeGetValueIn(acceleration.BaseUnit)
}

func (a *accelerationImplementation) valueInUnits(unit string) float64 {
	return a.unsafeGetValueIn(unit)
}

// Add returns the sum of the accelerations, expressed in the receiver's unit.
func (a *accelerationImplementation) Add(other Acceleration) Acceleration {
	return &accelerationImplementation{a.withValue(a.value + other.valueInUnits(a.unit))}
}

// Sub returns the difference of the accelerations, expressed in the receiver's unit.
func (a *accelerationImplementation) Sub(other Acceleration) Acceleration {
	return &accelerationImplementation{a.withValue(a.value - other.valueInUnits(a.unit))}
}

// Mul returns the acceleration multiplied by multiplier.
func (a *accelerationImplementation) Mul(multiplier float64) Acceleration {
	return &accelerationImplementation{a.withValue(a.value * multiplier)}
}

// Div returns the acceleration divided by divisor. It fails if divisor is zero.
func (a *accelerationImplementation) Div(divisor float64) (Acceleration, error) {
	if divisor == 0 {
		return nil, ErrZeroDivision
	}

	return &accelerationImplementation{a.withValue(a.value / divisor)}, nil
}

// MultiplyByTime returns the velocity a·t in m/s.
func (a *accelerationImplementation) MultiplyByTime(t Time) Velocity {
	return newVelocity(a.valueInBaseUnits()*t.valueInBaseUnits(), velocity.BaseUnit)
}

// MultiplyByMass returns the force m·a in N.
func (a *accelerationImplementation) MultiplyByMass(m Mass) Force {
	return newForce(m.valueInBaseUnits()*a.valueInBaseUnits(), force.BaseUnit)
}

// Compare compares the accelerations; see Comparer.
func (a *accelerationImplementation) Compare(other Acceleration) int {
	return cmp.Compare(a.commonValues(other.base()))
}

// Equal reports whether the accelerations are equal; see Comparer.
func (a *accelerationImplementation) Equal(other Acceleration) bool {
	return a.Compare(other) == 0
}

// EqualWithinRel reports whether the accelerations are equal within the relative
// tolerance tol; see Comparer.
func (a *accelerationImplementation) EqualWithinRel(other Acceleration, tol float64) bool {
	x, y := a.commonValues(other.base())

	return equalWithinRel(x, y, tol)
}

// NewAcceleration creates an acceleration of value in unit, one of the symbols defined in
// the acceleration package.
func NewAcceleration(value float64, unit string) (Acceleration, error) {
	msm, err := newBaseMeasurement(value, unit, accelerationRates)
	if err != nil {
		return nil, err
	}

	return &accelerationImplementation{msm}, nil
}

// newAcceleration is an unchecked constructor for internal use with valid, canonical
// units.
func newAcceleration(value float64, unit string) Acceleration {
	return &accelerationImplementation{baseMeasurement{value: value, unit: unit, conversionRates: accelerationRates}}
}

const errAccelerationConversion incompatibleError = "not an acceleration measure"

// ToAcceleration converts m, for example a result of ParseString, into an Acceleration.
// It fails if m isn't an acceleration.
func ToAcceleration(m Measurement) (Acceleration, error) {
	if m == nil || m.Type() != MeasureAcceleration {
		return nil, errAccelerationConversion
	}

	return &accelerationImplementation{*m.base()}, nil
}

// StringToAcceleration parses an acceleration from a string such as "9.8 m/s²".
// It fails if the unit isn't an acceleration unit.
func StringToAcceleration(s string) (Acceleration, error) {
	// Parse and construct directly: one allocation, where ParseString
	// followed by ToAcceleration would allocate twice.
	value, unit, err := parseValueAndUnit(s)
	if err != nil {
		return nil, err
	}

	return NewAcceleration(value, unit)
}
