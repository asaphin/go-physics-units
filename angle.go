package units

import (
	"cmp"

	"github.com/asaphin/go-physics-units/angle"
)

// Angle is an angle measurement. Its unit symbols are defined in the angle
// package.
type Angle interface {
	Measurement
	UnitConverter[Angle]
	Arithmetics[Angle]
	Comparer[Angle]
}

// angleImplementation is the concrete implementation of Angle.
type angleImplementation struct {
	baseMeasurement
}

// ConvertTo returns the angle expressed in unit.
func (a *angleImplementation) ConvertTo(unit string) (Angle, error) {
	msm, err := a.convertTo(unit)
	if err != nil {
		return nil, err
	}

	return &angleImplementation{msm}, nil
}

// MustConvertTo is like ConvertTo but panics if unit isn't an angle unit.
func (a *angleImplementation) MustConvertTo(unit string) Angle {
	msm, err := a.convertTo(unit)
	if err != nil {
		panic(err)
	}

	return &angleImplementation{msm}
}

// ConvertToBaseUnits returns the angle expressed in radians.
func (a *angleImplementation) ConvertToBaseUnits() Angle {
	return &angleImplementation{a.unsafeConvertTo(angle.BaseUnit)}
}

func (a *angleImplementation) valueInBaseUnits() float64 {
	return a.unsafeGetValueIn(angle.BaseUnit)
}

func (a *angleImplementation) valueInUnits(unit string) float64 {
	return a.unsafeGetValueIn(unit)
}

// Add returns the sum of the angles, expressed in the receiver's unit.
func (a *angleImplementation) Add(other Angle) Angle {
	return &angleImplementation{a.withValue(a.value + other.valueInUnits(a.unit))}
}

// Sub returns the difference of the angles, expressed in the receiver's unit.
func (a *angleImplementation) Sub(other Angle) Angle {
	return &angleImplementation{a.withValue(a.value - other.valueInUnits(a.unit))}
}

// Mul returns the angle multiplied by multiplier.
func (a *angleImplementation) Mul(multiplier float64) Angle {
	return &angleImplementation{a.withValue(a.value * multiplier)}
}

// Div returns the angle divided by divisor. It fails if divisor is zero.
func (a *angleImplementation) Div(divisor float64) (Angle, error) {
	if divisor == 0 {
		return nil, ErrZeroDivision
	}

	return &angleImplementation{a.withValue(a.value / divisor)}, nil
}

// Compare compares the angles; see Comparer.
func (a *angleImplementation) Compare(other Angle) int {
	return cmp.Compare(a.commonValues(other.base()))
}

// Equal reports whether the angles are equal; see Comparer.
func (a *angleImplementation) Equal(other Angle) bool {
	return a.Compare(other) == 0
}

// EqualWithinRel reports whether the angles are equal within the relative
// tolerance tol; see Comparer.
func (a *angleImplementation) EqualWithinRel(other Angle, tol float64) bool {
	x, y := a.commonValues(other.base())

	return equalWithinRel(x, y, tol)
}

// NewAngle creates an angle of value in unit, one of the symbols defined in
// the angle package.
func NewAngle(value float64, unit string) (Angle, error) {
	msm, err := newBaseMeasurement(value, unit, angleRates)
	if err != nil {
		return nil, err
	}

	return &angleImplementation{msm}, nil
}

const errAngleConversion incompatibleError = "not an angle measure"

// ToAngle converts m, for example a result of ParseString, into an Angle.
// It fails if m isn't an angle.
func ToAngle(m Measurement) (Angle, error) {
	if m == nil || m.Type() != MeasureAngle {
		return nil, errAngleConversion
	}

	return &angleImplementation{*m.base()}, nil
}

// StringToAngle parses an angle from a string such as "90 °".
// It fails if the unit isn't an angle unit.
func StringToAngle(s string) (Angle, error) {
	// Parse and construct directly: one allocation, where ParseString
	// followed by ToAngle would allocate twice.
	value, unit, err := parseValueAndUnit(s)
	if err != nil {
		return nil, err
	}

	return NewAngle(value, unit)
}
