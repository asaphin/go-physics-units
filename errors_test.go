package units_test

import (
	"errors"
	"strconv"
	"testing"

	units "github.com/asaphin/go-physics-units"
	"github.com/asaphin/go-physics-units/area"
	"github.com/asaphin/go-physics-units/conversion"
	"github.com/asaphin/go-physics-units/distance"
	"github.com/asaphin/go-physics-units/mass"
	"github.com/asaphin/go-physics-units/time"
)

// sentinels lists every sentinel error. A new sentinel needs an entry here.
var sentinels = []units.Error{
	units.ErrUnknownUnit,
	units.ErrIncompatibleMeasures,
	units.ErrZeroDivision,
	units.ErrInvalidFormat,
	units.ErrInvalidFactors,
	units.ErrInvalidMeasure,
}

// errOf discards a result and returns the error.
func errOf[T any](_ T, err error) error {
	return err
}

// TestErrors checks every kind of error the package returns: its sentinel
// (exactly one matches, and errors.As finds it) and its message, which must
// stay as detailed as before the sentinels existed.
func TestErrors(t *testing.T) {
	km, err := units.NewDistance(1, distance.Kilometer)
	if err != nil {
		t.Fatal(err)
	}

	kg, err := units.NewMass(1, mass.Kilogram)
	if err != nil {
		t.Fatal(err)
	}

	noTime, err := units.NewTime(0, time.Second)
	if err != nil {
		t.Fatal(err)
	}

	noDistance, err := units.NewDistance(0, distance.Meter)
	if err != nil {
		t.Fatal(err)
	}

	m2, err := units.NewArea(1, area.SquareMeter)
	if err != nil {
		t.Fatal(err)
	}

	mc := newCurrencyCreator(t)
	usd := newCurrency(t, mc, 1, "USD")

	// A measure with the same name as mc's, but other units.
	other, err := units.NewMeasurementCreator[Currency]("currency", "GBP", conversion.Factors{"GBP": 1.3})
	if err != nil {
		t.Fatal(err)
	}

	gbp, err := other.New(1, "GBP")
	if err != nil {
		t.Fatal(err)
	}

	cases := []struct {
		name string
		err  error
		kind units.Error
		msg  string
	}{
		{"NewDistance", errOf(units.NewDistance(1, "furlong")), units.ErrUnknownUnit, "unknown distance unit furlong"},
		{"ConvertTo", errOf(km.ConvertTo(mass.Kilogram)), units.ErrUnknownUnit, "unknown distance unit kg"},
		{
			"ConvertToMeasurement", errOf(km.ConvertToMeasurement(mass.Kilogram)),
			units.ErrUnknownUnit, "unknown distance unit kg",
		},
		{"StringToDistance", errOf(units.StringToDistance("1 kg")), units.ErrUnknownUnit, "unknown distance unit kg"},
		{"NewMeasurement", errOf(units.NewMeasurement(1, "furlong")), units.ErrUnknownUnit, "unknown unit furlong"},
		{"ParseString", errOf(units.ParseString("1 furlong")), units.ErrUnknownUnit, "unknown unit furlong"},
		{"creator New", errOf(mc.New(1, "GBP")), units.ErrUnknownUnit, "unknown currency unit GBP"},
		{
			"NewBaseMeasurement", errOf(units.NewBaseMeasurement(1, "B", conversion.Factors{"A": 1})),
			units.ErrUnknownUnit, "unknown custom unit B",
		},

		{"ToDistance", errOf(units.ToDistance(kg)), units.ErrIncompatibleMeasures, "not a distance measure"},
		{"ToArea", errOf(units.ToArea(nil)), units.ErrIncompatibleMeasures, "not an area measure"},
		{
			"AddMeasurement",
			errOf(km.AddMeasurement(kg)),
			units.ErrIncompatibleMeasures,
			"incompatible types distance and mass",
		},
		{"SubMeasurement nil", errOf(km.SubMeasurement(nil)), units.ErrIncompatibleMeasures, "measurement is nil"},
		{
			"CompareMeasurement", errOf(km.CompareMeasurement(kg)),
			units.ErrIncompatibleMeasures, "incompatible types distance and mass",
		},
		{
			"AddMeasurement of measures with the same name", errOf(usd.AddMeasurement(gbp)),
			units.ErrIncompatibleMeasures, "can't convert GBP to USD: the units belong to different currency measures",
		},

		{"Div", errOf(km.Div(0)), units.ErrZeroDivision, "division by zero"},
		{"DivMeasurement", errOf(km.DivMeasurement(0)), units.ErrZeroDivision, "division by zero"},
		{"DivideByTime", errOf(km.DivideByTime(noTime)), units.ErrZeroDivision, "division by zero"},
		{"DivideByDistance", errOf(m2.DivideByDistance(noDistance)), units.ErrZeroDivision, "division by zero"},

		{"one field", errOf(units.ParseString("10km")), units.ErrInvalidFormat, "invalid string format: 10km"},
		{"three fields", errOf(units.StringToMass("1 2 kg")), units.ErrInvalidFormat, "invalid string format: 1 2 kg"},
		{
			"value isn't a number", errOf(units.ParseString("ten km")), units.ErrInvalidFormat,
			`failed to parse value ten: strconv.ParseFloat: parsing "ten": invalid syntax`,
		},
		{
			"value isn't finite",
			errOf(units.ParseString("+Inf km")),
			units.ErrInvalidFormat,
			"value +Inf is not a finite number",
		},
		{"creator ParseString", errOf(mc.ParseString("USD")), units.ErrInvalidFormat, "invalid string format: USD"},

		{
			"empty factors", errOf(units.NewBaseMeasurement(1, "A", conversion.Factors{})),
			units.ErrInvalidFactors, "conversion factors are empty",
		},
		{
			"negative factor", errOf(units.NewBaseMeasurement(1, "A", conversion.Factors{"A": -1})),
			units.ErrInvalidFactors, "conversion factor of unit A must be positive and finite, got -1",
		},
		{
			"creator factors",
			errOf(units.NewMeasurementCreator[Currency]("currency", "USD", conversion.Factors{"": 1})),
			units.ErrInvalidFactors,
			"unit symbol must not be empty",
		},
		{
			"creator base unit",
			errOf(units.NewMeasurementCreator[Currency]("currency", "USD", conversion.Factors{"EUR": 1})),
			units.ErrInvalidFactors,
			"provided factors doesn't contain unit USD",
		},

		{
			"empty measure name", errOf(units.NewMeasurementCreator[Currency]("", "USD", currencyFactors())),
			units.ErrInvalidMeasure, "measure type must not be empty",
		},
		{
			"predefined measure name", errOf(units.NewMeasurementCreator[Currency]("mass", "USD", currencyFactors())),
			units.ErrInvalidMeasure, "measure type mass is reserved for a predefined measure",
		},
		{
			"predefined measure interface",
			errOf(units.NewMeasurementCreator[units.Distance]("currency", "USD", currencyFactors())),
			units.ErrInvalidMeasure,
			"unsupported measurement type units.Distance: it must be satisfied by units.Measurement alone",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if c.err == nil {
				t.Fatalf("no error, want %q", c.msg)
			}

			if got := c.err.Error(); got != c.msg {
				t.Errorf("message %q, want %q", got, c.msg)
			}

			for _, s := range sentinels {
				if is, want := errors.Is(c.err, s), s == c.kind; is != want {
					t.Errorf("errors.Is(err, %q) = %v, want %v", s, is, want)
				}
			}

			var kind units.Error
			if !errors.As(c.err, &kind) || kind != c.kind {
				t.Errorf("errors.As found the kind %q, want %q", kind, c.kind)
			}
		})
	}
}

