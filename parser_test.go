package units_test

import (
	"errors"
	"strconv"
	"sync"
	"testing"

	units "github.com/asaphin/go-physics-units"
	"github.com/asaphin/go-physics-units/conversion"
	"github.com/asaphin/go-physics-units/distance"
	"github.com/asaphin/go-physics-units/energy"
	"github.com/asaphin/go-physics-units/force"
	"github.com/asaphin/go-physics-units/velocity"
)

func TestParser_ZeroValueIsParseString(t *testing.T) {
	var zero units.Parser

	for _, s := range []string{"10 km", "1.5e3 m", "3 kg", "10km", "ten km", "1 furlong", "NaN m", ""} {
		got, gotErr := zero.Parse(s)
		want, wantErr := units.ParseString(s)

		if (gotErr == nil) != (wantErr == nil) || (gotErr != nil && gotErr.Error() != wantErr.Error()) {
			t.Errorf("Parse(%q) error = %v, ParseString error = %v", s, gotErr, wantErr)

			continue
		}

		if want != nil && (got.Value() != want.Value() || got.Symbol() != want.Symbol() || got.Type() != want.Type()) {
			t.Errorf("Parse(%q) = %v, ParseString = %v", s, got, want)
		}
	}

	if u, err := zero.Unit("\u00b5m"); err != nil || u != distance.Micrometer.String() {
		t.Errorf("zero Unit(\u00b5m) = %q, %v; want %q", u, err, distance.Micrometer)
	}

	if _, err := zero.Unit("furlong"); !errors.Is(err, units.ErrUnknownUnit) {
		t.Errorf("zero Unit(furlong): %v, want ErrUnknownUnit", err)
	}
}

func TestParser_Aliases(t *testing.T) {
	aliases := map[string]string{
		"metre":  distance.Meter.String(),
		"metres": distance.Meter.String(),
		"kph":    velocity.KilometerPerHour.String(),
	}
	p := mustParser(t, units.WithAliases(aliases))

	// Changing the map afterwards has no effect: WithAliases copies it.
	aliases["kph"] = distance.Kilometer.String()
	aliases["furlong"] = distance.Meter.String()

	m, err := p.Parse("5 metres")
	if err != nil {
		t.Fatal(err)
	}

	assertMeasurement(t, m, 5, distance.Meter.String())

	if m.Type() != units.MeasureDistance {
		t.Errorf("Type() = %s, want distance", m.Type())
	}

	v, err := p.Parse("100 kph")
	if err != nil {
		t.Fatal(err)
	}

	assertMeasurement(t, v, 100, velocity.KilometerPerHour.String())

	if _, parseErr := p.Parse("1 furlong"); !errors.Is(parseErr, units.ErrUnknownUnit) {
		t.Errorf("an alias added after NewParser parsed: %v", parseErr)
	}

	if u, unitErr := p.Unit("metre"); unitErr != nil || u != distance.Meter.String() {
		t.Errorf("Unit(metre) = %q, %v; want m", u, unitErr)
	}

	// Symbols keep working, and an alias denotes the unit itself.
	km, err := p.Parse("2 km")
	if err != nil {
		t.Fatal(err)
	}

	if sum, err := km.AddMeasurement(m); err != nil || sum.Value() != 2.005 {
		t.Errorf("2 km + 5 metres = %v, %v; want 2.005 km", sum, err)
	}
}

// metreMeasure is a user-defined measure with the predefined symbol m.
type metreMeasure struct{}

var metreUnits = units.MustNewUnitSet("metres", "m", conversion.Factors{"m": 1})

func (metreMeasure) Units() units.UnitSet { return metreUnits }

// dollarMeasure shares the symbol USD with currency.
type dollarMeasure struct{}

var dollarUnits = units.MustNewUnitSet("dollars", "USD", conversion.Factors{"USD": 1})

func (dollarMeasure) Units() units.UnitSet { return dollarUnits }

func TestParser_InvalidOptions(t *testing.T) {
	cases := map[string][]units.ParserOption{
		"alias of an unknown unit": {units.WithAliases(map[string]string{"furlong": "fur"})},
		"alias that is a symbol":   {units.WithAliases(map[string]string{"m": distance.Kilometer.String()})},
		"alias of an alias": {
			units.WithAliases(map[string]string{"metre": distance.Meter.String()}),
			units.WithAliases(map[string]string{"meter": "metre"}),
		},
		"alias defined twice": {
			units.WithAliases(map[string]string{"metre": distance.Meter.String()}),
			units.WithAliases(map[string]string{"metre": distance.Meter.String()}),
		},
		"empty alias":                 {units.WithAliases(map[string]string{"": distance.Meter.String()})},
		"alias with a space":          {units.WithAliases(map[string]string{"square metre": distance.Meter.String()})},
		"measure unit is predefined":  {units.WithMeasure[metreMeasure]()},
		"two measures share a symbol": {units.WithMeasure[currency](), units.WithMeasure[dollarMeasure]()},
		"measure added twice":         {units.WithMeasure[currency](), units.WithMeasure[currency]()},
		"measure without units":       {units.WithMeasure[noUnits]()},
	}

	for name, options := range cases {
		t.Run(name, func(t *testing.T) {
			p, err := units.NewParser(options...)
			if !errors.Is(err, units.ErrInvalidOption) || p != (units.Parser{}) {
				t.Errorf("NewParser = %v, %v; want the zero Parser and ErrInvalidOption", p, err)
			}
		})
	}
}

