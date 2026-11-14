package units_test

import (
	"cmp"
	"encoding"
	"errors"
	"fmt"
	"iter"
	"math"
	"strings"
	"testing"

	units "github.com/asaphin/go-physics-units"
	"github.com/asaphin/go-physics-units/acceleration"
	"github.com/asaphin/go-physics-units/angle"
	"github.com/asaphin/go-physics-units/angularvelocity"
	"github.com/asaphin/go-physics-units/area"
	"github.com/asaphin/go-physics-units/charge"
	"github.com/asaphin/go-physics-units/current"
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
	"github.com/asaphin/go-physics-units/resistance"
	"github.com/asaphin/go-physics-units/time"
	"github.com/asaphin/go-physics-units/velocity"
	"github.com/asaphin/go-physics-units/voltage"
	"github.com/asaphin/go-physics-units/volume"
)

// unitType is the API shared by the unit types: those of the unit packages
// and Unit[M].
type unitType interface {
	comparable
	fmt.Stringer
	Factor() float64
}

// predefinedUnit is the API of the unit types of the unit packages.
type predefinedUnit interface {
	unitType
	~uint8
	Valid() bool
}

// typedMeasure is the API shared by the predefined measure types.
type typedMeasure[T any, U unitType] interface {
	comparable
	units.Measurement
	encoding.TextMarshaler
	encoding.TextAppender
	units.Arithmetics[T]
	units.UnitConverter[T, U]
	units.Comparer[T]
	Unit() U
}

type typedCase[T typedMeasure[T, U], U unitType] struct {
	measureType   units.MeasureType
	newFn         func(value float64, unit U) T
	toFn          func(m units.Measurement) (T, error)
	parseFn       func(s string) (T, error)
	parseUnit     func(symbol string) (U, bool)
	allUnits      func() iter.Seq[U]
	rate          func(from, to U) float64
	unitA         U
	unitB         U
	ratio         float64 // 1 unitA = ratio unitB
	baseUnit      U
	baseRatio     float64 // 1 unitA = baseRatio baseUnit
	foreignSymbol string  // a unit symbol of another measure

	// parseGeneric parses a Measurement, ParseString if nil.
	parseGeneric func(s string) (units.Measurement, error)
}

