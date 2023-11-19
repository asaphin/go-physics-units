package units_test

import (
	units "github.com/asaphin/go-physics-units"
	"testing"
)

func BenchmarkDistance_DivideByTime(b *testing.B) {
	n := 10000

	distances := generateRandomDistances(n)

	times := generateRandomTimes(n)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		index := i % n

		_ = distances[index].DivideByTime(times[index])
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
