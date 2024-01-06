package units

import (
	"fmt"
	"github.com/asaphin/go-physics-units/conversion"
	"github.com/asaphin/go-physics-units/internal/immutable"
	"github.com/asaphin/go-physics-units/internal/rates"
	"strconv"
	"strings"
)

type MeasurementCreator[T Measurement] interface {
	New(value float64, unit string) (T, error)
	ParseString(s string) (T, error)
}

type measurementCreator[T Measurement] struct {
	baseUnit        string
	factors         conversion.Factors
	conversionRates immutable.Float64Map
	measureType     MeasureType
}

func NewMeasurementCreator[T Measurement](
	measureType, baseUnit string, factors conversion.Factors) (MeasurementCreator[T], error) {
	if _, ok := factors[baseUnit]; !ok {
		return nil, fmt.Errorf("provided factors doesn't contain unit %s", baseUnit)
	}

	f := conversion.CopyConversionFactors(factors)

	return &measurementCreator[T]{
		baseUnit:        baseUnit,
		factors:         f,
		conversionRates: immutable.MakeImmutable(rates.FactorsToRates(f)),
		measureType:     MeasureType(measureType),
	}, nil
}

func (c *measurementCreator[T]) New(value float64, unit string) (T, error) {
	bm, err := newBaseMeasurement(value, unit, c.conversionRates)

	return Measurement(bm).(T), err
}

func (c *measurementCreator[T]) ParseString(s string) (T, error) {
	var bm *baseMeasurement

	parts := strings.Fields(s)
	if len(parts) != 2 {
		return Measurement(bm).(T), fmt.Errorf("invalid string format: %s", s)
	}

	if _, ok := c.factors[parts[1]]; !ok {
		return Measurement(bm).(T), fmt.Errorf("unit %s isn't specified for measure %s", parts[1], c.measureType)
	}

	value, err := strconv.ParseFloat(parts[0], 64)
	if err != nil {
		return Measurement(bm).(T), fmt.Errorf("failed to parse value %v: %w", parts[0], err)
	}

	bm, err = newBaseMeasurement(value, parts[1], c.conversionRates)

	return Measurement(bm).(T), err
}
