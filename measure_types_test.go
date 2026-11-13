package units_test

import (
	"cmp"
	"errors"
	"math"
	"testing"

	units "github.com/asaphin/go-physics-units"
	"github.com/asaphin/go-physics-units/acceleration"
	"github.com/asaphin/go-physics-units/angle"
	"github.com/asaphin/go-physics-units/area"
	"github.com/asaphin/go-physics-units/datarate"
	"github.com/asaphin/go-physics-units/datasize"
	"github.com/asaphin/go-physics-units/density"
	"github.com/asaphin/go-physics-units/distance"
	"github.com/asaphin/go-physics-units/energy"
	"github.com/asaphin/go-physics-units/force"
	"github.com/asaphin/go-physics-units/frequency"
	"github.com/asaphin/go-physics-units/mass"
	"github.com/asaphin/go-physics-units/power"
	"github.com/asaphin/go-physics-units/pressure"
	"github.com/asaphin/go-physics-units/time"
	"github.com/asaphin/go-physics-units/velocity"
	"github.com/asaphin/go-physics-units/volume"
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
	Compare(other T) int
	Equal(other T) bool
	EqualWithinRel(other T, tol float64) bool
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
	volumeCase = typedCase[units.Volume]{
		measureType: units.MeasureVolume,
		newFn:       units.NewVolume, toFn: units.ToVolume, parseFn: units.StringToVolume,
		unitA: volume.Gallon, unitB: volume.Liter, ratio: 3.785411784,
		baseUnit: volume.BaseUnit, baseRatio: 0.003785411784, foreignUnit: mass.Kilogram,
	}
	accelerationCase = typedCase[units.Acceleration]{
		measureType: units.MeasureAcceleration,
		newFn:       units.NewAcceleration, toFn: units.ToAcceleration, parseFn: units.StringToAcceleration,
		unitA: acceleration.StandardGravity, unitB: acceleration.FootPerSecondSquared, ratio: 9.80665 / 0.3048,
		baseUnit: acceleration.BaseUnit, baseRatio: 9.80665, foreignUnit: velocity.MeterPerSecond,
	}
	forceCase = typedCase[units.Force]{
		measureType: units.MeasureForce,
		newFn:       units.NewForce, toFn: units.ToForce, parseFn: units.StringToForce,
		unitA: force.KilogramForce, unitB: force.Newton, ratio: 9.80665,
		baseUnit: force.BaseUnit, baseRatio: 9.80665, foreignUnit: mass.Kilogram,
	}
	energyCase = typedCase[units.Energy]{
		measureType: units.MeasureEnergy,
		newFn:       units.NewEnergy, toFn: units.ToEnergy, parseFn: units.StringToEnergy,
		unitA: energy.KilowattHour, unitB: energy.Megajoule, ratio: 3.6,
		baseUnit: energy.BaseUnit, baseRatio: 3.6e6, foreignUnit: power.Watt,
	}
	powerCase = typedCase[units.Power]{
		measureType: units.MeasurePower,
		newFn:       units.NewPower, toFn: units.ToPower, parseFn: units.StringToPower,
		unitA: power.Horsepower, unitB: power.Kilowatt, ratio: 0.74569987158227022,
		baseUnit: power.BaseUnit, baseRatio: 745.69987158227022, foreignUnit: energy.Joule,
	}
	pressureCase = typedCase[units.Pressure]{
		measureType: units.MeasurePressure,
		newFn:       units.NewPressure, toFn: units.ToPressure, parseFn: units.StringToPressure,
		unitA: pressure.Atmosphere, unitB: pressure.Kilopascal, ratio: 101.325,
		baseUnit: pressure.BaseUnit, baseRatio: 101325, foreignUnit: force.Newton,
	}
	frequencyCase = typedCase[units.Frequency]{
		measureType: units.MeasureFrequency,
		newFn:       units.NewFrequency, toFn: units.ToFrequency, parseFn: units.StringToFrequency,
		unitA: frequency.Kilohertz, unitB: frequency.RevolutionPerMinute, ratio: 60000,
		baseUnit: frequency.BaseUnit, baseRatio: 1000, foreignUnit: time.Second,
	}
	angleCase = typedCase[units.Angle]{
		measureType: units.MeasureAngle,
		newFn:       units.NewAngle, toFn: units.ToAngle, parseFn: units.StringToAngle,
		unitA: angle.Revolution, unitB: angle.Degree, ratio: 360,
		baseUnit: angle.BaseUnit, baseRatio: 2 * math.Pi, foreignUnit: time.Second,
	}
	dataSizeCase = typedCase[units.DataSize]{
		measureType: units.MeasureDataSize,
		newFn:       units.NewDataSize, toFn: units.ToDataSize, parseFn: units.StringToDataSize,
		unitA: datasize.Gibibyte, unitB: datasize.Megabyte, ratio: 1073.741824,
		baseUnit: datasize.BaseUnit, baseRatio: 1 << 30, foreignUnit: time.Second,
	}
	densityCase = typedCase[units.Density]{
		measureType: units.MeasureDensity,
		newFn:       units.NewDensity, toFn: units.ToDensity, parseFn: units.StringToDensity,
		unitA: density.GramPerCubicCentimeter, unitB: density.PoundPerCubicFoot, ratio: 62.42796057614462,
		baseUnit: density.BaseUnit, baseRatio: 1000, foreignUnit: mass.Kilogram,
	}
	dataRateCase = typedCase[units.DataRate]{
		measureType: units.MeasureDataRate,
		newFn:       units.NewDataRate, toFn: units.ToDataRate, parseFn: units.StringToDataRate,
		unitA: datarate.GigabitPerSecond, unitB: datarate.MegabytePerSecond, ratio: 125,
		baseUnit: datarate.BaseUnit, baseRatio: 1.25e8, foreignUnit: datasize.Byte,
	}
)

