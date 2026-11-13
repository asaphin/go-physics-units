package units_test

import (
	"errors"
	"math"
	"math/rand/v2"
	"testing"
	builtintime "time"

	units "github.com/asaphin/go-physics-units"
	"github.com/asaphin/go-physics-units/time"
)

func TestParseDuration(t *testing.T) {
	tm := units.ParseDuration(90*builtintime.Minute + 500*builtintime.Millisecond)

	assertMeasurement(t, tm, 5400.5, time.BaseUnit)
	assertMeasurement(t, tm.MustConvertTo(time.Hour), 5400.5/3600, time.Hour)
}

func TestTime_ToDuration(t *testing.T) {
	cases := []struct {
		value float64
		unit  string
		want  builtintime.Duration
	}{
		{1.5, time.Hour, 90 * builtintime.Minute},
		{1, time.Year, 8766 * builtintime.Hour}, // a Julian year: 365.25 days
		{1, time.Nanosecond, 1},
		{0.4, time.Nanosecond, 0},
		{0.5, time.Nanosecond, 1}, // halves round away from zero
		{-1.5, time.Nanosecond, -2},
		{-250, time.Millisecond, -250 * builtintime.Millisecond},
		{1e9, time.Second, 1e18}, // exact: the s→ns table rate would be 128 ns short
		{292, time.Year, 292 * 8766 * builtintime.Hour},
		{-292, time.Year, -292 * 8766 * builtintime.Hour},
	}

	for _, c := range cases {
		tm, err := units.NewTime(c.value, c.unit)
		if err != nil {
			t.Fatal(err)
		}

		if got, err := tm.ToDuration(); err != nil || got != c.want {
			t.Errorf("%v: ToDuration() = %v, %v; want %v", tm, got, err, c.want)
		}
	}
}

func TestTime_ToDuration_OutOfRange(t *testing.T) {
	for _, tm := range []units.Time{
		units.ParseDuration(math.MaxInt64), // 2⁶³-1 ns rounds to 2⁶³ ns in float64
		mustTime(t, 293, time.Year),
		mustTime(t, -293, time.Year),
		mustTime(t, math.NaN(), time.Second),
		mustTime(t, math.Inf(1), time.Second),
		mustTime(t, math.Inf(-1), time.Hour),
	} {
		if d, err := tm.ToDuration(); !errors.Is(err, units.ErrOutOfRange) || d != 0 {
			t.Errorf("%v: ToDuration() = %v, %v; want 0 and ErrOutOfRange", tm, d, err)
		}
	}

	// The lower bound, -2⁶³ ns, is exactly representable and in range.
	if d, err := units.ParseDuration(math.MinInt64).ToDuration(); err != nil || d != math.MinInt64 {
		t.Errorf("ToDuration of the minimum Duration = %v, %v; want %v", d, err, builtintime.Duration(math.MinInt64))
	}
}

// TestTime_ToDurationRoundTrip: within a day, float64 seconds are precise to a
// fraction of a nanosecond, so the round trip is exact.
func TestTime_ToDurationRoundTrip(t *testing.T) {
	rnd := rand.New(rand.NewPCG(3, 4))

	for range 10_000 {
		d := builtintime.Duration(rnd.Int64N(int64(48*builtintime.Hour))) - 24*builtintime.Hour

		if got, err := units.ParseDuration(d).ToDuration(); err != nil || got != d {
			t.Fatalf("ParseDuration(%v).ToDuration() = %v, %v", d, got, err)
		}
	}
}

func mustTime(t *testing.T, value float64, unit string) units.Time {
	t.Helper()

	tm, err := units.NewTime(value, unit)
	if err != nil {
		t.Fatal(err)
	}

	return tm
}
