package rates_test

import (
	"math"
	"strings"
	"testing"

	"github.com/asaphin/go-physics-units/internal/rates"
)

func TestValidate(t *testing.T) {
	if err := rates.Validate([]string{"m", "km", "\u03bcm"}, []float64{1, 1000, 1e-6}); err != nil {
		t.Errorf("valid units: %v", err)
	}

	cases := map[string]struct {
		symbols []string
		factors []float64
		want    string
	}{
		"empty":           {nil, nil, "empty"},
		"length mismatch": {[]string{"a"}, []float64{1, 2}, "1 symbols but 2 factors"},
		"empty symbol":    {[]string{"a", ""}, []float64{1, 1}, "must not be empty"},
		"zero":            {[]string{"a", "b"}, []float64{1, 0}, "unit b must be positive"},
		"negative":        {[]string{"a", "b"}, []float64{1, -1}, "positive"},
		"NaN":             {[]string{"a", "b"}, []float64{1, math.NaN()}, "positive"},
		"+Inf":            {[]string{"a", "b"}, []float64{1, math.Inf(1)}, "finite"},
		"duplicate":       {[]string{"a", "a"}, []float64{1, 2}, "unit symbol a is defined twice"},
		"micro signs":     {[]string{"\u03bcm", "\u00b5m"}, []float64{1, 2}, "differ only in the micro sign"},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			err := rates.Validate(c.symbols, c.factors)
			if err == nil || !strings.Contains(err.Error(), c.want) {
				t.Errorf("Validate = %v, want an error containing %q", err, c.want)
			}
		})
	}
}

func TestMustValidate_Panics(t *testing.T) {
	defer func() {
		r := recover()
		if s, ok := r.(string); !ok || !strings.Contains(s, "invalid distance units") {
			t.Errorf("panic value %v, want a message naming the measure", r)
		}
	}()

	rates.MustValidate("distance", []string{"m"}, []float64{0})
}

func TestMatrix(t *testing.T) {
	m := rates.Matrix([]float64{1, 1000, 0.001})

	want := []float64{
		1, 0.001, 1000,
		1000, 1, 1e6,
		0.001, 1e-6, 1,
	}

	if len(m) != len(want) {
		t.Fatalf("len = %d, want %d", len(m), len(want))
	}

	for i := range want {
		if math.Abs(m[i]-want[i]) > 1e-12*want[i] {
			t.Errorf("m[%d] = %v, want %v", i, m[i], want[i])
		}
	}

	// The diagonal is exactly 1, so converting to the same unit is exact.
	for i := range 3 {
		if m[i*3+i] != 1 {
			t.Errorf("rate %d→%d = %v, want exactly 1", i, i, m[i*3+i])
		}
	}
}

func TestCanonicalUnit(t *testing.T) {
	cases := map[string]string{
		"\u00b5m":      "\u03bcm",
		"\u03bcm":      "\u03bcm",
		"km":           "km",
		"":             "",
		"k\u2126":      "k\u03a9", // OHM SIGN
		"\u212b":       "\u00c5",  // ANGSTROM SIGN
		"\u212a":       "K",       // KELVIN SIGN
		"\u00b5\u2126": "\u03bc\u03a9",
	}

	for in, want := range cases {
		if got := rates.CanonicalUnit(in); got != want {
			t.Errorf("CanonicalUnit(%q) = %q, want %q", in, got, want)
		}
	}

	if n := testing.AllocsPerRun(100, func() { _ = rates.CanonicalUnit("km") }); n != 0 {
		t.Errorf("CanonicalUnit allocates %v times without a look-alike character", n)
	}
}
