package units

import (
	"errors"
	"fmt"
	"github.com/asaphin/go-physics-units/area"
	"github.com/asaphin/go-physics-units/distance"
)

// Area interface represents an area measurement.
type Area interface {
	Measurement
	UnitConverter[Area]
	Arithmetics[Area]

	DivideByDistance(d Distance) Distance
}

// areaImplementation is a concrete implementation of the Area interface.
type areaImplementation struct {
	baseMeasurement
}

// ConvertTo implements the ConvertTo method of the Area interface.
func (a *areaImplementation) ConvertTo(unit string) (Area, error) {
	msm, err := a.convertTo(unit)
	if err != nil {
		return nil, err
	}

	return &areaImplementation{*msm}, nil
}

// MustConvertTo implements the MustConvertTo method of the Area interface.
func (a *areaImplementation) MustConvertTo(unit string) Area {
	msm, err := a.convertTo(unit)
	if err != nil {
		panic(err)
	}

	return &areaImplementation{*msm}
}

// ConvertToBaseUnits implements the MustConvertToBaseUnits method of the Area interface.
func (a *areaImplementation) ConvertToBaseUnits() Area {
	msm := a.unsafeConvertTo(area.BaseUnit)

	return &areaImplementation{*msm}
}

func (a *areaImplementation) valueInBaseUnits() float64 {
	return a.unsafeGetValueIn(area.BaseUnit)
}

func (a *areaImplementation) valueInUnits(unit string) float64 {
	return a.unsafeGetValueIn(unit)
}

func (a *areaImplementation) Add(area Area) Area {
	a1 := a.value
	a2 := area.valueInUnits(a.unit)

	return newArea(a1+a2, a.unit)
}

func (a *areaImplementation) Sub(area Area) Area {
	a1 := a.value
	a2 := area.valueInUnits(a.unit)

	return newArea(a1-a2, a.unit)
}

func (a *areaImplementation) Mul(multiplier float64) Area {
	msm := newUnsafeBaseMeasurement(a.value*multiplier, a.unit, a.conversionRates)

	return &areaImplementation{*msm}
}

func (a *areaImplementation) Div(divisor float64) (Area, error) {
	if divisor == 0 {
		return nil, errZeroDivision
	}

	msm := newUnsafeBaseMeasurement(a.value/divisor, a.unit, a.conversionRates)

	return &areaImplementation{*msm}, nil
}

// DivideByDistance implements the DivideByDistance method of the Area interface.
func (a *areaImplementation) DivideByDistance(d Distance) Distance {
	if d.Value() == 0 {
		panic("division by zero")
	}

	baseA := a.valueInBaseUnits()
	baseD := d.valueInBaseUnits()

	return newDistance(baseA/baseD, distance.BaseUnit)
}

// NewArea creates a new Area instance.
func NewArea(value float64, unit string) (Area, error) {
	rates := area.ConversionRates()

	if _, ok := rates.Has(unit + unit); !ok {
		return nil, fmt.Errorf("unknown Area unit %s", unit)
	}

	msm := newUnsafeBaseMeasurement(value, unit, rates)

	return &areaImplementation{*msm}, nil
}

// newArea unsafe area constructor without checks for internal usage.
func newArea(value float64, unit string) Area {
	msm := newUnsafeBaseMeasurement(value, unit, area.ConversionRates())

	return &areaImplementation{*msm}
}

var errAreaConversion = errors.New("not an area measure")

func ToArea(m Measurement) (Area, error) {
	if m.Type() == MeasureArea {
		b, ok := m.(*baseMeasurement)
		if !ok {
			return nil, errBaseMeasurementConversion
		}

		return &areaImplementation{*b}, nil
	}

	return nil, errAreaConversion
}

func StringToArea(s string) (Area, error) {
	m, err := ParseString(s)
	if err != nil {
		return nil, err
	}

	return ToArea(m)
}
