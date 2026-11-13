package units

import (
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/asaphin/go-physics-units/conversion"
	"github.com/asaphin/go-physics-units/internal/rates"
)

// Measurement is a value together with its unit. It is implemented by every
// measure of this package (Distance, Time, …) and by the measurements produced
// by a MeasurementCreator.
//
// Measurements are immutable: every operation returns a new value.
// Measurement has unexported methods, so it can't be implemented outside this
// package.
type Measurement interface {
	// Value returns the numeric value, expressed in Unit.
	Value() float64
	// Unit returns the unit symbol, e.g. "km".
	Unit() string
	// Type returns the measure the measurement belongs to.
	Type() MeasureType
	// ConvertToMeasurement returns the measurement expressed in targetUnit.
	ConvertToMeasurement(targetUnit string) (Measurement, error)
	fmt.Stringer

	// AddMeasurement returns the sum of the measurement and m, expressed in the
	// receiver's unit. It fails if m can't be converted to that unit.
	AddMeasurement(m Measurement) (Measurement, error)
	// SubMeasurement returns the difference of the measurement and m, expressed
	// in the receiver's unit. It fails if m can't be converted to that unit.
	SubMeasurement(m Measurement) (Measurement, error)
	// MulMeasurement returns the measurement multiplied by multiplier.
	MulMeasurement(multiplier float64) Measurement
	// DivMeasurement returns the measurement divided by divisor. It fails if
	// divisor is zero.
	DivMeasurement(divisor float64) (Measurement, error)

	base() *baseMeasurement
}

// baseMeasurement is the representation shared by all measures. The measure
// a value belongs to is the measure of its conversion table.
type baseMeasurement struct {
	value           float64
	unit            string
	conversionRates rates.Table
}

var errNilMeasurement = errors.New("measurement is nil")

func (b *baseMeasurement) base() *baseMeasurement {
	return b
}

// Value returns the numeric value, expressed in Unit.
func (b *baseMeasurement) Value() float64 {
	return b.value
}

// Unit returns the unit symbol.
func (b *baseMeasurement) Unit() string {
	return b.unit
}

// Type returns the measure the measurement belongs to.
func (b *baseMeasurement) Type() MeasureType {
	return MeasureType(b.conversionRates.Measure())
}

// String formats the measurement as "<value> <unit>", the format ParseString
// accepts. The value is formatted like fmt's %v, in the shortest form that
// parses back to the same float64.
func (b *baseMeasurement) String() string {
	// strconv instead of fmt.Sprintf: same output, a third of the time.
	var buf [64]byte

	s := strconv.AppendFloat(buf[:0], b.value, 'g', -1, 64)
	s = append(s, ' ')
	s = append(s, b.unit...)

	return string(s)
}

// AddMeasurement returns the sum of b and m, expressed in b's unit.
func (b *baseMeasurement) AddMeasurement(m Measurement) (Measurement, error) {
	v, err := b.valueOf(m)
	if err != nil {
		return nil, err
	}

	bm := b.withValue(b.value + v)

	return &bm, nil
}

// SubMeasurement returns the difference of b and m, expressed in b's unit.
func (b *baseMeasurement) SubMeasurement(m Measurement) (Measurement, error) {
	v, err := b.valueOf(m)
	if err != nil {
		return nil, err
	}

	bm := b.withValue(b.value - v)

	return &bm, nil
}

// MulMeasurement returns b multiplied by multiplier.
func (b *baseMeasurement) MulMeasurement(multiplier float64) Measurement {
	bm := b.withValue(b.value * multiplier)

	return &bm
}

// DivMeasurement returns b divided by divisor.
func (b *baseMeasurement) DivMeasurement(divisor float64) (Measurement, error) {
	if divisor == 0 {
		return nil, errZeroDivision
	}

	bm := b.withValue(b.value / divisor)

	return &bm, nil
}

// ConvertToMeasurement returns b expressed in targetUnit.
func (b *baseMeasurement) ConvertToMeasurement(targetUnit string) (Measurement, error) {
	bm, err := b.convertTo(targetUnit)
	if err != nil {
		return nil, err
	}

	return &bm, nil
}

func (b *baseMeasurement) convertTo(targetUnit string) (baseMeasurement, error) {
	targetUnit = rates.CanonicalUnit(targetUnit)

	rate, ok := b.conversionRates.Rate(b.unit, targetUnit)
	if !ok {
		return baseMeasurement{}, unknownUnitError(b.Type(), targetUnit)
	}

	return baseMeasurement{value: b.value * rate, unit: targetUnit, conversionRates: b.conversionRates}, nil
}

