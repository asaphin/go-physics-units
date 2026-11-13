package units_test

import (
	"testing"
	builtintime "time"

	units "github.com/asaphin/go-physics-units"
	"github.com/asaphin/go-physics-units/time"
)

func TestParseDuration(t *testing.T) {
	tm := units.ParseDuration(90*builtintime.Minute + 500*builtintime.Millisecond)

	assertMeasurement(t, tm, 5400.5, time.BaseUnit)
	assertMeasurement(t, tm.MustConvertTo(time.Hour), 5400.5/3600, time.Hour)
}
