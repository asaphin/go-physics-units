//go:build !race

package units_test

import (
	"testing"
	builtintime "time"

	units "github.com/asaphin/go-physics-units"
	"github.com/asaphin/go-physics-units/conversion"
	"github.com/asaphin/go-physics-units/distance"
	"github.com/asaphin/go-physics-units/time"
)

// Sinks keep the compiler from optimizing the measured calls away.
var (
	sinkDistance    units.Distance
	sinkVelocity    units.Velocity
	sinkMeasurement units.Measurement
	sinkCurrency    Currency
	sinkString      string
	sinkRate        float64
	sinkOK          bool
	sinkInt         int
	sinkDuration    builtintime.Duration
	errSink         error
)

// TestAllocations guards the performance principle (P1). Unlike timings on
// shared CI runners, allocation counts are deterministic, so a regression fails
// the build. Raise a limit only on purpose and say why in the PR; lower it when
// an optimization saves an allocation. The race detector changes allocation
// counts, so this file is excluded from -race builds.
func TestAllocations(t *testing.T) {
	km, err := units.NewDistance(1, distance.Kilometer)
	if err != nil {
		t.Fatal(err)
	}

	m500, err := units.NewDistance(500, distance.Meter)
	if err != nil {
		t.Fatal(err)
	}

	hour, err := units.NewTime(1, time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	mc, err := units.NewMeasurementCreator[Currency]("currency", "USD", conversion.Factors{"USD": 1, "EUR": 1.1})
	if err != nil {
		t.Fatal(err)
	}

	parser, err := units.NewParser(units.AllowNoSpace(), units.CaseInsensitive())
	if err != nil {
		t.Fatal(err)
	}

	cases := []struct {
		name  string
		limit float64
		f     func()
	}{
		{"NewDistance", 1, func() { sinkDistance, errSink = units.NewDistance(1, distance.Kilometer) }},
		{"ConvertTo", 1, func() { sinkDistance, errSink = km.ConvertTo(distance.Mile) }},
		{"ConvertToBaseUnits", 1, func() { sinkDistance = km.ConvertToBaseUnits() }},
		{"Add", 1, func() { sinkDistance = km.Add(m500) }},
		{"Sub", 1, func() { sinkDistance = km.Sub(m500) }},
		{"Mul", 1, func() { sinkDistance = km.Mul(2) }},
		{"Div", 1, func() { sinkDistance, errSink = km.Div(2) }},
		{"DivideByTime", 1, func() { sinkVelocity, errSink = km.DivideByTime(hour) }},
		{"AddMeasurement", 1, func() { sinkMeasurement, errSink = km.AddMeasurement(m500) }},
		{"String", 1, func() { sinkString = km.String() }},
		{"ParseString", 1, func() { sinkMeasurement, errSink = units.ParseString("10 km") }},
		{"StringToDistance", 1, func() { sinkDistance, errSink = units.StringToDistance("10 km") }},
		{"MeasurementCreator.New", 1, func() { sinkCurrency, errSink = mc.New(1, "EUR") }},
		{"Compare", 0, func() { sinkInt = km.Compare(m500) }},
		{"EqualWithinRel", 0, func() { sinkOK = km.EqualWithinRel(m500, 1e-9) }},
		{"CompareMeasurement", 0, func() { sinkInt, errSink = km.CompareMeasurement(m500) }},
		{"ToDuration", 0, func() { sinkDuration, errSink = hour.ToDuration() }},
		{"Parser.Parse", 1, func() { sinkMeasurement, errSink = parser.Parse("10KM") }},
		{"distance.Rate", 0, func() { sinkRate, sinkOK = distance.Rate(distance.Kilometer, distance.Meter) }},
		{"unknown unit error", 1, func() { sinkDistance, errSink = units.NewDistance(1, "furlong") }},
		{"invalid format error", 2, func() { sinkMeasurement, errSink = units.ParseString("10km") }},
		{"zero division error", 0, func() { sinkDistance, errSink = km.Div(0) }},
		{"incompatible measures error", 0, func() { sinkDistance, errSink = units.ToDistance(nil) }},
	}

	for _, c := range cases {
		if got := testing.AllocsPerRun(100, c.f); got > c.limit {
			t.Errorf("%s: %v allocations per call, limit %v", c.name, got, c.limit)
		}
	}
}
