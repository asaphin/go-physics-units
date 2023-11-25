package units

import (
	"errors"
	"fmt"
	"github.com/asaphin/go-physics-units/conversion-factors"
	"strconv"
	"strings"
)

type Measurement interface {
	Value() float64
	Unit() string
	convertTo(targetUnit string) (*baseMeasurement, error)
	Type() MeasureType
	Mul(multiplier float64) *baseMeasurement
	Div(divisor float64) *baseMeasurement
	fmt.Stringer
}

type baseMeasurement struct {
	value             float64
	unit              string
	conversionFactors *immutableConversionFactors
}

var errBaseMeasurementConversion = errors.New("not a *baseMeasurement type")

func (b *baseMeasurement) Value() float64 {
	return b.value
}

func (b *baseMeasurement) Unit() string {
	return b.unit
}

func (b *baseMeasurement) convertTo(targetUnit string) (*baseMeasurement, error) {
	if b.unit == targetUnit {
		return newBaseMeasurement(b.value, targetUnit, b.conversionFactors)
	}

	initialFactor := b.conversionFactors.GetFactor(b.unit)

	if targetFactor, ok := b.conversionFactors.HasFactor(targetUnit); ok {
		factor := initialFactor / targetFactor

		return newBaseMeasurement(b.value*factor, targetUnit, b.conversionFactors)
	}

	return nil, fmt.Errorf("unit %s unspecified in conversionFactors", targetUnit)
}

func (b *baseMeasurement) Type() MeasureType {
	return DetectMeasureType(b.unit)
}

func (b *baseMeasurement) Mul(multiplier float64) *baseMeasurement {
	return &baseMeasurement{value: b.value * multiplier, unit: b.unit, conversionFactors: b.conversionFactors}
}

func (b *baseMeasurement) Div(divisor float64) *baseMeasurement {
	return &baseMeasurement{value: b.value / divisor, unit: b.unit, conversionFactors: b.conversionFactors}
}

func (b *baseMeasurement) String() string {
	return fmt.Sprintf("%v %s", b.value, b.unit)
}

func NewMeasurement(value float64, unit string) (Measurement, error) {
	mt := DetectMeasureType(unit)
	if mt == MeasureCustom {
		return nil, fmt.Errorf("unknown unit %s", unit)
	}

	return NewBaseMeasurement(value, unit, measureToConversionFactorsMapping[mt])
}

func ParseString(s string) (Measurement, error) {
	parts := strings.Fields(s)
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid string format: %s", s)
	}

	value, err := strconv.ParseFloat(parts[0], 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse value %v: %w", parts[0], err)
	}

	return NewMeasurement(value, parts[1])
}

// NewBaseMeasurement is intended to create measurements based on your custom conversion factors.
// To create common measurement use NewMeasurement.
func NewBaseMeasurement(value float64, unit string, conversionFactors conversion.Factors) (Measurement, error) {
	factor, ok := conversionFactors[unit]

	if !ok {
		return nil, fmt.Errorf("unit %s unspecified in conversionFactors", unit)
	}

	if factor == 0 {
		return nil, fmt.Errorf("unit conversion factor can not be zero (%s)", unit)
	}

	return &baseMeasurement{
		value:             value,
		unit:              unit,
		conversionFactors: newImmutableConversionFactors(conversionFactors),
	}, nil
}

func newBaseMeasurement(value float64, unit string,
	conversionFactors *immutableConversionFactors) (*baseMeasurement, error) {
	factor, ok := conversionFactors.HasFactor(unit)

	if !ok {
		return nil, fmt.Errorf("unit %s unspecified in conversionFactors", unit)
	}

	if factor == 0 {
		return nil, fmt.Errorf("unit conversion factor can not be zero (%s)", unit)
	}

	return &baseMeasurement{
		value:             value,
		unit:              unit,
		conversionFactors: conversionFactors,
	}, nil
}
