package units_test

import (
	"math"
	"testing"

	units "github.com/asaphin/go-physics-units"
	"github.com/asaphin/go-physics-units/conversion"
	"github.com/asaphin/go-physics-units/distance"
)

type Currency interface {
	units.Measurement
}

func currencyFactors() conversion.Factors {
	return conversion.Factors{
		"EUR": 1.1,
		"USD": 1,
		"UAH": 0.027,
	}
}

func newCurrencyCreator(t *testing.T) units.MeasurementCreator[Currency] {
	t.Helper()

	mc, err := units.NewMeasurementCreator[Currency]("currency", "USD", currencyFactors())
	if err != nil {
		t.Fatal(err)
	}

	return mc
}

func newCurrency(t *testing.T, mc units.MeasurementCreator[Currency], value float64, unit string) Currency {
	t.Helper()

	c, err := mc.New(value, unit)
	if err != nil {
		t.Fatal(err)
	}

	return c
}

func TestMeasurementCreator(t *testing.T) {
	mc := newCurrencyCreator(t)
	usd := newCurrency(t, mc, 1, "USD")

	uah, err := usd.ConvertToMeasurement("UAH")
	if err != nil {
		t.Fatal(err)
	}

	assertMeasurement(t, uah, 1/0.027, "UAH")

	if usd.Type() != "currency" {
		t.Errorf("Type() = %s, want currency", usd.Type())
	}
}

func TestMeasurementCreator_Arithmetic(t *testing.T) {
	mc := newCurrencyCreator(t)

	sum, err := newCurrency(t, mc, 1, "USD").AddMeasurement(newCurrency(t, mc, 1, "EUR"))
	if err != nil {
		t.Fatal(err)
	}

	assertMeasurement(t, sum, 2.1, "USD")
}

func TestMeasurementCreator_IncompatibleMeasures(t *testing.T) {
	type Weird interface{ units.Measurement }

	usd := newCurrency(t, newCurrencyCreator(t), 1, "USD")

	wc, err := units.NewMeasurementCreator[Weird]("weird", "foo", conversion.Factors{"foo": 1})
	if err != nil {
		t.Fatal(err)
	}

	foo, err := wc.New(5, "foo")
	if err != nil {
		t.Fatal(err)
	}

	if sum, err := usd.AddMeasurement(foo); err == nil || sum != nil {
		t.Errorf("1 USD + 5 foo = %v, %v; want nil and an error", sum, err)
	}

	// Two measures with the same name but disjoint units.
	other, err := units.NewMeasurementCreator[Currency]("currency", "GBP", conversion.Factors{"GBP": 1.3, "JPY": 0.007})
	if err != nil {
		t.Fatal(err)
	}

	gbp, err := other.New(1, "GBP")
	if err != nil {
		t.Fatal(err)
	}

	if sum, err := usd.AddMeasurement(gbp); err == nil || sum != nil {
		t.Errorf("1 USD + 1 GBP from unrelated creators = %v, %v; want nil and an error", sum, err)
	}
}

func TestMeasurementCreator_UnknownUnitReturnsNil(t *testing.T) {
	mc := newCurrencyCreator(t)

	c, err := mc.New(1, "GBP")
	if err == nil || c != nil {
		t.Errorf("New(1, GBP) = %v, %v; want nil and an error", c, err)
	}

	usd := newCurrency(t, mc, 1, "USD")
	if conv, err := usd.ConvertToMeasurement("GBP"); err == nil || conv != nil {
		t.Errorf("ConvertToMeasurement(GBP) = %v, %v; want nil and an error", conv, err)
	}
}

func TestMeasurementCreator_RejectsPredefinedMeasureInterface(t *testing.T) {
	mc, err := units.NewMeasurementCreator[units.Distance]("length", "m", distance.ConversionFactors())
	if err == nil || mc != nil {
		t.Errorf("NewMeasurementCreator[units.Distance] = %v, %v; want nil and an error", mc, err)
	}
}

func TestMeasurementCreator_AcceptsMeasurement(t *testing.T) {
	mc, err := units.NewMeasurementCreator[units.Measurement]("currency", "USD", currencyFactors())
	if err != nil {
		t.Fatal(err)
	}

	m, err := mc.New(2, "EUR")
	if err != nil {
		t.Fatal(err)
	}

	assertMeasurement(t, m, 2, "EUR")
}