// unsafeConvertTo converts b to targetUnit, which must belong to b's measure.
func (b *baseMeasurement) unsafeConvertTo(targetUnit string) baseMeasurement {
	return baseMeasurement{value: b.unsafeGetValueIn(targetUnit), unit: targetUnit, conversionRates: b.conversionRates}
}

// unsafeGetValueIn returns b's value expressed in targetUnit, which must belong
// to b's measure.
func (b *baseMeasurement) unsafeGetValueIn(targetUnit string) float64 {
	if b.unit == targetUnit {
		return b.value
	}

	rate, _ := b.conversionRates.Rate(b.unit, targetUnit)

	return b.value * rate
}

// valueOf returns m's value expressed in b's unit, or an error if m belongs to
// another measure or can't be converted to b's unit.
func (b *baseMeasurement) valueOf(m Measurement) (float64, error) {
	if m == nil {
		return 0, errNilMeasurement
	}

	if b.Type() != m.Type() {
		return 0, fmt.Errorf("incompatible types %s and %s", b.Type(), m.Type())
	}

	other := m.base()

	rate, ok := other.conversionRates.Rate(other.unit, b.unit)
	if !ok {
		return 0, fmt.Errorf("can't convert %s to %s: the units belong to different %s measures",
			other.unit, b.unit, b.Type())
	}

	return other.value * rate, nil
}

// withValue returns a measurement with b's unit and measure and the given value.
func (b *baseMeasurement) withValue(value float64) baseMeasurement {
	return baseMeasurement{value: value, unit: b.unit, conversionRates: b.conversionRates}
}

// NewMeasurement creates a measurement of one of the predefined measures,
// detected from unit. Use the measure-specific constructors (NewDistance, …)
// when the measure is known in advance.
func NewMeasurement(value float64, unit string) (Measurement, error) {
	mt := DetectMeasureType(unit)
	if mt == MeasureCustom {
		return nil, errors.New("unknown unit " + unit)
	}

	bm, err := newBaseMeasurement(value, unit, measureToConversionRatesMapping[mt])
	if err != nil {
		return nil, err
	}

	return &bm, nil
}

// ParseString parses a measurement of a predefined measure from a string of the
// form "<value> <unit>", e.g. "10 km" or "1.5e3 m". The value must be finite.
func ParseString(s string) (Measurement, error) {
	value, unit, err := parseValueAndUnit(s)
	if err != nil {
		return nil, err
	}

	return NewMeasurement(value, unit)
}

// NewBaseMeasurement creates a measurement from custom conversion factors.
// The measurement's Type is MeasureCustom.
//
// It builds a new conversion table on every call. To create many measurements
// of the same custom measure, use NewMeasurementCreator, which builds the table
// once and gives the measure a name.
func NewBaseMeasurement(value float64, unit string, conversionFactors conversion.Factors) (Measurement, error) {
	table, err := rates.NewTable(string(MeasureCustom), conversionFactors)
	if err != nil {
		return nil, err
	}

	bm, err := newBaseMeasurement(value, unit, table)
	if err != nil {
		return nil, err
	}

	return &bm, nil
}

// newBaseMeasurement validates unit against table and creates a measurement.
func newBaseMeasurement(value float64, unit string, table rates.Table) (baseMeasurement, error) {
	canonical := rates.CanonicalUnit(unit)
	if !table.Has(canonical) {
		return baseMeasurement{}, unknownUnitError(MeasureType(table.Measure()), unit)
	}

	return baseMeasurement{value: value, unit: canonical, conversionRates: table}, nil
}

// parseValueAndUnit splits a "<value> <unit>" string and parses the value.
func parseValueAndUnit(s string) (float64, string, error) {
	parts := strings.Fields(s)
	if len(parts) != 2 {
		return 0, "", fmt.Errorf("invalid string format: %s", s)
	}

	value, err := strconv.ParseFloat(parts[0], 64)
	if err != nil {
		return 0, "", fmt.Errorf("failed to parse value %v: %w", parts[0], err)
	}

	if math.IsNaN(value) || math.IsInf(value, 0) {
		return 0, "", fmt.Errorf("value %s is not a finite number", parts[0])
	}

	return value, parts[1], nil
}

// unknownUnitError builds the error for a unit that doesn't belong to measure mt.
// It concatenates instead of using fmt.Errorf, which formats the named string
// type MeasureType through reflection.
func unknownUnitError(mt MeasureType, unit string) error {
	return errors.New("unknown " + string(mt) + " unit " + unit)
}