func TestTypedMeasures(t *testing.T) {
	t.Run("distance", func(t *testing.T) { testTypedMeasure(t, distanceCase) })
	t.Run("area", func(t *testing.T) { testTypedMeasure(t, areaCase) })
	t.Run("mass", func(t *testing.T) { testTypedMeasure(t, massCase) })
	t.Run("time", func(t *testing.T) { testTypedMeasure(t, timeCase) })
	t.Run("velocity", func(t *testing.T) { testTypedMeasure(t, velocityCase) })
	t.Run("volume", func(t *testing.T) { testTypedMeasure(t, volumeCase) })
	t.Run("acceleration", func(t *testing.T) { testTypedMeasure(t, accelerationCase) })
	t.Run("force", func(t *testing.T) { testTypedMeasure(t, forceCase) })
	t.Run("energy", func(t *testing.T) { testTypedMeasure(t, energyCase) })
	t.Run("power", func(t *testing.T) { testTypedMeasure(t, powerCase) })
	t.Run("pressure", func(t *testing.T) { testTypedMeasure(t, pressureCase) })
	t.Run("frequency", func(t *testing.T) { testTypedMeasure(t, frequencyCase) })
	t.Run("angle", func(t *testing.T) { testTypedMeasure(t, angleCase) })
	t.Run("data size", func(t *testing.T) { testTypedMeasure(t, dataSizeCase) })
	t.Run("density", func(t *testing.T) { testTypedMeasure(t, densityCase) })
	t.Run("data rate", func(t *testing.T) { testTypedMeasure(t, dataRateCase) })
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

	t.Run("Compare", func(t *testing.T) {
		// x is 2 unitA; y is 1 unitB, which is 1/ratio unitA.
		want := cmp.Compare(2, 1/c.ratio)
		if got, rev := x.Compare(y), y.Compare(x); got != want || rev != -want {
			t.Errorf("x.Compare(y) = %d, y.Compare(x) = %d; want %d and %d", got, rev, want, -want)
		}

		bigger, same := x.Mul(2), x.Mul(1)
		if x.Compare(bigger) != -1 || bigger.Compare(x) != 1 || x.Compare(same) != 0 {
			t.Errorf("Compare: x vs 2x = %d, 2x vs x = %d, x vs x = %d; want -1, 1, 0",
				x.Compare(bigger), bigger.Compare(x), x.Compare(same))
		}
	})

	t.Run("Equal", func(t *testing.T) {
		if same := x.Mul(1); !x.Equal(same) || x.Equal(y) || y.Equal(x) {
			t.Errorf("Equal: x vs x = %v, x vs y = %v, y vs x = %v; want true, false, false",
				x.Equal(same), x.Equal(y), y.Equal(x))
		}

		converted := x.MustConvertTo(c.unitB)
		if x.Equal(converted) != converted.Equal(x) {
			t.Errorf("Equal isn't symmetric for %v and %v", x, converted)
		}
	})

	t.Run("EqualWithinRel", func(t *testing.T) {
		converted := x.MustConvertTo(c.unitB)
		if !x.EqualWithinRel(converted, 1e-12) || !converted.EqualWithinRel(x, 1e-12) {
			t.Errorf("%v and %v aren't equal within 1e-12", x, converted)
		}

		near := x.Mul(1 + 1e-6)
		if x.EqualWithinRel(near, 1e-9) || !x.EqualWithinRel(near, 1e-5) {
			t.Errorf("%v and %v: within 1e-9 = %v, within 1e-5 = %v; want false, true",
				x, near, x.EqualWithinRel(near, 1e-9), x.EqualWithinRel(near, 1e-5))
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
