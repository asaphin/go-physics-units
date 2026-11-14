package physics_test

import (
	"errors"
	"math"
	"testing"

	units "github.com/asaphin/go-physics-units"
	"github.com/asaphin/go-physics-units/acceleration"
	"github.com/asaphin/go-physics-units/amount"
	"github.com/asaphin/go-physics-units/angle"
	"github.com/asaphin/go-physics-units/capacitance"
	"github.com/asaphin/go-physics-units/constant"
	"github.com/asaphin/go-physics-units/current"
	"github.com/asaphin/go-physics-units/density"
	"github.com/asaphin/go-physics-units/distance"
	"github.com/asaphin/go-physics-units/frequency"
	"github.com/asaphin/go-physics-units/inductance"
	"github.com/asaphin/go-physics-units/mass"
	"github.com/asaphin/go-physics-units/physics"
	"github.com/asaphin/go-physics-units/pressure"
	"github.com/asaphin/go-physics-units/temperature"
	"github.com/asaphin/go-physics-units/time"
	"github.com/asaphin/go-physics-units/velocity"
	"github.com/asaphin/go-physics-units/voltage"
	"github.com/asaphin/go-physics-units/volume"
)

// value is what the results have in common.
type value interface {
	Value() float64
	Symbol() string
}

func near(a, b, tol float64) bool {
	return math.Abs(a-b) <= tol*math.Max(math.Abs(a), math.Abs(b))
}

func must[T any](v T, err error) T {
	if err != nil {
		panic(err)
	}

	return v
}

func kelvin(k float64) units.Temperature {
	return must(units.NewTemperature(k, temperature.Kelvin))
}

