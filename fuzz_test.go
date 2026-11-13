package units_test

import (
	"math"
	"testing"

	units "github.com/asaphin/go-physics-units"
)

func FuzzParseString(f *testing.F) {
	for _, seed := range []string{
		"10 km", "1.5e3 m", "-3 kg", "0 s", "1 \u00b5m", "5 m²", "2 km/h", "0x1p-2 ly",
		"NaN m", "+Inf m", "10km", "", "1 2 3", "1e400 m",
	} {
		f.Add(seed)
	}

	f.Fuzz(func(t *testing.T, s string) {
		m, err := units.ParseString(s)
		if err != nil {
			if m != nil {
				t.Fatalf("ParseString(%q) returned %v together with error %v", s, m, err)
			}

			return
		}

		if math.IsNaN(m.Value()) || math.IsInf(m.Value(), 0) {
			t.Fatalf("ParseString(%q) returned a non-finite value %v", s, m.Value())
		}

		if mt := units.DetectMeasureType(m.Unit()); mt != m.Type() {
			t.Fatalf("ParseString(%q): Type() = %s, but unit %s belongs to %s", s, m.Type(), m.Unit(), mt)
		}

		again, err := units.ParseString(m.String())
		if err != nil {
			t.Fatalf("re-parsing %q (from %q): %v", m.String(), s, err)
		}

		if again.Value() != m.Value() || again.Unit() != m.Unit() {
			t.Fatalf("String() doesn't round-trip: %q -> %v -> %v", s, m, again)
		}
	})
}
