package units_test

import (
	"encoding"
	"fmt"
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

		if mt := units.DetectMeasureType(m.Symbol()); mt != m.Type() {
			t.Fatalf("ParseString(%q): Type() = %s, but unit %s belongs to %s", s, m.Type(), m.Symbol(), mt)
		}

		//nolint:gocritic // deliberately formats through fmt, to compare Format with String
		if v := fmt.Sprintf("%v", m); v != m.String() {
			t.Fatalf("%%v of %q is %q, want String() %q", s, v, m.String())
		}

		again, err := units.ParseString(m.String())
		if err != nil {
			t.Fatalf("re-parsing %q (from %q): %v", m.String(), s, err)
		}

		if again != m {
			t.Fatalf("String() doesn't round-trip: %q -> %v -> %v", s, m, again)
		}

		marshaler, ok := m.(encoding.TextMarshaler)
		if !ok {
			t.Fatalf("%T isn't an encoding.TextMarshaler", m)
		}

		text, err := marshaler.MarshalText()
		if err != nil || string(text) != m.String() {
			t.Fatalf("MarshalText of %v = %q, %v; want String()", m, text, err)
		}
	})
}

// FuzzParser checks a parser with every option: it never panics, returns nil
// on error, accepts whatever ParseString accepts with the same result, and
// what it parses round-trips through String.
func FuzzParser(f *testing.F) {
	for _, seed := range []string{
		"10 km", "10km", "10KM", "1.5e3eV", "10eV", "-2.5m", "1 KN", "5 metres", "10", "km", "1e999km", "1 2 3",
	} {
		f.Add(seed)
	}

	p, err := units.NewParser(units.CaseInsensitive(), units.AllowNoSpace(),
		units.WithAliases(map[string]string{"metres": "m"}))
	if err != nil {
		f.Fatal(err)
	}

	f.Fuzz(func(t *testing.T, s string) {
		m, err := p.Parse(s)
		if err != nil {
			if m != nil {
				t.Fatalf("Parse(%q) returned %v together with error %v", s, m, err)
			}

			if strict, strictErr := units.ParseString(s); strictErr == nil {
				t.Fatalf("Parse(%q) failed with %v, but ParseString accepts it as %v", s, err, strict)
			}

			return
		}

		if math.IsNaN(m.Value()) || math.IsInf(m.Value(), 0) {
			t.Fatalf("Parse(%q) returned a non-finite value %v", s, m.Value())
		}

		if strict, strictErr := units.ParseString(s); strictErr == nil && strict != m {
			t.Fatalf("Parse(%q) = %v, but ParseString gives %v", s, m, strict)
		}

		again, err := p.Parse(m.String())
		if err != nil || again != m {
			t.Fatalf("String() doesn't round-trip: %q -> %v -> %v, %v", s, m, again, err)
		}
	})
}
