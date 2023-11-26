package units

import (
	"errors"
	"fmt"
	"github.com/asaphin/go-physics-units/time"
)

type Time interface {
	Measurement
	UnitConverter[Time]
	Arithmetics[Time]
}

type timeImplementation struct {
	baseMeasurement
}

func (t *timeImplementation) ConvertTo(unit string) (Time, error) {
	msm, err := t.convertTo(unit)
	if err != nil {
		return nil, err
	}

	return &timeImplementation{*msm}, nil
}

func (t *timeImplementation) MustConvertTo(unit string) Time {
	msm, err := t.convertTo(unit)
	if err != nil {
		panic(err)
	}

	return &timeImplementation{*msm}
}

func (t *timeImplementation) ConvertToBaseUnits() Time {
	msm := t.unsafeConvertTo(time.BaseUnit)

	return &timeImplementation{*msm}
}

func (t *timeImplementation) valueInBaseUnits() float64 {
	return t.unsafeGetValueIn(time.BaseUnit)
}

func (t *timeImplementation) valueInUnits(unit string) float64 {
	return t.unsafeGetValueIn(unit)
}

func (t *timeImplementation) Add(tm Time) Time {
	t1 := t.value
	t2 := tm.valueInUnits(t.unit)

	return newTime(t1+t2, t.unit)
}

func (t *timeImplementation) Sub(tm Time) Time {
	t1 := t.value
	t2 := tm.valueInUnits(t.unit)

	return newTime(t1-t2, t.unit)
}

func (t *timeImplementation) Mul(multiplier float64) Time {
	msm := newUnsafeBaseMeasurement(t.value*multiplier, t.unit, t.conversionRates)

	return &timeImplementation{*msm}
}

func (t *timeImplementation) Div(divisor float64) (Time, error) {
	if divisor == 0 {
		return nil, errZeroDivision
	}

	msm := newUnsafeBaseMeasurement(t.value/divisor, t.unit, t.conversionRates)

	return &timeImplementation{*msm}, nil
}

// newTime unsafe time constructor without checks for internal usage.
func newTime(value float64, unit string) Time {
	m := newUnsafeBaseMeasurement(value, unit, time.ConversionRates())

	return &timeImplementation{*m}
}

func NewTime(value float64, unit string) (Time, error) {
	rates := time.ConversionRates()

	if _, ok := rates.Has(unit + unit); !ok {
		return nil, fmt.Errorf("unknown Time unit %s", unit)
	}

	m := newUnsafeBaseMeasurement(value, unit, rates)

	return &timeImplementation{*m}, nil
}

var errTimeConversion = errors.New("not a time measure")

func ToTime(m Measurement) (Time, error) {
	b, ok := m.(*baseMeasurement)
	if !ok {
		return nil, errBaseMeasurementConversion
	}

	if m.Type() == MeasureTime {
		return &timeImplementation{*b}, nil
	}

	return nil, errTimeConversion
}
