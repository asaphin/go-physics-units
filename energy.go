package units

import (
	"cmp"

	"github.com/asaphin/go-physics-units/energy"
	"github.com/asaphin/go-physics-units/force"
	"github.com/asaphin/go-physics-units/power"
)

// Energy is an energy measurement. Its unit symbols are defined in the energy
// package.
type Energy interface {
	Measurement
	UnitConverter[Energy]
	Arithmetics[Energy]
	Comparer[Energy]

	// DivideByTime returns the power e/t in power.BaseUnit (W).
	// It fails if t is zero.
	DivideByTime(t Time) (Power, error)
	// DivideByDistance returns the force e/d in force.BaseUnit (N).
	// It fails if d is zero.
	DivideByDistance(d Distance) (Force, error)
}

// energyImplementation is the concrete implementation of Energy.
type energyImplementation struct {
	baseMeasurement
}

// ConvertTo returns the energy expressed in unit.
func (e *energyImplementation) ConvertTo(unit string) (Energy, error) {
	msm, err := e.convertTo(unit)
	if err != nil {
		return nil, err
	}

	return &energyImplementation{msm}, nil
}

// MustConvertTo is like ConvertTo but panics if unit isn't an energy unit.
func (e *energyImplementation) MustConvertTo(unit string) Energy {
	msm, err := e.convertTo(unit)
	if err != nil {
		panic(err)
	}

	return &energyImplementation{msm}
}

// ConvertToBaseUnits returns the energy expressed in joules.
func (e *energyImplementation) ConvertToBaseUnits() Energy {
	return &energyImplementation{e.unsafeConvertTo(energy.BaseUnit)}
}

func (e *energyImplementation) valueInBaseUnits() float64 {
	return e.unsafeGetValueIn(energy.BaseUnit)
}

func (e *energyImplementation) valueInUnits(unit string) float64 {
	return e.unsafeGetValueIn(unit)
}

// Add returns the sum of the energies, expressed in the receiver's unit.
func (e *energyImplementation) Add(other Energy) Energy {
	return &energyImplementation{e.withValue(e.value + other.valueInUnits(e.unit))}
}

// Sub returns the difference of the energies, expressed in the receiver's unit.
func (e *energyImplementation) Sub(other Energy) Energy {
	return &energyImplementation{e.withValue(e.value - other.valueInUnits(e.unit))}
}

// Mul returns the energy multiplied by multiplier.
func (e *energyImplementation) Mul(multiplier float64) Energy {
	return &energyImplementation{e.withValue(e.value * multiplier)}
}

// Div returns the energy divided by divisor. It fails if divisor is zero.
func (e *energyImplementation) Div(divisor float64) (Energy, error) {
	if divisor == 0 {
		return nil, ErrZeroDivision
	}

	return &energyImplementation{e.withValue(e.value / divisor)}, nil
}

// DivideByTime returns the power e/t in W, or an error if t is zero.
func (e *energyImplementation) DivideByTime(t Time) (Power, error) {
	if t.Value() == 0 {
		return nil, ErrZeroDivision
	}

	return newPower(e.valueInBaseUnits()/t.valueInBaseUnits(), power.BaseUnit), nil
}

// DivideByDistance returns the force e/d in N, or an error if d is zero.
func (e *energyImplementation) DivideByDistance(d Distance) (Force, error) {
	if d.Value() == 0 {
		return nil, ErrZeroDivision
	}

	return newForce(e.valueInBaseUnits()/d.valueInBaseUnits(), force.BaseUnit), nil
}

// Compare compares the energies; see Comparer.
func (e *energyImplementation) Compare(other Energy) int {
	return cmp.Compare(e.commonValues(other.base()))
}

// Equal reports whether the energies are equal; see Comparer.
func (e *energyImplementation) Equal(other Energy) bool {
	return e.Compare(other) == 0
}

// EqualWithinRel reports whether the energies are equal within the relative
// tolerance tol; see Comparer.
func (e *energyImplementation) EqualWithinRel(other Energy, tol float64) bool {
	x, y := e.commonValues(other.base())

	return equalWithinRel(x, y, tol)
}

// NewEnergy creates an energy of value in unit, one of the symbols defined in
// the energy package.
func NewEnergy(value float64, unit string) (Energy, error) {
	msm, err := newBaseMeasurement(value, unit, energyRates)
	if err != nil {
		return nil, err
	}

	return &energyImplementation{msm}, nil
}

// newEnergy is an unchecked constructor for internal use with valid, canonical
// units.
func newEnergy(value float64, unit string) Energy {
	return &energyImplementation{baseMeasurement{value: value, unit: unit, conversionRates: energyRates}}
}

const errEnergyConversion incompatibleError = "not an energy measure"

// ToEnergy converts m, for example a result of ParseString, into an Energy.
// It fails if m isn't an energy.
func ToEnergy(m Measurement) (Energy, error) {
	if m == nil || m.Type() != MeasureEnergy {
		return nil, errEnergyConversion
	}

	return &energyImplementation{*m.base()}, nil
}

// StringToEnergy parses an energy from a string such as "10 kWh".
// It fails if the unit isn't an energy unit.
func StringToEnergy(s string) (Energy, error) {
	// Parse and construct directly: one allocation, where ParseString
	// followed by ToEnergy would allocate twice.
	value, unit, err := parseValueAndUnit(s)
	if err != nil {
		return nil, err
	}

	return NewEnergy(value, unit)
}
