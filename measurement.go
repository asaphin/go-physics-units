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
	unsafeConvertTo(targetUnit string) *baseMeasurement
	unsafeGetValueIn(targetUnit string) float64
	Type() MeasureType
	fmt.Stringer

	AddMeasurement(m Measurement) (Measurement, error)
	SubMeasurement(m Measurement) (Measurement, error)
	MulMeasurement(multiplier float64) Measurement
	DivMeasurement(divisor float64) (Measurement, error)
}

type baseMeasurement struct {
	value           float64
	unit            string
	conversionRates immutable.Float64Map
}

func (b *baseMeasurement) AddMeasurement(m Measurement) (Measurement, error) {
	baseType := b.Type()
	adderType := m.Type()

	if baseType != adderType {
		return nil, fmt.Errorf("incompatible types %s and %s", baseType, adderType)
	}

	adder := m.unsafeGetValueIn(b.unit)

	return newUnsafeBaseMeasurement(b.value+adder, b.unit, b.conversionRates), nil
}

func (b *baseMeasurement) SubMeasurement(m Measurement) (Measurement, error) {
	baseType := b.Type()
	adderType := m.Type()

	if baseType != adderType {
		return nil, fmt.Errorf("incompatible types %s and %s", baseType, adderType)
	}

	subtract := m.unsafeGetValueIn(b.unit)

	return newUnsafeBaseMeasurement(b.value-subtract, b.unit, b.conversionRates), nil
}

func (b *baseMeasurement) MulMeasurement(multiplier float64) Measurement {
	return &baseMeasurement{value: b.value * multiplier, unit: b.unit, conversionRates: b.conversionRates}
}

func (b *baseMeasurement) DivMeasurement(divisor float64) (Measurement, error) {
	if divisor == 0 {
		return nil, errZeroDivision
	}

	return &baseMeasurement{value: b.value / divisor, unit: b.unit, conversionRates: b.conversionRates}, nil
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

func (b *baseMeasurement) unsafeConvertTo(targetUnit string) *baseMeasurement {
	if b.unit == targetUnit {
		return newUnsafeBaseMeasurement(b.value, targetUnit, b.conversionRates)
	}

	factor := b.conversionRates.Get(b.unit + targetUnit)

	return newUnsafeBaseMeasurement(b.value*factor, targetUnit, b.conversionRates)
}

func (b *baseMeasurement) unsafeGetValueIn(targetUnit string) float64 {
	if b.unit == targetUnit {
		return b.value
	}

	factor := b.conversionRates.Get(b.unit + targetUnit)

	return b.value * factor
}

func (b *baseMeasurement) Type() MeasureType {
	return DetectMeasureType(b.unit)
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

func newUnsafeBaseMeasurement(value float64, unit string, conversionRates immutable.Float64Map) *baseMeasurement {
	return &baseMeasurement{
		value:           value,
		unit:            unit,
		conversionRates: conversionRates,
	}
}
