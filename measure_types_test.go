package units_test

import (
	"errors"
	"testing"

	units "github.com/asaphin/go-physics-units"
	"github.com/asaphin/go-physics-units/area"
	"github.com/asaphin/go-physics-units/distance"
	"github.com/asaphin/go-physics-units/mass"
	"github.com/asaphin/go-physics-units/time"
	"github.com/asaphin/go-physics-units/velocity"
)

// typedMeasure is the API shared by the predefined measure interfaces.
type typedMeasure[T any] interface {
	units.Measurement
	ConvertTo(unit string) (T, error)
	MustConvertTo(unit string) T
	ConvertToBaseUnits() T
	Add(other T) T
	Sub(other T) T
	Mul(multiplier float64) T
	Div(divisor float64) (T, error)
}

type typedCase[T typedMeasure[T]] struct {
	measureType units.MeasureType
	newFn       func(value float64, unit string) (T, error)
	toFn        func(m units.Measurement) (T, error)
	parseFn     func(s string) (T, error)
	unitA       string
	unitB       string
	ratio       float64 // 1 unitA = ratio unitB
	baseUnit    string
	baseRatio   float64 // 1 unitA = baseRatio baseUnit
	foreignUnit string  // a unit of another measure
}

// One case per predefined measure, shared by TestTypedMeasures and
// BenchmarkMeasures. A new measure needs a case here.
var (
	distanceCase = typedCase[units.Distance]{
		measureType: units.MeasureDistance,
		newFn:       units.NewDistance, toFn: units.ToDistance, parseFn: units.StringToDistance,
		unitA: distance.Mile, unitB: distance.Kilometer, ratio: 1.609344,
		baseUnit: distance.BaseUnit, baseRatio: 1609.344, foreignUnit: mass.Kilogram,
	}
	areaCase = typedCase[units.Area]{
		measureType: units.MeasureArea,
		newFn:       units.NewArea, toFn: units.ToArea, parseFn: units.StringToArea,
		unitA: area.SquareFoot, unitB: area.SquareInch, ratio: 144,
		baseUnit: area.BaseUnit, baseRatio: 0.09290304, foreignUnit: time.Second,
	}
	massCase = typedCase[units.Mass]{
		measureType: units.MeasureMass,
		newFn:       units.NewMass, toFn: units.ToMass, parseFn: units.StringToMass,
		unitA: mass.Pound, unitB: mass.Ounce, ratio: 16,
		baseUnit: mass.BaseUnit, baseRatio: 0.45359237, foreignUnit: distance.Meter,
	}
	timeCase = typedCase[units.Time]{
		measureType: units.MeasureTime,
		newFn:       units.NewTime, toFn: units.ToTime, parseFn: units.StringToTime,
		unitA: time.Hour, unitB: time.Minute, ratio: 60,
		baseUnit: time.BaseUnit, baseRatio: 3600, foreignUnit: mass.Kilogram,
	}
	velocityCase = typedCase[units.Velocity]{
		measureType: units.MeasureVelocity,
		newFn:       units.NewVelocity, toFn: units.ToVelocity, parseFn: units.StringToVelocity,
		unitA: velocity.Knot, unitB: velocity.KilometerPerHour, ratio: 1.852,
		baseUnit: velocity.BaseUnit, baseRatio: 1852.0 / 3600, foreignUnit: distance.Meter,
	}
)

func TestTypedMeasures(t *testing.T) {
	t.Run("distance", func(t *testing.T) { testTypedMeasure(t, distanceCase) })
	t.Run("area", func(t *testing.T) { testTypedMeasure(t, areaCase) })
	t.Run("mass", func(t *testing.T) { testTypedMeasure(t, massCase) })
	t.Run("time", func(t *testing.T) { testTypedMeasure(t, timeCase) })
	t.Run("velocity", func(t *testing.T) { testTypedMeasure(t, velocityCase) })
}