func TestMeasurementCreator_InvalidArguments(t *testing.T) {
	cases := []struct {
		name        string
		measureType string
		baseUnit    string
		factors     conversion.Factors
	}{
		{"empty measure type", "", "USD", currencyFactors()},
		{"predefined measure type", string(units.MeasureDistance), "USD", currencyFactors()},
		{"custom measure type", string(units.MeasureCustom), "USD", currencyFactors()},
		{"missing base unit", "currency", "GBP", currencyFactors()},
		{"no factors", "currency", "USD", conversion.Factors{}},
		{"empty unit symbol", "currency", "USD", conversion.Factors{"USD": 1, "": 2}},
		{"zero factor", "currency", "USD", conversion.Factors{"USD": 1, "BTC": 0}},
		{"negative factor", "currency", "USD", conversion.Factors{"USD": 1, "BTC": -1}},
		{"NaN factor", "currency", "USD", conversion.Factors{"USD": 1, "BTC": math.NaN()}},
		{"infinite factor", "currency", "USD", conversion.Factors{"USD": 1, "BTC": math.Inf(1)}},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			mc, err := units.NewMeasurementCreator[Currency](c.measureType, c.baseUnit, c.factors)
			if err == nil || mc != nil {
				t.Errorf("NewMeasurementCreator = %v, %v; want nil and an error", mc, err)
			}
		})
	}
}

func TestMeasurementCreator_SymbolOfPredefinedUnit(t *testing.T) {
	type Money interface{ units.Measurement }

	mc, err := units.NewMeasurementCreator[Money]("money", "m", conversion.Factors{"m": 1, "k": 1000})
	if err != nil {
		t.Fatal(err)
	}

	money, err := mc.New(1, "m")
	if err != nil {
		t.Fatal(err)
	}

	if money.Type() != "money" {
		t.Errorf("Type() = %s, want money", money.Type())
	}

	if d, err := units.ToDistance(money); err == nil {
		t.Errorf("ToDistance(1 m of money) = %v, want an error", d)
	}
}

func TestMeasurementCreator_CopiesFactors(t *testing.T) {
	factors := currencyFactors()

	mc, err := units.NewMeasurementCreator[Currency]("currency", "USD", factors)
	if err != nil {
		t.Fatal(err)
	}

	factors["EUR"] = 1000

	eur, err := newCurrency(t, mc, 1, "EUR").ConvertToMeasurement("USD")
	if err != nil {
		t.Fatal(err)
	}

	assertMeasurement(t, eur, 1.1, "USD")
}

// A custom measure can extend the units of a predefined one without changing it.
func TestMeasurementCreator_ExtendsPredefinedFactors(t *testing.T) {
	type Surveying interface{ units.Measurement }

	factors := distance.ConversionFactors()
	factors["fur"] = 201.168 // furlong, 660 international feet

	mc, err := units.NewMeasurementCreator[Surveying]("surveying", distance.BaseUnit, factors)
	if err != nil {
		t.Fatal(err)
	}

	furlong, err := mc.New(1, "fur")
	if err != nil {
		t.Fatal(err)
	}

	m, err := furlong.ConvertToMeasurement(distance.Meter)
	if err != nil {
		t.Fatal(err)
	}

	assertMeasurement(t, m, 201.168, distance.Meter)

	if d, err := units.NewDistance(1, "fur"); err == nil || d != nil {
		t.Errorf("NewDistance(1, fur) = %v, %v; the custom unit leaked into the distance measure", d, err)
	}
}

func TestMeasurementCreator_ParseString(t *testing.T) {
	mc := newCurrencyCreator(t)

	c, err := mc.ParseString("2.5 EUR")
	if err != nil {
		t.Fatal(err)
	}

	assertMeasurement(t, c, 2.5, "EUR")

	for _, s := range []string{"2.5 GBP", "EUR", "x EUR", "NaN EUR", "1 2 EUR"} {
		if c, err := mc.ParseString(s); err == nil || c != nil {
			t.Errorf("ParseString(%q) = %v, %v; want nil and an error", s, c, err)
		}
	}
}

func TestNewBaseMeasurement(t *testing.T) {
	m, err := units.NewBaseMeasurement(1, "EUR", currencyFactors())
	if err != nil {
		t.Fatal(err)
	}

	if m.Type() != units.MeasureCustom {
		t.Errorf("Type() = %s, want %s", m.Type(), units.MeasureCustom)
	}

	usd, err := m.ConvertToMeasurement("USD")
	if err != nil {
		t.Fatal(err)
	}

	assertMeasurement(t, usd, 1.1, "USD")

	if m, err := units.NewBaseMeasurement(1, "GBP", currencyFactors()); err == nil || m != nil {
		t.Errorf("unknown unit: got %v, %v; want nil and an error", m, err)
	}

	if m, err := units.NewBaseMeasurement(1, "USD", conversion.Factors{"USD": 1, "BTC": 0}); err == nil || m != nil {
		t.Errorf("zero factor: got %v, %v; want nil and an error", m, err)
	}

	// Custom factors never produce a predefined measure, even with its symbols.
	custom, err := units.NewBaseMeasurement(1, "m", distance.ConversionFactors())
	if err != nil {
		t.Fatal(err)
	}

	if _, err := units.ToDistance(custom); err == nil {
		t.Error("ToDistance accepted a custom measurement")
	}
}
