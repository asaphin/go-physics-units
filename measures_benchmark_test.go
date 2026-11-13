package units_test

import (
	"fmt"
	"testing"

	units "github.com/asaphin/go-physics-units"
	"github.com/asaphin/go-physics-units/distance"
)

// BenchmarkMeasures runs the shared API of every predefined measure, using the
// cases of TestTypedMeasures.
func BenchmarkMeasures(b *testing.B) {
	b.Run("distance", func(b *testing.B) { benchTypedMeasure(b, distanceCase) })
	b.Run("area", func(b *testing.B) { benchTypedMeasure(b, areaCase) })
	b.Run("mass", func(b *testing.B) { benchTypedMeasure(b, massCase) })
	b.Run("time", func(b *testing.B) { benchTypedMeasure(b, timeCase) })
	b.Run("velocity", func(b *testing.B) { benchTypedMeasure(b, velocityCase) })
	b.Run("volume", func(b *testing.B) { benchTypedMeasure(b, volumeCase) })
	b.Run("acceleration", func(b *testing.B) { benchTypedMeasure(b, accelerationCase) })
	b.Run("force", func(b *testing.B) { benchTypedMeasure(b, forceCase) })
	b.Run("energy", func(b *testing.B) { benchTypedMeasure(b, energyCase) })
	b.Run("power", func(b *testing.B) { benchTypedMeasure(b, powerCase) })
	b.Run("pressure", func(b *testing.B) { benchTypedMeasure(b, pressureCase) })
	b.Run("frequency", func(b *testing.B) { benchTypedMeasure(b, frequencyCase) })
	b.Run("angle", func(b *testing.B) { benchTypedMeasure(b, angleCase) })
	b.Run("data size", func(b *testing.B) { benchTypedMeasure(b, dataSizeCase) })
}

func benchTypedMeasure[T typedMeasure[T]](b *testing.B, c typedCase[T]) {
	b.Helper()

	x, err := c.newFn(1234.5, c.unitA)
	if err != nil {
		b.Fatal(err)
	}

	y, err := c.newFn(2, c.unitB)
	if err != nil {
		b.Fatal(err)
	}

	s := x.String()

	b.Run("New", func(b *testing.B) {
		for b.Loop() {
			_, _ = c.newFn(1234.5, c.unitA)
		}
	})
	b.Run("ConvertTo", func(b *testing.B) {
		for b.Loop() {
			_, _ = x.ConvertTo(c.unitB)
		}
	})
	b.Run("Add", func(b *testing.B) {
		for b.Loop() {
			_ = x.Add(y)
		}
	})
	b.Run("Mul", func(b *testing.B) {
		for b.Loop() {
			_ = x.Mul(3)
		}
	})
	b.Run("Compare", func(b *testing.B) {
		for b.Loop() {
			_ = x.Compare(y)
		}
	})
	b.Run("String", func(b *testing.B) {
		for b.Loop() {
			_ = x.String()
		}
	})
	b.Run("Format", func(b *testing.B) {
		for b.Loop() {
			_ = fmt.Sprintf("%.2f", x)
		}
	})
	b.Run("StringTo", func(b *testing.B) {
		for b.Loop() {
			_, _ = c.parseFn(s)
		}
	})
}

// BenchmarkCustomMeasures covers user-defined measures, which must stay as
// fast as the predefined ones (P1, P3).
func BenchmarkCustomMeasures(b *testing.B) {
	mc, err := units.NewMeasurementCreator[Currency]("currency", "USD", currencyFactors())
	if err != nil {
		b.Fatal(err)
	}

	usd, err := mc.New(1, "USD")
	if err != nil {
		b.Fatal(err)
	}

	eur, err := mc.New(1, "EUR")
	if err != nil {
		b.Fatal(err)
	}

	threeUnits := currencyFactors()
	distanceUnits := distance.ConversionFactors()

	b.Run("MeasurementCreator.New", func(b *testing.B) {
		for b.Loop() {
			_, _ = mc.New(1, "EUR")
		}
	})
	b.Run("MeasurementCreator.ParseString", func(b *testing.B) {
		for b.Loop() {
			_, _ = mc.ParseString("2.5 EUR")
		}
	})
	b.Run("AddMeasurement", func(b *testing.B) {
		for b.Loop() {
			_, _ = usd.AddMeasurement(eur)
		}
	})
	b.Run("NewBaseMeasurement/3 units", func(b *testing.B) {
		for b.Loop() {
			_, _ = units.NewBaseMeasurement(1, "USD", threeUnits)
		}
	})
	b.Run("NewBaseMeasurement/19 units", func(b *testing.B) {
		for b.Loop() {
			_, _ = units.NewBaseMeasurement(1, distance.Meter, distanceUnits)
		}
	})
}

// BenchmarkErrors covers error paths, which P1 keeps cheap too.
func BenchmarkErrors(b *testing.B) {
	b.Run("unknown unit", func(b *testing.B) {
		for b.Loop() {
			_, _ = units.NewDistance(1, "furlong")
		}
	})
	b.Run("unit of another measure", func(b *testing.B) {
		for b.Loop() {
			_, _ = units.StringToDistance("1 kg")
		}
	})
	b.Run("invalid format", func(b *testing.B) {
		for b.Loop() {
			_, _ = units.ParseString("10km")
		}
	})
}

// BenchmarkRate covers the unit packages' Rate functions.
func BenchmarkRate(b *testing.B) {
	for b.Loop() {
		_, _ = distance.Rate(distance.Kilometer, distance.Meter)
	}
}
