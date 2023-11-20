package units_test

import (
	units "github.com/asaphin/go-physics-units"
	"testing"
)

func BenchmarkDistance_Mul(b *testing.B) {
	n := 10000

	distances := generateRandomDistances(n)

	multipliers := randFloats(-100, 100, n)

	var index int

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		index = i % n

		_ = distances[index].Mul(multipliers[index])
	}
}

func BenchmarkDistance_Div(b *testing.B) {
	n := 10000

	distances := generateRandomDistances(n)

	divisors := randFloats(-100, 100, n)

	var index int

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		index = i % n

		_ = distances[index].Div(divisors[index])
	}
}

func BenchmarkDistance_DivideByTime(b *testing.B) {
	n := 10000

	distances := generateRandomDistances(n)

	times := generateRandomTimes(n)

	var index int

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		index = i % n

		_ = distances[index].DivideByTime(times[index])
	}
}

func BenchmarkDistance_ConvertToBaseUnit(b *testing.B) {
	n := 10000

	distances := generateRandomDistances(n)

	var index int

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		index = i % n

		_ = distances[index].ConvertToBaseUnits()
	}
}

func BenchmarkDistance_ToString(b *testing.B) {
	n := 10000

	distances := generateRandomDistances(n)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		index := i % n

		_ = distances[index].String()
	}
}

func BenchmarkDistance_ParseString(b *testing.B) {
	n := 10000

	distances := generateRandomDistances(n)

	strs := make([]string, n)

	for i := range distances {
		strs[i] = distances[i].String()
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		index := i % n

		_, err := units.ParseString(strs[index])
		if err != nil {
			b.Error(err)
		}
	}
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
