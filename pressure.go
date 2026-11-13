package units

import (
	"cmp"

	"github.com/asaphin/go-physics-units/force"
	"github.com/asaphin/go-physics-units/pressure"
)

// Pressure is a pressure measurement. Its unit symbols are defined in the pressure
// package.
type Pressure interface {
	Measurement
	UnitConverter[Pressure]
	Arithmetics[Pressure]
	Comparer[Pressure]

	// MultiplyByArea returns the force p·a in force.BaseUnit (N).
	MultiplyByArea(a Area) Force
}

// pressureImplementation is the concrete implementation of Pressure.
type pressureImplementation struct {
	baseMeasurement
}

// ConvertTo returns the pressure expressed in unit.
func (p *pressureImplementation) ConvertTo(unit string) (Pressure, error) {
	msm, err := p.convertTo(unit)
	if err != nil {
		return nil, err
	}

	return &pressureImplementation{msm}, nil
}

// MustConvertTo is like ConvertTo but panics if unit isn't a pressure unit.
func (p *pressureImplementation) MustConvertTo(unit string) Pressure {
	msm, err := p.convertTo(unit)
	if err != nil {
		panic(err)
	}

	return &pressureImplementation{msm}
}

// ConvertToBaseUnits returns the pressure expressed in pascals.
func (p *pressureImplementation) ConvertToBaseUnits() Pressure {
	return &pressureImplementation{p.unsafeConvertTo(pressure.BaseUnit)}
}

func (p *pressureImplementation) valueInBaseUnits() float64 {
	return p.unsafeGetValueIn(pressure.BaseUnit)
}

func (p *pressureImplementation) valueInUnits(unit string) float64 {
	return p.unsafeGetValueIn(unit)
}

// Add returns the sum of the pressures, expressed in the receiver's unit.
func (p *pressureImplementation) Add(other Pressure) Pressure {
	return &pressureImplementation{p.withValue(p.value + other.valueInUnits(p.unit))}
}

// Sub returns the difference of the pressures, expressed in the receiver's unit.
func (p *pressureImplementation) Sub(other Pressure) Pressure {
	return &pressureImplementation{p.withValue(p.value - other.valueInUnits(p.unit))}
}

// Mul returns the pressure multiplied by multiplier.
func (p *pressureImplementation) Mul(multiplier float64) Pressure {
	return &pressureImplementation{p.withValue(p.value * multiplier)}
}

// Div returns the pressure divided by divisor. It fails if divisor is zero.
func (p *pressureImplementation) Div(divisor float64) (Pressure, error) {
	if divisor == 0 {
		return nil, ErrZeroDivision
	}

	return &pressureImplementation{p.withValue(p.value / divisor)}, nil
}

// MultiplyByArea returns the force p·a in N.
func (p *pressureImplementation) MultiplyByArea(a Area) Force {
	return newForce(p.valueInBaseUnits()*a.valueInBaseUnits(), force.BaseUnit)
}

// Compare compares the pressures; see Comparer.
func (p *pressureImplementation) Compare(other Pressure) int {
	return cmp.Compare(p.commonValues(other.base()))
}

// Equal reports whether the pressures are equal; see Comparer.
func (p *pressureImplementation) Equal(other Pressure) bool {
	return p.Compare(other) == 0
}

// EqualWithinRel reports whether the pressures are equal within the relative
// tolerance tol; see Comparer.
func (p *pressureImplementation) EqualWithinRel(other Pressure, tol float64) bool {
	x, y := p.commonValues(other.base())

	return equalWithinRel(x, y, tol)
}

// NewPressure creates a pressure of value in unit, one of the symbols defined in
// the pressure package.
func NewPressure(value float64, unit string) (Pressure, error) {
	msm, err := newBaseMeasurement(value, unit, pressureRates)
	if err != nil {
		return nil, err
	}

	return &pressureImplementation{msm}, nil
}

// newPressure is an unchecked constructor for internal use with valid, canonical
// units.
func newPressure(value float64, unit string) Pressure {
	return &pressureImplementation{baseMeasurement{value: value, unit: unit, conversionRates: pressureRates}}
}

const errPressureConversion incompatibleError = "not a pressure measure"

// ToPressure converts m, for example a result of ParseString, into a Pressure.
// It fails if m isn't a pressure.
func ToPressure(m Measurement) (Pressure, error) {
	if m == nil || m.Type() != MeasurePressure {
		return nil, errPressureConversion
	}

	return &pressureImplementation{*m.base()}, nil
}

// StringToPressure parses a pressure from a string such as "101.3 kPa".
// It fails if the unit isn't a pressure unit.
func StringToPressure(s string) (Pressure, error) {
	// Parse and construct directly: one allocation, where ParseString
	// followed by ToPressure would allocate twice.
	value, unit, err := parseValueAndUnit(s)
	if err != nil {
		return nil, err
	}

	return NewPressure(value, unit)
}
