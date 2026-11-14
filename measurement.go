package units

import (
	"cmp"
	"fmt"
	"io"
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
	// Format implements fmt.Formatter: %v prints String(), %.2f formats the
	// value and appends the unit ("5.00 km"). See the package documentation.
	fmt.Formatter

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
	// CompareMeasurement compares the measurement with m as Comparer.Compare
	// does: it returns -1, 0 or +1 as the measurement is less than, equal to or
	// greater than m. It fails if m can't be converted to the measurement's unit.
	CompareMeasurement(m Measurement) (int, error)

	base() *baseMeasurement
}

// baseMeasurement is the representation shared by all measures. The measure
// a value belongs to is the measure of its conversion table.
type baseMeasurement struct {
	value           float64
	unit            string
	conversionRates rates.Table
}

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

// Format implements fmt.Formatter, so fmt verbs format measurements sensibly:
//
//   - %v and %s print String(), honoring width and the '-' flag: "5 km".
//   - %b, %e, %E, %f, %F, %g, %G, %x and %X format the value with the given
//     flags, width and precision, then append the unit: %.2f gives "5.00 km".
//   - %q prints String() quoted: "\"5 km\"".
//   - %#v prints a Go-like form: units.Measurement{Value:5, Unit:"km", Type:"distance"}.
//   - Other verbs print fmt's usual error form: %!d(units.Measurement=5 km).
func (b *baseMeasurement) Format(f fmt.State, verb rune) {
	switch verb {
	case 'v', 's':
		if verb == 'v' && f.Flag('#') {
			_, _ = fmt.Fprintf(f, "units.Measurement{Value:%#v, Unit:%q, Type:%q}", b.value, b.unit, b.Type())

			return
		}

		if _, ok := f.Width(); !ok && !f.Flag('-') {
			_, _ = io.WriteString(f, b.String()) // fast path: fmt.Print, Println, plain %v

			return
		}

		_, _ = fmt.Fprintf(f, fmt.FormatString(f, 's'), b.String())
	case 'q':
		_, _ = fmt.Fprintf(f, fmt.FormatString(f, 'q'), b.String())
	case 'b', 'e', 'E', 'f', 'F', 'g', 'G', 'x', 'X':
		_, _ = fmt.Fprintf(f, fmt.FormatString(f, verb), b.value)
		_, _ = io.WriteString(f, " "+b.unit)
	default:
		_, _ = fmt.Fprintf(f, "%%!%c(units.Measurement=%s)", verb, b.String())
	}
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
		return nil, ErrZeroDivision
	}

	bm := b.withValue(b.value / divisor)

	return &bm, nil
}

// CompareMeasurement compares b with m. Like Comparer.Compare, it compares
// values in different units in the measure's base unit.
func (b *baseMeasurement) CompareMeasurement(m Measurement) (int, error) {
	if m != nil && m.base().conversionRates == b.conversionRates {
		return cmp.Compare(b.commonValues(m.base())), nil
	}

	// Another measure, or a measure with the same name but its own units: as in
	// AddMeasurement, compare in b's unit if m converts to it.
	v, err := b.valueOf(m)
	if err != nil {
		return 0, err
	}

	return cmp.Compare(b.value, v), nil
}

// commonValues returns the values of b and o, which must share b's conversion
// table, in a common unit: their own if they have the same unit, otherwise the
// reference unit of the table's factors (the base unit, for predefined
// measures). Either way the result doesn't depend on the order of b and o, so
// comparisons built on it are symmetric.
func (b *baseMeasurement) commonValues(o *baseMeasurement) (float64, float64) {
	if b.unit == o.unit {
		return b.value, o.value
	}

	bf, _ := b.conversionRates.Factor(b.unit)
	of, _ := b.conversionRates.Factor(o.unit)

	return b.value * bf, o.value * of
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
		return 0, newError(
			ErrIncompatibleMeasures,
			"incompatible types "+string(b.Type())+" and "+string(m.Type()),
			nil,
		)
	}

	other := m.base()

	rate, ok := other.conversionRates.Rate(other.unit, b.unit)
	if !ok {
		return 0, newError(ErrIncompatibleMeasures, "can't convert "+other.unit+" to "+b.unit+
			": the units belong to different "+string(b.Type())+" measures", nil)
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
	pu, ok := predefinedUnits[rates.CanonicalUnit(unit)]
	if !ok {
		return nil, unknownUnitError("", unit)
	}

	bm := baseMeasurement{value: value, unit: pu.unit, conversionRates: pu.table}

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
		return nil, newError(ErrInvalidFactors, err.Error(), nil)
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
	// FieldsSeq instead of Fields: no slice is allocated for the fields.
	var parts [2]string

	n := 0
	for field := range strings.FieldsSeq(s) {
		if n == len(parts) {
			n++ // a third field: the format is invalid

			break
		}

		parts[n] = field
		n++
	}

	if n != len(parts) {
		return 0, "", newError(ErrInvalidFormat, "invalid string format: "+s, nil)
	}

	value, err := strconv.ParseFloat(parts[0], 64)
	if err != nil {
		return 0, "", newError(ErrInvalidFormat, "failed to parse value "+parts[0]+": "+err.Error(), err)
	}

	if math.IsNaN(value) || math.IsInf(value, 0) {
		return 0, "", newError(ErrInvalidFormat, "value "+parts[0]+" is not a finite number", nil)
	}

	return value, parts[1], nil
}
