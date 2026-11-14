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
	b.Run("density", func(b *testing.B) { benchTypedMeasure(b, densityCase) })
	b.Run("data rate", func(b *testing.B) { benchTypedMeasure(b, dataRateCase) })
	b.Run("electric current", func(b *testing.B) { benchTypedMeasure(b, currentCase) })
	b.Run("voltage", func(b *testing.B) { benchTypedMeasure(b, voltageCase) })
	b.Run("electric charge", func(b *testing.B) { benchTypedMeasure(b, chargeCase) })
	b.Run("resistance", func(b *testing.B) { benchTypedMeasure(b, resistanceCase) })
	b.Run("angular velocity", func(b *testing.B) { benchTypedMeasure(b, angularVelocityCase) })
}

func benchTypedMeasure[T typedMeasure[T, U], U unitType](b *testing.B, c typedCase[T, U]) {
	b.Helper()

	x := c.newFn(1234.5, c.unitA)
	y := c.newFn(2, c.unitB)

	s := x.String()

	b.Run("New", func(b *testing.B) {
		for b.Loop() {
			_ = c.newFn(1234.5, c.unitA)
		}
	})
	b.Run("ConvertTo", func(b *testing.B) {
		for b.Loop() {
			_ = x.ConvertTo(c.unitB)
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
	usd, eur := units.NewQuantity(1, usdUnit), units.NewQuantity(1, eurUnit)
	threeUnits := currencyFactors()
	distanceUnits := distance.ConversionFactors()

	b.Run("NewQuantity", func(b *testing.B) {
		for b.Loop() {
			_ = units.NewQuantity(1, eurUnit)
		}
	})
	b.Run("ParseQuantity", func(b *testing.B) {
		for b.Loop() {
			_, _ = units.ParseQuantity[currency]("2.5 EUR")
		}
	})
	b.Run("Add", func(b *testing.B) {
		for b.Loop() {
			_ = usd.Add(eur)
		}
	})
	b.Run("ConvertTo", func(b *testing.B) {
		for b.Loop() {
			_ = usd.ConvertTo(eurUnit)
		}
	})
	b.Run("Compare", func(b *testing.B) {
		for b.Loop() {
			_ = usd.Compare(eur)
		}
	})
	b.Run("AddMeasurement", func(b *testing.B) {
		for b.Loop() {
			_, _ = usd.AddMeasurement(eur)
		}
	})
	b.Run("NewUnitSet/3 units", func(b *testing.B) {
		for b.Loop() {
			_, _ = units.NewUnitSet("currency", "USD", threeUnits)
		}
	})
	b.Run("NewUnitSet/19 units", func(b *testing.B) {
		for b.Loop() {
			_, _ = units.NewUnitSet("long distance", distance.BaseUnit.String(), distanceUnits)
		}
	})
}

// BenchmarkSum adds up mixed distances, the loop the value types were designed
// for.
func BenchmarkSum(b *testing.B) {
	ds := generateRandomDistances(1000)

	for b.Loop() {
		var total units.Distance
		for _, d := range ds {
			total = total.Add(d)
		}

		sinkDistance = total
	}
}

// BenchmarkErrors covers error paths, which P1 keeps cheap too.
func BenchmarkErrors(b *testing.B) {
	b.Run("unknown unit", func(b *testing.B) {
		for b.Loop() {
			_, _ = units.StringToDistance("1 furlong")
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

// BenchmarkParser covers the configurable parser, against ParseString.
func BenchmarkParser(b *testing.B) {
	strict, _ := units.NewParser()
	all, _ := units.NewParser(units.CaseInsensitive(), units.AllowNoSpace(),
		units.WithAliases(map[string]string{"metres": distance.Meter.String()}))

	b.Run("ParseString", func(b *testing.B) {
		for b.Loop() {
			_, _ = units.ParseString("10 km")
		}
	})
	b.Run("default options", func(b *testing.B) {
		for b.Loop() {
			_, _ = strict.Parse("10 km")
		}
	})
	b.Run("alias", func(b *testing.B) {
		for b.Loop() {
			_, _ = all.Parse("10 metres")
		}
	})
	b.Run("no space, other case", func(b *testing.B) {
		for b.Loop() {
			_, _ = all.Parse("10KM")
		}
	})
}

// BenchmarkRate covers the unit packages' Rate functions.
func BenchmarkRate(b *testing.B) {
	for b.Loop() {
		_ = distance.Rate(distance.Kilometer, distance.Meter)
	}
}
