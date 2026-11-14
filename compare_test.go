package units_test

import (
	"errors"
	"math"
	"math/rand/v2"
	"slices"
	"testing"

	units "github.com/asaphin/go-physics-units"
	"github.com/asaphin/go-physics-units/distance"
)

func TestCompare_SpecialValues(t *testing.T) {
	nan := units.NewDistance(math.NaN(), distance.Kilometer)
	inf := units.NewDistance(math.Inf(1), distance.Kilometer)
	infM := units.NewDistance(math.Inf(1), distance.Meter)
	negInf := units.NewDistance(math.Inf(-1), distance.Meter)
	zero := units.NewDistance(0, distance.Kilometer)
	negZero := units.NewDistance(math.Copysign(0, -1), distance.Meter)
	one := units.NewDistance(1, distance.Kilometer)

	cases := []struct {
		name    string
		x, y    units.Distance
		compare int
		equal   bool
	}{
		{"NaN is less than a number", nan, one, -1, false},
		{"NaN equals NaN, as in cmp.Compare", nan, nan, 0, true},
		{"-0 equals +0", negZero, zero, 0, true},
		{"+Inf in km equals +Inf in m", inf, infM, 0, true},
		{"-Inf is less than +Inf", negInf, inf, -1, false},
		{"+Inf is greater than a number", inf, one, 1, false},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got, rev := c.x.Compare(c.y), c.y.Compare(c.x); got != c.compare || rev != -c.compare {
				t.Errorf("Compare = %d, reversed %d; want %d, %d", got, rev, c.compare, -c.compare)
			}

			if got := c.x.Equal(c.y); got != c.equal {
				t.Errorf("Equal = %v, want %v", got, c.equal)
			}

			// A tolerance never makes infinities or NaN equal to anything else.
			if got := c.x.EqualWithinRel(c.y, 1); got != c.equal {
				t.Errorf("EqualWithinRel(1) = %v, want %v", got, c.equal)
			}
		})
	}
}

func TestEqualWithinRel_Tolerance(t *testing.T) {
	x := units.NewDistance(1000, distance.Meter)
	near := units.NewDistance(1.000001, distance.Kilometer) // 1 mm more

	for _, c := range []struct {
		tol  float64
		want bool
	}{
		{0, false},
		{1e-7, false},
		{1e-6, true}, // |a-b| = 1e-3 m <= 1e-6 · 1000.001 m
		{-1, false},
		{math.NaN(), false},
		{math.Inf(1), true},
	} {
		if got, rev := x.EqualWithinRel(near, c.tol), near.EqualWithinRel(x, c.tol); got != c.want || rev != c.want {
			t.Errorf("EqualWithinRel(tol %v) = %v, reversed %v; want %v", c.tol, got, rev, c.want)
		}
	}
}

// TestCompare_RoundingIsSymmetric: 12 in and 1 ft differ in the last bit once
// converted to metres. Whatever the outcome, it must not depend on the order of
// the operands.
func TestCompare_RoundingIsSymmetric(t *testing.T) {
	foot := units.NewDistance(1, distance.Feet)
	inches := units.NewDistance(12, distance.Inch)

	if foot.Equal(inches) || inches.Equal(foot) {
		t.Errorf("1 ft and 12 in are equal, but they differ in the last bit in metres; update the test")
	}

	if foot.Compare(inches) != -inches.Compare(foot) {
		t.Errorf("Compare isn't antisymmetric: %d and %d", foot.Compare(inches), inches.Compare(foot))
	}

	if !foot.EqualWithinRel(inches, 1e-15) || !inches.EqualWithinRel(foot, 1e-15) {
		t.Error("1 ft and 12 in aren't equal within 1e-15")
	}
}

func TestCompare_Symmetric(t *testing.T) {
	all := slices.Collect(distance.Units())
	rnd := rand.New(rand.NewPCG(1, 2))

	for range 10_000 {
		// Equal values in two units, nudged apart by a few ulps half the time.
		v := rnd.Float64() * 1e6
		a := units.NewDistance(v, all[rnd.IntN(len(all))])
		b := a.ConvertTo(all[rnd.IntN(len(all))])

		if rnd.IntN(2) == 0 {
			b = units.NewDistance(math.Nextafter(b.Value(), math.Inf(1)), b.Unit())
		}

		if a.Compare(b) != -b.Compare(a) || a.Equal(b) != b.Equal(a) {
			t.Fatalf("%v and %v: Compare %d and %d, Equal %v and %v",
				a, b, a.Compare(b), b.Compare(a), a.Equal(b), b.Equal(a))
		}
	}
}

func TestCompareMeasurement(t *testing.T) {
	km, err := units.ParseString("1 km")
	if err != nil {
		t.Fatal(err)
	}

	m := units.NewDistance(999, distance.Meter)

	if got, cmpErr := km.CompareMeasurement(m); cmpErr != nil || got != 1 {
		t.Errorf("1 km vs 999 m = %d, %v; want 1", got, cmpErr)
	}

	if got, cmpErr := m.CompareMeasurement(km); cmpErr != nil || got != -1 {
		t.Errorf("999 m vs 1 km = %d, %v; want -1", got, cmpErr)
	}

	usd, eur := units.NewQuantity(1, usdUnit), units.NewQuantity(1, eurUnit)

	if got, cmpErr := usd.CompareMeasurement(eur); cmpErr != nil || got != -1 {
		t.Errorf("1 USD vs 1 EUR = %d, %v; want -1", got, cmpErr)
	}

	// A measure with the same name but its own unit set.
	gbp := units.NewQuantity(1, units.MustUnitOf[otherCurrency]("GBP"))
	if got, cmpErr := usd.CompareMeasurement(gbp); !errors.Is(cmpErr, units.ErrIncompatibleMeasures) || got != 0 {
		t.Errorf("1 USD vs 1 GBP of another set = %d, %v; want 0 and ErrIncompatibleMeasures", got, cmpErr)
	}

	if got, cmpErr := usd.CompareMeasurement(nil); !errors.Is(cmpErr, units.ErrIncompatibleMeasures) || got != 0 {
		t.Errorf("CompareMeasurement(nil) = %d, %v; want 0 and ErrIncompatibleMeasures", got, cmpErr)
	}
}
