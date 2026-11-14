package units

import (
	"cmp"
	"fmt"

	"github.com/asaphin/go-physics-units/ratio"
)

// Measure is implemented by the marker type of a user-defined measure. Its
// Units method returns the measure's units, and must return the same UnitSet
// on every call: build it once, in a package-level variable.
//
// A measure takes four declarations, and then has the same typed API as the
// predefined measures:
//
//	type currency struct{}
//
//	var currencyUnits = units.MustNewUnitSet("currency", "USD", conversion.Factors{"USD": 1, "EUR": 1.1})
//
//	func (currency) Units() units.UnitSet { return currencyUnits }
//
//	type Money = units.Quantity[currency]
//
// Any type can be a marker; an empty struct costs nothing. UnitOf and
// ParseQuantity fail if Units returns the zero UnitSet, and the methods of a
// Quantity of such a measure panic.
type Measure interface {
	Units() UnitSet
}

// Unit is a unit of the user-defined measure M. Look one up by symbol with
// UnitOf, once, and use it for any number of values. The zero Unit is M's base
// unit.
type Unit[M Measure] struct {
	i uint16
}

// UnitOf returns the unit of M with the given symbol. It fails with
// [ErrUnknownUnit] if M has no such unit, and with [ErrInvalidMeasure] if M's
// Units method returns the zero UnitSet. Both forms of the micro sign are
// accepted.
func UnitOf[M Measure](symbol string) (Unit[M], error) {
	s, err := setOf[M]()
	if err != nil {
		return Unit[M]{}, err
	}

	i, err := s.unitOf(symbol)
	if err != nil {
		return Unit[M]{}, err
	}

	return Unit[M]{i: i}, nil
}

// MustUnitOf is like UnitOf but panics if the unit doesn't exist. It is
// intended for package-level variables:
//
//	var usd = units.MustUnitOf[currency]("USD")
func MustUnitOf[M Measure](symbol string) Unit[M] {
	u, err := UnitOf[M](symbol)
	if err != nil {
		panic(err)
	}

	return u
}

// setOf returns M's unit set, or an error if M has none.
func setOf[M Measure]() (*unitSet, error) {
	var m M

	s := m.Units().s
	if s == nil {
		return nil, newError(ErrInvalidMeasure,
			"the Units method of the measure returns the zero UnitSet: build it with NewUnitSet", nil)
	}

	return s, nil
}

// set returns M's unit set. M's Units method must return a non-zero set;
// UnitOf and ParseQuantity check it.
func set[M Measure]() *unitSet {
	var m M

	return m.Units().s
}

// String returns the unit's symbol.
func (u Unit[M]) String() string {
	return set[M]().symbols[u.i]
}

// Factor returns how many reference units there are in one u: the factor u
// was given to NewUnitSet.
func (u Unit[M]) Factor() float64 {
	return set[M]().factors[u.i]
}

// Quantity is a value of the user-defined measure M: a value and a Unit[M].
// It has the typed API of the predefined measures, and implements Measurement.
// See Measure for how to declare M.
//
// Like the predefined measures, a Quantity is a small immutable value: every
// method returns a new value. The zero Quantity is 0 in M's base unit.
type Quantity[M Measure] struct {
	value float64
	unit  Unit[M]
}

// NewQuantity returns a quantity of value in unit.
func NewQuantity[M Measure](value float64, unit Unit[M]) Quantity[M] {
	return Quantity[M]{value: value, unit: unit}
}

// ParseQuantity parses a quantity of M from a string such as "10 USD". It
// fails with [ErrInvalidFormat] unless the string has the form
// "<value> <symbol>" with a finite value, and with [ErrUnknownUnit] if the
// symbol isn't a unit of M.
func ParseQuantity[M Measure](s string) (Quantity[M], error) {
	value, symbol, err := parseValueAndUnit(s)
	if err != nil {
		return Quantity[M]{}, err
	}

	u, err := UnitOf[M](symbol)
	if err != nil {
		return Quantity[M]{}, err
	}

	return Quantity[M]{value: value, unit: u}, nil
}

