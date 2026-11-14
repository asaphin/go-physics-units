package units_test

import (
	"errors"
	"math"
	"strconv"
	"testing"

	units "github.com/asaphin/go-physics-units"
	"github.com/asaphin/go-physics-units/area"
	"github.com/asaphin/go-physics-units/conversion"
	"github.com/asaphin/go-physics-units/distance"
	"github.com/asaphin/go-physics-units/mass"
	"github.com/asaphin/go-physics-units/temperature"
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
	units.ErrOutOfRange,
	units.ErrAmbiguousUnit,
	units.ErrInvalidOption,
}

// errOf discards a result and returns the error.
func errOf[T any](_ T, err error) error {
	return err
}

// TestErrors checks every kind of error the package returns: its sentinel
// (exactly one matches, and errors.As finds it) and its message, which must
// stay as detailed as before the sentinels existed.
func TestErrors(t *testing.T) {
	km := units.NewDistance(1, distance.Kilometer)
	kg := units.NewMass(1, mass.Kilogram)
	noTime := units.NewTime(0, time.Second)
	noDistance := units.NewDistance(0, distance.Meter)
	m2 := units.NewArea(1, area.SquareMeter)
	usd := units.NewQuantity(1, usdUnit)
	gbp := units.NewQuantity(1, units.MustUnitOf[otherCurrency]("GBP"))

	cases := []struct {
		name string
		err  error
		kind units.Error
		msg  string
	}{
		{"StringToDistance", errOf(units.StringToDistance("1 kg")), units.ErrUnknownUnit, "unknown distance unit kg"},
		{
			"ConvertToMeasurement", errOf(km.ConvertToMeasurement(mass.Kilogram.String())),
			units.ErrUnknownUnit, "unknown distance unit kg",
		},
		{"NewMeasurement", errOf(units.NewMeasurement(1, "furlong")), units.ErrUnknownUnit, "unknown unit furlong"},
		{"ParseString", errOf(units.ParseString("1 furlong")), units.ErrUnknownUnit, "unknown unit furlong"},
		{"UnitOf", errOf(units.UnitOf[currency]("GBP")), units.ErrUnknownUnit, "unknown currency unit GBP"},
		{
			"Quantity.ConvertToMeasurement", errOf(usd.ConvertToMeasurement("GBP")),
			units.ErrUnknownUnit, "unknown currency unit GBP",
		},
		{
			"AppendText of an invalid unit", errOf(units.NewDistance(1, distance.Unit(200)).MarshalText()),
			units.ErrUnknownUnit, "unknown distance unit distance.Unit(200)",
		},

		{
			"StringToTemperature", errOf(units.StringToTemperature("1 C")),
			units.ErrUnknownUnit, "unknown temperature unit C",
		},

		{"ToDistance", errOf(units.ToDistance(kg)), units.ErrIncompatibleMeasures, "not a distance measure"},
		{
			"ParseString of a temperature", errOf(units.ParseString("20 °C")), units.ErrIncompatibleMeasures,
			"unit °C is a temperature unit, and temperatures aren't Measurements: use StringToTemperature",
		},
		{"ToArea", errOf(units.ToArea(nil)), units.ErrIncompatibleMeasures, "not an area measure"},
		{"ToQuantity", errOf(units.ToQuantity[currency](km)), units.ErrIncompatibleMeasures, "not a currency measure"},
		{
			"AddMeasurement",
			errOf(km.AddMeasurement(kg)),
			units.ErrIncompatibleMeasures,
			"incompatible types distance and mass",
		},
		{"SubMeasurement nil", errOf(km.SubMeasurement(nil)), units.ErrIncompatibleMeasures, "measurement is nil"},
		{
			"AddMeasurement of a measure without units",
			errOf(km.AddMeasurement(units.NewQuantity(1, units.Unit[noUnits]{}))),
			units.ErrIncompatibleMeasures,
			"the measurement's measure has no units: its Units method returns the zero UnitSet",
		},
		{
			"CompareMeasurement", errOf(km.CompareMeasurement(kg)),
			units.ErrIncompatibleMeasures, "incompatible types distance and mass",
		},
		{
			"AddMeasurement of measures with the same name", errOf(usd.AddMeasurement(gbp)),
			units.ErrIncompatibleMeasures, "incompatible currency measures: they have different unit sets",
		},

		{
			"ToDuration", errOf(units.NewTime(293, time.Year).ToDuration()),
			units.ErrOutOfRange, "time 293 yr can't be represented as a time.Duration",
		},
		{
			"below absolute zero", errOf(units.NewTemperature(-300, temperature.Celsius)),
			units.ErrOutOfRange, "temperature -300 °C is below absolute zero, -273.15 °C",
		},
		{
			"MarshalText of NaN", errOf(units.NewDistance(math.NaN(), distance.Meter).MarshalText()),
			units.ErrOutOfRange, "value NaN m can't be encoded as text: it isn't finite",
		},

		{
			"ambiguous unit", errOf(mustParser(t, units.CaseInsensitive()).Parse("1 KN")),
			units.ErrAmbiguousUnit, "ambiguous unit KN: it could be any of kN, kn",
		},
		{
			"invalid alias", errOf(units.NewParser(units.WithAliases(map[string]string{"fur": "furlong"}))),
			units.ErrInvalidOption, "alias fur: furlong isn't a unit symbol",
		},
		{
			"measure without units", errOf(units.NewParser(units.WithMeasure[noUnits]())), units.ErrInvalidOption,
			"WithMeasure: the Units method of the measure returns the zero UnitSet: build it with NewUnitSet",
		},
		{
			"Parser no space",
			errOf(mustParser(t, units.AllowNoSpace()).Parse("km")),
			units.ErrInvalidFormat,
			"invalid string format: km",
		},

		{"Div", errOf(km.Div(0)), units.ErrZeroDivision, "division by zero"},
		{"DivMeasurement", errOf(km.DivMeasurement(0)), units.ErrZeroDivision, "division by zero"},
		{"DivideByTime", errOf(km.DivideByTime(noTime)), units.ErrZeroDivision, "division by zero"},
		{"DivideByDistance", errOf(m2.DivideByDistance(noDistance)), units.ErrZeroDivision, "division by zero"},
		{"Quantity.Div", errOf(usd.Div(0)), units.ErrZeroDivision, "division by zero"},

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
		{
			"ParseQuantity",
			errOf(units.ParseQuantity[currency]("USD")),
			units.ErrInvalidFormat,
			"invalid string format: USD",
		},

		{
			"empty factors", errOf(units.NewUnitSet("currency", "A", conversion.Factors{})),
			units.ErrInvalidFactors, "provided factors don't contain the base unit A",
		},
		{
			"negative factor", errOf(units.NewUnitSet("currency", "A", conversion.Factors{"A": -1})),
			units.ErrInvalidFactors, "conversion factor of unit A must be positive and finite, got -1",
		},
		{
			"empty symbol",
			errOf(units.NewUnitSet("currency", "USD", conversion.Factors{"USD": 1, "": 1})),
			units.ErrInvalidFactors,
			"unit symbol must not be empty",
		},
		{
			"base unit",
			errOf(units.NewUnitSet("currency", "USD", conversion.Factors{"EUR": 1})),
			units.ErrInvalidFactors,
			"provided factors don't contain the base unit USD",
		},

		{
			"empty measure name", errOf(units.NewUnitSet("", "USD", currencyFactors())),
			units.ErrInvalidMeasure, "measure type must not be empty",
		},
		{
			"predefined measure name", errOf(units.NewUnitSet("mass", "USD", currencyFactors())),
			units.ErrInvalidMeasure, "measure type mass is reserved for a predefined measure",
		},
		{
			"measure without units", errOf(units.UnitOf[noUnits]("x")), units.ErrInvalidMeasure,
			"the Units method of the measure returns the zero UnitSet: build it with NewUnitSet",
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
		{"custom measure", errOf(units.UnitOf[currency]("GBP")), "currency", "GBP"},
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

func mustParser(t *testing.T, options ...units.ParserOption) units.Parser {
	t.Helper()

	p, err := units.NewParser(options...)
	if err != nil {
		t.Fatal(err)
	}

	return p
}
