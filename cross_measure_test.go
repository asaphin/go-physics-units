package units_test

import (
	"errors"
	"testing"

	units "github.com/asaphin/go-physics-units"
	"github.com/asaphin/go-physics-units/acceleration"
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

// must returns a function that returns v, or fails the test if err isn't nil.
// Call it as must(units.NewDistance(2, distance.Meter))(t): Go passes the
// results of a call on only when they are the only arguments.
func must[T any](v T, err error) func(*testing.T) T {
	return func(t *testing.T) T {
		t.Helper()

		if err != nil {
			t.Fatal(err)
		}

		return v
	}
}

// ok adapts an operation that can't fail to the shape of the table below.
func ok[T units.Measurement](m T) (units.Measurement, error) {
	return m, nil
}

// fallible adapts an operation that can fail, keeping a nil result nil.
func fallible[T units.Measurement](m T, err error) (units.Measurement, error) {
	if err != nil {
		return nil, err
	}

	return m, nil
}

// TestCrossMeasureOperations checks each operation that combines measures
// against an independently known result, in the result's base unit.
func TestCrossMeasureOperations(t *testing.T) {
	m2 := must(units.NewDistance(2, distance.Meter))(t)
	cm10 := must(units.NewDistance(10, distance.Centimeter))(t)
	cm50 := must(units.NewDistance(50, distance.Centimeter))(t)
	sqm2 := must(units.NewArea(2, area.SquareMeter))(t)
	sqin := must(units.NewArea(1, area.SquareInch))(t)
	cubicM := must(units.NewVolume(1, volume.CubicMeter))(t)
	liter := must(units.NewVolume(1, volume.Liter))(t)
	kmh100 := must(units.NewVelocity(100, velocity.KilometerPerHour))(t)
	s2 := must(units.NewTime(2, time.Second))(t)
	s10 := must(units.NewTime(10, time.Second))(t)
	hour := must(units.NewTime(1, time.Hour))(t)
	gn := must(units.NewAcceleration(1, acceleration.StandardGravity))(t)
	kg1 := must(units.NewMass(1, mass.Kilogram))(t)
	kg2 := must(units.NewMass(2, mass.Kilogram))(t)
	lb := must(units.NewMass(1, mass.Pound))(t)
	n10 := must(units.NewForce(10, force.Newton))(t)
	kgf := must(units.NewForce(1, force.KilogramForce))(t)
	lbf := must(units.NewForce(1, force.PoundForce))(t)
	j20 := must(units.NewEnergy(20, energy.Joule))(t)
	kwh := must(units.NewEnergy(1, energy.KilowattHour))(t)
	kw := must(units.NewPower(1, power.Kilowatt))(t)
	atm := must(units.NewPressure(1, pressure.Atmosphere))(t)
	hz50 := must(units.NewFrequency(50, frequency.Hertz))(t)
	rpm60 := must(units.NewFrequency(60, frequency.RevolutionPerMinute))(t)
	water := must(units.NewDensity(1, density.GramPerCubicCentimeter))(t)
	liter2 := must(units.NewVolume(2, volume.Liter))(t)
	gb := must(units.NewDataSize(1, datasize.Gigabyte))(t)
	s8 := must(units.NewTime(8, time.Second))(t)
	mbit100 := must(units.NewDataRate(100, datarate.MegabitPerSecond))(t)
	minute := must(units.NewTime(1, time.Minute))(t)

	cases := []struct {
		name  string
		op    func() (units.Measurement, error)
		value float64
		unit  string
	}{
		{"2 m cubed", func() (units.Measurement, error) { return ok(m2.Cube()) }, 8, volume.BaseUnit},
		{"10 cm cubed is a litre", func() (units.Measurement, error) { return ok(cm10.Cube()) }, 1e-3, volume.BaseUnit},
		{
			"2 m² × 50 cm",
			func() (units.Measurement, error) { return ok(sqm2.MultiplyByDistance(cm50)) },
			1,
			volume.BaseUnit,
		},
		{
			"1 m³ / 2 m²",
			func() (units.Measurement, error) { return fallible(cubicM.DivideByArea(sqm2)) },
			0.5,
			distance.BaseUnit,
		},
		{
			"1 L / 10 cm",
			func() (units.Measurement, error) { return fallible(liter.DivideByDistance(cm10)) },
			0.01,
			area.BaseUnit,
		},
		{
			"100 km/h / 10 s", func() (units.Measurement, error) { return fallible(kmh100.DivideByTime(s10)) },
			100.0 / 36, acceleration.BaseUnit,
		},
		{
			"gₙ × 2 s",
			func() (units.Measurement, error) { return ok(gn.MultiplyByTime(s2)) },
			19.6133,
			velocity.BaseUnit,
		},
		{"gₙ × 1 kg", func() (units.Measurement, error) { return ok(gn.MultiplyByMass(kg1)) }, 9.80665, force.BaseUnit},
		{
			"1 lb × gₙ is a pound-force",
			func() (units.Measurement, error) { return ok(lb.MultiplyByAcceleration(gn)) },
			4.4482216152605,
			force.BaseUnit,
		},
		{
			"10 N / 2 kg",
			func() (units.Measurement, error) { return fallible(n10.DivideByMass(kg2)) },
			5,
			acceleration.BaseUnit,
		},
		{
			"1 kgf / gₙ",
			func() (units.Measurement, error) { return fallible(kgf.DivideByAcceleration(gn)) },
			1,
			mass.BaseUnit,
		},
		{
			"1 lbf / 1 in² is a psi", func() (units.Measurement, error) { return fallible(lbf.DivideByArea(sqin)) },
			6894.757293168361, pressure.BaseUnit,
		},
		{
			"10 N × 2 m",
			func() (units.Measurement, error) { return ok(n10.MultiplyByDistance(m2)) },
			20,
			energy.BaseUnit,
		},
		{
			"1 kWh / 1 h",
			func() (units.Measurement, error) { return fallible(kwh.DivideByTime(hour)) },
			1000,
			power.BaseUnit,
		},
		{
			"20 J / 2 m",
			func() (units.Measurement, error) { return fallible(j20.DivideByDistance(m2)) },
			10,
			force.BaseUnit,
		},
		{
			"1 kW × 1 h",
			func() (units.Measurement, error) { return ok(kw.MultiplyByTime(hour)) },
			3.6e6,
			energy.BaseUnit,
		},
		{
			"1 atm × 2 m²",
			func() (units.Measurement, error) { return ok(atm.MultiplyByArea(sqm2)) },
			202650,
			force.BaseUnit,
		},
		{"period of 50 Hz", func() (units.Measurement, error) { return fallible(hz50.Period()) }, 0.02, time.BaseUnit},
		{"period of 60 rpm", func() (units.Measurement, error) { return fallible(rpm60.Period()) }, 1, time.BaseUnit},
		{
			"1 kg / 1 L",
			func() (units.Measurement, error) { return fallible(kg1.DivideByVolume(liter)) },
			1000,
			density.BaseUnit,
		},
		{
			"1 g/cm³ × 2 L",
			func() (units.Measurement, error) { return ok(water.MultiplyByVolume(liter2)) },
			2,
			mass.BaseUnit,
		},
		{
			"1 kg / 1 g/cm³",
			func() (units.Measurement, error) { return fallible(kg1.DivideByDensity(water)) },
			1e-3,
			volume.BaseUnit,
		},
		{
			"1 GB / 8 s",
			func() (units.Measurement, error) { return fallible(gb.DivideByTime(s8)) },
			1.25e8,
			datarate.BaseUnit,
		},
		{
			"1 GB at 100 Mbit/s", func() (units.Measurement, error) { return fallible(gb.DivideByDataRate(mbit100)) },
			80, time.BaseUnit,
		},
		{
			"100 Mbit/s × 1 min",
			func() (units.Measurement, error) { return ok(mbit100.MultiplyByTime(minute)) },
			7.5e8,
			datasize.BaseUnit,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			m, err := c.op()
			if err != nil {
				t.Fatal(err)
			}

			assertMeasurement(t, m, c.value, c.unit)
		})
	}
}