// ToQuantity returns m as a Quantity[M]. It fails with [ErrIncompatibleMeasures]
// if m isn't one, and with [ErrInvalidMeasure] if M's Units method returns the
// zero UnitSet.
func ToQuantity[M Measure](m Measurement) (Quantity[M], error) {
	q, ok := m.(Quantity[M])
	if !ok {
		s, err := setOf[M]()
		if err != nil {
			return Quantity[M]{}, err
		}

		return Quantity[M]{}, newError(ErrIncompatibleMeasures, "not a "+string(s.name)+" measure", nil)
	}

	return q, nil
}

// Value returns the numeric value, expressed in Unit.
func (q Quantity[M]) Value() float64 {
	return q.value
}

// Unit returns the unit of the value.
func (q Quantity[M]) Unit() Unit[M] {
	return q.unit
}

// Symbol returns the symbol of the unit.
func (q Quantity[M]) Symbol() string {
	return q.unit.String()
}

// Type returns the name of the measure, as given to NewUnitSet.
func (q Quantity[M]) Type() MeasureType {
	return set[M]().name
}

// String formats the quantity as "<value> <symbol>", such as "10 USD", the
// form ParseQuantity accepts.
func (q Quantity[M]) String() string {
	return formatString(q.value, q.unit.String())
}

// Format implements fmt.Formatter, as the predefined measures do. See the
// package documentation.
func (q Quantity[M]) Format(f fmt.State, verb rune) {
	formatMeasurement(f, verb, "Quantity", q.value, q.unit.String())
}

// AppendText implements encoding.TextAppender: it appends the String form,
// such as "10 USD". It fails with [ErrOutOfRange] if the value isn't finite.
func (q Quantity[M]) AppendText(b []byte) ([]byte, error) {
	return appendText(b, q.value, q.unit.String())
}

// MarshalText implements encoding.TextMarshaler, so quantities encode in JSON
// as strings such as "10 USD". It fails like AppendText.
func (q Quantity[M]) MarshalText() ([]byte, error) {
	return q.AppendText(nil)
}

// UnmarshalText implements encoding.TextUnmarshaler: it parses text as
// ParseQuantity does and stores the result in q. Unmarshalling needs a pointer
// receiver; it changes only the value it is called on. It fails like
// [ParseQuantity], and leaves q unchanged when it does.
func (q *Quantity[M]) UnmarshalText(text []byte) error {
	value, err := ParseQuantity[M](string(text))
	if err != nil {
		return err
	}

	*q = value

	return nil
}

// ConvertTo returns the quantity expressed in unit.
func (q Quantity[M]) ConvertTo(unit Unit[M]) Quantity[M] {
	return Quantity[M]{value: q.value * set[M]().rate(q.unit.i, unit.i), unit: unit}
}

// ConvertToBaseUnits returns the quantity expressed in M's base unit.
func (q Quantity[M]) ConvertToBaseUnits() Quantity[M] {
	return q.ConvertTo(Unit[M]{})
}

// Add returns the sum of the quantities, expressed in the receiver's unit.
func (q Quantity[M]) Add(other Quantity[M]) Quantity[M] {
	return Quantity[M]{value: q.value + other.value*set[M]().rate(other.unit.i, q.unit.i), unit: q.unit}
}

// Sub returns the difference of the quantities, expressed in the receiver's
// unit.
func (q Quantity[M]) Sub(other Quantity[M]) Quantity[M] {
	return Quantity[M]{value: q.value - other.value*set[M]().rate(other.unit.i, q.unit.i), unit: q.unit}
}

// Mul returns the quantity multiplied by multiplier.
func (q Quantity[M]) Mul(multiplier float64) Quantity[M] {
	return Quantity[M]{value: q.value * multiplier, unit: q.unit}
}

