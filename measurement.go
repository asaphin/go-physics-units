package units

import (
	"fmt"
	"io"
	"math"
	"strconv"
	"strings"

	"github.com/asaphin/go-physics-units/internal/rates"
)

// Measurement is the interface that every value of this package implements:
// the predefined measures (Distance, Time, …) and Quantity, the values of
// user-defined measures. Use it for code that handles measurements of any
// measure, such as the results of ParseString; the concrete types have typed
// methods that can't fail where these can.
//
// A value of a concrete type converts to Measurement and back:
//
//	m, _ := units.ParseString("10 km") // m holds a units.Distance
//	d, ok := m.(units.Distance)
//
// Measurement has an unexported method, so it can't be implemented outside
// this package.
type Measurement interface {
	// Value returns the numeric value, expressed in the unit.
	Value() float64
	// Symbol returns the unit symbol, such as "km".
	Symbol() string
	// Type returns the measure the measurement belongs to.
	Type() MeasureType
	fmt.Stringer
	// Format implements fmt.Formatter: %v prints String(), %.2f formats the
	// value and appends the symbol ("5.00 km"). See the package documentation.
	fmt.Formatter

	// ConvertToMeasurement returns the measurement expressed in the unit with
	// the given symbol. It fails with ErrUnknownUnit if the symbol isn't a
	// unit of the measure.
	ConvertToMeasurement(symbol string) (Measurement, error)
	// AddMeasurement returns the sum of the measurement and m, expressed in the
	// receiver's unit. It fails with ErrIncompatibleMeasures if m is nil or of
	// another measure.
	AddMeasurement(m Measurement) (Measurement, error)
	// SubMeasurement returns the difference of the measurement and m, expressed
	// in the receiver's unit. It fails like AddMeasurement.
	SubMeasurement(m Measurement) (Measurement, error)
	// MulMeasurement returns the measurement multiplied by multiplier.
	MulMeasurement(multiplier float64) Measurement
	// DivMeasurement returns the measurement divided by divisor. It fails with
	// ErrZeroDivision if divisor is zero.
	DivMeasurement(divisor float64) (Measurement, error)
	// CompareMeasurement compares the measurement with m as Comparer.Compare
	// does: it returns -1, 0 or +1 as the measurement is less than, equal to or
	// greater than m. It fails like AddMeasurement.
	CompareMeasurement(m Measurement) (int, error)

	generic() genericValue
}

// genericValue is a measurement seen by the generic code: its value, the
// index of its unit in its measure's unit set, and that set.
type genericValue struct {
	value float64
	unit  uint16
	set   *unitSet
}

// NewMeasurement creates a measurement of one of the predefined measures,
// detected from the unit symbol: NewMeasurement(10, "km") returns a Distance.
// It fails with ErrUnknownUnit if no predefined measure has the unit. Use the
// measure-specific constructors (NewDistance, …) when the measure is known in
// advance.
func NewMeasurement(value float64, symbol string) (Measurement, error) {
	pu, ok := predefinedUnits[rates.CanonicalUnit(symbol)]
	if !ok {
		return nil, unknownUnitError("", symbol)
	}

	return pu.box(value, pu.unit), nil
}

// ParseString parses a measurement of a predefined measure from a string of the
// form "<value> <symbol>", such as "10 km" or "1.5e3 m". The value must be
// finite. The result has the measure's concrete type: "10 km" gives a Distance.
func ParseString(s string) (Measurement, error) {
	value, symbol, err := parseValueAndUnit(s)
	if err != nil {
		return nil, err
	}

	return NewMeasurement(value, symbol)
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

// formatString formats a measurement as "<value> <symbol>". The value is
// formatted like fmt's %v, in the shortest form that parses back to the same
// float64.
func formatString(value float64, symbol string) string {
	// strconv instead of fmt.Sprintf: same output, a third of the time.
	var buf [64]byte

	return string(appendMeasurement(buf[:0], value, symbol))
}

// appendMeasurement appends the String form of a measurement to b.
func appendMeasurement(b []byte, value float64, symbol string) []byte {
	b = strconv.AppendFloat(b, value, 'g', -1, 64)
	b = append(b, ' ')

	return append(b, symbol...)
}

// appendText appends the text form of a measurement to b, for AppendText. It
// fails for a value that isn't finite, which ParseString wouldn't accept.
func appendText(b []byte, value float64, symbol string) ([]byte, error) {
	if math.IsNaN(value) || math.IsInf(value, 0) {
		return b, newError(ErrOutOfRange, "value "+strconv.FormatFloat(value, 'g', -1, 64)+
			" "+symbol+" can't be encoded as text: it isn't finite", nil)
	}

	return appendMeasurement(b, value, symbol), nil
}

// formatMeasurement implements fmt.Formatter for every measurement type:
//
//   - %v and %s print String(), honoring width and the '-' flag: "5 km".
//   - %b, %e, %E, %f, %F, %g, %G, %x and %X format the value with the given
//     flags, width and precision, then append the symbol: %.2f gives "5.00 km".
//   - %q prints String() quoted: "\"5 km\"".
//   - %#v prints a Go-like form: units.Distance{Value:5, Unit:"km"}.
//   - Other verbs print fmt's usual error form: %!d(units.Distance=5 km).
func formatMeasurement(f fmt.State, verb rune, typeName string, value float64, symbol string) {
	switch verb {
	case 'v', 's':
		if verb == 'v' && f.Flag('#') {
			_, _ = fmt.Fprintf(f, "units.%s{Value:%#v, Unit:%q}", typeName, value, symbol)

			return
		}

		if _, ok := f.Width(); !ok && !f.Flag('-') {
			_, _ = io.WriteString(f, formatString(value, symbol)) // fast path: fmt.Print, Println, plain %v

			return
		}

		_, _ = fmt.Fprintf(f, fmt.FormatString(f, 's'), formatString(value, symbol))
	case 'q':
		_, _ = fmt.Fprintf(f, fmt.FormatString(f, 'q'), formatString(value, symbol))
	case 'b', 'e', 'E', 'f', 'F', 'g', 'G', 'x', 'X':
		_, _ = fmt.Fprintf(f, fmt.FormatString(f, verb), value)
		_, _ = io.WriteString(f, " "+symbol)
	default:
		_, _ = fmt.Fprintf(f, "%%!%c(units.%s=%s)", verb, typeName, formatString(value, symbol))
	}
}
