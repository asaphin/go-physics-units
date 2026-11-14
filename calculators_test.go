package units_test

import (
	"errors"
	"math"
	"testing"

	units "github.com/asaphin/go-physics-units"
	"github.com/asaphin/go-physics-units/absorbeddose"
	"github.com/asaphin/go-physics-units/amount"
	"github.com/asaphin/go-physics-units/area"
	"github.com/asaphin/go-physics-units/distance"
	"github.com/asaphin/go-physics-units/energy"
	"github.com/asaphin/go-physics-units/equivalentdose"
	"github.com/asaphin/go-physics-units/force"
	"github.com/asaphin/go-physics-units/frequency"
	"github.com/asaphin/go-physics-units/fuelconsumption"
	"github.com/asaphin/go-physics-units/fueleconomy"
	"github.com/asaphin/go-physics-units/mass"
	"github.com/asaphin/go-physics-units/pace"
	"github.com/asaphin/go-physics-units/radioactivity"
	"github.com/asaphin/go-physics-units/ratio"
	"github.com/asaphin/go-physics-units/time"
	"github.com/asaphin/go-physics-units/velocity"
	"github.com/asaphin/go-physics-units/volume"
)

// TestCalculators checks the methods that aren't plain products or quotients
// of two measures against values computed from their formulas.
func TestCalculators(t *testing.T) {
	m1 := units.NewDistance(1, distance.Meter)
	m2 := units.NewDistance(2, distance.Meter)
	m3 := units.NewDistance(300, distance.Centimeter)
	c := units.NewVelocity(1, velocity.SpeedOfLight)

	cases := []struct {
		name  string
		op    func() (units.Measurement, error)
		value float64
		unit  string
	}{
		{"1 m CircleArea", func() (units.Measurement, error) { return ok(m1.CircleArea()) }, math.Pi, "m²"},
		{"2 m SphereArea", func() (units.Measurement, error) { return ok(m2.SphereArea()) }, 16 * math.Pi, "m²"},
		{"3 m SphereVolume", func() (units.Measurement, error) { return ok(m3.SphereVolume()) }, 36 * math.Pi, "m³"},
		{
			"9 m² SquareSide",
			func() (units.Measurement, error) { return fallible(units.NewArea(9, area.SquareMeter).SquareSide()) },
			3, "m",
		},
		{
			"π m² CircleRadius",
			func() (units.Measurement, error) {
				return fallible(units.NewArea(math.Pi*1e4, area.SquareCentimeter).CircleRadius())
			},
			1, "m",
		},
		{
			"8 L CubeEdge",
			func() (units.Measurement, error) { return ok(units.NewVolume(8, volume.Liter).CubeEdge()) },
			0.2,
			"m",
		},
		{
			"-8 m³ CubeEdge",
			func() (units.Measurement, error) { return ok(units.NewVolume(-8, volume.CubicMeter).CubeEdge()) },
			-2, "m",
		},
		{
			"4π/3 m³ SphereRadius",
			func() (units.Measurement, error) {
				return ok(units.NewVolume(4*math.Pi/3, volume.CubicMeter).SphereRadius())
			},
			1, "m",
		},
		{
			"2 kg at 3 m/s",
			func() (units.Measurement, error) {
				return ok(units.NewMass(2, mass.Kilogram).KineticEnergy(units.NewVelocity(3, velocity.MeterPerSecond)))
			},
			9, "J",
		},
		{
			"1 kg Weight", func() (units.Measurement, error) { return ok(units.NewMass(1, mass.Kilogram).Weight()) },
			units.NewForce(1, force.KilogramForce).ConvertToBaseUnits().Value(), "N",
		},
		{
			"1 g RestEnergy", func() (units.Measurement, error) { return ok(units.NewMass(1, mass.Gram).RestEnergy()) },
			299792458.0 * 299792458 / 1000, "J",
		},
		{
			"1 Mt TNT MassEquivalent",
			func() (units.Measurement, error) {
				return ok(units.NewEnergy(1, energy.MegatonOfTNT).MassEquivalent())
			},
			4.184e15 / (299792458.0 * 299792458), "kg",
		},
		{
			"1 eV PhotonFrequency",
			func() (units.Measurement, error) {
				return ok(units.NewEnergy(1, energy.Electronvolt).PhotonFrequency())
			},
			2.417989242084918e14, "Hz", // CODATA 2022: 2.417 989 242... e14 Hz
		},
		{
			"1 PHz PhotonEnergy",
			func() (units.Measurement, error) {
				return ok(units.NewFrequency(1000, frequency.Terahertz).PhotonEnergy())
			},
			6.62607015e-19, "J",
		},
		{
			"5 m/s Pace", func() (units.Measurement, error) { return fallible(c.Mul(5 / 299792458.0).Pace()) },
			0.2, "s/m",
		},
		{
			"5 min/km Velocity",
			func() (units.Measurement, error) {
				return fallible(units.NewPace(5, pace.MinutePerKilometer).Velocity())
			},
			1000.0 / 300, "m/s",
		},
		{
			"20 ms Frequency",
			func() (units.Measurement, error) { return fallible(units.NewTime(20, time.Millisecond).Frequency()) },
			50, "Hz",
		},
		{
			"20 km/L Consumption",
			func() (units.Measurement, error) {
				return fallible(units.NewFuelEconomy(20, fueleconomy.KilometerPerLiter).Consumption())
			},
			5e-8, "m³/m", // 5 L/100km
		},
		{
			"30 mpg Consumption",
			func() (units.Measurement, error) {
				return fallible(units.NewFuelEconomy(30, fueleconomy.MilePerGallon).Consumption())
			},
			235.2145833333333 / 30 * 1e-8, "m³/m", // L/100km = 235.2145833… / mpg
		},
		{
			"4 L/100km Economy",
			func() (units.Measurement, error) {
				return fallible(units.NewFuelConsumption(4, fuelconsumption.LiterPer100Kilometers).Economy())
			},
			2.5e7, "m/m³", // 25 km/L
		},
		{
			"two half-lives",
			func() (units.Measurement, error) {
				return fallible(units.NewRadioactivity(1000, radioactivity.Becquerel).DecayAfter(
					units.NewTime(10, time.Day), units.NewTime(5, time.Day)))
			},
			250, "Bq",
		},
		{
			"one half-life before, in Ci",
			func() (units.Measurement, error) {
				return fallible(units.NewRadioactivity(1, radioactivity.Curie).DecayAfter(
					units.NewTime(-30, time.Year), units.NewTime(30, time.Year)))
			},
			2, "Ci",
		},
		{
			"2 mGy of alpha",
			func() (units.Measurement, error) {
				return ok(units.NewAbsorbedDose(2, absorbeddose.Milligray).Equivalent(20))
			},
			0.04, "Sv",
		},
		{
			"1 km / 250 m",
			func() (units.Measurement, error) {
				return fallible(units.NewDistance(1, distance.Kilometer).Ratio(units.NewDistance(250, distance.Meter)))
			},
			4, "1",
		},
		{
			"3 mSv / 1 rem, same measure in two units",
			func() (units.Measurement, error) {
				return fallible(units.NewEquivalentDose(3, equivalentdose.Millisievert).Ratio(
					units.NewEquivalentDose(1, equivalentdose.Rem)))
			},
			0.3, "1",
		},
	}

	for _, c := range cases {
		m, err := c.op()
		if err != nil {
			t.Fatalf("%s: %v", c.name, err)
		}

		assertMeasurement(t, m, c.value, c.unit)
	}
}

