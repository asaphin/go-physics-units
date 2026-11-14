package units_test

import (
	"fmt"
	"math"
	"testing"

	units "github.com/asaphin/go-physics-units"
	"github.com/asaphin/go-physics-units/acceleration"
	"github.com/asaphin/go-physics-units/angle"
	"github.com/asaphin/go-physics-units/angularvelocity"
	"github.com/asaphin/go-physics-units/area"
	"github.com/asaphin/go-physics-units/charge"
	"github.com/asaphin/go-physics-units/conversion"
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

var predefinedFactors = map[units.MeasureType]conversion.Factors{
	units.MeasureDistance:        distance.ConversionFactors(),
	units.MeasureArea:            area.ConversionFactors(),
	units.MeasureMass:            mass.ConversionFactors(),
	units.MeasureTime:            time.ConversionFactors(),
	units.MeasureVelocity:        velocity.ConversionFactors(),
	units.MeasureVolume:          volume.ConversionFactors(),
	units.MeasureAcceleration:    acceleration.ConversionFactors(),
	units.MeasureForce:           force.ConversionFactors(),
	units.MeasureEnergy:          energy.ConversionFactors(),
	units.MeasurePower:           power.ConversionFactors(),
	units.MeasurePressure:        pressure.ConversionFactors(),
	units.MeasureFrequency:       frequency.ConversionFactors(),
	units.MeasureAngle:           angle.ConversionFactors(),
	units.MeasureDataSize:        datasize.ConversionFactors(),
	units.MeasureDensity:         density.ConversionFactors(),
	units.MeasureDataRate:        datarate.ConversionFactors(),
	units.MeasureCurrent:         current.ConversionFactors(),
	units.MeasureVoltage:         voltage.ConversionFactors(),
	units.MeasureCharge:          charge.ConversionFactors(),
	units.MeasureResistance:      resistance.ConversionFactors(),
	units.MeasureAngularVelocity: angularvelocity.ConversionFactors(),
}

// TestConversion_AllUnitPairs converts between every ordered pair of units of
// every predefined measure, and back again.
func TestConversion_AllUnitPairs(t *testing.T) {
	for mt, factors := range predefinedFactors {
		for from, fromFactor := range factors {
			for to, toFactor := range factors {
				m, err := units.NewMeasurement(1, from)
				if err != nil {
					t.Fatalf("%s: %v", mt, err)
				}

				if m.Type() != mt {
					t.Errorf("unit %s has type %s, want %s", from, m.Type(), mt)
				}

				got, err := m.ConvertToMeasurement(to)
				if err != nil {
					t.Fatalf("%s: converting %s to %s: %v", mt, from, to, err)
				}

				if want := fromFactor / toFactor; !almostEqual(got.Value(), want, 1e-12) {
					t.Errorf("%s: 1 %s = %v %s, want %v", mt, from, got.Value(), to, want)
				}

				back, err := got.ConvertToMeasurement(from)
				if err != nil {
					t.Fatalf("%s: converting %s back to %s: %v", mt, to, from, err)
				}

				if !almostEqual(back.Value(), 1, 1e-12) {
					t.Errorf("%s: round trip %s -> %s -> %s gives %v, want 1", mt, from, to, from, back.Value())
				}
			}
		}
	}
}

// TestConversion_ReferenceValues checks conversions against independently known
// equivalences, so that a wrong factor can't go unnoticed.
func TestConversion_ReferenceValues(t *testing.T) {
	cases := []struct {
		from fmt.Stringer
		to   fmt.Stringer
		want float64 // 1 from = want to
	}{
		{distance.Meter, distance.Millimeter, 1000},
		{distance.Millimeter, distance.Meter, 0.001},
		{distance.Kilometer, distance.Meter, 1000},
		{distance.Mile, distance.Kilometer, 1.609344},
		{distance.NauticalMile, distance.Meter, 1852},
		{distance.Inch, distance.Centimeter, 2.54},
		{distance.Feet, distance.Inch, 12},
		{distance.Yard, distance.Feet, 3},
		{distance.Angstrom, distance.Nanometer, 0.1},
		{distance.Micrometer, distance.Nanometer, 1000},
		{distance.AstronomicalUnit, distance.Kilometer, 149597870.7},
		{distance.LightYear, distance.Kilometer, 9460730472580.8},
		{distance.Parsec, distance.LightYear, 3.2615637771674},
		{area.SquareKilometer, area.SquareMeter, 1e6},
		{area.SquareMeter, area.SquareCentimeter, 1e4},
		{area.SquareCentimeter, area.SquareMillimeter, 100},
		{area.SquareFoot, area.SquareInch, 144},
		{area.SquareMile, area.SquareKilometer, 2.589988110336},
		{mass.Kilogram, mass.Gram, 1000},
		{mass.Gram, mass.Microgram, 1e6},
		{mass.Milligram, mass.Microgram, 1000},
		{mass.Pound, mass.Kilogram, 0.45359237},
		{mass.Pound, mass.Ounce, 16},
		{mass.Quintal, mass.Kilogram, 100},
		{mass.Tonne, mass.Quintal, 10},
		{time.Second, time.Millisecond, 1000},
		{time.Millisecond, time.Microsecond, 1000},
		{time.Microsecond, time.Nanosecond, 1000},
		{time.Hour, time.Minute, 60},
		{time.Day, time.Hour, 24},
		{time.Week, time.Day, 7},
		{time.Year, time.Day, 365.25},
		{velocity.MeterPerSecond, velocity.KilometerPerHour, 3.6},
		{velocity.Knot, velocity.KilometerPerHour, 1.852},
		{velocity.MilePerHour, velocity.KilometerPerHour, 1.609344},
		{velocity.FeetPerSecond, velocity.MeterPerSecond, 0.3048},
		{velocity.SpeedOfLight, velocity.KilometerPerHour, 1079252848.8},
		{volume.CubicMeter, volume.Liter, 1000},
		{volume.Liter, volume.Milliliter, 1000},
		{volume.Milliliter, volume.CubicCentimeter, 1},
		{volume.Microliter, volume.CubicMillimeter, 1},
		{volume.Gallon, volume.Liter, 3.785411784},
		{volume.Gallon, volume.Quart, 4},
		{volume.Quart, volume.Pint, 2},
		{volume.Gallon, volume.CubicInch, 231},
		{volume.CubicFoot, volume.CubicInch, 1728},
		{volume.CubicYard, volume.CubicFoot, 27},
		{volume.Barrel, volume.Gallon, 42},
		{volume.CubicKilometer, volume.CubicMeter, 1e9},
		{acceleration.StandardGravity, acceleration.MeterPerSecondSquared, 9.80665},
		{acceleration.MeterPerSecondSquared, acceleration.CentimeterPerSecondSquared, 100},
		{acceleration.FootPerSecondSquared, acceleration.MeterPerSecondSquared, 0.3048},
		{force.KilogramForce, force.Newton, 9.80665},
		{force.PoundForce, force.Newton, 4.4482216152605},
		{force.Kilonewton, force.Newton, 1000},
		{force.Newton, force.Dyne, 1e5},
		{force.Meganewton, force.Kilonewton, 1000},
		{force.Newton, force.Millinewton, 1000},
		{energy.KilowattHour, energy.Megajoule, 3.6},
		{energy.WattHour, energy.Joule, 3600},
		{energy.MegawattHour, energy.KilowattHour, 1000},
		{energy.Kilocalorie, energy.Kilojoule, 4.184},
		{energy.Kilocalorie, energy.Calorie, 1000},
		{energy.BTU, energy.Joule, 1055.05585262},
		{energy.Electronvolt, energy.Joule, 1.602176634e-19},
		{energy.Joule, energy.Erg, 1e7},
		{energy.FootPound, energy.Joule, 1.3558179483314004},
		{energy.Gigajoule, energy.Megajoule, 1000},
		{power.Horsepower, power.Watt, 745.69987158227022},
		{power.MetricHorsepower, power.Watt, 735.49875},
		{power.Kilowatt, power.Watt, 1000},
		{power.Watt, power.Milliwatt, 1000},
		{power.Gigawatt, power.Megawatt, 1000},
		{power.BTUPerHour, power.Watt, 0.29307107017222222},
		{pressure.Atmosphere, pressure.Pascal, 101325},
		{pressure.Atmosphere, pressure.Torr, 760},
		{pressure.Bar, pressure.Kilopascal, 100},
		{pressure.Hectopascal, pressure.Millibar, 1},
		{pressure.PoundPerSquareInch, pressure.Pascal, 6894.757293168361},
		{pressure.MillimeterOfMercury, pressure.Pascal, 133.322387415},
		{pressure.InchOfMercury, pressure.MillimeterOfMercury, 25.4},
		{pressure.Megapascal, pressure.Bar, 10},
		{pressure.Gigapascal, pressure.Megapascal, 1000},
		{frequency.Kilohertz, frequency.Hertz, 1000},
		{frequency.Hertz, frequency.RevolutionPerMinute, 60},
		{frequency.Gigahertz, frequency.Megahertz, 1000},
		{frequency.Terahertz, frequency.Gigahertz, 1000},
		{angle.Revolution, angle.Degree, 360},
		{angle.Degree, angle.Arcminute, 60},
		{angle.Arcminute, angle.Arcsecond, 60},
		{angle.Revolution, angle.Gradian, 400},
		{angle.Radian, angle.Milliradian, 1000},
		{angle.Degree, angle.Radian, math.Pi / 180},
		{datasize.Byte, datasize.Bit, 8},
		{datasize.Kilobyte, datasize.Byte, 1000},
		{datasize.Kibibyte, datasize.Byte, 1024},
		{datasize.Mebibyte, datasize.Kibibyte, 1024},
		{datasize.Gibibyte, datasize.Mebibyte, 1024},
		{datasize.Tebibyte, datasize.Gibibyte, 1024},
		{datasize.Pebibyte, datasize.Tebibyte, 1024},
		{datasize.Gigabyte, datasize.Megabyte, 1000},
		{datasize.Petabyte, datasize.Terabyte, 1000},
		{datasize.Megabit, datasize.Kilobit, 1000},
		{datasize.Gigabit, datasize.Megabyte, 125},
		{datasize.Megabyte, datasize.Megabit, 8},
		{density.GramPerCubicCentimeter, density.KilogramPerCubicMeter, 1000},
		{density.KilogramPerLiter, density.GramPerMilliliter, 1},
		{density.GramPerLiter, density.KilogramPerCubicMeter, 1},
		{density.PoundPerCubicFoot, density.KilogramPerCubicMeter, 16.018463373960138},
		{density.PoundPerGallon, density.KilogramPerCubicMeter, 119.82642731689663},
		{datarate.MegabitPerSecond, datarate.KilobitPerSecond, 1000},
		{datarate.BytePerSecond, datarate.BitPerSecond, 8},
		{datarate.MegabytePerSecond, datarate.MegabitPerSecond, 8},
		{datarate.GibibytePerSecond, datarate.MebibytePerSecond, 1024},
		{datarate.GigabitPerSecond, datarate.MegabytePerSecond, 125},
		{datarate.KibibytePerSecond, datarate.BytePerSecond, 1024},
		{datarate.GigabytePerSecond, datarate.MegabytePerSecond, 1000},
		{current.Ampere, current.Milliampere, 1000},
		{current.Milliampere, current.Microampere, 1000},
		{current.Kiloampere, current.Ampere, 1000},
		{voltage.Kilovolt, voltage.Volt, 1000},
		{voltage.Megavolt, voltage.Kilovolt, 1000},
		{voltage.Millivolt, voltage.Microvolt, 1000},
		{charge.AmpereHour, charge.Coulomb, 3600},
		{charge.MilliampereHour, charge.Coulomb, 3.6},
		{charge.Coulomb, charge.Millicoulomb, 1000},
		{charge.Millicoulomb, charge.Microcoulomb, 1000},
		{resistance.Kiloohm, resistance.Ohm, 1000},
		{resistance.Gigaohm, resistance.Megaohm, 1000},
		{resistance.Ohm, resistance.Milliohm, 1000},
		{angularvelocity.RevolutionPerSecond, angularvelocity.RevolutionPerMinute, 60},
		{angularvelocity.RevolutionPerSecond, angularvelocity.DegreePerSecond, 360},
		{angularvelocity.RevolutionPerMinute, angularvelocity.RadianPerSecond, 0.10471975511965977},
	}

	for _, c := range cases {
		m, err := units.NewMeasurement(1, c.from.String())
		if err != nil {
			t.Fatal(err)
		}

		got, err := m.ConvertToMeasurement(c.to.String())
		if err != nil {
			t.Fatal(err)
		}

		if !almostEqual(got.Value(), c.want, 1e-12) {
			t.Errorf("1 %s = %v %s, want %v", c.from, got.Value(), c.to, c.want)
		}
	}
}

func TestConversion_MicroSignVariants(t *testing.T) {
	const (
		microSign = "\u00b5" // MICRO SIGN
		greekMu   = "\u03bc" // GREEK SMALL LETTER MU
	)

	for _, unit := range []fmt.Stringer{distance.Micrometer, mass.Microgram, time.Microsecond} {
		if unit.String()[:len(greekMu)] != greekMu {
			t.Errorf("predefined symbol %s doesn't use U+03BC", unit)
		}
	}

	if u, ok := distance.ParseUnit(microSign + "m"); !ok || u != distance.Micrometer {
		t.Errorf("ParseUnit(%sm) = %v, %v; want %v", microSign, u, ok, distance.Micrometer)
	}

	d, err := units.StringToDistance("1 " + microSign + "m")
	if err != nil {
		t.Fatal(err)
	}

	assertMeasurement(t, d, 1, distance.Micrometer.String())

	nm, err := d.ConvertToMeasurement(microSign + "m")
	if err != nil {
		t.Fatal(err)
	}

	assertMeasurement(t, nm, 1, distance.Micrometer.String())

	g, err := units.ParseString("5 " + microSign + "g")
	if err != nil {
		t.Fatal(err)
	}

	assertMeasurement(t, g, 5, mass.Microgram.String())

	if mt := units.DetectMeasureType(microSign + "s"); mt != units.MeasureTime {
		t.Errorf("DetectMeasureType(\u00b5s) = %s, want %s", mt, units.MeasureTime)
	}
}

// TestRate checks each unit package's Rate against the examples that its
// documentation used to give. testUnitPackage checks every rate against the
// factors.
func TestRate(t *testing.T) {
	examples := []struct {
		name string
		rate float64
		want float64
	}{
		{"distance", distance.Rate(distance.Kilometer, distance.Meter), 1000},
		{"area", area.Rate(area.SquareKilometer, area.SquareMeter), 1e6},
		{"mass", mass.Rate(mass.Kilogram, mass.Gram), 1000},
		{"time", time.Rate(time.Hour, time.Minute), 60},
		{"velocity", velocity.Rate(velocity.MeterPerSecond, velocity.KilometerPerHour), 3.6},
		{"volume", volume.Rate(volume.Liter, volume.Milliliter), 1000},
		{"acceleration", acceleration.Rate(acceleration.StandardGravity, acceleration.MeterPerSecondSquared), 9.80665},
		{"force", force.Rate(force.Kilonewton, force.Newton), 1000},
		{"energy", energy.Rate(energy.KilowattHour, energy.Megajoule), 3.6},
		{"power", power.Rate(power.Kilowatt, power.Watt), 1000},
		{"pressure", pressure.Rate(pressure.Bar, pressure.Kilopascal), 100},
		{"frequency", frequency.Rate(frequency.Kilohertz, frequency.Hertz), 1000},
		{"angle", angle.Rate(angle.Revolution, angle.Degree), 360},
		{"data size", datasize.Rate(datasize.Kibibyte, datasize.Byte), 1024},
		{"density", density.Rate(density.KilogramPerLiter, density.KilogramPerCubicMeter), 1000},
		{"data rate", datarate.Rate(datarate.MegabytePerSecond, datarate.MegabitPerSecond), 8},
		{"current", current.Rate(current.Ampere, current.Milliampere), 1000},
		{"voltage", voltage.Rate(voltage.Kilovolt, voltage.Volt), 1000},
		{"charge", charge.Rate(charge.AmpereHour, charge.MilliampereHour), 1000},
		{"resistance", resistance.Rate(resistance.Megaohm, resistance.Kiloohm), 1000},
		{
			"angular velocity",
			angularvelocity.Rate(angularvelocity.RevolutionPerSecond, angularvelocity.RevolutionPerMinute),
			60,
		},
	}

	for _, e := range examples {
		if !almostEqual(e.rate, e.want, 1e-15) {
			t.Errorf("%s: Rate = %v, want %v", e.name, e.rate, e.want)
		}
	}
}

func TestConversionFactors_ReturnsCopy(t *testing.T) {
	for mt, factorsFn := range map[units.MeasureType]func() conversion.Factors{
		units.MeasureDistance:        distance.ConversionFactors,
		units.MeasureArea:            area.ConversionFactors,
		units.MeasureMass:            mass.ConversionFactors,
		units.MeasureTime:            time.ConversionFactors,
		units.MeasureVelocity:        velocity.ConversionFactors,
		units.MeasureVolume:          volume.ConversionFactors,
		units.MeasureAcceleration:    acceleration.ConversionFactors,
		units.MeasureForce:           force.ConversionFactors,
		units.MeasureEnergy:          energy.ConversionFactors,
		units.MeasurePower:           power.ConversionFactors,
		units.MeasurePressure:        pressure.ConversionFactors,
		units.MeasureFrequency:       frequency.ConversionFactors,
		units.MeasureAngle:           angle.ConversionFactors,
		units.MeasureDataSize:        datasize.ConversionFactors,
		units.MeasureDensity:         density.ConversionFactors,
		units.MeasureDataRate:        datarate.ConversionFactors,
		units.MeasureCurrent:         current.ConversionFactors,
		units.MeasureVoltage:         voltage.ConversionFactors,
		units.MeasureCharge:          charge.ConversionFactors,
		units.MeasureResistance:      resistance.ConversionFactors,
		units.MeasureAngularVelocity: angularvelocity.ConversionFactors,
	} {
		f := factorsFn()
		for unit := range f {
			f[unit] = -1
		}

		for unit, factor := range factorsFn() {
			if factor == -1 {
				t.Errorf("%s: mutating the result of ConversionFactors changed factor of %s", mt, unit)
			}
		}
	}
}