// Div returns the quantity divided by divisor. It fails with [ErrZeroDivision]
// if divisor is zero.
func (q Quantity[M]) Div(divisor float64) (Quantity[M], error) {
	if divisor == 0 {
		return Quantity[M]{}, ErrZeroDivision
	}

	return Quantity[M]{value: q.value / divisor, unit: q.unit}, nil
}

// Compare compares the quantities; see Comparer. Quantities in different
// units are compared through their factors.
func (q Quantity[M]) Compare(other Quantity[M]) int {
	return cmp.Compare(q.common(other))
}

// Equal reports whether the quantities are equal; see Comparer.
func (q Quantity[M]) Equal(other Quantity[M]) bool {
	return q.Compare(other) == 0
}

// EqualWithinRel reports whether the quantities are equal within the relative
// tolerance tol; see Comparer.
func (q Quantity[M]) EqualWithinRel(other Quantity[M], tol float64) bool {
	x, y := q.common(other)

	return equalWithinRel(x, y, tol)
}

// Ratio returns q/other, the ratio of the quantities, in [ratio.BaseUnit]:
// convert it with ConvertTo(ratio.Percent) for a percentage. It fails with
// [ErrZeroDivision] if other is zero.
func (q Quantity[M]) Ratio(other Quantity[M]) (Ratio, error) {
	x, y := q.common(other)
	if y == 0 {
		return Ratio{}, ErrZeroDivision
	}

	return Ratio{value: x / y, unit: ratio.BaseUnit}, nil
}

// common returns the values of the quantities in a common unit, as the
// predefined measures do: their own if they share it, otherwise the reference
// unit of the factors.
func (q Quantity[M]) common(other Quantity[M]) (float64, float64) {
	if q.unit == other.unit {
		return q.value, other.value
	}

	f := set[M]().factors

	return q.value * f[q.unit.i], other.value * f[other.unit.i]
}

// ConvertToMeasurement is ConvertTo for generic code: it takes a unit symbol.
// It fails with [ErrUnknownUnit] if symbol isn't a unit of M.
func (q Quantity[M]) ConvertToMeasurement(symbol string) (Measurement, error) {
	i, err := set[M]().unitOf(symbol)
	if err != nil {
		return nil, err
	}

	return q.ConvertTo(Unit[M]{i: i}), nil
}

// AddMeasurement is Add for generic code. It fails with
// [ErrIncompatibleMeasures] if m isn't a Quantity[M].
func (q Quantity[M]) AddMeasurement(m Measurement) (Measurement, error) {
	v, err := set[M]().valueIn(m, q.unit.i)
	if err != nil {
		return nil, err
	}

	return Quantity[M]{value: q.value + v, unit: q.unit}, nil
}

// SubMeasurement is Sub for generic code. It fails with
// [ErrIncompatibleMeasures] if m isn't a Quantity[M].
func (q Quantity[M]) SubMeasurement(m Measurement) (Measurement, error) {
	v, err := set[M]().valueIn(m, q.unit.i)
	if err != nil {
		return nil, err
	}

	return Quantity[M]{value: q.value - v, unit: q.unit}, nil
}

// MulMeasurement is Mul for generic code.
func (q Quantity[M]) MulMeasurement(multiplier float64) Measurement {
	return q.Mul(multiplier)
}

// DivMeasurement is Div for generic code. It fails with [ErrZeroDivision] if
// divisor is zero.
func (q Quantity[M]) DivMeasurement(divisor float64) (Measurement, error) {
	if divisor == 0 {
		return nil, ErrZeroDivision
	}

	return Quantity[M]{value: q.value / divisor, unit: q.unit}, nil
}

// CompareMeasurement is Compare for generic code. It fails with
// [ErrIncompatibleMeasures] if m isn't a Quantity[M].
func (q Quantity[M]) CompareMeasurement(m Measurement) (int, error) {
	return set[M]().compare(q.value, q.unit.i, m)
}

func (q Quantity[M]) generic() genericValue {
	return genericValue{value: q.value, unit: q.unit.i, set: set[M]()}
}
