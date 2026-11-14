package units

import (
	"cmp"
	"fmt"

	"github.com/asaphin/go-physics-units/temperature"
)

// Temperature is a temperature, a point on a scale such as 20 °C: a value in a
// temperature.Unit. Create one with NewTemperature or StringToTemperature.
//
// Temperature scales have offsets, so temperatures follow different rules
// from the other measures (see docs/rfc-temperature.md):
//
//   - The difference of two temperatures is a TemperatureDifference, and a
//     temperature plus or minus a difference is a temperature. Adding two
//     temperatures, or multiplying one, isn't offered: 20 °C + 20 °C has no
//     physical meaning.
//   - A temperature is never below absolute zero: the constructors, the
//     parsers and the operations that could go below it fail with
//     [ErrOutOfRange] instead.
//   - Temperature doesn't implement Measurement, whose AddMeasurement and
//     MulMeasurement it couldn't support; ParseString rejects temperature
//     units.
//
// A Temperature is a small immutable value. The zero Temperature is 0 K.
type Temperature struct {
	value float64
	unit  temperature.Unit
}

// TemperatureDifference is a difference of temperatures, such as a rise of
// 10 °C, which is a rise of 18 °F: unlike a temperature, it converts without
// an offset, and has the arithmetic of the other measures. It can be
// negative. The zero TemperatureDifference is 0 K.
type TemperatureDifference struct {
	value float64
	unit  temperature.Unit
}

// NewTemperature returns a temperature of value in unit, such as
// NewTemperature(20, temperature.Celsius). It fails with [ErrOutOfRange] if the
// value is below absolute zero (-273.15 °C) or NaN, and with [ErrUnknownUnit] if
// unit isn't valid.
func NewTemperature(value float64, unit temperature.Unit) (Temperature, error) {
	if !unit.Valid() {
		return Temperature{}, unknownUnitError(MeasureTemperature, unit.String())
	}

	if !(value >= unit.AbsoluteZero()) {
		return Temperature{}, belowAbsoluteZero(value, unit)
	}

	return Temperature{value: value, unit: unit}, nil
}

// belowAbsoluteZero returns the error for a temperature below absolute zero.
func belowAbsoluteZero(value float64, unit temperature.Unit) error {
	return newError(ErrOutOfRange, "temperature "+formatString(value, unit.String())+
		" is below absolute zero, "+formatString(unit.AbsoluteZero(), unit.String()), nil)
}

// StringToTemperature parses a temperature from a string such as "20 °C". It
// fails with [ErrInvalidFormat] unless the string has the form "<value> <symbol>"
// with a finite value, with [ErrUnknownUnit] if the symbol isn't a temperature
// unit, and with [ErrOutOfRange] if the temperature is below absolute zero.
func StringToTemperature(s string) (Temperature, error) {
	value, unit, err := parseTemperature(s)
	if err != nil {
		return Temperature{}, err
	}

	return NewTemperature(value, unit)
}

// parseTemperature parses a temperature or a temperature difference.
func parseTemperature(s string) (float64, temperature.Unit, error) {
	value, symbol, err := parseValueAndUnit(s)
	if err != nil {
		return 0, 0, err
	}

	unit, ok := temperature.ParseUnit(symbol)
	if !ok {
		return 0, 0, unknownUnitError(MeasureTemperature, symbol)
	}

	return value, unit, nil
}

// Value returns the numeric value, expressed in Unit.
func (t Temperature) Value() float64 {
	return t.value
}

// Unit returns the unit of the value.
func (t Temperature) Unit() temperature.Unit {
	return t.unit
}

// Symbol returns the symbol of the unit, as Unit().String() does.
func (t Temperature) Symbol() string {
	return t.unit.String()
}

// Kelvin returns the temperature in kelvins, the absolute scale that physical
// formulas such as the ideal gas law need.
func (t Temperature) Kelvin() float64 {
	return temperature.Convert(t.value, t.unit, temperature.Kelvin)
}

// String formats the temperature as "<value> <symbol>", such as "20 °C", the
// form StringToTemperature accepts.
func (t Temperature) String() string {
	return formatString(t.value, t.unit.String())
}