func TestCalculators_Scalars(t *testing.T) {
	c := units.NewVelocity(1, velocity.SpeedOfLight)

	for _, tc := range []struct {
		name string
		v    units.Velocity
		want float64
	}{
		{"0.6 c", c.Mul(0.6), 1.25},
		{"-0.6 c", c.Mul(-0.6), 1.25},
		{"0.8 c", c.Mul(0.8), 5.0 / 3},
		{"at rest", c.Mul(0), 1},
		{"100 km/h", units.NewVelocity(100, velocity.KilometerPerHour), 1 + 4.290e-15},
	} {
		got, err := tc.v.LorentzFactor()
		if err != nil || !almostEqual(got, tc.want, 1e-14) {
			t.Errorf("LorentzFactor(%s) = %v, %v; want %v", tc.name, got, err, tc.want)
		}
	}

	for _, v := range []units.Velocity{c, c.Mul(-1.5), c.Mul(math.NaN())} {
		if got, err := v.LorentzFactor(); !errors.Is(err, units.ErrOutOfRange) || got != 0 {
			t.Errorf("LorentzFactor(%v) = %v, %v; want 0 and ErrOutOfRange", v, got, err)
		}
	}

	if got := units.NewAmount(2, amount.Millimole).Entities(); !almostEqual(got, 1.204428152e21, 1e-15) {
		t.Errorf("Entities(2 mmol) = %v, want 1.204428152e21", got)
	}

	if got := units.NewRatio(25, ratio.Percent).Fraction(); got != 0.25 {
		t.Errorf("Fraction(25 %%) = %v, want 0.25", got)
	}
}

func TestCalculators_Errors(t *testing.T) {
	zeroDivision(t, "Velocity.Pace", units.NewVelocity(0, velocity.MeterPerSecond).Pace)
	zeroDivision(t, "Pace.Velocity", units.NewPace(0, pace.SecondPerMeter).Velocity)
	zeroDivision(t, "Time.Frequency", units.NewTime(0, time.Second).Frequency)
	zeroDivision(t, "FuelEconomy.Consumption", units.NewFuelEconomy(0, fueleconomy.KilometerPerLiter).Consumption)
	zeroDivision(t, "FuelConsumption.Economy", units.NewFuelConsumption(0, fuelconsumption.LiterPerKilometer).Economy)
	zeroDivision(t, "Radioactivity.DecayAfter", func() (units.Radioactivity, error) {
		return units.NewRadioactivity(1, radioactivity.Becquerel).DecayAfter(units.NewTime(1, time.Second),
			units.NewTime(0, time.Second))
	})
	zeroDivision(t, "Distance.Ratio", func() (units.Ratio, error) {
		return units.NewDistance(1, distance.Meter).Ratio(units.NewDistance(0, distance.Kilometer))
	})

	negative := units.NewArea(-1, area.SquareMeter)
	for name, op := range map[string]func() (units.Distance, error){
		"SquareSide":   negative.SquareSide,
		"CircleRadius": negative.CircleRadius,
	} {
		if d, err := op(); !errors.Is(err, units.ErrOutOfRange) || d != (units.Distance{}) {
			t.Errorf("%s of a negative area = %v, %v; want the zero Distance and ErrOutOfRange", name, d, err)
		}
	}
}
