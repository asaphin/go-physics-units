//go:build !race

package units_test

import (
	"testing"

	units "github.com/asaphin/go-physics-units"
	"github.com/asaphin/go-physics-units/distance"
	"github.com/asaphin/go-physics-units/temperature"
	"github.com/asaphin/go-physics-units/time"
)

// TestAllocations guards the performance principle (P1). Unlike timings on
// shared CI runners, allocation counts are deterministic, so a regression fails
// the build. Raise a limit only on purpose and say why in the PR; lower it when
// an optimization saves an allocation. The race detector changes allocation
// counts, so this file is excluded from -race builds.
//
// Values are plain structs, so the typed API doesn't allocate at all. What
// remains is the generic API, which returns a Measurement interface (one
// allocation to box the value), String, and the error paths that build a
// message.
func TestAllocations(t *testing.T) {
	km := units.NewDistance(1, distance.Kilometer)
	m500 := units.NewDistance(500, distance.Meter)

	// Boxing a value into a Measurement allocates; do it once, outside the
	// measured calls, as code that holds Measurements already has.
	var m500m units.Measurement = m500
	hour := units.NewTime(1, time.Hour)
	eur := units.NewQuantity(1, eurUnit)
	room := units.Temperature{} // 0 K
	rise := units.NewTemperatureDifference(10, temperature.Celsius)

	parser, err := units.NewParser(units.AllowNoSpace(), units.CaseInsensitive())
	if err != nil {
		t.Fatal(err)
	}

	buf := make([]byte, 0, 64)

	cases := []struct {
		name  string
		limit float64
		f     func()
	}{
		// The typed API: no allocations.
		{"NewDistance", 0, func() { sinkDistance = units.NewDistance(1, distance.Kilometer) }},
		{"ConvertTo", 0, func() { sinkDistance = km.ConvertTo(distance.Mile) }},
		{"ConvertToBaseUnits", 0, func() { sinkDistance = km.ConvertToBaseUnits() }},
		{"Add", 0, func() { sinkDistance = km.Add(m500) }},
		{"Sub", 0, func() { sinkDistance = km.Sub(m500) }},
		{"Mul", 0, func() { sinkDistance = km.Mul(2) }},
		{"Div", 0, func() { sinkDistance, errSink = km.Div(2) }},
		{"Compare", 0, func() { sinkInt = km.Compare(m500) }},
		{"EqualWithinRel", 0, func() { sinkOK = km.EqualWithinRel(m500, 1e-9) }},
		{"DivideByTime", 0, func() { sinkVelocity, errSink = km.DivideByTime(hour) }},
		{"StringToDistance", 0, func() { sinkDistance, errSink = units.StringToDistance("10 km") }},
		{"ToDuration", 0, func() { sinkDuration, errSink = hour.ToDuration() }},
		{"distance.Rate", 0, func() { sinkRate = distance.Rate(distance.Kilometer, distance.Meter) }},
		{"distance.ParseUnit", 0, func() { sinkUnit, sinkOK = distance.ParseUnit("km") }},
		{"Quantity.Add", 0, func() { sinkMoney = eur.Add(eur) }},
		{"Quantity.ConvertTo", 0, func() { sinkMoney = eur.ConvertTo(usdUnit) }},
		{"ParseQuantity", 0, func() { sinkMoney, errSink = units.ParseQuantity[currency]("1 EUR") }},
		{"Temperature.ConvertTo", 0, func() { sinkTemperature = room.ConvertTo(temperature.Fahrenheit) }},
		{"Temperature.Add", 0, func() { sinkTemperature, errSink = room.Add(rise) }},
		{"Temperature.Sub", 0, func() { sinkDifference = room.Sub(room) }},
		{"CompareMeasurement", 0, func() { sinkInt, errSink = km.CompareMeasurement(m500m) }},
		{"AppendText", 0, func() { sinkBytes, errSink = km.AppendText(buf[:0]) }},

		// One allocation: the result string, or the Measurement interface.
		{"String", 1, func() { sinkString = km.String() }},
		{"AddMeasurement", 1, func() { sinkMeasurement, errSink = km.AddMeasurement(m500m) }},
		{"ParseString", 1, func() { sinkMeasurement, errSink = units.ParseString("10 km") }},
		{"Parser.Parse", 1, func() { sinkMeasurement, errSink = parser.Parse("10KM") }},

		// Error paths.
		{"unknown unit error", 1, func() { sinkDistance, errSink = units.StringToDistance("1 furlong") }},
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