// One case per predefined measure, shared by TestTypedMeasures and
// BenchmarkMeasures. A new measure needs a case here.
var (
	distanceCase = typedCase[units.Distance, distance.Unit]{
		measureType: units.MeasureDistance,
		newFn:       units.NewDistance, toFn: units.ToDistance, parseFn: units.StringToDistance,
		parseUnit: distance.ParseUnit, allUnits: distance.Units, rate: distance.Rate,
		unitA: distance.Mile, unitB: distance.Kilometer, ratio: 1.609344,
		baseUnit: distance.BaseUnit, baseRatio: 1609.344, foreignSymbol: mass.Kilogram.String(),
	}
	areaCase = typedCase[units.Area, area.Unit]{
		measureType: units.MeasureArea,
		newFn:       units.NewArea, toFn: units.ToArea, parseFn: units.StringToArea,
		parseUnit: area.ParseUnit, allUnits: area.Units, rate: area.Rate,
		unitA: area.SquareFoot, unitB: area.SquareInch, ratio: 144,
		baseUnit: area.BaseUnit, baseRatio: 0.09290304, foreignSymbol: time.Second.String(),
	}
	massCase = typedCase[units.Mass, mass.Unit]{
		measureType: units.MeasureMass,
		newFn:       units.NewMass, toFn: units.ToMass, parseFn: units.StringToMass,
		parseUnit: mass.ParseUnit, allUnits: mass.Units, rate: mass.Rate,
		unitA: mass.Pound, unitB: mass.Ounce, ratio: 16,
		baseUnit: mass.BaseUnit, baseRatio: 0.45359237, foreignSymbol: distance.Meter.String(),
	}
	timeCase = typedCase[units.Time, time.Unit]{
		measureType: units.MeasureTime,
		newFn:       units.NewTime, toFn: units.ToTime, parseFn: units.StringToTime,
		parseUnit: time.ParseUnit, allUnits: time.Units, rate: time.Rate,
		unitA: time.Hour, unitB: time.Minute, ratio: 60,
		baseUnit: time.BaseUnit, baseRatio: 3600, foreignSymbol: mass.Kilogram.String(),
	}
	velocityCase = typedCase[units.Velocity, velocity.Unit]{
		measureType: units.MeasureVelocity,
		newFn:       units.NewVelocity, toFn: units.ToVelocity, parseFn: units.StringToVelocity,
		parseUnit: velocity.ParseUnit, allUnits: velocity.Units, rate: velocity.Rate,
		unitA: velocity.Knot, unitB: velocity.KilometerPerHour, ratio: 1.852,
		baseUnit: velocity.BaseUnit, baseRatio: 1852.0 / 3600, foreignSymbol: distance.Meter.String(),
	}
	volumeCase = typedCase[units.Volume, volume.Unit]{
		measureType: units.MeasureVolume,
		newFn:       units.NewVolume, toFn: units.ToVolume, parseFn: units.StringToVolume,
		parseUnit: volume.ParseUnit, allUnits: volume.Units, rate: volume.Rate,
		unitA: volume.Gallon, unitB: volume.Liter, ratio: 3.785411784,
		baseUnit: volume.BaseUnit, baseRatio: 0.003785411784, foreignSymbol: mass.Kilogram.String(),
	}
	accelerationCase = typedCase[units.Acceleration, acceleration.Unit]{
		measureType: units.MeasureAcceleration,
		newFn:       units.NewAcceleration, toFn: units.ToAcceleration, parseFn: units.StringToAcceleration,
		parseUnit: acceleration.ParseUnit, allUnits: acceleration.Units, rate: acceleration.Rate,
		unitA: acceleration.StandardGravity, unitB: acceleration.FootPerSecondSquared, ratio: 9.80665 / 0.3048,
		baseUnit: acceleration.BaseUnit, baseRatio: 9.80665, foreignSymbol: velocity.MeterPerSecond.String(),
	}
	forceCase = typedCase[units.Force, force.Unit]{
		measureType: units.MeasureForce,
		newFn:       units.NewForce, toFn: units.ToForce, parseFn: units.StringToForce,
		parseUnit: force.ParseUnit, allUnits: force.Units, rate: force.Rate,
		unitA: force.KilogramForce, unitB: force.Newton, ratio: 9.80665,
		baseUnit: force.BaseUnit, baseRatio: 9.80665, foreignSymbol: mass.Kilogram.String(),
	}
	energyCase = typedCase[units.Energy, energy.Unit]{
		measureType: units.MeasureEnergy,
		newFn:       units.NewEnergy, toFn: units.ToEnergy, parseFn: units.StringToEnergy,
		parseUnit: energy.ParseUnit, allUnits: energy.Units, rate: energy.Rate,
		unitA: energy.KilowattHour, unitB: energy.Megajoule, ratio: 3.6,
		baseUnit: energy.BaseUnit, baseRatio: 3.6e6, foreignSymbol: power.Watt.String(),
	}
	powerCase = typedCase[units.Power, power.Unit]{
		measureType: units.MeasurePower,
		newFn:       units.NewPower, toFn: units.ToPower, parseFn: units.StringToPower,
		parseUnit: power.ParseUnit, allUnits: power.Units, rate: power.Rate,
		unitA: power.Horsepower, unitB: power.Kilowatt, ratio: 0.74569987158227022,
		baseUnit: power.BaseUnit, baseRatio: 745.69987158227022, foreignSymbol: energy.Joule.String(),
	}
	pressureCase = typedCase[units.Pressure, pressure.Unit]{
		measureType: units.MeasurePressure,
		newFn:       units.NewPressure, toFn: units.ToPressure, parseFn: units.StringToPressure,
		parseUnit: pressure.ParseUnit, allUnits: pressure.Units, rate: pressure.Rate,
		unitA: pressure.Atmosphere, unitB: pressure.Kilopascal, ratio: 101.325,
		baseUnit: pressure.BaseUnit, baseRatio: 101325, foreignSymbol: force.Newton.String(),
	}
	frequencyCase = typedCase[units.Frequency, frequency.Unit]{
		measureType: units.MeasureFrequency,
		newFn:       units.NewFrequency, toFn: units.ToFrequency, parseFn: units.StringToFrequency,
		parseUnit: frequency.ParseUnit, allUnits: frequency.Units, rate: frequency.Rate,
		unitA: frequency.Kilohertz, unitB: frequency.RevolutionPerMinute, ratio: 60000,
		baseUnit: frequency.BaseUnit, baseRatio: 1000, foreignSymbol: time.Second.String(),
	}
	angleCase = typedCase[units.Angle, angle.Unit]{
		measureType: units.MeasureAngle,
		newFn:       units.NewAngle, toFn: units.ToAngle, parseFn: units.StringToAngle,
		parseUnit: angle.ParseUnit, allUnits: angle.Units, rate: angle.Rate,
		unitA: angle.Revolution, unitB: angle.Degree, ratio: 360,
		baseUnit: angle.BaseUnit, baseRatio: 2 * math.Pi, foreignSymbol: time.Second.String(),
	}
	dataSizeCase = typedCase[units.DataSize, datasize.Unit]{
		measureType: units.MeasureDataSize,
		newFn:       units.NewDataSize, toFn: units.ToDataSize, parseFn: units.StringToDataSize,
		parseUnit: datasize.ParseUnit, allUnits: datasize.Units, rate: datasize.Rate,
		unitA: datasize.Gibibyte, unitB: datasize.Megabyte, ratio: 1073.741824,
		baseUnit: datasize.BaseUnit, baseRatio: 1 << 30, foreignSymbol: time.Second.String(),
	}
	densityCase = typedCase[units.Density, density.Unit]{
		measureType: units.MeasureDensity,
		newFn:       units.NewDensity, toFn: units.ToDensity, parseFn: units.StringToDensity,
		parseUnit: density.ParseUnit, allUnits: density.Units, rate: density.Rate,
		unitA: density.GramPerCubicCentimeter, unitB: density.PoundPerCubicFoot, ratio: 62.42796057614462,
		baseUnit: density.BaseUnit, baseRatio: 1000, foreignSymbol: mass.Kilogram.String(),
	}
	dataRateCase = typedCase[units.DataRate, datarate.Unit]{
		measureType: units.MeasureDataRate,
		newFn:       units.NewDataRate, toFn: units.ToDataRate, parseFn: units.StringToDataRate,
		parseUnit: datarate.ParseUnit, allUnits: datarate.Units, rate: datarate.Rate,
		unitA: datarate.GigabitPerSecond, unitB: datarate.MegabytePerSecond, ratio: 125,
		baseUnit: datarate.BaseUnit, baseRatio: 1.25e8, foreignSymbol: datasize.Byte.String(),
	}
	currentCase = typedCase[units.Current, current.Unit]{
		measureType: units.MeasureCurrent,
		newFn:       units.NewCurrent, toFn: units.ToCurrent, parseFn: units.StringToCurrent,
		parseUnit: current.ParseUnit, allUnits: current.Units, rate: current.Rate,
		unitA: current.Kiloampere, unitB: current.Milliampere, ratio: 1e6,
		baseUnit: current.BaseUnit, baseRatio: 1000, foreignSymbol: voltage.Volt.String(),
	}
	voltageCase = typedCase[units.Voltage, voltage.Unit]{
		measureType: units.MeasureVoltage,
		newFn:       units.NewVoltage, toFn: units.ToVoltage, parseFn: units.StringToVoltage,
		parseUnit: voltage.ParseUnit, allUnits: voltage.Units, rate: voltage.Rate,
		unitA: voltage.Kilovolt, unitB: voltage.Millivolt, ratio: 1e6,
		baseUnit: voltage.BaseUnit, baseRatio: 1000, foreignSymbol: current.Ampere.String(),
	}
	chargeCase = typedCase[units.Charge, charge.Unit]{
		measureType: units.MeasureCharge,
		newFn:       units.NewCharge, toFn: units.ToCharge, parseFn: units.StringToCharge,
		parseUnit: charge.ParseUnit, allUnits: charge.Units, rate: charge.Rate,
		unitA: charge.AmpereHour, unitB: charge.MilliampereHour, ratio: 1000,
		baseUnit: charge.BaseUnit, baseRatio: 3600, foreignSymbol: current.Ampere.String(),
	}
	resistanceCase = typedCase[units.Resistance, resistance.Unit]{
		measureType: units.MeasureResistance,
		newFn:       units.NewResistance, toFn: units.ToResistance, parseFn: units.StringToResistance,
		parseUnit: resistance.ParseUnit, allUnits: resistance.Units, rate: resistance.Rate,
		unitA: resistance.Megaohm, unitB: resistance.Kiloohm, ratio: 1000,
		baseUnit: resistance.BaseUnit, baseRatio: 1e6, foreignSymbol: voltage.Volt.String(),
	}
	angularVelocityCase = typedCase[units.AngularVelocity, angularvelocity.Unit]{
		measureType: units.MeasureAngularVelocity,
		newFn:       units.NewAngularVelocity, toFn: units.ToAngularVelocity, parseFn: units.StringToAngularVelocity,
		parseUnit: angularvelocity.ParseUnit, allUnits: angularvelocity.Units, rate: angularvelocity.Rate,
		unitA: angularvelocity.RevolutionPerSecond, unitB: angularvelocity.RevolutionPerMinute, ratio: 60,
		baseUnit: angularvelocity.BaseUnit, baseRatio: 2 * math.Pi, foreignSymbol: frequency.Hertz.String(),
	}
)

