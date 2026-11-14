package rates_test

import (
	"math"
	"slices"
	"testing"

	"github.com/asaphin/go-physics-units/conversion"
	"github.com/asaphin/go-physics-units/internal/rates"
)

func TestNewTable_PairsDontCollide(t *testing.T) {
	// "m"+"mm" and "mm"+"m" concatenate to the same string.
	table, err := rates.NewTable("distance", conversion.Factors{"m": 1, "mm": 0.001, "mmm": 1e-6})
	if err != nil {
		t.Fatal(err)
	}

	cases := []struct {
		from, to string
		want     float64
	}{
		{"m", "mm", 1000},
		{"mm", "m", 0.001},
		{"m", "mmm", 1e6},
		{"mmm", "m", 1e-6},
		{"mm", "mmm", 1000},
		{"m", "m", 1},
	}

	for _, c := range cases {
		got, ok := table.Rate(c.from, c.to)
		if !ok || math.Abs(got-c.want) > 1e-9*c.want {
			t.Errorf("Rate(%s, %s) = %v, %v; want %v, true", c.from, c.to, got, ok, c.want)
		}
	}

	if _, ok := table.Rate("m", "km"); ok {
		t.Error("Rate to an unknown unit reported ok")
	}
}

func TestNewTable_Measure(t *testing.T) {
	table := rates.MustNewTable("currency", conversion.Factors{"USD": 1})

	if table.Measure() != "currency" {
		t.Errorf("Measure() = %s, want currency", table.Measure())
	}

	if !table.Has("USD") || table.Has("EUR") {
		t.Error("Has reports the wrong units")
	}
}

func TestNewTable_RejectsInvalidFactors(t *testing.T) {
	cases := map[string]conversion.Factors{
		"nil":              nil,
		"empty":            {},
		"empty symbol":     {"a": 1, "": 1},
		"zero":             {"a": 1, "b": 0},
		"negative":         {"a": 1, "b": -2},
		"NaN":              {"a": 1, "b": math.NaN()},
		"+Inf":             {"a": 1, "b": math.Inf(1)},
		"micro sign clash": {"\u00b5g": 1e-9, "\u03bcg": 1e-9},
	}

	for name, f := range cases {
		if table, err := rates.NewTable("test", f); err == nil || table != (rates.Table{}) {
			t.Errorf("%s: NewTable = %v, %v; want the zero Table and an error", name, table, err)
		}
	}
}

func TestNewTable_DoesNotKeepFactors(t *testing.T) {
	f := conversion.Factors{"a": 1, "b": 2}

	table, err := rates.NewTable("test", f)
	if err != nil {
		t.Fatal(err)
	}

	f["b"] = 4
	f["c"] = 8

	delete(f, "a")

	if got, ok := table.Rate("b", "a"); !ok || got != 2 {
		t.Errorf("Rate(b, a) = %v, %v after changing the factors; want 2, true", got, ok)
	}

	if table.Has("c") {
		t.Error("a unit added to the factors after NewTable appeared in the table")
	}

	if got, ok := table.Factor("b"); !ok || got != 2 {
		t.Errorf("Factor(b) = %v, %v after changing the factors; want 2, true", got, ok)
	}
}

func TestTable_Factor(t *testing.T) {
	table, err := rates.NewTable("distance", conversion.Factors{"m": 1, "km": 1000, "\u00b5m": 1e-6})
	if err != nil {
		t.Fatal(err)
	}

	for unit, want := range map[string]float64{"m": 1, "km": 1000, "\u03bcm": 1e-6} {
		if got, ok := table.Factor(unit); !ok || got != want {
			t.Errorf("Factor(%s) = %v, %v; want %v, true", unit, got, ok, want)
		}
	}

	if got := slices.Sorted(table.Units()); !slices.Equal(got, []string{"km", "m", "\u03bcm"}) {
		t.Errorf("Units() = %q, want km, m and \u03bcm", got)
	}

	// Factor, like Rate, takes canonical symbols only.
	for _, unit := range []string{"ft", "\u00b5m", ""} {
		if got, ok := table.Factor(unit); ok || got != 0 {
			t.Errorf("Factor(%q) = %v, %v; want 0, false", unit, got, ok)
		}
	}
}

func TestFactorRate(t *testing.T) {
	f := conversion.Factors{"m": 1, "km": 1000, "in": 0.0254}

	table, err := rates.NewTable("distance", f)
	if err != nil {
		t.Fatal(err)
	}

	for from := range f {
		for to := range f {
			got, ok := rates.FactorRate(f, from, to)
			want, _ := table.Rate(from, to)

			if !ok || got != want {
				t.Errorf("FactorRate(%s, %s) = %v, %v; want %v, true, as the table", from, to, got, ok, want)
			}
		}
	}

	for _, pair := range [][2]string{{"m", "ft"}, {"ft", "m"}, {"ft", "ft"}} {
		if got, ok := rates.FactorRate(f, pair[0], pair[1]); ok || got != 0 {
			t.Errorf("FactorRate(%s, %s) = %v, %v; want 0, false", pair[0], pair[1], got, ok)
		}
	}
}

func TestTable_Zero(t *testing.T) {
	var zero rates.Table

	if zero.Measure() != "" || zero.Has("a") {
		t.Errorf("zero Table: Measure() = %q, Has(a) = %v; want empty and false", zero.Measure(), zero.Has("a"))
	}

	if r, ok := zero.Rate("a", "a"); ok || r != 0 {
		t.Errorf("zero Table: Rate(a, a) = %v, %v; want 0, false", r, ok)
	}

	if f, ok := zero.Factor("a"); ok || f != 0 {
		t.Errorf("zero Table: Factor(a) = %v, %v; want 0, false", f, ok)
	}

	for u := range zero.Units() {
		t.Errorf("zero Table: Units() yields %q", u)
	}
}

func TestMustNewTable_PanicsOnInvalidFactors(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Error("MustNewTable didn't panic")
		}
	}()

	rates.MustNewTable("test", conversion.Factors{"a": 0})
}

func TestNewTable_StoresCanonicalSymbols(t *testing.T) {
	table, err := rates.NewTable("mass", conversion.Factors{"\u00b5g": 1e-9, "g": 1e-3})
	if err != nil {
		t.Fatal(err)
	}

	if !table.Has("\u03bcg") {
		t.Error("the micro sign wasn't canonicalised to U+03BC")
	}

	if got, ok := table.Rate("g", "\u03bcg"); !ok || math.Abs(got-1e6) > 1e-3 {
		t.Errorf("Rate(g, \u03bcg) = %v, %v; want 1e6, true", got, ok)
	}
}

func TestCanonicalUnit(t *testing.T) {
	cases := map[string]string{
		"\u00b5m": "\u03bcm",
		"\u03bcm": "\u03bcm",
		"m":       "m",
		"":        "",
	}

	for in, want := range cases {
		if got := rates.CanonicalUnit(in); got != want {
			t.Errorf("CanonicalUnit(%q) = %q, want %q", in, got, want)
		}
	}
}
