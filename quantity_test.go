package units_test

import (
	"errors"
	"fmt"
	"math"
	"slices"
	"strconv"
	"testing"

	units "github.com/asaphin/go-physics-units"
	"github.com/asaphin/go-physics-units/conversion"
	"github.com/asaphin/go-physics-units/distance"
)

// currency is the user-defined measure of the tests: a marker type and its
// unit set.
type currency struct{}

func currencyFactors() conversion.Factors {
	return conversion.Factors{
		"EUR": 1.1,
		"USD": 1,
		"UAH": 0.027,
	}
}

var currencyUnits = units.MustNewUnitSet("currency", "USD", currencyFactors())

func (currency) Units() units.UnitSet { return currencyUnits }

// Money is the typed currency value.
type Money = units.Quantity[currency]

var (
	usdUnit = units.MustUnitOf[currency]("USD")
	eurUnit = units.MustUnitOf[currency]("EUR")
	uahUnit = units.MustUnitOf[currency]("UAH")
)

// noUnits is a measure whose Units method returns the zero UnitSet.
type noUnits struct{}

func (noUnits) Units() units.UnitSet { return units.UnitSet{} }

// otherCurrency has the name of currency but its own units.
type otherCurrency struct{}

var otherCurrencyUnits = units.MustNewUnitSet("currency", "GBP", conversion.Factors{"GBP": 1.3, "JPY": 0.007})

func (otherCurrency) Units() units.UnitSet { return otherCurrencyUnits }

func TestQuantity_TypedAPI(t *testing.T) {
	testTypedMeasure(t, typedCase[Money, units.Unit[currency]]{
		measureType:   "currency",
		newFn:         units.NewQuantity[currency],
		toFn:          units.ToQuantity[currency],
		parseFn:       units.ParseQuantity[currency],
		unitA:         eurUnit,
		unitB:         uahUnit,
		ratio:         1.1 / 0.027,
		baseUnit:      usdUnit,
		baseRatio:     1.1,
		foreignSymbol: distance.Meter.String(),
		parseGeneric:  mustParser(t, units.WithMeasure[currency]()).Parse,
	})
}

func TestQuantity_GenericMethods(t *testing.T) {
	x, y := units.NewQuantity(2, eurUnit), units.NewQuantity(1, uahUnit)
	testGenericMethods(t, x, y, "UAH", 1.1/0.027, distance.Meter.String())
}

func TestQuantity_ZeroValue(t *testing.T) {
	var zero Money

	if zero.Value() != 0 || zero.Unit() != usdUnit || zero.String() != "0 USD" {
		t.Errorf("the zero Money is %v, want 0 USD", zero)
	}
}

func TestUnitOf(t *testing.T) {
	if usdUnit.String() != "USD" || usdUnit.Factor() != 1 || eurUnit.Factor() != 1.1 {
		t.Errorf("units: %v (%v), %v (%v)", usdUnit, usdUnit.Factor(), eurUnit, eurUnit.Factor())
	}

	if u, err := units.UnitOf[currency]("GBP"); !errors.Is(err, units.ErrUnknownUnit) || u != (units.Unit[currency]{}) {
		t.Errorf("UnitOf(GBP) = %v, %v; want the zero Unit and ErrUnknownUnit", u, err)
	}

	if _, err := units.UnitOf[noUnits]("x"); !errors.Is(err, units.ErrInvalidMeasure) {
		t.Errorf("UnitOf on a measure without units: %v, want ErrInvalidMeasure", err)
	}

	assertPanics(t, func() { units.MustUnitOf[currency]("GBP") })
}

func TestParseQuantity(t *testing.T) {
	q, err := units.ParseQuantity[currency]("10 EUR")
	if err != nil || q != units.NewQuantity(10, eurUnit) {
		t.Fatalf("ParseQuantity(10 EUR) = %v, %v", q, err)
	}

	for s, want := range map[string]error{
		"10 GBP":  units.ErrUnknownUnit,
		"10EUR":   units.ErrInvalidFormat,
		"NaN EUR": units.ErrInvalidFormat,
	} {
		if bad, err := units.ParseQuantity[currency](s); !errors.Is(err, want) || bad != (Money{}) {
			t.Errorf("ParseQuantity(%q) = %v, %v; want the zero value and %v", s, bad, err, want)
		}
	}

	if _, err := units.ParseQuantity[noUnits]("1 x"); !errors.Is(err, units.ErrInvalidMeasure) {
		t.Errorf("ParseQuantity on a measure without units: %v, want ErrInvalidMeasure", err)
	}
}

func TestQuantity_DifferentSetsWithTheSameName(t *testing.T) {
	usd := units.NewQuantity(1, usdUnit)
	gbp := units.NewQuantity(1, units.MustUnitOf[otherCurrency]("GBP"))

	if gbp.Type() != usd.Type() {
		t.Fatalf("types %s and %s, want both currency", gbp.Type(), usd.Type())
	}

	if sum, err := usd.AddMeasurement(gbp); !errors.Is(err, units.ErrIncompatibleMeasures) || sum != nil {
		t.Errorf("1 USD + 1 GBP of another set = %v, %v; want nil and ErrIncompatibleMeasures", sum, err)
	}

	if _, err := units.ToQuantity[currency](gbp); !errors.Is(err, units.ErrIncompatibleMeasures) {
		t.Errorf("ToQuantity of another measure: %v, want ErrIncompatibleMeasures", err)
	}

	if _, err := units.ToQuantity[noUnits](gbp); !errors.Is(err, units.ErrIncompatibleMeasures) {
		t.Errorf("ToQuantity to a measure without units: %v, want ErrIncompatibleMeasures", err)
	}
}