func TestUnitError(t *testing.T) {
	cases := []struct {
		name    string
		err     error
		measure units.MeasureType
		unit    string
	}{
		{"typed measure", errOf(units.StringToDistance("1 kg")), units.MeasureDistance, "kg"},
		{"no predefined measure", errOf(units.ParseString("1 furlong")), "", "furlong"},
		{"custom measure", errOf(newCurrencyCreator(t).New(1, "GBP")), "currency", "GBP"},
		{
			"NewBaseMeasurement",
			errOf(units.NewBaseMeasurement(1, "B", conversion.Factors{"A": 1})),
			units.MeasureCustom,
			"B",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			var ue *units.UnitError
			if !errors.As(c.err, &ue) {
				t.Fatalf("errors.As(%v) found no *UnitError", c.err)
			}

			if ue.Measure != c.measure || ue.Unit != c.unit {
				t.Errorf("UnitError{%q, %q}, want {%q, %q}", ue.Measure, ue.Unit, c.measure, c.unit)
			}
		})
	}
}

func TestErrors_WrapTheCause(t *testing.T) {
	err := errOf(units.ParseString("1e999 km"))

	var numErr *strconv.NumError
	if !errors.As(err, &numErr) || !errors.Is(err, strconv.ErrRange) {
		t.Errorf("%v doesn't wrap the *strconv.NumError of the value", err)
	}

	if !errors.Is(err, units.ErrInvalidFormat) {
		t.Errorf("%v isn't an ErrInvalidFormat", err)
	}
}

func TestSentinelMessages(t *testing.T) {
	seen := make(map[string]bool, len(sentinels))

	for _, s := range sentinels {
		if s.Error() == "" || s.Error() != string(s) || seen[s.Error()] {
			t.Errorf("sentinel %q: its message must be its own, non-empty and unique", s)
		}

		seen[s.Error()] = true
	}
}