// TestFormulas checks each formula against a value known independently of it,
// in the result's base unit.
func TestFormulas(t *testing.T) {
	gn := units.NewAcceleration(1, acceleration.StandardGravity)
	g10 := units.NewAcceleration(10, acceleration.MeterPerSecondSquared)
	earthMass := units.NewMass(1, mass.EarthMass)
	sunMass := units.NewMass(1, mass.SolarMass)
	earthRadius := units.NewDistance(constant.NominalEarthEquatorialRadius, distance.Meter)
	molarVolume := units.NewVolume(22.41396954, volume.Liter) // CODATA 2022: an ideal gas at 0 °C and 1 atm
	atm := units.NewPressure(1, pressure.Atmosphere)
	gmEarth := constant.NominalEarthMassParameter
	gmSun := constant.NominalSolarMassParameter

	cases := []struct {
		name   string
		got    value
		want   float64
		symbol string
		tol    float64
	}{
		{
			"fall 2·gₙ m", must(physics.FreeFallTime(units.NewDistance(2*9.80665, distance.Meter), gn)),
			2, "s", 1e-15,
		},
		{
			"speed after 5 m at 9.8 m/s²",
			must(physics.FreeFallSpeed(units.NewDistance(5, distance.Meter),
				units.NewAcceleration(9.8, acceleration.MeterPerSecondSquared))),
			math.Sqrt(98), "m/s", 1e-15,
		},
		{"fallen in 2 s", physics.FallDistance(units.NewTime(2, time.Second), gn), 4 * 9.80665 / 2, "m", 1e-15},
		{
			"10 kg at 10 m",
			physics.PotentialEnergy(units.NewMass(10, mass.Kilogram), gn, units.NewDistance(10, distance.Meter)),
			980.665, "J", 1e-15,
		},
		{
			"20 m/s at 45°",
			must(physics.ProjectileRange(units.NewVelocity(20, velocity.MeterPerSecond),
				units.NewAngle(45, angle.Degree), g10)),
			40, "m", 1e-14,
		},
		{
			"seconds pendulum",
			must(physics.PendulumPeriod(units.NewDistance(9.80665/(math.Pi*math.Pi), distance.Meter), gn)),
			2, "s", 1e-15,
		},
		{
			"10 m of water",
			physics.HydrostaticPressure(units.NewDensity(1000, density.KilogramPerCubicMeter), gn,
				units.NewDistance(10, distance.Meter)),
			98066.5, "Pa", 1e-15, // 1 kgf/cm²
		},
		{
			"1 kg and 1 kg at 1 m",
			must(physics.GravitationalForce(units.NewMass(1, mass.Kilogram), units.NewMass(1, mass.Kilogram),
				units.NewDistance(1, distance.Meter))),
			constant.GravitationalConstant, "N", 1e-15,
		},
		{
			"the Earth's surface gravity", must(physics.SurfaceGravity(earthMass, earthRadius)),
			gmEarth / (6.3781e6 * 6.3781e6), "m/s²", 1e-4,
		},
		{
			"the Earth's escape velocity", must(physics.EscapeVelocity(earthMass, earthRadius)),
			math.Sqrt(2 * gmEarth / 6.3781e6), "m/s", 1e-4,
		},
		{
			"a low Earth orbit", must(physics.OrbitalVelocity(earthMass, earthRadius)),
			math.Sqrt(gmEarth / 6.3781e6), "m/s", 1e-4,
		},
		{
			"one Gaussian year",
			must(physics.OrbitalPeriod(sunMass, units.NewDistance(1, distance.AstronomicalUnit))),
			2 * math.Pi / 0.01720209895 * 86400, "s", 1e-4,
		},
		{"the Sun's horizon", physics.SchwarzschildRadius(sunMass), 2 * gmSun / (299792458.0 * 299792458), "m", 1e-4},
		{
			"100 MHz in vacuum", must(physics.LightWavelength(units.NewFrequency(100, frequency.Megahertz))),
			2.99792458, "m", 1e-15,
		},
		{
			"500 nm in vacuum", must(physics.LightFrequency(units.NewDistance(500, distance.Nanometer))),
			5.99584916e14, "Hz", 1e-15,
		},
		{
			"a 1239.84198 nm photon", must(physics.PhotonEnergy(units.NewDistance(1239.84198, distance.Nanometer))),
			1.602176634e-19, "J", 1e-8, // hc = 1239.841 98… eV·nm
		},
		{"sound at 0 °C", physics.SpeedOfSoundInAir(kelvin(273.15)), 331.3, "m/s", 1e-3},
		{"sound at 20 °C", physics.SpeedOfSoundInAir(kelvin(293.15)), 343.2, "m/s", 1e-3},
		{
			"1 s at 0.6 c",
			must(physics.TimeDilation(units.NewTime(1, time.Second), units.NewVelocity(0.6, velocity.SpeedOfLight))),
			1.25, "s", 1e-15,
		},
		{
			"1 mol at 0 °C in 22.4 L",
			must(physics.IdealGasPressure(units.NewAmount(1, amount.Mole), kelvin(273.15), molarVolume)),
			101325, "Pa", 1e-9,
		},
		{
			"1 mol at 0 °C and 1 atm",
			must(physics.IdealGasVolume(units.NewAmount(1, amount.Mole), kelvin(273.15), atm)),
			22.41396954e-3, "m³", 1e-9,
		},
		{
			"22.4 L at 0 °C and 1 atm",
			must(physics.IdealGasAmount(atm, molarVolume, kelvin(273.15))),
			1, "mol", 1e-9,
		},
		{
			"the Sun's luminosity",
			physics.BlackBodyPower(units.NewDistance(constant.NominalSolarRadius, distance.Meter).SphereArea(),
				kelvin(5772)),
			constant.NominalSolarLuminosity, "W", 1e-3, // IAU 2015: from R☉, T_eff = 5772 K and σ
		},
		{"the Sun's peak", must(physics.WienPeakWavelength(kelvin(5772))), 502.04e-9, "m", 1e-4},
		{
			"1 F at 2 V",
			physics.CapacitorEnergy(units.NewCapacitance(1, capacitance.Farad), units.NewVoltage(2, voltage.Volt)),
			2, "J", 1e-15,
		},
		{
			"2 H at 3 A",
			physics.InductorEnergy(units.NewInductance(2, inductance.Henry), units.NewCurrent(3, current.Ampere)),
			9, "J", 1e-15,
		},
		{
			"1 mH and 1 μF",
			must(physics.ResonantFrequency(units.NewInductance(1, inductance.Millihenry),
				units.NewCapacitance(1, capacitance.Microfarad))),
			5032.9212104487, "Hz", 1e-12,
		},
	}

	for _, c := range cases {
		if c.got.Symbol() != c.symbol || !near(c.got.Value(), c.want, c.tol) {
			t.Errorf("%s = %v %s, want %v %s", c.name, c.got.Value(), c.got.Symbol(), c.want, c.symbol)
		}
	}

	mach, err := physics.MachNumber(units.NewVelocity(2*343.2, velocity.MeterPerSecond), kelvin(293.15))
	if err != nil || !near(mach, 2, 1e-3) {
		t.Errorf("MachNumber = %v, %v; want 2", mach, err)
	}
}

