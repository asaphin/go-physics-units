package units_test

import (
	"math"
	"math/rand/v2"
	"testing"

	units "github.com/asaphin/go-physics-units"
	"github.com/asaphin/go-physics-units/conversion"
	"github.com/asaphin/go-physics-units/distance"
	"github.com/asaphin/go-physics-units/time"
)

// almostEqual reports whether a and b are equal within the relative error epsilon.
func almostEqual(a, b, epsilon float64) bool {
	if a == b {
		return true
	}

	absA := math.Abs(a)
	absB := math.Abs(b)
	diff := math.Abs(a - b)

	if a == 0 || b == 0 || absA+absB < math.SmallestNonzeroFloat64 {
		return diff < epsilon
	}

	return diff/(absA+absB) < epsilon
}

// assertMeasurement fails the test unless m has the given value (within a
// relative error of 1e-12) and unit.
func assertMeasurement(t *testing.T, m units.Measurement, value float64, unit string) {
	t.Helper()

	if m == nil {
		t.Fatalf("measurement is nil, want %v %s", value, unit)
	}

	if m.Unit() != unit {
		t.Errorf("unit is %s, want %s", m.Unit(), unit)
	}

	if !almostEqual(m.Value(), value, 1e-12) {
		t.Errorf("value is %v, want %v", m.Value(), value)
	}
}

// assertPanics fails the test unless f panics.
func assertPanics(t *testing.T, f func()) {
	t.Helper()

	defer func() {
		if recover() == nil {
			t.Error("expected a panic")
		}
	}()

	f()
}

func getValidUnits(factors conversion.Factors) []string {
	keys := make([]string, 0, len(factors))
	for key := range factors {
		keys = append(keys, key)
	}

	return keys
}

func randFloat(lo, hi float64) float64 {
	return lo + rand.Float64()*(hi-lo)
}

func randFloats(lo, hi float64, n int) []float64 {
	floats := make([]float64, n)

	for i := range floats {
		floats[i] = randFloat(lo, hi)
	}

	return floats
}

func randStringElement(arr []string) string {
	if len(arr) == 0 {
		panic("empty array")
	}

	return arr[rand.IntN(len(arr))]
}

func generateRandomDistances(n int) []units.Distance {
	validDistanceUnits := getValidUnits(distance.ConversionFactors())

	distances := make([]units.Distance, n)

	for i := range n {
		dist, err := units.NewDistance(randFloat(-100, 100), randStringElement(validDistanceUnits))
		if err != nil {
			panic(err)
		}

		distances[i] = dist
	}

	return distances
}

func generateRandomTimes(n int) []units.Time {
	validTimeUnits := getValidUnits(time.ConversionFactors())

	times := make([]units.Time, n)

	for i := range n {
		tm, err := units.NewTime(randFloat(-100, 100), randStringElement(validTimeUnits))
		if err != nil {
			panic(err)
		}

		times[i] = tm
	}

	return times
}
