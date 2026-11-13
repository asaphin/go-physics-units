package units_test

import (
	"testing"

	units "github.com/asaphin/go-physics-units"
	"github.com/asaphin/go-physics-units/area"
	"github.com/asaphin/go-physics-units/conversion"
	"github.com/asaphin/go-physics-units/distance"
	"github.com/asaphin/go-physics-units/mass"
	"github.com/asaphin/go-physics-units/time"
	"github.com/asaphin/go-physics-units/velocity"
)

var predefinedFactors = map[units.MeasureType]conversion.Factors{
	units.MeasureDistance: distance.ConversionFactors(),
	units.MeasureArea:     area.ConversionFactors(),
	units.MeasureMass:     mass.ConversionFactors(),
	units.MeasureTime:     time.ConversionFactors(),
	units.MeasureVelocity: velocity.ConversionFactors(),
}

// TestConversion_AllUnitPairs converts between every ordered pair of units of
// every predefined measure, and back again.
func TestConversion_AllUnitPairs(t *testing.T) {
	for mt, factors := range predefinedFactors {
		for from, fromFactor := range factors {
			for to, toFactor := range factors {
				m, err := units.NewMeasurement(1, from)
				if err != nil {
					t.Fatalf("%s: %v", mt, err)
				}

				if m.Type() != mt {
					t.Errorf("unit %s has type %s, want %s", from, m.Type(), mt)
				}

				got, err := m.ConvertToMeasurement(to)
				if err != nil {
					t.Fatalf("%s: converting %s to %s: %v", mt, from, to, err)
				}

				if want := fromFactor / toFactor; !almostEqual(got.Value(), want, 1e-12) {
					t.Errorf("%s: 1 %s = %v %s, want %v", mt, from, got.Value(), to, want)
				}

				back, err := got.ConvertToMeasurement(from)
				if err != nil {
					t.Fatalf("%s: converting %s back to %s: %v", mt, to, from, err)
				}

				if !almostEqual(back.Value(), 1, 1e-12) {
					t.Errorf("%s: round trip %s -> %s -> %s gives %v, want 1", mt, from, to, from, back.Value())
				}
			}
		}
	}
}

// TestConversion_ReferenceValues checks conversions against independently known
// equivalences, so that a wrong factor can't go unnoticed.
func TestConversion_ReferenceValues(t *testing.T) {
	cases := []struct {
		from string
		to   string
		want float64 // 1 from = want to
	}{
		{distance.Meter, distance.Millimeter, 1000},
		{distance.Millimeter, distance.Meter, 0.001},
		{distance.Kilometer, distance.Meter, 1000},
		{distance.Mile, distance.Kilometer, 1.609344},
		{distance.NauticalMile, distance.Meter, 1852},
		{distance.Inch, distance.Centimeter, 2.54},
		{distance.Feet, distance.Inch, 12},
		{distance.Yard, distance.Feet, 3},
		{distance.Angstrom, distance.Nanometer, 0.1},
		{distance.Micrometer, distance.Nanometer, 1000},
		{distance.AstronomicalUnit, distance.Kilometer, 149597870.7},
		{distance.LightYear, distance.Kilometer, 9460730472580.8},
		{distance.Parsec, distance.LightYear, 3.2615637771674},
		{area.SquareKilometer, area.SquareMeter, 1e6},
		{area.SquareMeter, area.SquareCentimeter, 1e4},
		{area.SquareCentimeter, area.SquareMillimeter, 100},
		{area.SquareFoot, area.SquareInch, 144},
		{area.SquareMile, area.SquareKilometer, 2.589988110336},
		{mass.Kilogram, mass.Gram, 1000},
		{mass.Gram, mass.Microgram, 1e6},
		{mass.Milligram, mass.Microgram, 1000},
		{mass.Pound, mass.Kilogram, 0.45359237},
		{mass.Pound, mass.Ounce, 16},
		{mass.Quintal, mass.Kilogram, 100},
		{mass.Tonne, mass.Quintal, 10},
		{time.Second, time.Millisecond, 1000},
		{time.Millisecond, time.Microsecond, 1000},
		{time.Microsecond, time.Nanosecond, 1000},
		{time.Hour, time.Minute, 60},
		{time.Day, time.Hour, 24},
		{time.Week, time.Day, 7},
		{time.Year, time.Day, 365.25},
		{velocity.MeterPerSecond, velocity.KilometerPerHour, 3.6},
		{velocity.Knot, velocity.KilometerPerHour, 1.852},
		{velocity.MilePerHour, velocity.KilometerPerHour, 1.609344},
		{velocity.FeetPerSecond, velocity.MeterPerSecond, 0.3048},
		{velocity.SpeedOfLight, velocity.KilometerPerHour, 1079252848.8},
	}

	for _, c := range cases {
		m, err := units.NewMeasurement(1, c.from)
		if err != nil {
			t.Fatal(err)
		}

		got, err := m.ConvertToMeasurement(c.to)
		if err != nil {
			t.Fatal(err)
		}

		if !almostEqual(got.Value(), c.want, 1e-12) {
			t.Errorf("1 %s = %v %s, want %v", c.from, got.Value(), c.to, c.want)
		}
	}
}

func TestConversion_MicroSignVariants(t *testing.T) {
	const (
		microSign = "\u00b5" // MICRO SIGN
		greekMu   = "\u03bc" // GREEK SMALL LETTER MU
	)

	for _, unit := range []string{distance.Micrometer, mass.Microgram, time.Microsecond} {
		if unit[:len(greekMu)] != greekMu {
			t.Errorf("predefined symbol %s doesn't use U+03BC", unit)
		}
	}

	d, err := units.NewDistance(1, microSign+"m")
	if err != nil {
		t.Fatal(err)
	}

	assertMeasurement(t, d, 1, distance.Micrometer)

	nm, err := d.ConvertTo("nm")
	if err != nil {
		t.Fatal(err)
	}

	assertMeasurement(t, nm, 1000, distance.Nanometer)

	back, err := nm.ConvertTo(microSign + "m")
	if err != nil {
		t.Fatal(err)
	}

	assertMeasurement(t, back, 1, distance.Micrometer)

	g, err := units.StringToMass("5 " + microSign + "g")
	if err != nil {
		t.Fatal(err)
	}

	assertMeasurement(t, g, 5, mass.Microgram)

	if mt := units.DetectMeasureType(microSign + "s"); mt != units.MeasureTime {
		t.Errorf("DetectMeasureType(\u00b5s) = %s, want %s", mt, units.MeasureTime)
	}
}

func TestConversionFactors_ReturnsCopy(t *testing.T) {
	for mt, factorsFn := range map[units.MeasureType]func() conversion.Factors{
		units.MeasureDistance: distance.ConversionFactors,
		units.MeasureArea:     area.ConversionFactors,
		units.MeasureMass:     mass.ConversionFactors,
		units.MeasureTime:     time.ConversionFactors,
		units.MeasureVelocity: velocity.ConversionFactors,
	} {
		f := factorsFn()
		for unit := range f {
			f[unit] = -1
		}

		for unit, factor := range factorsFn() {
			if factor == -1 {
				t.Errorf("%s: mutating the result of ConversionFactors changed factor of %s", mt, unit)
			}
		}
	}
}
