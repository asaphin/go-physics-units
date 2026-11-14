package physics_test

import (
	"fmt"

	units "github.com/asaphin/go-physics-units"
	"github.com/asaphin/go-physics-units/acceleration"
	"github.com/asaphin/go-physics-units/amount"
	"github.com/asaphin/go-physics-units/angle"
	"github.com/asaphin/go-physics-units/area"
	"github.com/asaphin/go-physics-units/capacitance"
	"github.com/asaphin/go-physics-units/current"
	"github.com/asaphin/go-physics-units/density"
	"github.com/asaphin/go-physics-units/distance"
	"github.com/asaphin/go-physics-units/energy"
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

// earth is the Earth's standard gravity.
var earth = units.NewAcceleration(1, acceleration.StandardGravity)

func celsius(c float64) units.Temperature {
	t, err := units.NewTemperature(c, temperature.Celsius)
	if err != nil {
		panic(err) // below absolute zero
	}

	return t
}

func ExampleFreeFallTime() {
	// A stone dropped into a 45 m well.
	t, err := physics.FreeFallTime(units.NewDistance(45, distance.Meter), earth)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%.2f\n", t)
	// Output: 3.03 s
}

func ExampleFreeFallSpeed() {
	// Diving from the 10 m platform.
	v, err := physics.FreeFallSpeed(units.NewDistance(10, distance.Meter), earth)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%.1f\n", v.ConvertTo(velocity.KilometerPerHour))
	// Output: 50.4 km/h
}

func ExampleFallDistance() {
	fmt.Printf("%.1f\n", physics.FallDistance(units.NewTime(3, time.Second), earth))
	// Output: 44.1 m
}

func ExamplePotentialEnergy() {
	// A 1000 kg lift raised 30 m.
	e := physics.PotentialEnergy(units.NewMass(1000, mass.Kilogram), earth, units.NewDistance(30, distance.Meter))

	fmt.Printf("%.0f\n", e.ConvertTo(energy.Kilojoule))
	// Output: 294 kJ
}

func ExampleProjectileRange() {
	r, err := physics.ProjectileRange(units.NewVelocity(30, velocity.MeterPerSecond),
		units.NewAngle(30, angle.Degree), earth)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%.1f\n", r)
	// Output: 79.5 m
}

func ExamplePendulumPeriod() {
	// The Foucault pendulum of the Panthéon in Paris is 67 m long.
	t, err := physics.PendulumPeriod(units.NewDistance(67, distance.Meter), earth)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%.1f\n", t)
	// Output: 16.4 s
}

func ExampleHydrostaticPressure() {
	// The bottom of the Mariana Trench, ignoring that sea water compresses.
	seawater := units.NewDensity(1025, density.KilogramPerCubicMeter)
	p := physics.HydrostaticPressure(seawater, earth, units.NewDistance(10935, distance.Meter))

	fmt.Printf("%.0f\n", p.ConvertTo(pressure.Bar))
	// Output: 1099 bar
}

func ExampleGravitationalForce() {
	moon := units.NewMass(7.342e22, mass.Kilogram)

	f, err := physics.GravitationalForce(units.NewMass(1, mass.EarthMass), moon,
		units.NewDistance(384400, distance.Kilometer))
	if err != nil {
		panic(err)
	}

	fmt.Printf("%.2e\n", f)
	// Output: 1.98e+20 N
}

func ExampleSurfaceGravity() {
	moon := units.NewMass(7.342e22, mass.Kilogram)

	g, err := physics.SurfaceGravity(moon, units.NewDistance(1737.4, distance.Kilometer))
	if err != nil {
		panic(err)
	}

	fmt.Printf("%.2f\n", g)
	// Output: 1.62 m/s²
}

func ExampleEscapeVelocity() {
	v, err := physics.EscapeVelocity(units.NewMass(1, mass.EarthMass), units.NewDistance(6371, distance.Kilometer))
	if err != nil {
		panic(err)
	}

	fmt.Printf("%.1f\n", v.ConvertTo(velocity.KilometerPerSecond))
	// Output: 11.2 km/s
}

func ExampleOrbitalVelocity() {
	// The International Space Station, about 420 km above the surface.
	iss := units.NewDistance(6371+420, distance.Kilometer)

	v, err := physics.OrbitalVelocity(units.NewMass(1, mass.EarthMass), iss)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%.2f\n", v.ConvertTo(velocity.KilometerPerSecond))
	// Output: 7.66 km/s
}

func ExampleOrbitalPeriod() {
	// A geostationary orbit goes round once a sidereal day.
	t, err := physics.OrbitalPeriod(units.NewMass(1, mass.EarthMass), units.NewDistance(42164, distance.Kilometer))
	if err != nil {
		panic(err)
	}

	fmt.Printf("%.2f\n", t.ConvertTo(time.Hour))
	// Output: 23.93 h
}

func ExampleSchwarzschildRadius() {
	// The Earth squeezed into a black hole.
	r := physics.SchwarzschildRadius(units.NewMass(1, mass.EarthMass))

	fmt.Printf("%.2f\n", r.ConvertTo(distance.Millimeter))
	// Output: 8.87 mm
}

func ExampleLightWavelength() {
	// An FM radio station.
	l, err := physics.LightWavelength(units.NewFrequency(100, frequency.Megahertz))
	if err != nil {
		panic(err)
	}

	fmt.Printf("%.2f\n", l)
	// Output: 3.00 m
}

func ExampleLightFrequency() {
	// A red laser pointer.
	f, err := physics.LightFrequency(units.NewDistance(650, distance.Nanometer))
	if err != nil {
		panic(err)
	}

	fmt.Printf("%.0f\n", f.ConvertTo(frequency.Terahertz))
	// Output: 461 THz
}

func ExamplePhotonEnergy() {
	// A violet photon.
	e, err := physics.PhotonEnergy(units.NewDistance(405, distance.Nanometer))
	if err != nil {
		panic(err)
	}

	fmt.Printf("%.2f\n", e.ConvertTo(energy.Electronvolt))
	// Output: 3.06 eV
}

func ExampleSpeedOfSoundInAir() {
	fmt.Printf("%.1f\n", physics.SpeedOfSoundInAir(celsius(20)))
	// Output: 343.2 m/s
}

func ExampleMachNumber() {
	// Concorde cruised at about 2180 km/h, where the air is at -56.5 °C.
	mach, err := physics.MachNumber(units.NewVelocity(2180, velocity.KilometerPerHour), celsius(-56.5))
	if err != nil {
		panic(err)
	}

	fmt.Printf("Mach %.2f\n", mach)
	// Output: Mach 2.05
}

func ExampleTimeDilation() {
	// A muon lives 2.2 μs in its own frame; at 0.998 c it lives far longer.
	life, err := physics.TimeDilation(units.NewTime(2.2, time.Microsecond),
		units.NewVelocity(0.998, velocity.SpeedOfLight))
	if err != nil {
		panic(err)
	}

	fmt.Printf("%.1f\n", life)
	// Output: 34.8 μs
}

func ExampleIdealGasPressure() {
	// A 12 L scuba tank holding 100 mol of air at 15 °C.
	p, err := physics.IdealGasPressure(units.NewAmount(100, amount.Mole), celsius(15),
		units.NewVolume(12, volume.Liter))
	if err != nil {
		panic(err)
	}

	fmt.Printf("%.0f\n", p.ConvertTo(pressure.Bar))
	// Output: 200 bar
}

func ExampleIdealGasVolume() {
	v, err := physics.IdealGasVolume(units.NewAmount(1, amount.Mole), celsius(25),
		units.NewPressure(1, pressure.Atmosphere))
	if err != nil {
		panic(err)
	}

	fmt.Printf("%.2f\n", v.ConvertTo(volume.Liter))
	// Output: 24.47 L
}

func ExampleIdealGasAmount() {
	// A 50 L balloon at 1 atm and 20 °C.
	n, err := physics.IdealGasAmount(units.NewPressure(1, pressure.Atmosphere), units.NewVolume(50, volume.Liter),
		celsius(20))
	if err != nil {
		panic(err)
	}

	fmt.Printf("%.2f\n", n)
	// Output: 2.08 mol
}

func ExampleBlackBodyPower() {
	// A person's 1.8 m² of skin at 33 °C radiates this much, and absorbs most
	// of it back from a room at a similar temperature.
	p := physics.BlackBodyPower(units.NewArea(1.8, area.SquareMeter), celsius(33))

	fmt.Printf("%.0f\n", p)
	// Output: 897 W
}

func ExampleWienPeakWavelength() {
	// The human body glows in the infrared.
	l, err := physics.WienPeakWavelength(celsius(37))
	if err != nil {
		panic(err)
	}

	fmt.Printf("%.2f\n", l.ConvertTo(distance.Micrometer))
	// Output: 9.34 μm
}

func ExampleCapacitorEnergy() {
	// A camera flash capacitor.
	e := physics.CapacitorEnergy(units.NewCapacitance(330, capacitance.Microfarad), units.NewVoltage(300, voltage.Volt))

	fmt.Printf("%.2f\n", e)
	// Output: 14.85 J
}

func ExampleInductorEnergy() {
	e := physics.InductorEnergy(units.NewInductance(10, inductance.Millihenry), units.NewCurrent(2, current.Ampere))

	fmt.Printf("%.0f\n", e.ConvertTo(energy.Millijoule))
	// Output: 20 mJ
}

func ExampleResonantFrequency() {
	// The tuning circuit of an AM radio.
	f, err := physics.ResonantFrequency(units.NewInductance(250, inductance.Microhenry),
		units.NewCapacitance(100, capacitance.Picofarad))
	if err != nil {
		panic(err)
	}

	fmt.Printf("%.0f\n", f.ConvertTo(frequency.Kilohertz))
	// Output: 1007 kHz
}