func TestCrossMeasureOperations_ZeroDivision(t *testing.T) {
	noDistance := must(units.NewDistance(0, distance.Meter))(t)
	noArea := must(units.NewArea(0, area.SquareMeter))(t)
	noTime := must(units.NewTime(0, time.Second))(t)
	noMass := must(units.NewMass(0, mass.Kilogram))(t)
	noAcceleration := must(units.NewAcceleration(0, acceleration.StandardGravity))(t)
	cubicM := must(units.NewVolume(1, volume.CubicMeter))(t)
	speed := must(units.NewVelocity(1, velocity.MeterPerSecond))(t)
	n := must(units.NewForce(1, force.Newton))(t)
	j := must(units.NewEnergy(1, energy.Joule))(t)
	noFrequency := must(units.NewFrequency(0, frequency.Hertz))(t)
	kg := must(units.NewMass(1, mass.Kilogram))(t)
	noVolume := must(units.NewVolume(0, volume.Liter))(t)
	noDensity := must(units.NewDensity(0, density.GramPerLiter))(t)
	gb := must(units.NewDataSize(1, datasize.Gigabyte))(t)
	noRate := must(units.NewDataRate(0, datarate.BitPerSecond))(t)

	cases := map[string]func() (units.Measurement, error){
		"Volume.DivideByArea":        func() (units.Measurement, error) { return fallible(cubicM.DivideByArea(noArea)) },
		"Volume.DivideByDistance":    func() (units.Measurement, error) { return fallible(cubicM.DivideByDistance(noDistance)) },
		"Velocity.DivideByTime":      func() (units.Measurement, error) { return fallible(speed.DivideByTime(noTime)) },
		"Force.DivideByMass":         func() (units.Measurement, error) { return fallible(n.DivideByMass(noMass)) },
		"Force.DivideByAcceleration": func() (units.Measurement, error) { return fallible(n.DivideByAcceleration(noAcceleration)) },
		"Force.DivideByArea":         func() (units.Measurement, error) { return fallible(n.DivideByArea(noArea)) },
		"Energy.DivideByTime":        func() (units.Measurement, error) { return fallible(j.DivideByTime(noTime)) },
		"Energy.DivideByDistance":    func() (units.Measurement, error) { return fallible(j.DivideByDistance(noDistance)) },
		"Frequency.Period":           func() (units.Measurement, error) { return fallible(noFrequency.Period()) },
		"Mass.DivideByVolume":        func() (units.Measurement, error) { return fallible(kg.DivideByVolume(noVolume)) },
		"Mass.DivideByDensity":       func() (units.Measurement, error) { return fallible(kg.DivideByDensity(noDensity)) },
		"DataSize.DivideByTime":      func() (units.Measurement, error) { return fallible(gb.DivideByTime(noTime)) },
		"DataSize.DivideByDataRate":  func() (units.Measurement, error) { return fallible(gb.DivideByDataRate(noRate)) },
	}

	for name, op := range cases {
		if m, err := op(); !errors.Is(err, units.ErrZeroDivision) || m != nil {
			t.Errorf("%s by zero = %v, %v; want nil and ErrZeroDivision", name, m, err)
		}
	}

	// The typed results are nil interfaces, not interfaces wrapping nil.
	if v, err := cubicM.DivideByArea(noArea); err == nil || v != nil {
		t.Errorf("DivideByArea by zero = %v, %v; want a nil Distance", v, err)
	}
}
