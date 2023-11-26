package units_test

import (
	units "github.com/asaphin/go-physics-units"
	"testing"
)

func BenchmarkDistance_AddSub(b *testing.B) {
	n := 10000

	distances1 := generateRandomDistances(n)
	distances2 := generateRandomDistances(n)

	var index int

	b.Run("Addition of two random distances", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			index = i % n

			_ = distances1[index].Add(distances2[index])
		}
	})

	b.Run("Subtraction of two random distances", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			index = i % n

			_ = distances1[index].Add(distances2[index])
		}
	})
}

func BenchmarkDistance_MulDiv(b *testing.B) {
	n := 10000

	distances := generateRandomDistances(n)
	multipliers := randFloats(-100, 100, n)

	var index int

	b.Run("Multiplying of distance", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			index = i % n

			_ = distances[index].Mul(multipliers[index])
		}
	})

	b.Run("Division of distance", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			index = i % n

			_, _ = distances[index].Div(multipliers[index])
		}
	})
}

func BenchmarkDistance_DivideByTime(b *testing.B) {
	n := 10000

	distances := generateRandomDistances(n)
	times := generateRandomTimes(n)

	var index int

	b.Run("Division of distance by time to get a velocity", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			index = i % n

			_ = distances[index].DivideByTime(times[index])
		}
	})
}

func BenchmarkDistance_ConvertToBaseUnit(b *testing.B) {
	n := 10000

	distances := generateRandomDistances(n)

	var index int

	b.Run("Conversion of distances to basic units - meters", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			index = i % n

			_ = distances[index].ConvertToBaseUnits()
		}
	})
}

func BenchmarkDistance_ToString(b *testing.B) {
	n := 10000

	distances := generateRandomDistances(n)

	var index int

	b.Run("Conversion of distances to strings", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			index = i % n

			_ = distances[index].String()
		}
	})
}

func BenchmarkDistance_ParseString(b *testing.B) {
	n := 10000

	distances := generateRandomDistances(n)

	str := make([]string, n)

	for i := range distances {
		str[i] = distances[i].String()
	}

	b.Run("Parse strings to distances", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			index := i % n

			_, err := units.ParseString(str[index])
			if err != nil {
				b.Error(err)
			}
		}
	})
}

func BenchmarkDistance_Creation(b *testing.B) {
	b.Run("NewDistanceMeters", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, _ = units.NewDistance(10.0, "m")
		}
	})
}

func BenchmarkDistance_ErrorHandling(b *testing.B) {
	b.Run("NewDistanceUnknownUnit", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, _ = units.NewDistance(10.0, "unk")
		}
	})
}