func TestTypedMeasures(t *testing.T) {
	t.Run("distance", func(t *testing.T) { testTypedMeasure(t, distanceCase); testUnitPackage(t, distanceCase) })
	t.Run("area", func(t *testing.T) { testTypedMeasure(t, areaCase); testUnitPackage(t, areaCase) })
	t.Run("mass", func(t *testing.T) { testTypedMeasure(t, massCase); testUnitPackage(t, massCase) })
	t.Run("time", func(t *testing.T) { testTypedMeasure(t, timeCase); testUnitPackage(t, timeCase) })
	t.Run("velocity", func(t *testing.T) { testTypedMeasure(t, velocityCase); testUnitPackage(t, velocityCase) })
	t.Run("volume", func(t *testing.T) { testTypedMeasure(t, volumeCase); testUnitPackage(t, volumeCase) })
	t.Run(
		"acceleration",
		func(t *testing.T) { testTypedMeasure(t, accelerationCase); testUnitPackage(t, accelerationCase) },
	)
	t.Run("force", func(t *testing.T) { testTypedMeasure(t, forceCase); testUnitPackage(t, forceCase) })
	t.Run("energy", func(t *testing.T) { testTypedMeasure(t, energyCase); testUnitPackage(t, energyCase) })
	t.Run("power", func(t *testing.T) { testTypedMeasure(t, powerCase); testUnitPackage(t, powerCase) })
	t.Run("pressure", func(t *testing.T) { testTypedMeasure(t, pressureCase); testUnitPackage(t, pressureCase) })
	t.Run("frequency", func(t *testing.T) { testTypedMeasure(t, frequencyCase); testUnitPackage(t, frequencyCase) })
	t.Run("angle", func(t *testing.T) { testTypedMeasure(t, angleCase); testUnitPackage(t, angleCase) })
	t.Run("data size", func(t *testing.T) { testTypedMeasure(t, dataSizeCase); testUnitPackage(t, dataSizeCase) })
	t.Run("density", func(t *testing.T) { testTypedMeasure(t, densityCase); testUnitPackage(t, densityCase) })
	t.Run("data rate", func(t *testing.T) { testTypedMeasure(t, dataRateCase); testUnitPackage(t, dataRateCase) })
	t.Run("electric current", func(t *testing.T) { testTypedMeasure(t, currentCase); testUnitPackage(t, currentCase) })
	t.Run("voltage", func(t *testing.T) { testTypedMeasure(t, voltageCase); testUnitPackage(t, voltageCase) })
	t.Run("electric charge", func(t *testing.T) { testTypedMeasure(t, chargeCase); testUnitPackage(t, chargeCase) })
	t.Run("resistance", func(t *testing.T) { testTypedMeasure(t, resistanceCase); testUnitPackage(t, resistanceCase) })
	t.Run(
		"angular velocity",
		func(t *testing.T) { testTypedMeasure(t, angularVelocityCase); testUnitPackage(t, angularVelocityCase) },
	)
}

