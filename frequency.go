package units

import (
	"cmp"

	"github.com/asaphin/go-physics-units/frequency"
	"github.com/asaphin/go-physics-units/time"
)

// Frequency is a frequency measurement. Its unit symbols are defined in the frequency
// package.
type Frequency interface {
	Measurement
	UnitConverter[Frequency]
	Arithmetics[Frequency]
	Comparer[Frequency]

	// Period returns the period 1/f in time.BaseUnit (s). It fails if f is
	// zero.
	Period() (Time, error)
}

// frequencyImplementation is the concrete implementation of Frequency.
type frequencyImplementation struct {
	baseMeasurement
}

// ConvertTo returns the frequency expressed in unit.
func (f *frequencyImplementation) ConvertTo(unit string) (Frequency, error) {
	msm, err := f.convertTo(unit)
	if err != nil {
		return nil, err
	}

	return &frequencyImplementation{msm}, nil
}

// MustConvertTo is like ConvertTo but panics if unit isn't a frequency unit.
func (f *frequencyImplementation) MustConvertTo(unit string) Frequency {
	msm, err := f.convertTo(unit)
	if err != nil {
		panic(err)
	}

	return &frequencyImplementation{msm}
}

// ConvertToBaseUnits returns the frequency expressed in hertz.
func (f *frequencyImplementation) ConvertToBaseUnits() Frequency {
	return &frequencyImplementation{f.unsafeConvertTo(frequency.BaseUnit)}
}

func (f *frequencyImplementation) valueInBaseUnits() float64 {
	return f.unsafeGetValueIn(frequency.BaseUnit)
}

func (f *frequencyImplementation) valueInUnits(unit string) float64 {
	return f.unsafeGetValueIn(unit)
}

// Add returns the sum of the frequencies, expressed in the receiver's unit.
func (f *frequencyImplementation) Add(other Frequency) Frequency {
	return &frequencyImplementation{f.withValue(f.value + other.valueInUnits(f.unit))}
}

// Sub returns the difference of the frequencies, expressed in the receiver's unit.
func (f *frequencyImplementation) Sub(other Frequency) Frequency {
	return &frequencyImplementation{f.withValue(f.value - other.valueInUnits(f.unit))}
}

// Mul returns the frequency multiplied by multiplier.
func (f *frequencyImplementation) Mul(multiplier float64) Frequency {
	return &frequencyImplementation{f.withValue(f.value * multiplier)}
}

// Div returns the frequency divided by divisor. It fails if divisor is zero.
func (f *frequencyImplementation) Div(divisor float64) (Frequency, error) {
	if divisor == 0 {
		return nil, ErrZeroDivision
	}

	return &frequencyImplementation{f.withValue(f.value / divisor)}, nil
}

// Period returns the period 1/f in s, or an error if f is zero.
func (f *frequencyImplementation) Period() (Time, error) {
	if f.value == 0 {
		return nil, ErrZeroDivision
	}

	return newTime(1/f.valueInBaseUnits(), time.BaseUnit), nil
}

// Compare compares the frequencies; see Comparer.
func (f *frequencyImplementation) Compare(other Frequency) int {
	return cmp.Compare(f.commonValues(other.base()))
}

// Equal reports whether the frequencies are equal; see Comparer.
func (f *frequencyImplementation) Equal(other Frequency) bool {
	return f.Compare(other) == 0
}

// EqualWithinRel reports whether the frequencies are equal within the relative
// tolerance tol; see Comparer.
func (f *frequencyImplementation) EqualWithinRel(other Frequency, tol float64) bool {
	x, y := f.commonValues(other.base())

	return equalWithinRel(x, y, tol)
}

// NewFrequency creates a frequency of value in unit, one of the symbols defined in
// the frequency package.
func NewFrequency(value float64, unit string) (Frequency, error) {
	msm, err := newBaseMeasurement(value, unit, frequencyRates)
	if err != nil {
		return nil, err
	}

	return &frequencyImplementation{msm}, nil
}

const errFrequencyConversion incompatibleError = "not a frequency measure"

// ToFrequency converts m, for example a result of ParseString, into a Frequency.
// It fails if m isn't a frequency.
func ToFrequency(m Measurement) (Frequency, error) {
	if m == nil || m.Type() != MeasureFrequency {
		return nil, errFrequencyConversion
	}

	return &frequencyImplementation{*m.base()}, nil
}

// StringToFrequency parses a frequency from a string such as "50 Hz".
// It fails if the unit isn't a frequency unit.
func StringToFrequency(s string) (Frequency, error) {
	// Parse and construct directly: one allocation, where ParseString
	// followed by ToFrequency would allocate twice.
	value, unit, err := parseValueAndUnit(s)
	if err != nil {
		return nil, err
	}

	return NewFrequency(value, unit)
}
