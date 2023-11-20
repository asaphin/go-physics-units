package units_test

import (
	units "github.com/asaphin/go-physics-units"
	"github.com/asaphin/go-physics-units/conversion-factors"
	"github.com/asaphin/go-physics-units/distance"
	"github.com/asaphin/go-physics-units/time"
	"math"
	"math/rand"
)

func almostEqual(a, b, epsilon float64) bool {
	if a == b {
		return true
	}

	absA := math.Abs(a)
	absB := math.Abs(b)
	diff := math.Abs(a - b)

	if a == 0 || b == 0 || absA+absB < math.SmallestNonzeroFloat64 {
		return diff < epsilon
	} else {
		return diff/(absA+absB) < epsilon
	}
}

func getValidUnits(factors conversion.Factors) []string {
	keys := make([]string, 0, len(factors))
	for key := range factors {
		keys = append(keys, key)
	}

	return keys
}

func randFloat(min, max float64) float64 {
	return min + rand.Float64()*(max-min)
}

func randFloats(min, max float64, n int) []float64 {
	floats := make([]float64, n)

	for i := range floats {
		floats[i] = randFloat(min, max)
	}

	return floats
}

func randStringElement(arr []string) string {
	if len(arr) == 0 {
		panic("empty array")
	}

	index := rand.Intn(len(arr))

	return arr[index]
}

func generateRandomDistances(n int) []*units.Distance {
	validDistanceUnits := getValidUnits(distance.ConversionFactors())

	distances := make([]*units.Distance, n)

	for i := 0; i < n; i++ {
		dist, err := units.NewDistance(randFloat(-100, 100), randStringElement(validDistanceUnits))
		if err != nil {
			panic(err)
		}

		distances[i] = dist
	}

	return distances
}

func generateRandomTimes(n int) []*units.Time {
	validTimeUnits := getValidUnits(time.ConversionFactors())

	times := make([]*units.Time, n)

	for i := 0; i < n; i++ {
		tm, err := units.NewTime(randFloat(-100, 100), randStringElement(validTimeUnits))
		if err != nil {
			panic(err)
		}

		times[i] = tm
	}

	return times
}
