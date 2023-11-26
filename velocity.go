package units

import (
	"errors"
	"fmt"
	"github.com/asaphin/go-physics-units/velocity"
)

type Velocity interface {
	Measurement
	UnitConverter[Velocity]
	Arithmetics[Velocity]
}

type velocityImplenentation struct {
	baseMeasurement
}

func (v *velocityImplenentation) ConvertTo(unit string) (Velocity, error) {
	msm, err := v.convertTo(unit)
	if err != nil {
		return nil, err
	}

	return &velocityImplenentation{*msm}, nil
}

func (v *velocityImplenentation) MustConvertTo(unit string) Velocity {
	msm, err := v.convertTo(unit)
	if err != nil {
		panic(err)
	}

	return &velocityImplenentation{*msm}
}

func (v *velocityImplenentation) ConvertToBaseUnits() Velocity {
	msm := v.unsafeConvertTo(velocity.BaseUnit)

	return &velocityImplenentation{*msm}
}

func (v *velocityImplenentation) valueInBaseUnits() float64 {
	return v.unsafeGetValueIn(velocity.BaseUnit)
}

func (v *velocityImplenentation) valueInUnits(unit string) float64 {
	return v.unsafeGetValueIn(unit)
}

func (v *velocityImplenentation) Add(vl Velocity) Velocity {
	v1 := v.value
	v2 := vl.valueInUnits(v.unit)

	return newVelocity(v1+v2, v.unit)
}

func (v *velocityImplenentation) Sub(vl Velocity) Velocity {
	v1 := v.value
	v2 := vl.valueInUnits(v.unit)

	return newVelocity(v1-v2, v.unit)
}

func (v *velocityImplenentation) Mul(multiplier float64) Velocity {
	msm := newUnsafeBaseMeasurement(v.value*multiplier, v.unit, v.conversionRates)

	return &velocityImplenentation{*msm}
}

func (v *velocityImplenentation) Div(divisor float64) (Velocity, error) {
	if divisor == 0 {
		return nil, errZeroDivision
	}

	msm := newUnsafeBaseMeasurement(v.value/divisor, v.unit, v.conversionRates)

	return &velocityImplenentation{*msm}, nil
}

func NewVelocity(value float64, unit string) (Velocity, error) {
	rates := velocity.ConversionRates()

	if _, ok := rates.Has(unit + unit); !ok {
		return nil, fmt.Errorf("unknown Velocity unit %s", unit)
	}

	m := newUnsafeBaseMeasurement(value, unit, rates)

	return &velocityImplenentation{*m}, nil
}

// newVelocity unsafe velocity constructor without checks for internal usage.
func newVelocity(value float64, unit string) Velocity {
	m := newUnsafeBaseMeasurement(value, unit, velocity.ConversionRates())

	return &velocityImplenentation{*m}
}

var errVelocityConversion = errors.New("not a velocity measure")

func ToVelocity(m Measurement) (Velocity, error) {
	b, ok := m.(*baseMeasurement)
	if !ok {
		return nil, errBaseMeasurementConversion
	}

	if m.Type() == MeasureVelocity {
		return &velocityImplenentation{*b}, nil
	}

	return nil, errVelocityConversion
}
