package units

import (
	"cmp"

	"github.com/asaphin/go-physics-units/energy"
	"github.com/asaphin/go-physics-units/power"
)

// Power is a power measurement. Its unit symbols are defined in the power
// package.
type Power interface {
	Measurement
	UnitConverter[Power]
	Arithmetics[Power]
	Comparer[Power]

	// MultiplyByTime returns the energy p·t in energy.BaseUnit (J).
	MultiplyByTime(t Time) Energy
}

// powerImplementation is the concrete implementation of Power.
type powerImplementation struct {
	baseMeasurement
}

// ConvertTo returns the power expressed in unit.
func (p *powerImplementation) ConvertTo(unit string) (Power, error) {
	msm, err := p.convertTo(unit)
	if err != nil {
		return nil, err
	}

	return &powerImplementation{msm}, nil
}

// MustConvertTo is like ConvertTo but panics if unit isn't a power unit.
func (p *powerImplementation) MustConvertTo(unit string) Power {
	msm, err := p.convertTo(unit)
	if err != nil {
		panic(err)
	}

	return &powerImplementation{msm}
}

// ConvertToBaseUnits returns the power expressed in watts.
func (p *powerImplementation) ConvertToBaseUnits() Power {
	return &powerImplementation{p.unsafeConvertTo(power.BaseUnit)}
}

func (p *powerImplementation) valueInBaseUnits() float64 {
	return p.unsafeGetValueIn(power.BaseUnit)
}

func (p *powerImplementation) valueInUnits(unit string) float64 {
	return p.unsafeGetValueIn(unit)
}

// Add returns the sum of the powers, expressed in the receiver's unit.
func (p *powerImplementation) Add(other Power) Power {
	return &powerImplementation{p.withValue(p.value + other.valueInUnits(p.unit))}
}

// Sub returns the difference of the powers, expressed in the receiver's unit.
func (p *powerImplementation) Sub(other Power) Power {
	return &powerImplementation{p.withValue(p.value - other.valueInUnits(p.unit))}
}

// Mul returns the power multiplied by multiplier.
func (p *powerImplementation) Mul(multiplier float64) Power {
	return &powerImplementation{p.withValue(p.value * multiplier)}
}

// Div returns the power divided by divisor. It fails if divisor is zero.
func (p *powerImplementation) Div(divisor float64) (Power, error) {
	if divisor == 0 {
		return nil, ErrZeroDivision
	}

	return &powerImplementation{p.withValue(p.value / divisor)}, nil
}

// MultiplyByTime returns the energy p·t in J.
func (p *powerImplementation) MultiplyByTime(t Time) Energy {
	return newEnergy(p.valueInBaseUnits()*t.valueInBaseUnits(), energy.BaseUnit)
}

// Compare compares the powers; see Comparer.
func (p *powerImplementation) Compare(other Power) int {
	return cmp.Compare(p.commonValues(other.base()))
}

// Equal reports whether the powers are equal; see Comparer.
func (p *powerImplementation) Equal(other Power) bool {
	return p.Compare(other) == 0
}

// EqualWithinRel reports whether the powers are equal within the relative
// tolerance tol; see Comparer.
func (p *powerImplementation) EqualWithinRel(other Power, tol float64) bool {
	x, y := p.commonValues(other.base())

	return equalWithinRel(x, y, tol)
}

// NewPower creates a power of value in unit, one of the symbols defined in
// the power package.
func NewPower(value float64, unit string) (Power, error) {
	msm, err := newBaseMeasurement(value, unit, powerRates)
	if err != nil {
		return nil, err
	}

	return &powerImplementation{msm}, nil
}

// newPower is an unchecked constructor for internal use with valid, canonical
// units.
func newPower(value float64, unit string) Power {
	return &powerImplementation{baseMeasurement{value: value, unit: unit, conversionRates: powerRates}}
}

const errPowerConversion incompatibleError = "not a power measure"

// ToPower converts m, for example a result of ParseString, into a Power.
// It fails if m isn't a power.
func ToPower(m Measurement) (Power, error) {
	if m == nil || m.Type() != MeasurePower {
		return nil, errPowerConversion
	}

	return &powerImplementation{*m.base()}, nil
}

// StringToPower parses a power from a string such as "100 W".
// It fails if the unit isn't a power unit.
func StringToPower(s string) (Power, error) {
	// Parse and construct directly: one allocation, where ParseString
	// followed by ToPower would allocate twice.
	value, unit, err := parseValueAndUnit(s)
	if err != nil {
		return nil, err
	}

	return NewPower(value, unit)
}
