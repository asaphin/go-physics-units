package units

import (
	"cmp"

	"github.com/asaphin/go-physics-units/density"
	"github.com/asaphin/go-physics-units/force"
	"github.com/asaphin/go-physics-units/mass"
	"github.com/asaphin/go-physics-units/volume"
)

// Mass is a mass measurement. Its unit symbols are defined in the mass package.
type Mass interface {
	Measurement
	UnitConverter[Mass]
	Arithmetics[Mass]
	Comparer[Mass]

	// MultiplyByAcceleration returns the force m·a in force.BaseUnit (N).
	MultiplyByAcceleration(a Acceleration) Force
	// DivideByVolume returns the density m/v in density.BaseUnit (kg/m³).
	// It fails if v is zero.
	DivideByVolume(v Volume) (Density, error)
	// DivideByDensity returns the volume m/d in volume.BaseUnit (m³).
	// It fails if d is zero.
	DivideByDensity(d Density) (Volume, error)
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
		return nil, ErrZeroDivision
	}

	return &massImplementation{m.withValue(m.value / divisor)}, nil
}

// MultiplyByAcceleration returns the force m·a in N.
func (m *massImplementation) MultiplyByAcceleration(a Acceleration) Force {
	return newForce(m.valueInBaseUnits()*a.valueInBaseUnits(), force.BaseUnit)
}

// DivideByVolume returns the density m/v in kg/m³, or an error if v is zero.
func (m *massImplementation) DivideByVolume(v Volume) (Density, error) {
	if v.Value() == 0 {
		return nil, ErrZeroDivision
	}

	return newDensity(m.valueInBaseUnits()/v.valueInBaseUnits(), density.BaseUnit), nil
}

// DivideByDensity returns the volume m/d in m³, or an error if d is zero.
func (m *massImplementation) DivideByDensity(d Density) (Volume, error) {
	if d.Value() == 0 {
		return nil, ErrZeroDivision
	}

	return newVolume(m.valueInBaseUnits()/d.valueInBaseUnits(), volume.BaseUnit), nil
}

// Compare compares the masses; see Comparer.
func (m *massImplementation) Compare(other Mass) int {
	return cmp.Compare(m.commonValues(other.base()))
}

// Equal reports whether the masses are equal; see Comparer.
func (m *massImplementation) Equal(other Mass) bool {
	return m.Compare(other) == 0
}

// EqualWithinRel reports whether the masses are equal within the relative
// tolerance tol; see Comparer.
func (m *massImplementation) EqualWithinRel(other Mass, tol float64) bool {
	x, y := m.commonValues(other.base())

	return equalWithinRel(x, y, tol)
}

// NewMass creates a mass of value in unit, one of the symbols defined in the
// mass package.
func NewMass(value float64, unit string) (Mass, error) {
	msm, err := newBaseMeasurement(value, unit, massRates)
	if err != nil {
		return nil, err
	}

	return &massImplementation{msm}, nil
}

// newMass is an unchecked constructor for internal use with valid, canonical
// units.
func newMass(value float64, unit string) Mass {
	return &massImplementation{baseMeasurement{value: value, unit: unit, conversionRates: massRates}}
}

const errMassConversion incompatibleError = "not a mass measure"

// ToMass converts m, for example a result of ParseString, into a Mass.
// It fails if m isn't a mass.
func ToMass(m Measurement) (Mass, error) {
	if m == nil || m.Type() != MeasureMass {
		return nil, errMassConversion
	}

	return &massImplementation{*m.base()}, nil
}

// StringToMass parses a mass from a string such as "10 kg".
// It fails if the unit isn't a mass unit.
func StringToMass(s string) (Mass, error) {
	// Parse and construct directly: one allocation, where ParseString
	// followed by ToMass would allocate twice.
	value, unit, err := parseValueAndUnit(s)
	if err != nil {
		return nil, err
	}

	return NewMass(value, unit)
}
