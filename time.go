package units

import (
	"math"
	builtintime "time"

	"github.com/asaphin/go-physics-units/frequency"
	"github.com/asaphin/go-physics-units/pace"
	"github.com/asaphin/go-physics-units/time"
)

// Time is a duration: a value in a time.Unit (of this module's time package),
// such as 10 min. Create one with NewTime, StringToTime or ParseDuration.
//
// A Time is a small immutable value: every method returns a new value, so copy
// and share it freely. The zero Time is 0 s.
type Time struct {
	value float64
	unit  time.Unit
}

// ParseDuration converts a standard library time.Duration into a Time in
// seconds. Time.ToDuration converts back.
func ParseDuration(d builtintime.Duration) Time {
	return Time{value: d.Seconds(), unit: time.BaseUnit}
}

// ToDuration converts the time into a standard library time.Duration, rounded
// to the nearest nanosecond. It fails with ErrOutOfRange if the time is NaN or
// outside the range of time.Duration, about ±292 years.
func (t Time) ToDuration() (builtintime.Duration, error) {
	// Through seconds: 1e9 is exact, while the s→ns rate, 1/1e-9 in binary
	// floating point, is not, and its error grows with the value.
	ns := math.Round(t.baseValue() * float64(builtintime.Second))

	// Duration is an int64 count of nanoseconds, whose range [-2⁶³, 2⁶³) has
	// exact float64 bounds. NaN fails both comparisons.
	if !(ns >= math.MinInt64 && ns < -math.MinInt64) {
		return 0, newError(ErrOutOfRange, "time "+t.String()+" can't be represented as a time.Duration", nil)
	}

	return builtintime.Duration(ns), nil
}

// DivideByDistance returns the pace t/d (s/m) in pace.BaseUnit. It fails with ErrZeroDivision if d is
// zero.
func (t Time) DivideByDistance(d Distance) (Pace, error) {
	if d.value == 0 {
		return Pace{}, ErrZeroDivision
	}

	return Pace{value: t.baseValue() / d.baseValue(), unit: pace.BaseUnit}, nil
}

// Frequency returns the frequency 1/t (Hz) in frequency.BaseUnit of something
// that happens once every t. It fails with ErrZeroDivision if t is zero.
func (t Time) Frequency() (Frequency, error) {
	if t.value == 0 {
		return Frequency{}, ErrZeroDivision
	}

	return Frequency{value: 1 / t.baseValue(), unit: frequency.BaseUnit}, nil
}
