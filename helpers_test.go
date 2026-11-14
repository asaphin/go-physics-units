package units_test

import (
	"math"
	"math/rand/v2"
	"slices"
	"testing"
	builtintime "time"

	units "github.com/asaphin/go-physics-units"
	"github.com/asaphin/go-physics-units/distance"
	"github.com/asaphin/go-physics-units/time"
)

// Sinks keep the compiler from optimizing the measured calls away.
var (
	sinkDistance    units.Distance
	sinkVelocity    units.Velocity
	sinkMeasurement units.Measurement
	sinkMoney       Money
	sinkTemperature units.Temperature
	sinkDifference  units.TemperatureDifference
	sinkUnit        distance.Unit
	sinkString      string
	sinkBytes       []byte
	sinkRate        float64
	sinkOK          bool
	sinkInt         int
	sinkDuration    builtintime.Duration
	errSink         error
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

// valueWithSymbol is what assertMeasurement checks: a Measurement, or a
// temperature, which isn't one.
type valueWithSymbol interface {
	Value() float64
	Symbol() string
}

// assertMeasurement fails the test unless m has the given value (within a
// relative error of 1e-12) and unit.
func assertMeasurement(t *testing.T, m valueWithSymbol, value float64, unit string) {
	t.Helper()

	if m == nil {
		t.Fatalf("measurement is nil, want %v %s", value, unit)
	}

	if m.Symbol() != unit {
		t.Errorf("unit is %s, want %s", m.Symbol(), unit)
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

func randElement[T any](arr []T) T {
	if len(arr) == 0 {
		panic("empty array")
	}

	return arr[rand.IntN(len(arr))]
}

func generateRandomDistances(n int) []units.Distance {
	all := slices.Collect(distance.Units())
	distances := make([]units.Distance, n)

	for i := range n {
		distances[i] = units.NewDistance(randFloat(-100, 100), randElement(all))
	}

	return distances
}

func generateRandomTimes(n int) []units.Time {
	all := slices.Collect(time.Units())
	times := make([]units.Time, n)

	for i := range n {
		times[i] = units.NewTime(randFloat(-100, 100), randElement(all))
	}

	return times
}
