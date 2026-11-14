package units

import (
	"cmp"

	"github.com/asaphin/go-physics-units/acceleration"
	"github.com/asaphin/go-physics-units/energy"
	"github.com/asaphin/go-physics-units/force"
	"github.com/asaphin/go-physics-units/mass"
	"github.com/asaphin/go-physics-units/pressure"
)

// Force is a force measurement. Its unit symbols are defined in the force
// package.
type Force interface {
	Measurement
	UnitConverter[Force]
	Arithmetics[Force]
	Comparer[Force]

	// DivideByMass returns the acceleration f/m in acceleration.BaseUnit
	// (m/s²). It fails if m is zero.
	DivideByMass(m Mass) (Acceleration, error)
	// DivideByAcceleration returns the mass f/a in mass.BaseUnit (kg).
	// It fails if a is zero.
	DivideByAcceleration(a Acceleration) (Mass, error)
	// DivideByArea returns the pressure f/a in pressure.BaseUnit (Pa).
	// It fails if a is zero.
	DivideByArea(a Area) (Pressure, error)
	// MultiplyByDistance returns the energy (work) f·d in energy.BaseUnit (J).
	MultiplyByDistance(d Distance) Energy
}

// forceImplementation is the concrete implementation of Force.
type forceImplementation struct {
	baseMeasurement
}

// ConvertTo returns the force expressed in unit.
func (f *forceImplementation) ConvertTo(unit string) (Force, error) {
	msm, err := f.convertTo(unit)
	if err != nil {
		return nil, err
	}

	return &forceImplementation{msm}, nil
}

// MustConvertTo is like ConvertTo but panics if unit isn't a force unit.
func (f *forceImplementation) MustConvertTo(unit string) Force {
	msm, err := f.convertTo(unit)
	if err != nil {
		panic(err)
	}

	return &forceImplementation{msm}
}

// ConvertToBaseUnits returns the force expressed in newtons.
func (f *forceImplementation) ConvertToBaseUnits() Force {
	return &forceImplementation{f.unsafeConvertTo(force.BaseUnit)}
}

func (f *forceImplementation) valueInBaseUnits() float64 {
	return f.unsafeGetValueIn(force.BaseUnit)
}

func (f *forceImplementation) valueInUnits(unit string) float64 {
	return f.unsafeGetValueIn(unit)
}

// Add returns the sum of the forces, expressed in the receiver's unit.
func (f *forceImplementation) Add(other Force) Force {
	return &forceImplementation{f.withValue(f.value + other.valueInUnits(f.unit))}
}

// Sub returns the difference of the forces, expressed in the receiver's unit.
func (f *forceImplementation) Sub(other Force) Force {
	return &forceImplementation{f.withValue(f.value - other.valueInUnits(f.unit))}
}

// Mul returns the force multiplied by multiplier.
func (f *forceImplementation) Mul(multiplier float64) Force {
	return &forceImplementation{f.withValue(f.value * multiplier)}
}

// Div returns the force divided by divisor. It fails if divisor is zero.
func (f *forceImplementation) Div(divisor float64) (Force, error) {
	if divisor == 0 {
		return nil, ErrZeroDivision
	}

	return &forceImplementation{f.withValue(f.value / divisor)}, nil
}

// DivideByMass returns the acceleration f/m in m/s², or an error if m is zero.
func (f *forceImplementation) DivideByMass(m Mass) (Acceleration, error) {
	if m.Value() == 0 {
		return nil, ErrZeroDivision
	}

	return newAcceleration(f.valueInBaseUnits()/m.valueInBaseUnits(), acceleration.BaseUnit), nil
}

// DivideByAcceleration returns the mass f/a in kg, or an error if a is zero.
func (f *forceImplementation) DivideByAcceleration(a Acceleration) (Mass, error) {
	if a.Value() == 0 {
		return nil, ErrZeroDivision
	}

	return newMass(f.valueInBaseUnits()/a.valueInBaseUnits(), mass.BaseUnit), nil
}

// DivideByArea returns the pressure f/a in Pa, or an error if a is zero.
func (f *forceImplementation) DivideByArea(a Area) (Pressure, error) {
	if a.Value() == 0 {
		return nil, ErrZeroDivision
	}

	return newPressure(f.valueInBaseUnits()/a.valueInBaseUnits(), pressure.BaseUnit), nil
}

// MultiplyByDistance returns the energy f·d in J.
func (f *forceImplementation) MultiplyByDistance(d Distance) Energy {
	return newEnergy(f.valueInBaseUnits()*d.valueInBaseUnits(), energy.BaseUnit)
}

// Compare compares the forces; see Comparer.
func (f *forceImplementation) Compare(other Force) int {
	return cmp.Compare(f.commonValues(other.base()))
}

// Equal reports whether the forces are equal; see Comparer.
func (f *forceImplementation) Equal(other Force) bool {
	return f.Compare(other) == 0
}

// EqualWithinRel reports whether the forces are equal within the relative
// tolerance tol; see Comparer.
func (f *forceImplementation) EqualWithinRel(other Force, tol float64) bool {
	x, y := f.commonValues(other.base())

	return equalWithinRel(x, y, tol)
}

// NewForce creates a force of value in unit, one of the symbols defined in
// the force package.
func NewForce(value float64, unit string) (Force, error) {
	msm, err := newBaseMeasurement(value, unit, forceRates)
	if err != nil {
		return nil, err
	}

	return &forceImplementation{msm}, nil
}

// newForce is an unchecked constructor for internal use with valid, canonical
// units.
//
//nolint:unparam // every caller passes the base unit, but the signature matches the other unchecked constructors
func newForce(value float64, unit string) Force {
	return &forceImplementation{baseMeasurement{value: value, unit: unit, conversionRates: forceRates}}
}

const errForceConversion incompatibleError = "not a force measure"

// ToForce converts m, for example a result of ParseString, into a Force.
// It fails if m isn't a force.
func ToForce(m Measurement) (Force, error) {
	if m == nil || m.Type() != MeasureForce {
		return nil, errForceConversion
	}

	return &forceImplementation{*m.base()}, nil
}

// StringToForce parses a force from a string such as "10 N".
// It fails if the unit isn't a force unit.
func StringToForce(s string) (Force, error) {
	// Parse and construct directly: one allocation, where ParseString
	// followed by ToForce would allocate twice.
	value, unit, err := parseValueAndUnit(s)
	if err != nil {
		return nil, err
	}

	return NewForce(value, unit)
}