func testTypedMeasure[T typedMeasure[T]](t *testing.T, c typedCase[T]) {
	t.Helper()

	x, err := c.newFn(2, c.unitA)
	if err != nil {
		t.Fatal(err)
	}

	y, err := c.newFn(1, c.unitB)
	if err != nil {
		t.Fatal(err)
	}

	t.Run("Type", func(t *testing.T) {
		if x.Type() != c.measureType {
			t.Errorf("Type() = %s, want %s", x.Type(), c.measureType)
		}
	})

	t.Run("String", func(t *testing.T) {
		if want := "2 " + c.unitA; x.String() != want {
			t.Errorf("String() = %q, want %q", x.String(), want)
		}
	})

	t.Run("ConvertTo", func(t *testing.T) {
		got, err := x.ConvertTo(c.unitB)
		if err != nil {
			t.Fatal(err)
		}

		assertMeasurement(t, got, 2*c.ratio, c.unitB)

		bad, err := x.ConvertTo(c.foreignUnit)
		if !errors.Is(err, units.ErrUnknownUnit) || any(bad) != nil {
			t.Errorf("ConvertTo(%s) = %v, %v; want nil and ErrUnknownUnit", c.foreignUnit, bad, err)
		}
	})

	t.Run("MustConvertTo", func(t *testing.T) {
		assertMeasurement(t, x.MustConvertTo(c.unitB), 2*c.ratio, c.unitB)
		assertPanics(t, func() { x.MustConvertTo(c.foreignUnit) })
	})

	t.Run("ConvertToBaseUnits", func(t *testing.T) {
		assertMeasurement(t, x.ConvertToBaseUnits(), 2*c.baseRatio, c.baseUnit)
	})

	t.Run("Add", func(t *testing.T) {
		assertMeasurement(t, x.Add(y), 2+1/c.ratio, c.unitA)
	})

	t.Run("Sub", func(t *testing.T) {
		assertMeasurement(t, x.Sub(y), 2-1/c.ratio, c.unitA)
	})

	t.Run("Mul", func(t *testing.T) {
		assertMeasurement(t, x.Mul(3), 6, c.unitA)
	})

	t.Run("Div", func(t *testing.T) {
		got, err := x.Div(4)
		if err != nil {
			t.Fatal(err)
		}

		assertMeasurement(t, got, 0.5, c.unitA)

		bad, err := x.Div(0)
		if !errors.Is(err, units.ErrZeroDivision) || any(bad) != nil {
			t.Errorf("Div(0) = %v, %v; want nil and ErrZeroDivision", bad, err)
		}
	})

	t.Run("operations don't mutate the receiver", func(t *testing.T) {
		assertMeasurement(t, x, 2, c.unitA)
	})

	t.Run("New with a unit of another measure", func(t *testing.T) {
		bad, err := c.newFn(1, c.foreignUnit)
		if !errors.Is(err, units.ErrUnknownUnit) || any(bad) != nil {
			t.Errorf("New(1, %s) = %v, %v; want nil and ErrUnknownUnit", c.foreignUnit, bad, err)
		}
	})

	t.Run("To", func(t *testing.T) {
		same, err := c.toFn(x)
		if err != nil {
			t.Fatalf("converting a value of the same measure: %v", err)
		}

		assertMeasurement(t, same, 2, c.unitA)

		parsed, err := units.ParseString("2 " + c.unitA)
		if err != nil {
			t.Fatal(err)
		}

		narrowed, err := c.toFn(parsed)
		if err != nil {
			t.Fatal(err)
		}

		assertMeasurement(t, narrowed.Add(y), 2+1/c.ratio, c.unitA)

		foreign, err := units.NewMeasurement(1, c.foreignUnit)
		if err != nil {
			t.Fatal(err)
		}

		if _, err := c.toFn(foreign); !errors.Is(err, units.ErrIncompatibleMeasures) {
			t.Errorf("converting a %s measurement: %v, want ErrIncompatibleMeasures", foreign.Type(), err)
		}

		if _, err := c.toFn(nil); !errors.Is(err, units.ErrIncompatibleMeasures) {
			t.Errorf("converting nil: %v, want ErrIncompatibleMeasures", err)
		}
	})

	t.Run("StringTo", func(t *testing.T) {
		got, err := c.parseFn("2 " + c.unitA)
		if err != nil {
			t.Fatal(err)
		}

		assertMeasurement(t, got, 2, c.unitA)

		for _, s := range []string{"2 " + c.foreignUnit, "2" + c.unitA, "two " + c.unitA} {
			if _, err := c.parseFn(s); err == nil {
				t.Errorf("parsing %q succeeded", s)
			}
		}
	})
}