func TestQuantity_Format(t *testing.T) {
	q := units.NewQuantity(1.5, eurUnit)

	for format, want := range map[string]string{
		"%v":   "1.5 EUR",
		"%.2f": "1.50 EUR",
		"%#v":  `units.Quantity{Value:1.5, Unit:"EUR"}`,
	} {
		if got := fmt.Sprintf(format, q); got != want {
			t.Errorf("Sprintf(%q) = %q, want %q", format, got, want)
		}
	}
}

func TestNewUnitSet(t *testing.T) {
	s, err := units.NewUnitSet("currency", "USD", currencyFactors())
	if err != nil {
		t.Fatal(err)
	}

	if s.Type() != "currency" {
		t.Errorf("Type() = %s, want currency", s.Type())
	}

	// The base unit first, then the others sorted: the order doesn't depend on
	// map iteration.
	if got := slices.Collect(s.Symbols()); !slices.Equal(got, []string{"USD", "EUR", "UAH"}) {
		t.Errorf("Symbols() = %v, want [USD EUR UAH]", got)
	}

	for range s.Symbols() {
		break // the iterator stops when asked to
	}

	if got := s.ConversionFactors(); len(got) != 3 || got["EUR"] != 1.1 {
		t.Errorf("ConversionFactors() = %v", got)
	}

	var zero units.UnitSet
	if zero.Type() != "" || len(slices.Collect(zero.Symbols())) != 0 || len(zero.ConversionFactors()) != 0 {
		t.Error("the zero UnitSet isn't empty")
	}
}

func TestNewUnitSet_InvalidArguments(t *testing.T) {
	cases := map[string]struct {
		measure, base string
		factors       conversion.Factors
		want          error
	}{
		"empty measure": {"", "USD", currencyFactors(), units.ErrInvalidMeasure},
		"predefined measure": {
			string(units.MeasureDistance),
			"m",
			conversion.Factors{"m": 1},
			units.ErrInvalidMeasure,
		},
		"custom measure":       {string(units.MeasureCustom), "m", conversion.Factors{"m": 1}, units.ErrInvalidMeasure},
		"base unit is missing": {"currency", "GBP", currencyFactors(), units.ErrInvalidFactors},
		"no factors":           {"currency", "USD", nil, units.ErrInvalidFactors},
		"empty symbol":         {"currency", "USD", conversion.Factors{"USD": 1, "": 2}, units.ErrInvalidFactors},
		"zero factor":          {"currency", "USD", conversion.Factors{"USD": 1, "EUR": 0}, units.ErrInvalidFactors},
		"infinite factor":      {"currency", "USD", conversion.Factors{"USD": math.Inf(1)}, units.ErrInvalidFactors},
		"NaN factor":           {"currency", "USD", conversion.Factors{"USD": math.NaN()}, units.ErrInvalidFactors},
		"micro signs": {
			"tiny", "\u00b5m", conversion.Factors{"\u00b5m": 1, "\u03bcm": 1}, units.ErrInvalidFactors,
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			s, err := units.NewUnitSet(c.measure, c.base, c.factors)
			if !errors.Is(err, c.want) || s != (units.UnitSet{}) {
				t.Errorf("NewUnitSet = %v, %v; want the zero UnitSet and %v", s, err, c.want)
			}
		})
	}

	assertPanics(t, func() { units.MustNewUnitSet("", "USD", currencyFactors()) })

	tooMany := conversion.Factors{}
	for i := range 1<<16 + 1 {
		tooMany["u"+strconv.Itoa(i)] = 1
	}

	if _, err := units.NewUnitSet("huge", "u0", tooMany); !errors.Is(err, units.ErrInvalidFactors) {
		t.Errorf("a set of 65537 units: %v, want ErrInvalidFactors", err)
	}
}

func TestNewUnitSet_CopiesFactors(t *testing.T) {
	factors := currencyFactors()
	s := units.MustNewUnitSet("currency", "USD", factors)

	factors["EUR"] = 1000
	factors["GBP"] = 1.3

	if got := s.ConversionFactors(); got["EUR"] != 1.1 || len(got) != 3 {
		t.Errorf("changing the caller's map changed the set: %v", got)
	}
}

func TestNewUnitSet_MicroSign(t *testing.T) {
	s := units.MustNewUnitSet("tiny", "\u00b5m", conversion.Factors{"\u00b5m": 1, "mm": 1000})

	if got := slices.Collect(s.Symbols()); got[0] != "\u03bcm" {
		t.Errorf("Symbols() = %q, want the canonical micro sign first", got)
	}
}

// longDistance extends the predefined distance units with the furlong.
type longDistance struct{}

var longDistanceUnits = func() units.UnitSet {
	f := distance.ConversionFactors()
	f["fur"] = 201.168

	return units.MustNewUnitSet("long distance", distance.BaseUnit.String(), f)
}()

func (longDistance) Units() units.UnitSet { return longDistanceUnits }

func TestNewUnitSet_ExtendsPredefinedFactors(t *testing.T) {
	furlongs, err := units.ParseQuantity[longDistance]("8 fur")
	if err != nil {
		t.Fatal(err)
	}

	mi := furlongs.ConvertTo(units.MustUnitOf[longDistance]("mi"))
	assertMeasurement(t, mi, 1, "mi")

	// The predefined distance measure is unchanged.
	if _, ok := distance.ParseUnit("fur"); ok {
		t.Error("the furlong leaked into the distance package")
	}

	if _, err := units.ParseString("1 fur"); !errors.Is(err, units.ErrUnknownUnit) {
		t.Errorf("ParseString(1 fur): %v, want ErrUnknownUnit", err)
	}
}
