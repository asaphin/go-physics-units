package units

import (
	"errors"
	"fmt"
	"github.com/asaphin/go-physics-units/conversion"
	"github.com/asaphin/go-physics-units/internal/immutable"
	"github.com/asaphin/go-physics-units/internal/rates"
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
	value           float64
	unit            string
	conversionRates immutable.Float64Map
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
		return newBaseMeasurement(b.value, targetUnit, b.conversionRates)
	}

	factor, ok := b.conversionRates.Has(b.unit + targetUnit)
	if !ok {
		measureFactors := getMeasureFactors(b.Type())

		if _, ok = measureFactors[b.unit]; !ok {
			return nil, fmt.Errorf("initial unit %s unspecified in conversionRates", b.unit)
		}

		if _, ok = measureFactors[targetUnit]; !ok {
			return nil, fmt.Errorf("target unit %s unspecified in conversionRates", targetUnit)
		}

		return nil, errors.New("global conversion units error")
	}

	return newBaseMeasurement(b.value*factor, targetUnit, b.conversionRates)
}

func (b *baseMeasurement) Type() MeasureType {
	return DetectMeasureType(b.unit)
}

func (b *baseMeasurement) Mul(multiplier float64) *baseMeasurement {
	return &baseMeasurement{value: b.value * multiplier, unit: b.unit, conversionRates: b.conversionRates}
}

func (b *baseMeasurement) Div(divisor float64) *baseMeasurement {
	return &baseMeasurement{value: b.value / divisor, unit: b.unit, conversionRates: b.conversionRates}
}

func (b *baseMeasurement) String() string {
	return fmt.Sprintf("%v %s", b.value, b.unit)
}

func NewMeasurement(value float64, unit string) (Measurement, error) {
	mt := DetectMeasureType(unit)
	if mt == MeasureCustom {
		return nil, fmt.Errorf("unknown unit %s", unit)
	}

	return newBaseMeasurement(value, unit, measureToConversionRatesMapping[mt])
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
		return nil, fmt.Errorf("unit %s unspecified in conversionRates", unit)
	}

	if factor == 0 {
		return nil, fmt.Errorf("unit conversion factor can not be zero (%s)", unit)
	}

	return &baseMeasurement{
		value:           value,
		unit:            unit,
		conversionRates: immutable.MakeImmutable(rates.FactorsToRates(conversionFactors)),
	}, nil
}

func newBaseMeasurement(value float64, unit string, conversionRates immutable.Float64Map) (*baseMeasurement, error) {
	factor, ok := conversionRates.Has(unit + unit)
	if !ok {
		return nil, fmt.Errorf("unit %s unspecified in conversionRates", unit)
	}

	if factor == 0 {
		return nil, fmt.Errorf("unit conversion factor can not be zero (%s)", unit)
	}

	return &baseMeasurement{
		value:           value,
		unit:            unit,
		conversionRates: conversionRates,
	}, nil
}
