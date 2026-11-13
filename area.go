package units

import (
	"errors"

	"github.com/asaphin/go-physics-units/area"
	"github.com/asaphin/go-physics-units/distance"
)

// Area is an area measurement. Its unit symbols are defined in the area package.
type Area interface {
	Measurement
	UnitConverter[Area]
	Arithmetics[Area]

	// DivideByDistance returns the distance a/d in distance.BaseUnit (m).
	// It panics if d is zero.
	DivideByDistance(d Distance) Distance
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
		return nil, errZeroDivision
	}

	return &areaImplementation{a.withValue(a.value / divisor)}, nil
}

// DivideByDistance returns the distance a/d in m. It panics if d is zero.
func (a *areaImplementation) DivideByDistance(d Distance) Distance {
	if d.Value() == 0 {
		panic(errZeroDivision)
	}

	return newDistance(a.valueInBaseUnits()/d.valueInBaseUnits(), distance.BaseUnit)
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

var errAreaConversion = errors.New("not an area measure")

// ToArea converts m, for example a result of ParseString, into an Area.
// It fails if m isn't an area.
func ToArea(m Measurement) (Area, error) {
	if m == nil || m.Type() != MeasureArea {
		return nil, errAreaConversion
	}

	return &areaImplementation{*m.base()}, nil
}

// StringToArea parses an area from a string such as "10 m²".
func StringToArea(s string) (Area, error) {
	m, err := ParseString(s)
	if err != nil {
		return nil, err
	}

	return ToArea(m)
}
