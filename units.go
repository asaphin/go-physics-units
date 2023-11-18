package units

import (
	"fmt"
)

type Measurement interface {
	Value() float64
	Unit() string
	ConvertTo(targetUnit string) (Measurement, error)
	fmt.Stringer
}

type BaseMeasurement struct {
	value             float64
	unit              string
	conversionFactors map[string]float64
}

func (b *BaseMeasurement) Value() float64 {
	return b.value
}

func (b *BaseMeasurement) Unit() string {
	return b.unit
}

func (b *BaseMeasurement) ConvertTo(targetUnit string) (Measurement, error) {
	initialFactor := b.conversionFactors[b.unit]

	if targetFactor, ok := b.conversionFactors[targetUnit]; ok {
		factor := targetFactor / initialFactor

		return NewBaseMeasurement(b.value*factor, targetUnit, b.conversionFactors)
	}

	return nil, fmt.Errorf("unit %s unspecified in conversionFactors", targetUnit)
}

func (b *BaseMeasurement) String() string {
	return fmt.Sprintf("%v %s", b.value, b.unit)
}

func NewBaseMeasurement(value float64, unit string, conversionFactors map[string]float64) (Measurement, error) {
	factor, ok := conversionFactors[unit]

	if !ok {
		return nil, fmt.Errorf("unit %s unspecified in conversionFactors", unit)
	}

	if factor == 0 {
		return nil, fmt.Errorf("unit conversion factor can not be zero (%s)", unit)
	}

	return &BaseMeasurement{value: value, unit: unit, conversionFactors: conversionFactors}, nil
}