// Format implements fmt.Formatter, as the other measures do. See the package
// documentation.
func (t Temperature) Format(state fmt.State, verb rune) {
	formatMeasurement(state, verb, "Temperature", t.value, t.unit.String())
}

// ConvertTo returns the temperature expressed in unit: 100 °C is 212 °F.
func (t Temperature) ConvertTo(unit temperature.Unit) Temperature {
	// Rounding can't take a temperature below absolute zero.
	return Temperature{value: max(temperature.Convert(t.value, t.unit, unit), unit.AbsoluteZero()), unit: unit}
}

// Sub returns the difference t - other, in t's unit: 20 °C - 50 °F is 10 °C.
func (t Temperature) Sub(other Temperature) TemperatureDifference {
	return TemperatureDifference{
		value: t.value - temperature.Convert(other.value, other.unit, t.unit),
		unit:  t.unit,
	}
}

// Add returns t raised by d, in t's unit: 20 °C + 18 °F is 30 °C. It fails
// with [ErrOutOfRange] if the result is below absolute zero.
func (t Temperature) Add(d TemperatureDifference) (Temperature, error) {
	return NewTemperature(t.value+d.value*temperature.DifferenceRate(d.unit, t.unit), t.unit)
}

// SubDifference returns t lowered by d, in t's unit. It fails with
// [ErrOutOfRange] if the result is below absolute zero.
func (t Temperature) SubDifference(d TemperatureDifference) (Temperature, error) {
	return NewTemperature(t.value-d.value*temperature.DifferenceRate(d.unit, t.unit), t.unit)
}

// Compare compares the temperatures, as Comparer.Compare does: exactly in the
// same unit, otherwise in kelvins.
func (t Temperature) Compare(other Temperature) int {
	return cmp.Compare(t.common(other))
}

// Equal reports whether the temperatures are equal; see Compare.
func (t Temperature) Equal(other Temperature) bool {
	return t.Compare(other) == 0
}

// EqualWithinRel reports whether the temperatures are equal within the
// relative tolerance tol, compared in kelvins unless they share a unit; see
// Comparer.
func (t Temperature) EqualWithinRel(other Temperature, tol float64) bool {
	x, y := t.common(other)

	return equalWithinRel(x, y, tol)
}

func (t Temperature) common(other Temperature) (float64, float64) {
	if t.unit == other.unit {
		return t.value, other.value
	}

	return t.Kelvin(), other.Kelvin()
}

// AppendText implements encoding.TextAppender: it appends the String form,
// such as "20 °C". It fails with [ErrOutOfRange] if the value is +Inf.
func (t Temperature) AppendText(b []byte) ([]byte, error) {
	return appendText(b, t.value, t.unit.String())
}

// MarshalText implements encoding.TextMarshaler, so temperatures encode in
// JSON as strings such as "20 °C". It fails like AppendText.
func (t Temperature) MarshalText() ([]byte, error) {
	return t.AppendText(nil)
}

// UnmarshalText implements encoding.TextUnmarshaler: it parses text as
// StringToTemperature does and stores the result in t. It fails like
// [StringToTemperature], and leaves t unchanged when it does.
func (t *Temperature) UnmarshalText(text []byte) error {
	value, err := StringToTemperature(string(text))
	if err != nil {
		return err
	}

	*t = value

	return nil
}

// NewTemperatureDifference returns a temperature difference of value in unit,
// such as NewTemperatureDifference(10, temperature.Celsius).
func NewTemperatureDifference(value float64, unit temperature.Unit) TemperatureDifference {
	return TemperatureDifference{value: value, unit: unit}
}

// StringToTemperatureDifference parses a temperature difference from a string
// such as "10 °C". It fails with [ErrInvalidFormat] unless the string has the
// form "<value> <symbol>" with a finite value, and with [ErrUnknownUnit] if the
// symbol isn't a temperature unit.
func StringToTemperatureDifference(s string) (TemperatureDifference, error) {
	value, unit, err := parseTemperature(s)
	if err != nil {
		return TemperatureDifference{}, err
	}

	return TemperatureDifference{value: value, unit: unit}, nil
}

// Value returns the numeric value, expressed in Unit.
func (d TemperatureDifference) Value() float64 {
	return d.value
}

