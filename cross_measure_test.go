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

// ok adapts an operation that can't fail to the shape of the table below.
func ok[T units.Measurement](m T) (units.Measurement, error) {
	return m, nil
}

// fallible adapts an operation that can fail, so that the result is nil on
// error.
func fallible[T units.Measurement](m T, err error) (units.Measurement, error) {
	if err != nil {
		return nil, err
	}

	return m, nil
}

// TestCrossMeasureOperations checks each operation that combines measures
// against an independently known result, in the result's base unit.
func TestCrossMeasureOperations(t *testing.T) {
	m2 := units.NewDistance(2, distance.Meter)
	cm10 := units.NewDistance(10, distance.Centimeter)
	cm50 := units.NewDistance(50, distance.Centimeter)
	sqm2 := units.NewArea(2, area.SquareMeter)
	sqin := units.NewArea(1, area.SquareInch)
	cubicM := units.NewVolume(1, volume.CubicMeter)
	liter := units.NewVolume(1, volume.Liter)
	kmh100 := units.NewVelocity(100, velocity.KilometerPerHour)
	s2 := units.NewTime(2, time.Second)
	s10 := units.NewTime(10, time.Second)
	hour := units.NewTime(1, time.Hour)
	gn := units.NewAcceleration(1, acceleration.StandardGravity)
	kg1 := units.NewMass(1, mass.Kilogram)
	kg2 := units.NewMass(2, mass.Kilogram)
	lb := units.NewMass(1, mass.Pound)
	n10 := units.NewForce(10, force.Newton)
	kgf := units.NewForce(1, force.KilogramForce)
	lbf := units.NewForce(1, force.PoundForce)
	j20 := units.NewEnergy(20, energy.Joule)
	kwh := units.NewEnergy(1, energy.KilowattHour)
	kw := units.NewPower(1, power.Kilowatt)
	atm := units.NewPressure(1, pressure.Atmosphere)
	hz50 := units.NewFrequency(50, frequency.Hertz)
	rpm60 := units.NewFrequency(60, frequency.RevolutionPerMinute)
	water := units.NewDensity(1, density.GramPerCubicCentimeter)
	liter2 := units.NewVolume(2, volume.Liter)
	gb := units.NewDataSize(1, datasize.Gigabyte)
	s8 := units.NewTime(8, time.Second)
	mbit100 := units.NewDataRate(100, datarate.MegabitPerSecond)
	minute := units.NewTime(1, time.Minute)
	km100 := units.NewDistance(100, distance.Kilometer)
	h2 := units.NewTime(2, time.Hour)
	km3 := units.NewDistance(3, distance.Kilometer)
	sqkm := units.NewArea(1, area.SquareKilometer)
	m500 := units.NewDistance(500, distance.Meter)
	kmh20 := units.NewVelocity(20, velocity.KilometerPerHour)
	min30 := units.NewTime(30, time.Minute)

	cases := []struct {
		name  string
		op    func() (units.Measurement, error)
		value float64
		unit  string
	}{
		{
			"100 km / 2 h",
			func() (units.Measurement, error) { return fallible(km100.DivideByTime(h2)) },
			100000.0 / 7200,
			velocity.BaseUnit.String(),
		},
		{"3 km squared", func() (units.Measurement, error) { return ok(km3.Square()) }, 9e6, area.BaseUnit.String()},
		{
			"1 km² / 500 m",
			func() (units.Measurement, error) { return fallible(sqkm.DivideByDistance(m500)) },
			2000,
			distance.BaseUnit.String(),
		},
		{
			"20 km/h × 30 min",
			func() (units.Measurement, error) { return ok(kmh20.MultiplyByTime(min30)) },
			10000,
			distance.BaseUnit.String(),
		},
		{"2 m cubed", func() (units.Measurement, error) { return ok(m2.Cube()) }, 8, volume.BaseUnit.String()},
		{
			"10 cm cubed is a litre",
			func() (units.Measurement, error) { return ok(cm10.Cube()) },
			1e-3,
			volume.BaseUnit.String(),
		},
		{
			"2 m² × 50 cm",
			func() (units.Measurement, error) { return ok(sqm2.MultiplyByDistance(cm50)) },
			1,
			volume.BaseUnit.String(),
		},
		{
			"1 m³ / 2 m²",
			func() (units.Measurement, error) { return fallible(cubicM.DivideByArea(sqm2)) },
			0.5,
			distance.BaseUnit.String(),
		},
		{
			"1 L / 10 cm",
			func() (units.Measurement, error) { return fallible(liter.DivideByDistance(cm10)) },
			0.01,
			area.BaseUnit.String(),
		},
		{
			"100 km/h / 10 s", func() (units.Measurement, error) { return fallible(kmh100.DivideByTime(s10)) },
			100.0 / 36, acceleration.BaseUnit.String(),
		},
		{
			"gₙ × 2 s",
			func() (units.Measurement, error) { return ok(gn.MultiplyByTime(s2)) },
			19.6133,
			velocity.BaseUnit.String(),
		},
		{
			"gₙ × 1 kg",
			func() (units.Measurement, error) { return ok(gn.MultiplyByMass(kg1)) },
			9.80665,
			force.BaseUnit.String(),
		},
		{
			"1 lb × gₙ is a pound-force",
			func() (units.Measurement, error) { return ok(lb.MultiplyByAcceleration(gn)) },
			4.4482216152605,
			force.BaseUnit.String(),
		},
		{
			"10 N / 2 kg",
			func() (units.Measurement, error) { return fallible(n10.DivideByMass(kg2)) },
			5,
			acceleration.BaseUnit.String(),
		},
		{
			"1 kgf / gₙ",
			func() (units.Measurement, error) { return fallible(kgf.DivideByAcceleration(gn)) },
			1,
			mass.BaseUnit.String(),
		},
		{
			"1 lbf / 1 in² is a psi", func() (units.Measurement, error) { return fallible(lbf.DivideByArea(sqin)) },
			6894.757293168361, pressure.BaseUnit.String(),
		},
		{
			"10 N × 2 m",
			func() (units.Measurement, error) { return ok(n10.MultiplyByDistance(m2)) },
			20,
			energy.BaseUnit.String(),
		},
		{
			"1 kWh / 1 h",
			func() (units.Measurement, error) { return fallible(kwh.DivideByTime(hour)) },
			1000,
			power.BaseUnit.String(),
		},
		{
			"20 J / 2 m",
			func() (units.Measurement, error) { return fallible(j20.DivideByDistance(m2)) },
			10,
			force.BaseUnit.String(),
		},
		{
			"1 kW × 1 h",
			func() (units.Measurement, error) { return ok(kw.MultiplyByTime(hour)) },
			3.6e6,
			energy.BaseUnit.String(),
		},
		{
			"1 atm × 2 m²",
			func() (units.Measurement, error) { return ok(atm.MultiplyByArea(sqm2)) },
			202650,
			force.BaseUnit.String(),
		},
		{
			"period of 50 Hz",
			func() (units.Measurement, error) { return fallible(hz50.Period()) },
			0.02,
			time.BaseUnit.String(),
		},
		{
			"period of 60 rpm",
			func() (units.Measurement, error) { return fallible(rpm60.Period()) },
			1,
			time.BaseUnit.String(),
		},
		{
			"1 kg / 1 L",
			func() (units.Measurement, error) { return fallible(kg1.DivideByVolume(liter)) },
			1000,
			density.BaseUnit.String(),
		},
		{
			"1 g/cm³ × 2 L",
			func() (units.Measurement, error) { return ok(water.MultiplyByVolume(liter2)) },
			2,
			mass.BaseUnit.String(),
		},
		{
			"1 kg / 1 g/cm³",
			func() (units.Measurement, error) { return fallible(kg1.DivideByDensity(water)) },
			1e-3,
			volume.BaseUnit.String(),
		},
		{
			"1 GB / 8 s",
			func() (units.Measurement, error) { return fallible(gb.DivideByTime(s8)) },
			1.25e8,
			datarate.BaseUnit.String(),
		},
		{
			"1 GB at 100 Mbit/s", func() (units.Measurement, error) { return fallible(gb.DivideByDataRate(mbit100)) },
			80, time.BaseUnit.String(),
		},
		{
			"100 Mbit/s × 1 min",
			func() (units.Measurement, error) { return ok(mbit100.MultiplyByTime(minute)) },
			7.5e8,
			datasize.BaseUnit.String(),
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
	noDistance := units.NewDistance(0, distance.Meter)
	noArea := units.NewArea(0, area.SquareMeter)
	noTime := units.NewTime(0, time.Second)
	noMass := units.NewMass(0, mass.Kilogram)
	noAcceleration := units.NewAcceleration(0, acceleration.StandardGravity)
	cubicM := units.NewVolume(1, volume.CubicMeter)
	speed := units.NewVelocity(1, velocity.MeterPerSecond)
	n := units.NewForce(1, force.Newton)
	j := units.NewEnergy(1, energy.Joule)
	noFrequency := units.NewFrequency(0, frequency.Hertz)
	kg := units.NewMass(1, mass.Kilogram)
	noVolume := units.NewVolume(0, volume.Liter)
	noDensity := units.NewDensity(0, density.GramPerLiter)
	gb := units.NewDataSize(1, datasize.Gigabyte)
	noRate := units.NewDataRate(0, datarate.BitPerSecond)

	cases := map[string]func() (units.Measurement, error){
		"Distance.DivideByTime":      func() (units.Measurement, error) { return fallible(noDistance.DivideByTime(noTime)) },
		"Area.DivideByDistance":      func() (units.Measurement, error) { return fallible(noArea.DivideByDistance(noDistance)) },
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

	// The typed results are zero values.
	if d, err := cubicM.DivideByArea(noArea); err == nil || d != (units.Distance{}) {
		t.Errorf("DivideByArea by zero = %v, %v; want the zero Distance", d, err)
	}

	if v, err := units.NewDistance(1, distance.Meter).DivideByTime(noTime); err == nil || v != (units.Velocity{}) {
		t.Errorf("DivideByTime by zero = %v, %v; want the zero Velocity", v, err)
	}
}
