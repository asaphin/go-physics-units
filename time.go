package units

import (
	"cmp"
	"math"
	builtintime "time"

	"github.com/asaphin/go-physics-units/time"
)

// Time is a time (duration) measurement. Its unit symbols are defined in the
// time package of this module.
type Time interface {
	Measurement
	UnitConverter[Time]
	Arithmetics[Time]
	Comparer[Time]

	// ToDuration converts the time into a standard library time.Duration,
	// rounded to the nearest nanosecond. It fails with ErrOutOfRange if the time
	// is NaN or outside the range of time.Duration, about ±292 years.
	ToDuration() (builtintime.Duration, error)
}

// timeImplementation is the concrete implementation of Time.
type timeImplementation struct {
	baseMeasurement
}

// ConvertTo returns the time expressed in unit.
func (t *timeImplementation) ConvertTo(unit string) (Time, error) {
	msm, err := t.convertTo(unit)
	if err != nil {
		return nil, err
	}

	return &timeImplementation{msm}, nil
}

// MustConvertTo is like ConvertTo but panics if unit isn't a time unit.
func (t *timeImplementation) MustConvertTo(unit string) Time {
	msm, err := t.convertTo(unit)
	if err != nil {
		panic(err)
	}

	return &timeImplementation{msm}
}

// ConvertToBaseUnits returns the time expressed in seconds.
func (t *timeImplementation) ConvertToBaseUnits() Time {
	return &timeImplementation{t.unsafeConvertTo(time.BaseUnit)}
}

func (t *timeImplementation) valueInBaseUnits() float64 {
	return t.unsafeGetValueIn(time.BaseUnit)
}

func (t *timeImplementation) valueInUnits(unit string) float64 {
	return t.unsafeGetValueIn(unit)
}

// Add returns the sum of the times, expressed in the receiver's unit.
func (t *timeImplementation) Add(other Time) Time {
	return &timeImplementation{t.withValue(t.value + other.valueInUnits(t.unit))}
}

// Sub returns the difference of the times, expressed in the receiver's unit.
func (t *timeImplementation) Sub(other Time) Time {
	return &timeImplementation{t.withValue(t.value - other.valueInUnits(t.unit))}
}

// Mul returns the time multiplied by multiplier.
func (t *timeImplementation) Mul(multiplier float64) Time {
	return &timeImplementation{t.withValue(t.value * multiplier)}
}

// Div returns the time divided by divisor. It fails if divisor is zero.
func (t *timeImplementation) Div(divisor float64) (Time, error) {
	if divisor == 0 {
		return nil, ErrZeroDivision
	}

	return &timeImplementation{t.withValue(t.value / divisor)}, nil
}

// Compare compares the times; see Comparer.
func (t *timeImplementation) Compare(other Time) int {
	return cmp.Compare(t.commonValues(other.base()))
}

// Equal reports whether the times are equal; see Comparer.
func (t *timeImplementation) Equal(other Time) bool {
	return t.Compare(other) == 0
}

// EqualWithinRel reports whether the times are equal within the relative
// tolerance tol; see Comparer.
func (t *timeImplementation) EqualWithinRel(other Time, tol float64) bool {
	x, y := t.commonValues(other.base())

	return equalWithinRel(x, y, tol)
}

// NewTime creates a time of value in unit, one of the symbols defined in the
// time package of this module.
func NewTime(value float64, unit string) (Time, error) {
	msm, err := newBaseMeasurement(value, unit, timeRates)
	if err != nil {
		return nil, err
	}

	return &timeImplementation{msm}, nil
}

// newTime is an unchecked constructor for internal use with valid, canonical
// units.
func newTime(value float64, unit string) Time {
	return &timeImplementation{baseMeasurement{value: value, unit: unit, conversionRates: timeRates}}
}

const errTimeConversion incompatibleError = "not a time measure"

// ToTime converts m, for example a result of ParseString, into a Time.
// It fails if m isn't a time.
func ToTime(m Measurement) (Time, error) {
	if m == nil || m.Type() != MeasureTime {
		return nil, errTimeConversion
	}

	return &timeImplementation{*m.base()}, nil
}

// StringToTime parses a time from a string such as "10 min".
// It fails if the unit isn't a time unit.
func StringToTime(s string) (Time, error) {
	// Parse and construct directly: one allocation, where ParseString
	// followed by ToTime would allocate twice.
	value, unit, err := parseValueAndUnit(s)
	if err != nil {
		return nil, err
	}

	return NewTime(value, unit)
}

// ParseDuration converts a standard library time.Duration into a Time in seconds.
// Time.ToDuration converts back.
func ParseDuration(d builtintime.Duration) Time {
	return newTime(d.Seconds(), time.BaseUnit)
}

// ToDuration converts the time into a time.Duration, rounded to the nearest
// nanosecond.
func (t *timeImplementation) ToDuration() (builtintime.Duration, error) {
	// Through seconds: 1e9 is exact, while the s→ns rate, 1/1e-9 in binary
	// floating point, is not, and its error grows with the value.
	ns := math.Round(t.valueInBaseUnits() * float64(builtintime.Second))

	// Duration is an int64 count of nanoseconds, whose range [-2⁶³, 2⁶³) has
	// exact float64 bounds. NaN fails both comparisons.
	if !(ns >= math.MinInt64 && ns < -math.MinInt64) {
		return 0, newError(ErrOutOfRange, "time "+t.String()+" can't be represented as a time.Duration", nil)
	}

	return builtintime.Duration(ns), nil
}
