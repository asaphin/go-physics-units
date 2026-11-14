package prototype_test

import (
	"testing"

	units "github.com/asaphin/go-physics-units"
	"github.com/asaphin/go-physics-units/distance"
	prototype "github.com/asaphin/go-physics-units/docs/rfc-value-types"
)

// Sinks keep the compiler from removing the benchmarked calls.
var (
	sinkCurrent  units.Distance
	sinkS        prototype.DistanceS
	sinkI        prototype.Distance
	sinkQ        Money
	sinkInt      int
	sinkErr      error
	sinkBoxed    any
	sinkString   string
	sinkInterval float64
)

func BenchmarkDistance(b *testing.B) {
	cur, _ := units.NewDistance(1234.5, distance.Mile)
	curY, _ := units.NewDistance(2, distance.Kilometer)
	s, _ := prototype.NewDistanceS(1234.5, prototype.SMile)
	sY, _ := prototype.NewDistanceS(2, prototype.SKilometer)
	i := prototype.NewDistance(1234.5, prototype.Mile)
	iY := prototype.NewDistance(2, prototype.Kilometer)
	usd, _ := prototype.UnitOf[currency]("USD")
	eur, _ := prototype.UnitOf[currency]("EUR")
	q := prototype.New(1234.5, usd)
	qY := prototype.New(2, eur)

	b.Run("New/current", func(b *testing.B) {
		for b.Loop() {
			sinkCurrent, sinkErr = units.NewDistance(1234.5, distance.Mile)
		}
	})
	b.Run("New/S", func(b *testing.B) {
		for b.Loop() {
			sinkS, sinkErr = prototype.NewDistanceS(1234.5, prototype.SMile)
		}
	})
	b.Run("New/I", func(b *testing.B) {
		for b.Loop() {
			sinkI = prototype.NewDistance(1234.5, prototype.Mile)
		}
	})
	b.Run("New/Quantity", func(b *testing.B) {
		for b.Loop() {
			sinkQ = prototype.New(1234.5, usd)
		}
	})

	b.Run("ConvertTo/current", func(b *testing.B) {
		for b.Loop() {
			sinkCurrent, sinkErr = cur.ConvertTo(distance.Kilometer)
		}
	})
	b.Run("ConvertTo/S", func(b *testing.B) {
		for b.Loop() {
			sinkS, sinkErr = s.ConvertTo(prototype.SKilometer)
		}
	})
	b.Run("ConvertTo/I", func(b *testing.B) {
		for b.Loop() {
			sinkI = i.ConvertTo(prototype.Kilometer)
		}
	})
	b.Run("ConvertTo/Quantity", func(b *testing.B) {
		for b.Loop() {
			sinkQ = q.ConvertTo(eur)
		}
	})

	b.Run("Add/current", func(b *testing.B) {
		for b.Loop() {
			sinkCurrent = cur.Add(curY)
		}
	})
	b.Run("Add/S", func(b *testing.B) {
		for b.Loop() {
			sinkS = s.Add(sY)
		}
	})
	b.Run("Add/I", func(b *testing.B) {
		for b.Loop() {
			sinkI = i.Add(iY)
		}
	})
	b.Run("Add/Quantity", func(b *testing.B) {
		for b.Loop() {
			sinkQ = q.Add(qY)
		}
	})

	b.Run("Compare/current", func(b *testing.B) {
		for b.Loop() {
			sinkInt = cur.Compare(curY)
		}
	})
	b.Run("Compare/S", func(b *testing.B) {
		for b.Loop() {
			sinkInt = s.Compare(sY)
		}
	})
	b.Run("Compare/I", func(b *testing.B) {
		for b.Loop() {
			sinkInt = i.Compare(iY)
		}
	})
	b.Run("Compare/Quantity", func(b *testing.B) {
		for b.Loop() {
			sinkInt = q.Compare(qY)
		}
	})

	// The cost of the generic path in a value-type design: boxing a value into
	// an interface, as passing it as a units.Measurement would.
	b.Run("Boxing/I", func(b *testing.B) {
		for b.Loop() {
			sinkBoxed = i.Add(iY)
		}
	})
	b.Run("String/I", func(b *testing.B) {
		for b.Loop() {
			sinkString = i.String()
		}
	})
}

// BenchmarkLoop sums 1000 distances in mixed units: the workload P1 is about.
func BenchmarkLoop(b *testing.B) {
	const n = 1000

	cur := make([]units.Distance, n)
	vals := make([]prototype.Distance, n)
	us := []prototype.Unit{prototype.Meter, prototype.Kilometer, prototype.Mile}
	cus := []string{distance.Meter, distance.Kilometer, distance.Mile}

	for k := range n {
		cur[k], _ = units.NewDistance(float64(k), cus[k%3])
		vals[k] = prototype.NewDistance(float64(k), us[k%3])
	}

	b.Run("current", func(b *testing.B) {
		for b.Loop() {
			sum, _ := units.NewDistance(0, distance.Meter)
			for _, d := range cur {
				sum = sum.Add(d)
			}

			sinkInterval = sum.Value()
		}
	})
	b.Run("I", func(b *testing.B) {
		for b.Loop() {
			sum := prototype.NewDistance(0, prototype.Meter)
			for _, d := range vals {
				sum = sum.Add(d)
			}

			sinkInterval = sum.Value()
		}
	})
}
