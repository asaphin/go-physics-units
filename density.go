package units

import (
	"cmp"

	"github.com/asaphin/go-physics-units/density"
	"github.com/asaphin/go-physics-units/mass"
)

// Density is a density measurement. Its unit symbols are defined in the density
// package.
type Density interface {
	Measurement
	UnitConverter[Density]
	Arithmetics[Density]
	Comparer[Density]

	// MultiplyByVolume returns the mass d·v in mass.BaseUnit (kg).
	MultiplyByVolume(v Volume) Mass
}

// densityImplementation is the concrete implementation of Density.
type densityImplementation struct {
	baseMeasurement
}

// ConvertTo returns the density expressed in unit.
func (d *densityImplementation) ConvertTo(unit string) (Density, error) {
	msm, err := d.convertTo(unit)
	if err != nil {
		return nil, err
	}

	return &densityImplementation{msm}, nil
}

// MustConvertTo is like ConvertTo but panics if unit isn't a density unit.
func (d *densityImplementation) MustConvertTo(unit string) Density {
	msm, err := d.convertTo(unit)
	if err != nil {
		panic(err)
	}

	return &densityImplementation{msm}
}

// ConvertToBaseUnits returns the density expressed in kilograms per cubic metre.
func (d *densityImplementation) ConvertToBaseUnits() Density {
	return &densityImplementation{d.unsafeConvertTo(density.BaseUnit)}
}

func (d *densityImplementation) valueInBaseUnits() float64 {
	return d.unsafeGetValueIn(density.BaseUnit)
}

func (d *densityImplementation) valueInUnits(unit string) float64 {
	return d.unsafeGetValueIn(unit)
}

// Add returns the sum of the densities, expressed in the receiver's unit.
func (d *densityImplementation) Add(other Density) Density {
	return &densityImplementation{d.withValue(d.value + other.valueInUnits(d.unit))}
}

// Sub returns the difference of the densities, expressed in the receiver's unit.
func (d *densityImplementation) Sub(other Density) Density {
	return &densityImplementation{d.withValue(d.value - other.valueInUnits(d.unit))}
}

// Mul returns the density multiplied by multiplier.
func (d *densityImplementation) Mul(multiplier float64) Density {
	return &densityImplementation{d.withValue(d.value * multiplier)}
}

// Div returns the density divided by divisor. It fails if divisor is zero.
func (d *densityImplementation) Div(divisor float64) (Density, error) {
	if divisor == 0 {
		return nil, ErrZeroDivision
	}

	return &densityImplementation{d.withValue(d.value / divisor)}, nil
}

// MultiplyByVolume returns the mass d·v in kg.
func (d *densityImplementation) MultiplyByVolume(v Volume) Mass {
	return newMass(d.valueInBaseUnits()*v.valueInBaseUnits(), mass.BaseUnit)
}

// Compare compares the densities; see Comparer.
func (d *densityImplementation) Compare(other Density) int {
	return cmp.Compare(d.commonValues(other.base()))
}

// Equal reports whether the densities are equal; see Comparer.
func (d *densityImplementation) Equal(other Density) bool {
	return d.Compare(other) == 0
}

// EqualWithinRel reports whether the densities are equal within the relative
// tolerance tol; see Comparer.
func (d *densityImplementation) EqualWithinRel(other Density, tol float64) bool {
	x, y := d.commonValues(other.base())

	return equalWithinRel(x, y, tol)
}

// NewDensity creates a density of value in unit, one of the symbols defined in
// the density package.
func NewDensity(value float64, unit string) (Density, error) {
	msm, err := newBaseMeasurement(value, unit, densityRates)
	if err != nil {
		return nil, err
	}

	return &densityImplementation{msm}, nil
}

// newDensity is an unchecked constructor for internal use with valid, canonical
// units.
func newDensity(value float64, unit string) Density {
	return &densityImplementation{baseMeasurement{value: value, unit: unit, conversionRates: densityRates}}
}

const errDensityConversion incompatibleError = "not a density measure"

// ToDensity converts m, for example a result of ParseString, into a Density.
// It fails if m isn't a density.
func ToDensity(m Measurement) (Density, error) {
	if m == nil || m.Type() != MeasureDensity {
		return nil, errDensityConversion
	}

	return &densityImplementation{*m.base()}, nil
}

// StringToDensity parses a density from a string such as "1.2 kg/m³".
// It fails if the unit isn't a density unit.
func StringToDensity(s string) (Density, error) {
	// Parse and construct directly: one allocation, where ParseString
	// followed by ToDensity would allocate twice.
	value, unit, err := parseValueAndUnit(s)
	if err != nil {
		return nil, err
	}

	return NewDensity(value, unit)
}
