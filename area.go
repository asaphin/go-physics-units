package units

import (
	"cmp"

	"github.com/asaphin/go-physics-units/area"
	"github.com/asaphin/go-physics-units/distance"
)

// Area is an area measurement. Its unit symbols are defined in the area package.
type Area interface {
	Measurement
	UnitConverter[Area]
	Arithmetics[Area]
	Comparer[Area]

	// DivideByDistance returns the distance a/d in distance.BaseUnit (m).
	// It fails if d is zero.
	DivideByDistance(d Distance) (Distance, error)
}

// areaImplementation is the concrete implementation of Area.
type areaImplementation struct {
	baseMeasurement
}

// ConvertTo returns the area expressed in unit.
func (a *areaImplementation) ConvertTo(unit string) (Area, error) {
	msm, err := a.convertTo(unit)
	if err != nil {
		return nil, err
	}

	return &areaImplementation{msm}, nil
}

// MustConvertTo is like ConvertTo but panics if unit isn't an area unit.
func (a *areaImplementation) MustConvertTo(unit string) Area {
	msm, err := a.convertTo(unit)
	if err != nil {
		panic(err)
	}

	return &areaImplementation{msm}
}

// ConvertToBaseUnits returns the area expressed in square metres.
func (a *areaImplementation) ConvertToBaseUnits() Area {
	return &areaImplementation{a.unsafeConvertTo(area.BaseUnit)}
}

func (a *areaImplementation) valueInBaseUnits() float64 {
	return a.unsafeGetValueIn(area.BaseUnit)
}

func (a *areaImplementation) valueInUnits(unit string) float64 {
	return a.unsafeGetValueIn(unit)
}

// Add returns the sum of the areas, expressed in the receiver's unit.
func (a *areaImplementation) Add(other Area) Area {
	return &areaImplementation{a.withValue(a.value + other.valueInUnits(a.unit))}
}

// Sub returns the difference of the areas, expressed in the receiver's unit.
func (a *areaImplementation) Sub(other Area) Area {
	return &areaImplementation{a.withValue(a.value - other.valueInUnits(a.unit))}
}

// Mul returns the area multiplied by multiplier.
func (a *areaImplementation) Mul(multiplier float64) Area {
	return &areaImplementation{a.withValue(a.value * multiplier)}
}

// Div returns the area divided by divisor. It fails if divisor is zero.
func (a *areaImplementation) Div(divisor float64) (Area, error) {
	if divisor == 0 {
		return nil, ErrZeroDivision
	}

	return &areaImplementation{a.withValue(a.value / divisor)}, nil
}

// DivideByDistance returns the distance a/d in m, or an error if d is zero.
func (a *areaImplementation) DivideByDistance(d Distance) (Distance, error) {
	if d.Value() == 0 {
		return nil, ErrZeroDivision
	}

	return newDistance(a.valueInBaseUnits()/d.valueInBaseUnits(), distance.BaseUnit), nil
}

// Compare compares the areas; see Comparer.
func (a *areaImplementation) Compare(other Area) int {
	return cmp.Compare(a.commonValues(other.base()))
}

// Equal reports whether the areas are equal; see Comparer.
func (a *areaImplementation) Equal(other Area) bool {
	return a.Compare(other) == 0
}

// EqualWithinRel reports whether the areas are equal within the relative
// tolerance tol; see Comparer.
func (a *areaImplementation) EqualWithinRel(other Area, tol float64) bool {
	x, y := a.commonValues(other.base())

	return equalWithinRel(x, y, tol)
}

// NewArea creates an area of value in unit, one of the symbols defined in the
// area package.
func NewArea(value float64, unit string) (Area, error) {
	msm, err := newBaseMeasurement(value, unit, areaRates)
	if err != nil {
		return nil, err
	}

	return &areaImplementation{msm}, nil
}

// newArea is an unchecked constructor for internal use with valid, canonical
// units.
func newArea(value float64, unit string) Area {
	return &areaImplementation{baseMeasurement{value: value, unit: unit, conversionRates: areaRates}}
}

const errAreaConversion incompatibleError = "not an area measure"

// ToArea converts m, for example a result of ParseString, into an Area.
// It fails if m isn't an area.
func ToArea(m Measurement) (Area, error) {
	if m == nil || m.Type() != MeasureArea {
		return nil, errAreaConversion
	}

	return &areaImplementation{*m.base()}, nil
}

// StringToArea parses an area from a string such as "10 m²".
// It fails if the unit isn't a area unit.
func StringToArea(s string) (Area, error) {
	// Parse and construct directly: one allocation, where ParseString
	// followed by ToArea would allocate twice.
	value, unit, err := parseValueAndUnit(s)
	if err != nil {
		return nil, err
	}

	return NewArea(value, unit)
}