func testTypedMeasure[T typedMeasure[T, U], U unitType](t *testing.T, c typedCase[T, U]) {
	t.Helper()

	x := c.newFn(2, c.unitA)
	y := c.newFn(1, c.unitB)
	symbolA, symbolB := c.unitA.String(), c.unitB.String()

	t.Run("accessors", func(t *testing.T) {
		if x.Type() != c.measureType || x.Unit() != c.unitA || x.Symbol() != symbolA || x.Value() != 2 {
			t.Errorf("x = %v of type %s, unit %v, symbol %s; want 2 %s of type %s",
				x.Value(), x.Type(), x.Unit(), x.Symbol(), symbolA, c.measureType)
		}

		if want := "2 " + symbolA; x.String() != want {
			t.Errorf("String() = %q, want %q", x.String(), want)
		}
	})

	t.Run("zero value", func(t *testing.T) {
		var zero T
		if zero.Value() != 0 || zero.Unit() != c.baseUnit || zero.Type() != c.measureType {
			t.Errorf("the zero value is %v of type %s, want 0 %v", zero, zero.Type(), c.baseUnit)
		}
	})

	t.Run("ConvertTo", func(t *testing.T) {
		assertMeasurement(t, x.ConvertTo(c.unitB), 2*c.ratio, symbolB)
		if must, ok := any(x).(interface{ MustConvertTo(unit U) T }); ok {
			assertMeasurement(t, must.MustConvertTo(c.unitB), 2*c.ratio, symbolB)
		}
		assertMeasurement(t, x.ConvertToBaseUnits(), 2*c.baseRatio, c.baseUnit.String())

		if got := x.ConvertTo(c.unitA); got != x {
			t.Errorf("converting to the same unit gave %v, want %v exactly", got, x)
		}
	})

	t.Run("Add", func(t *testing.T) {
		assertMeasurement(t, x.Add(y), 2+1/c.ratio, symbolA)
	})

	t.Run("Sub", func(t *testing.T) {
		assertMeasurement(t, x.Sub(y), 2-1/c.ratio, symbolA)
	})

	t.Run("Mul", func(t *testing.T) {
		assertMeasurement(t, x.Mul(3), 6, symbolA)
	})

	t.Run("Div", func(t *testing.T) {
		got, err := x.Div(4)
		if err != nil {
			t.Fatal(err)
		}

		assertMeasurement(t, got, 0.5, symbolA)

		var zero T
		if bad, err := x.Div(0); !errors.Is(err, units.ErrZeroDivision) || bad != zero {
			t.Errorf("Div(0) = %v, %v; want the zero value and ErrZeroDivision", bad, err)
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

		converted := x.ConvertTo(c.unitB)
		if x.Equal(converted) != converted.Equal(x) {
			t.Errorf("Equal isn't symmetric for %v and %v", x, converted)
		}
	})

	t.Run("EqualWithinRel", func(t *testing.T) {
		converted := x.ConvertTo(c.unitB)
		if !x.EqualWithinRel(converted, 1e-12) || !converted.EqualWithinRel(x, 1e-12) {
			t.Errorf("%v and %v aren't equal within 1e-12", x, converted)
		}

		near := x.Mul(1 + 1e-6)
		if x.EqualWithinRel(near, 1e-9) || !x.EqualWithinRel(near, 1e-5) {
			t.Errorf("%v and %v: within 1e-9 = %v, within 1e-5 = %v; want false, true",
				x, near, x.EqualWithinRel(near, 1e-9), x.EqualWithinRel(near, 1e-5))
		}
	})

	t.Run("Format", func(t *testing.T) {
		got := fmt.Sprintf("%.1f|%v|%6.2e", x, x, x)
		if want := "2.0 " + symbolA + "|2 " + symbolA + "|2.00e+00 " + symbolA; got != want {
			t.Errorf("Sprintf = %q, want %q", got, want)
		}
	})

	t.Run("text", func(t *testing.T) {
		testText(t, x, c.parseFn)
	})

	t.Run("operations don't change the receiver", func(t *testing.T) {
		assertMeasurement(t, x, 2, symbolA)
	})

	t.Run("generic methods", func(t *testing.T) {
		testGenericMethods(t, x, y, symbolB, c.ratio, c.foreignSymbol)
	})

	t.Run("To", func(t *testing.T) {
		same, err := c.toFn(x)
		if err != nil || same != x {
			t.Fatalf("converting a value of the same measure: %v, %v", same, err)
		}

		parse := c.parseGeneric
		if parse == nil {
			parse = units.ParseString
		}

		parsed, err := parse("2 " + symbolA)
		if err != nil {
			t.Fatal(err)
		}

		if narrowed, ok := parsed.(T); !ok || narrowed != x {
			t.Errorf("ParseString(2 %s) holds %T %v, want %v", symbolA, parsed, parsed, x)
		}

		foreign, err := units.NewMeasurement(1, c.foreignSymbol)
		if err != nil {
			t.Fatal(err)
		}

		var zero T
		if bad, err := c.toFn(foreign); !errors.Is(err, units.ErrIncompatibleMeasures) || bad != zero {
			t.Errorf("converting a %s measurement: %v, %v; want the zero value and ErrIncompatibleMeasures",
				foreign.Type(), bad, err)
		}

		if _, err := c.toFn(nil); !errors.Is(err, units.ErrIncompatibleMeasures) {
			t.Errorf("converting nil: %v, want ErrIncompatibleMeasures", err)
		}
	})

	t.Run("StringTo", func(t *testing.T) {
		got, err := c.parseFn("2 " + symbolA)
		if err != nil || got != x {
			t.Fatalf("StringTo(2 %s) = %v, %v; want %v", symbolA, got, err, x)
		}

		var zero T

		for s, want := range map[string]error{
			"2 " + c.foreignSymbol: units.ErrUnknownUnit,
			"2" + symbolA:          units.ErrInvalidFormat,
			"two " + symbolA:       units.ErrInvalidFormat,
		} {
			if bad, err := c.parseFn(s); !errors.Is(err, want) || bad != zero {
				t.Errorf("StringTo(%q) = %v, %v; want the zero value and %v", s, bad, err, want)
			}
		}
	})
}

// testText checks the text encoding of x, which must be finite: MarshalText
// and AppendText give String, and UnmarshalText parses it back as parse does.
func testText[T typedMeasure[T, U], U unitType](t *testing.T, x T, parse func(string) (T, error)) {
	t.Helper()

	text, err := x.MarshalText()
	if err != nil || string(text) != x.String() {
		t.Fatalf("MarshalText = %q, %v; want %q", text, err, x.String())
	}

	appended, err := x.AppendText([]byte("x = "))
	if err != nil || string(appended) != "x = "+x.String() {
		t.Errorf("AppendText = %q, %v", appended, err)
	}

	var back T

	unmarshaler, ok := any(&back).(encoding.TextUnmarshaler)
	if !ok {
		t.Fatalf("*%T isn't an encoding.TextUnmarshaler", back)
	}

	if err := unmarshaler.UnmarshalText(text); err != nil || back != x {
		t.Errorf("UnmarshalText(%q) = %v, %v; want %v", text, back, err, x)
	}

	before := back
	if err := unmarshaler.UnmarshalText([]byte("2 furlong")); err == nil || back != before {
		t.Errorf("UnmarshalText of an unknown unit: %v, and the value became %v", err, back)
	}

	if _, err := parse("NaN " + x.Symbol()); err == nil {
		t.Error("parsing NaN succeeded")
	}

	nan := x.Mul(math.NaN())
	if text, err := nan.MarshalText(); !errors.Is(err, units.ErrOutOfRange) {
		t.Errorf("MarshalText of NaN = %q, %v; want ErrOutOfRange", text, err)
	}
}

// testGenericMethods checks the Measurement methods of x, which is 2 in some
// unit, and y, which is 1 in unitB = 1/ratio of x's unit.
func testGenericMethods(t *testing.T, x, y units.Measurement, unitB string, ratio float64, foreignSymbol string) {
	t.Helper()

	converted, err := x.ConvertToMeasurement(unitB)
	if err != nil {
		t.Fatal(err)
	}

	assertMeasurement(t, converted, 2*ratio, unitB)

	if bad, convErr := x.ConvertToMeasurement(foreignSymbol); !errors.Is(convErr, units.ErrUnknownUnit) || bad != nil {
		t.Errorf("ConvertToMeasurement(%s) = %v, %v; want nil and ErrUnknownUnit", foreignSymbol, bad, convErr)
	}

	sum, err := x.AddMeasurement(y)
	if err != nil {
		t.Fatal(err)
	}

	assertMeasurement(t, sum, 2+1/ratio, x.Symbol())

	diff, err := x.SubMeasurement(y)
	if err != nil {
		t.Fatal(err)
	}

	assertMeasurement(t, diff, 2-1/ratio, x.Symbol())
	assertMeasurement(t, x.MulMeasurement(3), 6, x.Symbol())

	quot, err := x.DivMeasurement(4)
	if err != nil {
		t.Fatal(err)
	}

	assertMeasurement(t, quot, 0.5, x.Symbol())

	if bad, divErr := x.DivMeasurement(0); !errors.Is(divErr, units.ErrZeroDivision) || bad != nil {
		t.Errorf("DivMeasurement(0) = %v, %v; want nil and ErrZeroDivision", bad, divErr)
	}

	if c, cmpErr := x.CompareMeasurement(x.MulMeasurement(1)); cmpErr != nil || c != 0 {
		t.Errorf("CompareMeasurement with itself = %d, %v; want 0", c, cmpErr)
	}

	if c, cmpErr := x.CompareMeasurement(y); cmpErr != nil || c != cmp.Compare(2, 1/ratio) {
		t.Errorf("CompareMeasurement = %d, %v; want %d", c, cmpErr, cmp.Compare(2, 1/ratio))
	}

	foreign, err := units.NewMeasurement(1, foreignSymbol)
	if err != nil {
		t.Fatal(err)
	}

	for name, op := range map[string]func(units.Measurement) (any, error){
		"AddMeasurement":     func(m units.Measurement) (any, error) { return x.AddMeasurement(m) },
		"SubMeasurement":     func(m units.Measurement) (any, error) { return x.SubMeasurement(m) },
		"CompareMeasurement": func(m units.Measurement) (any, error) { return x.CompareMeasurement(m) },
	} {
		for _, m := range []units.Measurement{foreign, nil} {
			if _, err := op(m); !errors.Is(err, units.ErrIncompatibleMeasures) {
				t.Errorf("%s(%v): %v, want ErrIncompatibleMeasures", name, m, err)
			}
		}
	}
}

// testUnitPackage checks the functions of a unit package: every unit has a
// symbol and a factor, ParseUnit inverts String, and the rates agree with the
// factors.
func testUnitPackage[T typedMeasure[T, U], U predefinedUnit](t *testing.T, c typedCase[T, U]) {
	t.Helper()

	invalid := U(200)
	if invalid.Valid() || !c.unitA.Valid() {
		t.Errorf("Valid: %v for U(200), %v for %v", invalid.Valid(), c.unitA.Valid(), c.unitA)
	}

	if s := invalid.String(); !strings.HasSuffix(s, ".Unit(200)") {
		t.Errorf("U(200).String() = %q, want <package>.Unit(200)", s)
	}

	assertPanics(t, func() { c.newFn(1, c.unitA).ConvertTo(invalid) })

	if text, err := c.newFn(1, invalid).MarshalText(); !errors.Is(err, units.ErrUnknownUnit) {
		t.Errorf("MarshalText in an invalid unit = %q, %v; want ErrUnknownUnit", text, err)
	}

	n := 0
	for u := range c.allUnits() {
		if n == 0 && u != c.baseUnit {
			t.Errorf("Units starts with %v, want the base unit %v", u, c.baseUnit)
		}

		if U(n) != u {
			t.Errorf("unit %d is %v: Units must yield the constants in order", n, u)
		}

		n++

		if f := u.Factor(); !(f > 0) || math.IsInf(f, 0) {
			t.Errorf("%v has the factor %v", u, f)
		}

		if parsed, ok := c.parseUnit(u.String()); !ok || parsed != u {
			t.Errorf("ParseUnit(%q) = %v, %v; want %v", u.String(), parsed, ok, u)
		}

		if r := c.rate(u, u); r != 1 {
			t.Errorf("Rate(%v, %v) = %v, want exactly 1", u, u, r)
		}

		if r := c.rate(u, c.baseUnit); !almostEqual(r, u.Factor(), 1e-15) {
			t.Errorf("Rate(%v, base) = %v, want its factor %v", u, r, u.Factor())
		}
	}

	if n < 2 {
		t.Errorf("only %d units", n)
	}

	for range c.allUnits() {
		break // the iterator stops when asked to
	}

	if _, ok := c.parseUnit(c.foreignSymbol); ok {
		t.Errorf("ParseUnit(%s) found a unit of another measure", c.foreignSymbol)
	}

	if r := c.rate(c.unitA, c.unitB); !almostEqual(r, c.ratio, 1e-12) {
		t.Errorf("Rate(%v, %v) = %v, want %v", c.unitA, c.unitB, r, c.ratio)
	}
}