// Unit returns the unit of the value.
func (d TemperatureDifference) Unit() temperature.Unit {
	return d.unit
}

// Symbol returns the symbol of the unit, as Unit().String() does.
func (d TemperatureDifference) Symbol() string {
	return d.unit.String()
}

// String formats the difference as "<value> <symbol>", such as "10 °C".
func (d TemperatureDifference) String() string {
	return formatString(d.value, d.unit.String())
}

// Format implements fmt.Formatter, as the other measures do.
func (d TemperatureDifference) Format(state fmt.State, verb rune) {
	formatMeasurement(state, verb, "TemperatureDifference", d.value, d.unit.String())
}

// ConvertTo returns the difference expressed in unit: a rise of 10 °C is a
// rise of 18 °F.
func (d TemperatureDifference) ConvertTo(unit temperature.Unit) TemperatureDifference {
	return TemperatureDifference{value: d.value * temperature.DifferenceRate(d.unit, unit), unit: unit}
}

// ConvertToBaseUnits returns the difference expressed in kelvins.
func (d TemperatureDifference) ConvertToBaseUnits() TemperatureDifference {
	return d.ConvertTo(temperature.BaseUnit)
}

// Add returns the sum of the differences, expressed in the receiver's unit.
func (d TemperatureDifference) Add(other TemperatureDifference) TemperatureDifference {
	return TemperatureDifference{value: d.value + other.ConvertTo(d.unit).value, unit: d.unit}
}

// Sub returns the difference of the differences, expressed in the receiver's
// unit.
func (d TemperatureDifference) Sub(other TemperatureDifference) TemperatureDifference {
	return TemperatureDifference{value: d.value - other.ConvertTo(d.unit).value, unit: d.unit}
}

// Mul returns the difference multiplied by multiplier.
func (d TemperatureDifference) Mul(multiplier float64) TemperatureDifference {
	return TemperatureDifference{value: d.value * multiplier, unit: d.unit}
}

// Div returns the difference divided by divisor. It fails with
// [ErrZeroDivision] if divisor is zero.
func (d TemperatureDifference) Div(divisor float64) (TemperatureDifference, error) {
	if divisor == 0 {
		return TemperatureDifference{}, ErrZeroDivision
	}

	return TemperatureDifference{value: d.value / divisor, unit: d.unit}, nil
}

// Compare compares the differences; see Comparer.
func (d TemperatureDifference) Compare(other TemperatureDifference) int {
	return cmp.Compare(d.common(other))
}

// Equal reports whether the differences are equal; see Comparer.
func (d TemperatureDifference) Equal(other TemperatureDifference) bool {
	return d.Compare(other) == 0
}

// EqualWithinRel reports whether the differences are equal within the
// relative tolerance tol; see Comparer.
func (d TemperatureDifference) EqualWithinRel(other TemperatureDifference, tol float64) bool {
	x, y := d.common(other)

	return equalWithinRel(x, y, tol)
}

func (d TemperatureDifference) common(other TemperatureDifference) (float64, float64) {
	if d.unit == other.unit {
		return d.value, other.value
	}

	return d.ConvertToBaseUnits().value, other.ConvertToBaseUnits().value
}

// AppendText implements encoding.TextAppender: it appends the String form,
// such as "10 °C". It fails with [ErrOutOfRange] if the value isn't finite, and
// with [ErrUnknownUnit] if the unit isn't valid.
func (d TemperatureDifference) AppendText(b []byte) ([]byte, error) {
	if !d.unit.Valid() {
		return b, unknownUnitError(MeasureTemperature, d.unit.String())
	}

	return appendText(b, d.value, d.unit.String())
}

// MarshalText implements encoding.TextMarshaler. It fails like AppendText.
func (d TemperatureDifference) MarshalText() ([]byte, error) {
	return d.AppendText(nil)
}

// UnmarshalText implements encoding.TextUnmarshaler: it parses text as
// StringToTemperatureDifference does and stores the result in d. It fails like
// [StringToTemperatureDifference], and leaves d unchanged when it does.
func (d *TemperatureDifference) UnmarshalText(text []byte) error {
	value, err := StringToTemperatureDifference(string(text))
	if err != nil {
		return err
	}

	*d = value

	return nil
}
