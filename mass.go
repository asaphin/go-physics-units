package units

import (
	"errors"

	"github.com/asaphin/go-physics-units/mass"
)

// Mass is a mass measurement. Its unit symbols are defined in the mass package.
type Mass interface {
	Measurement
	UnitConverter[Mass]
	Arithmetics[Mass]
}

// massImplementation is the concrete implementation of Mass.
type massImplementation struct {
	baseMeasurement
}

// ConvertTo returns the mass expressed in unit.
func (m *massImplementation) ConvertTo(unit string) (Mass, error) {
	msm, err := m.convertTo(unit)
	if err != nil {
		return nil, err
	}

	return &massImplementation{msm}, nil
}

// MustConvertTo is like ConvertTo but panics if unit isn't a mass unit.
func (m *massImplementation) MustConvertTo(unit string) Mass {
	msm, err := m.convertTo(unit)
	if err != nil {
		panic(err)
	}

	return &massImplementation{msm}
}

// ConvertToBaseUnits returns the mass expressed in kilograms.
func (m *massImplementation) ConvertToBaseUnits() Mass {
	return &massImplementation{m.unsafeConvertTo(mass.BaseUnit)}
}

func (m *massImplementation) valueInBaseUnits() float64 {
	return m.unsafeGetValueIn(mass.BaseUnit)
}

func (m *massImplementation) valueInUnits(unit string) float64 {
	return m.unsafeGetValueIn(unit)
}

// Add returns the sum of the masses, expressed in the receiver's unit.
func (m *massImplementation) Add(other Mass) Mass {
	return &massImplementation{m.withValue(m.value + other.valueInUnits(m.unit))}
}

// Sub returns the difference of the masses, expressed in the receiver's unit.
func (m *massImplementation) Sub(other Mass) Mass {
	return &massImplementation{m.withValue(m.value - other.valueInUnits(m.unit))}
}

// Mul returns the mass multiplied by multiplier.
func (m *massImplementation) Mul(multiplier float64) Mass {
	return &massImplementation{m.withValue(m.value * multiplier)}
}

// Div returns the mass divided by divisor. It fails if divisor is zero.
func (m *massImplementation) Div(divisor float64) (Mass, error) {
	if divisor == 0 {
		return nil, errZeroDivision
	}

	return &massImplementation{m.withValue(m.value / divisor)}, nil
}

// NewMass creates a mass of value in unit, one of the symbols defined in the
// mass package.
func NewMass(value float64, unit string) (Mass, error) {
	msm, err := newBaseMeasurement(value, unit, mass.ConversionRates())
	if err != nil {
		return nil, err
	}

	return &massImplementation{msm}, nil
}

var errMassConversion = errors.New("not a mass measure")

// ToMass converts m, for example a result of ParseString, into a Mass.
// It fails if m isn't a mass.
func ToMass(m Measurement) (Mass, error) {
	if m == nil || m.Type() != MeasureMass {
		return nil, errMassConversion
	}

	return &massImplementation{*m.base()}, nil
}

// StringToMass parses a mass from a string such as "10 kg".
func StringToMass(s string) (Mass, error) {
	m, err := ParseString(s)
	if err != nil {
		return nil, err
	}

	return ToMass(m)
}