func TestFormulas_Errors(t *testing.T) {
	m := units.NewDistance(1, distance.Meter)
	up := units.NewAcceleration(-1, acceleration.StandardGravity)
	noG := units.NewAcceleration(0, acceleration.MeterPerSecondSquared)
	kg := units.NewMass(1, mass.Kilogram)
	noM := units.NewDistance(0, distance.Meter)
	zeroK := kelvin(0)

	cases := []struct {
		name string
		err  error
		want units.Error
	}{
		{"FreeFallTime, no gravity", errOf(physics.FreeFallTime(m, noG)), units.ErrZeroDivision},
		{"FreeFallTime, upwards", errOf(physics.FreeFallTime(m, up)), units.ErrOutOfRange},
		{"FreeFallSpeed, upwards", errOf(physics.FreeFallSpeed(m, up)), units.ErrOutOfRange},
		{
			"ProjectileRange, no gravity",
			errOf(physics.ProjectileRange(units.NewVelocity(1, velocity.MeterPerSecond),
				units.NewAngle(1, angle.Radian), noG)),
			units.ErrZeroDivision,
		},
		{"PendulumPeriod, no gravity", errOf(physics.PendulumPeriod(m, noG)), units.ErrZeroDivision},
		{"PendulumPeriod, upwards", errOf(physics.PendulumPeriod(m, up)), units.ErrOutOfRange},
		{"GravitationalForce at 0 m", errOf(physics.GravitationalForce(kg, kg, noM)), units.ErrZeroDivision},
		{"SurfaceGravity at 0 m", errOf(physics.SurfaceGravity(kg, noM)), units.ErrZeroDivision},
		{"EscapeVelocity at 0 m", errOf(physics.EscapeVelocity(kg, noM)), units.ErrZeroDivision},
		{"EscapeVelocity, negative mass", errOf(physics.EscapeVelocity(kg.Mul(-1), m)), units.ErrOutOfRange},
		{"OrbitalVelocity at 0 m", errOf(physics.OrbitalVelocity(kg, noM)), units.ErrZeroDivision},
		{"OrbitalPeriod, no mass", errOf(physics.OrbitalPeriod(kg.Mul(0), m)), units.ErrZeroDivision},
		{"OrbitalPeriod, negative mass", errOf(physics.OrbitalPeriod(kg.Mul(-1), m)), units.ErrOutOfRange},
		{
			"LightWavelength of 0 Hz", errOf(physics.LightWavelength(units.NewFrequency(0, frequency.Hertz))),
			units.ErrZeroDivision,
		},
		{"LightFrequency of 0 m", errOf(physics.LightFrequency(noM)), units.ErrZeroDivision},
		{"PhotonEnergy of 0 m", errOf(physics.PhotonEnergy(noM)), units.ErrZeroDivision},
		{
			"MachNumber at 0 K",
			errOf(physics.MachNumber(units.NewVelocity(1, velocity.MeterPerSecond), zeroK)),
			units.ErrZeroDivision,
		},
		{
			"TimeDilation at c",
			errOf(physics.TimeDilation(units.NewTime(1, time.Second), units.NewVelocity(1, velocity.SpeedOfLight))),
			units.ErrOutOfRange,
		},
		{
			"IdealGasPressure in 0 L",
			errOf(physics.IdealGasPressure(units.NewAmount(1, amount.Mole), kelvin(300),
				units.NewVolume(0, volume.Liter))),
			units.ErrZeroDivision,
		},
		{
			"IdealGasVolume at 0 Pa",
			errOf(physics.IdealGasVolume(units.NewAmount(1, amount.Mole), kelvin(300),
				units.NewPressure(0, pressure.Pascal))),
			units.ErrZeroDivision,
		},
		{
			"IdealGasAmount at 0 K",
			errOf(physics.IdealGasAmount(units.NewPressure(1, pressure.Pascal), units.NewVolume(1, volume.Liter),
				zeroK)),
			units.ErrZeroDivision,
		},
		{"WienPeakWavelength at 0 K", errOf(physics.WienPeakWavelength(zeroK)), units.ErrZeroDivision},
		{
			"ResonantFrequency, no capacitance",
			errOf(physics.ResonantFrequency(units.NewInductance(1, inductance.Henry),
				units.NewCapacitance(0, capacitance.Farad))),
			units.ErrZeroDivision,
		},
		{
			"ResonantFrequency, negative capacitance",
			errOf(physics.ResonantFrequency(units.NewInductance(1, inductance.Henry),
				units.NewCapacitance(-1, capacitance.Farad))),
			units.ErrOutOfRange,
		},
	}

	if err := errOf(physics.FreeFallTime(m, up)); err == nil ||
		err.Error() != "the distance and the acceleration have opposite signs: the result isn't real" {
		t.Errorf("FreeFallTime upwards: %v", err)
	}

	sentinels := []units.Error{
		units.ErrZeroDivision, units.ErrOutOfRange, units.ErrUnknownUnit,
		units.ErrIncompatibleMeasures, units.ErrInvalidFormat,
	}

	for _, c := range cases {
		if c.err == nil {
			t.Errorf("%s: no error, want %v", c.name, c.want)

			continue
		}

		if c.err.Error() == "" {
			t.Errorf("%s: the error has no message", c.name)
		}

		for _, s := range sentinels {
			if errors.Is(c.err, s) != (s == c.want) {
				t.Errorf("%s: %v; errors.Is(%v) = %v", c.name, c.err, s, errors.Is(c.err, s))
			}
		}
	}
}

func errOf[T any](_ T, err error) error {
	return err
}