func TestParser_CaseInsensitive(t *testing.T) {
	p := mustParser(
		t,
		units.CaseInsensitive(),
		units.WithAliases(map[string]string{"Metre": distance.Meter.String(), "KM": distance.Kilometer.String()}),
	)

	cases := []struct {
		in    string
		value float64
		unit  string
	}{
		{"10 KM", 10, distance.Kilometer.String()},
		{"10 Km", 10, distance.Kilometer.String()},
		{"1 kn", 1, velocity.Knot.String()},    // exact matches win
		{"1 kN", 1, force.Kilonewton.String()}, // exact matches win
		{"3 METRE", 3, distance.Meter.String()},
		{"2 kwh", 2, energy.KilowattHour.String()},
		{"4 \u00e5", 4, distance.Angstrom.String()}, // non-ASCII: å folds to Å
	}

	for _, c := range cases {
		m, err := p.Parse(c.in)
		if err != nil {
			t.Errorf("Parse(%q): %v", c.in, err)

			continue
		}

		assertMeasurement(t, m, c.value, c.unit)
	}

	_, err := p.Parse("1 KN")
	if !errors.Is(err, units.ErrAmbiguousUnit) || err.Error() != "ambiguous unit KN: it could be any of kN, kn" {
		t.Errorf("Parse(1 KN): %v, want ErrAmbiguousUnit naming kN and kn", err)
	}

	if _, err := p.Unit("KN"); !errors.Is(err, units.ErrAmbiguousUnit) {
		t.Errorf("Unit(KN): %v, want ErrAmbiguousUnit", err)
	}

	// Without the option, case matters.
	if _, err := mustParser(t).Parse("10 KM"); !errors.Is(err, units.ErrUnknownUnit) {
		t.Errorf("a case-sensitive parser accepted KM: %v", err)
	}
}

func TestParser_AllowNoSpace(t *testing.T) {
	p := mustParser(t, units.AllowNoSpace())

	cases := []struct {
		in    string
		value float64
		unit  string
	}{
		{"10km", 10, distance.Kilometer.String()},
		{"  -2.5m ", -2.5, distance.Meter.String()},
		{"1.5e3eV", 1500, energy.Electronvolt.String()},
		{"10eV", 10, energy.Electronvolt.String()}, // "10e" isn't a number
		{"3e2 km", 300, distance.Kilometer.String()},
		{"10 km", 10, distance.Kilometer.String()},
	}

	for _, c := range cases {
		m, err := p.Parse(c.in)
		if err != nil {
			t.Errorf("Parse(%q): %v", c.in, err)

			continue
		}

		assertMeasurement(t, m, c.value, c.unit)
	}

	for in, want := range map[string]error{
		"km":        units.ErrInvalidFormat,
		"10":        units.ErrInvalidFormat,
		"1e999km":   units.ErrInvalidFormat,
		"NaNkm":     units.ErrInvalidFormat,
		"10furlong": units.ErrUnknownUnit,
	} {
		if m, err := p.Parse(in); !errors.Is(err, want) || m != nil {
			t.Errorf("Parse(%q) = %v, %v; want nil and %v", in, m, err, want)
		}
	}

	var numErr *strconv.NumError
	if _, err := p.Parse("1e999km"); !errors.As(err, &numErr) {
		t.Errorf("Parse(1e999km): %v doesn't wrap the *strconv.NumError", err)
	}

	if _, err := mustParser(t).Parse("10km"); !errors.Is(err, units.ErrInvalidFormat) {
		t.Errorf("without AllowNoSpace, 10km parsed: %v", err)
	}
}

func TestParser_WithMeasure(t *testing.T) {
	p := mustParser(t, units.WithMeasure[currency](), units.WithAliases(map[string]string{"dollars": "USD"}))

	m, err := p.Parse("10 dollars")
	if err != nil {
		t.Fatal(err)
	}

	usd, ok := m.(Money)
	if !ok || usd != units.NewQuantity(10, usdUnit) {
		t.Fatalf("Parse(10 dollars) = %v of type %T; want 10 USD of type Money", m, m)
	}

	// Parsed values belong to the measure's unit set.
	eur := units.NewQuantity(10, eurUnit)

	sum, err := eur.AddMeasurement(usd)
	if err != nil {
		t.Fatal(err)
	}

	assertMeasurement(t, sum, 10+10/1.1, "EUR")

	if c, err := usd.CompareMeasurement(eur); err != nil || c != -1 {
		t.Errorf("10 USD vs 10 EUR = %d, %v; want -1", c, err)
	}

	if u, err := p.Unit("dollars"); err != nil || u != "USD" {
		t.Errorf("Unit(dollars) = %q, %v; want USD", u, err)
	}

	// Predefined symbols still parse.
	if d, err := p.Parse("2 km"); err != nil || d != units.NewDistance(2, distance.Kilometer) {
		t.Errorf("Parse(2 km) = %v, %v", d, err)
	}
}

func TestParser_ConcurrentUse(t *testing.T) {
	p := mustParser(
		t,
		units.CaseInsensitive(),
		units.AllowNoSpace(),
		units.WithAliases(map[string]string{"metre": distance.Meter.String()}),
	)

	var wg sync.WaitGroup

	for range 8 {
		wg.Go(func() {
			for range 1000 {
				if _, err := p.Parse("10KM"); err != nil {
					t.Error(err)

					return
				}

				if _, err := p.Parse("3 metre"); err != nil {
					t.Error(err)

					return
				}
			}
		})
	}

	wg.Wait()
}
