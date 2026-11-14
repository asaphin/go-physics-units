// Package physics provides formulas that combine several measures, rest on an
// approximation, or use a measured constant: falling bodies, projectiles and
// pendulums, gravitation and orbits, light and sound, relativity, ideal gases,
// thermal radiation and circuits.
//
// Relations between two measures that hold exactly, such as
// Distance.DivideByTime or Mass.RestEnergy, are methods of the measures in
// package units. The functions here take and return the same typed values:
//
//	g := units.NewAcceleration(1, acceleration.StandardGravity)
//	t, err := physics.FreeFallTime(units.NewDistance(20, distance.Meter), g)
//
// They compute in SI base units and are safe for concurrent use. A function
// fails with [units.ErrZeroDivision] when a divisor is zero, and with
// [units.ErrOutOfRange] when the result has no real value, such as the time to
// fall a height against gravity.
package physics

import (
	"math"

	units "github.com/asaphin/go-physics-units"
	"github.com/asaphin/go-physics-units/acceleration"
	"github.com/asaphin/go-physics-units/amount"
	"github.com/asaphin/go-physics-units/constant"
	"github.com/asaphin/go-physics-units/distance"
	"github.com/asaphin/go-physics-units/energy"
	"github.com/asaphin/go-physics-units/force"
	"github.com/asaphin/go-physics-units/frequency"
	"github.com/asaphin/go-physics-units/power"
	"github.com/asaphin/go-physics-units/pressure"
	"github.com/asaphin/go-physics-units/time"
	"github.com/asaphin/go-physics-units/velocity"
	"github.com/asaphin/go-physics-units/volume"
)

// si returns the value of x in its SI base unit.
func si[T interface {
	ConvertToBaseUnits() T
	Value() float64
}](x T) float64 {
	return x.ConvertToBaseUnits().Value()
}

// rangeError is a [units.ErrOutOfRange] error with a fixed message. Its values
// are constants, so returning one doesn't allocate.
type rangeError string

const (
	errOppositeSigns rangeError = "the distance and the acceleration have opposite signs: the result isn't real"
	errNegativeMass  rangeError = "the mass and the distance have opposite signs: the result isn't real"
	errOppositeLC    rangeError = "the inductance and the capacitance have opposite signs: the result isn't real"
)

func (e rangeError) Error() string {
	return string(e)
}

func (rangeError) Unwrap() error {
	return units.ErrOutOfRange
}

// Falling bodies, projectiles and pendulums. Pass the Earth's standard gravity
// as units.NewAcceleration(1, acceleration.StandardGravity); none of these
// formulas includes air resistance.

// FreeFallTime returns the time √(2h/g) (s) in [time.BaseUnit] to fall the
// height h from rest under the gravity g. It fails with [units.ErrZeroDivision]
// if g is zero, and with [units.ErrOutOfRange] if h and g have opposite signs.
func FreeFallTime(h units.Distance, g units.Acceleration) (units.Time, error) {
	gv := si(g)
	if gv == 0 {
		return units.Time{}, units.ErrZeroDivision
	}

	x := 2 * si(h) / gv
	if x < 0 {
		return units.Time{}, errOppositeSigns
	}

	return units.NewTime(math.Sqrt(x), time.BaseUnit), nil
}

// FreeFallSpeed returns the speed √(2gh) (m/s) in [velocity.BaseUnit] after
// falling the height h from rest under the gravity g. It fails with
// [units.ErrOutOfRange] if h and g have opposite signs.
func FreeFallSpeed(h units.Distance, g units.Acceleration) (units.Velocity, error) {
	x := 2 * si(g) * si(h)
	if x < 0 {
		return units.Velocity{}, errOppositeSigns
	}

	return units.NewVelocity(math.Sqrt(x), velocity.BaseUnit), nil
}

// FallDistance returns the distance g·t²/2 (m) in [distance.BaseUnit] fallen from
// rest in the time t under the gravity g.
func FallDistance(t units.Time, g units.Acceleration) units.Distance {
	s := si(t)

	return units.NewDistance(si(g)*s*s/2, distance.BaseUnit)
}

// PotentialEnergy returns the gravitational potential energy m·g·h (J) in
// [energy.BaseUnit] of the mass m at the height h under the gravity g.
func PotentialEnergy(m units.Mass, g units.Acceleration, h units.Distance) units.Energy {
	return units.NewEnergy(si(m)*si(g)*si(h), energy.BaseUnit)
}

// ProjectileRange returns the distance v²·sin(2θ)/g (m) in [distance.BaseUnit]
// that a projectile launched at the speed v and the elevation θ covers before
// it comes back to its launch height. It fails with [units.ErrZeroDivision] if g
// is zero.
func ProjectileRange(v units.Velocity, elevation units.Angle, g units.Acceleration) (units.Distance, error) {
	gv := si(g)
	if gv == 0 {
		return units.Distance{}, units.ErrZeroDivision
	}

	s := si(v)

	return units.NewDistance(s*s*math.Sin(2*si(elevation))/gv, distance.BaseUnit), nil
}

// PendulumPeriod returns the period 2π√(L/g) (s) in [time.BaseUnit] of a simple
// pendulum of the length L under the gravity g, for small swings. It fails
// with [units.ErrZeroDivision] if g is zero, and with [units.ErrOutOfRange] if L
// and g have opposite signs.
func PendulumPeriod(length units.Distance, g units.Acceleration) (units.Time, error) {
	gv := si(g)
	if gv == 0 {
		return units.Time{}, units.ErrZeroDivision
	}

	x := si(length) / gv
	if x < 0 {
		return units.Time{}, errOppositeSigns
	}

	return units.NewTime(2*math.Pi*math.Sqrt(x), time.BaseUnit), nil
}

// HydrostaticPressure returns the pressure ρ·g·h (Pa) in [pressure.BaseUnit] at
// the depth h in a fluid of the density ρ under the gravity g, above the
// pressure at its surface.
func HydrostaticPressure(rho units.Density, g units.Acceleration, depth units.Distance) units.Pressure {
	return units.NewPressure(si(rho)*si(g)*si(depth), pressure.BaseUnit)
}

// Gravitation and orbits, with the Newtonian constant of gravitation G of
// CODATA 2022, known to a relative uncertainty of 2.2e-5. The masses of the
// Sun and the Earth in package mass carry that uncertainty too.

// GravitationalForce returns the force G·m₁·m₂/r² (N) in [force.BaseUnit]
// between two masses at the distance r, by Newton's law of gravitation. It
// fails with [units.ErrZeroDivision] if r is zero.
func GravitationalForce(m1, m2 units.Mass, r units.Distance) (units.Force, error) {
	d := si(r)
	if d == 0 {
		return units.Force{}, units.ErrZeroDivision
	}

	return units.NewForce(constant.GravitationalConstant*si(m1)*si(m2)/(d*d), force.BaseUnit), nil
}

// SurfaceGravity returns the acceleration G·m/r² (m/s²) in
// [acceleration.BaseUnit] at the distance r from the centre of a spherical body
// of the mass m: at its surface if r is its radius. It ignores the body's
// rotation. It fails with [units.ErrZeroDivision] if r is zero.
func SurfaceGravity(m units.Mass, r units.Distance) (units.Acceleration, error) {
	d := si(r)
	if d == 0 {
		return units.Acceleration{}, units.ErrZeroDivision
	}

	return units.NewAcceleration(constant.GravitationalConstant*si(m)/(d*d), acceleration.BaseUnit), nil
}

// EscapeVelocity returns the speed √(2G·m/r) (m/s) in [velocity.BaseUnit] needed
// to escape a body of the mass m from the distance r from its centre. It
// fails with [units.ErrZeroDivision] if r is zero, and with [units.ErrOutOfRange]
// if m and r have opposite signs.
func EscapeVelocity(m units.Mass, r units.Distance) (units.Velocity, error) {
	return orbitalSpeed(m, r, 2)
}

// OrbitalVelocity returns the speed √(G·m/r) (m/s) in [velocity.BaseUnit] of a
// circular orbit of the radius r around a body of the mass m. It fails with
// [units.ErrZeroDivision] if r is zero, and with [units.ErrOutOfRange] if m and r
// have opposite signs.
func OrbitalVelocity(m units.Mass, r units.Distance) (units.Velocity, error) {
	return orbitalSpeed(m, r, 1)
}

func orbitalSpeed(m units.Mass, r units.Distance, k float64) (units.Velocity, error) {
	d := si(r)
	if d == 0 {
		return units.Velocity{}, units.ErrZeroDivision
	}

	x := k * constant.GravitationalConstant * si(m) / d
	if x < 0 {
		return units.Velocity{}, errNegativeMass
	}

	return units.NewVelocity(math.Sqrt(x), velocity.BaseUnit), nil
}

// OrbitalPeriod returns the period 2π√(a³/(G·m)) (s) in [time.BaseUnit] of an
// orbit with the semi-major axis a around a body of the mass m, by Kepler's
// third law, for an orbiting body much lighter than m. It fails with
// [units.ErrZeroDivision] if m is zero, and with [units.ErrOutOfRange] if m and a
// have opposite signs.
func OrbitalPeriod(m units.Mass, a units.Distance) (units.Time, error) {
	gm := constant.GravitationalConstant * si(m)
	if gm == 0 {
		return units.Time{}, units.ErrZeroDivision
	}

	d := si(a)

	x := d * d * d / gm
	if x < 0 {
		return units.Time{}, errNegativeMass
	}

	return units.NewTime(2*math.Pi*math.Sqrt(x), time.BaseUnit), nil
}

// SchwarzschildRadius returns the radius 2G·m/c² (m) in [distance.BaseUnit] of
// the event horizon of a black hole of the mass m.
func SchwarzschildRadius(m units.Mass) units.Distance {
	return units.NewDistance(2*constant.GravitationalConstant*si(m)/
		(constant.SpeedOfLight*constant.SpeedOfLight), distance.BaseUnit)
}

// Light, sound and relativity.

// LightWavelength returns the wavelength c/f (m) in [distance.BaseUnit], in
// vacuum, of light or any electromagnetic wave of the frequency f. It fails
// with [units.ErrZeroDivision] if f is zero.
func LightWavelength(f units.Frequency) (units.Distance, error) {
	hz := si(f)
	if hz == 0 {
		return units.Distance{}, units.ErrZeroDivision
	}

	return units.NewDistance(constant.SpeedOfLight/hz, distance.BaseUnit), nil
}

// LightFrequency returns the frequency c/λ (Hz) in [frequency.BaseUnit] of light
// or any electromagnetic wave of the wavelength λ in vacuum. It fails with
// [units.ErrZeroDivision] if the wavelength is zero.
func LightFrequency(wavelength units.Distance) (units.Frequency, error) {
	m := si(wavelength)
	if m == 0 {
		return units.Frequency{}, units.ErrZeroDivision
	}

	return units.NewFrequency(constant.SpeedOfLight/m, frequency.BaseUnit), nil
}

// PhotonEnergy returns the energy h·c/λ (J) in [energy.BaseUnit] of a photon of
// the wavelength λ in vacuum. It fails with [units.ErrZeroDivision] if the
// wavelength is zero.
func PhotonEnergy(wavelength units.Distance) (units.Energy, error) {
	m := si(wavelength)
	if m == 0 {
		return units.Energy{}, units.ErrZeroDivision
	}

	return units.NewEnergy(constant.PlanckConstant*constant.SpeedOfLight/m, energy.BaseUnit), nil
}

// Dry air as an ideal gas: its ratio of specific heats and its molar mass (ISO
// 2533, the standard atmosphere).
const (
	airHeatCapacityRatio = 1.4
	airMolarMass         = 28.96442e-3 // kg/mol
)

// SpeedOfSoundInAir returns the speed √(γ·R·T/M) (m/s) in [velocity.BaseUnit] of
// sound in dry air at the temperature t, as an ideal gas with γ = 1.4 and M =
// 28.96442 g/mol: about 331.3 m/s at 0 °C and 343.2 m/s at 20 °C. Humidity
// raises it by up to about 0.5 %.
func SpeedOfSoundInAir(t units.Temperature) units.Velocity {
	return units.NewVelocity(math.Sqrt(airHeatCapacityRatio*constant.MolarGasConstant*t.Kelvin()/airMolarMass),
		velocity.BaseUnit)
}

// MachNumber returns the speed v as a multiple of the speed of sound in dry air
// at the temperature t (see SpeedOfSoundInAir). It fails with
// [units.ErrZeroDivision] at absolute zero.
func MachNumber(v units.Velocity, t units.Temperature) (float64, error) {
	sound := si(SpeedOfSoundInAir(t))
	if sound == 0 {
		return 0, units.ErrZeroDivision
	}

	return si(v) / sound, nil
}

// TimeDilation returns the time γ·τ in the unit of properTime that passes for
// an observer at rest while a clock moving at the speed v measures the proper
// time τ (see Velocity.LorentzFactor). It fails with [units.ErrOutOfRange] unless
// the speed is below the speed of light.
func TimeDilation(properTime units.Time, v units.Velocity) (units.Time, error) {
	gamma, err := v.LorentzFactor()
	if err != nil {
		return units.Time{}, err
	}

	return properTime.Mul(gamma), nil
}

// Thermodynamics.

// IdealGasPressure returns the pressure n·R·T/V (Pa) in [pressure.BaseUnit] of the
// amount n of an ideal gas at the temperature t in the volume v. It fails with
// [units.ErrZeroDivision] if v is zero.
func IdealGasPressure(n units.Amount, t units.Temperature, v units.Volume) (units.Pressure, error) {
	m3 := si(v)
	if m3 == 0 {
		return units.Pressure{}, units.ErrZeroDivision
	}

	return units.NewPressure(si(n)*constant.MolarGasConstant*t.Kelvin()/m3, pressure.BaseUnit), nil
}

// IdealGasVolume returns the volume n·R·T/p (m³) in [volume.BaseUnit] of the
// amount n of an ideal gas at the temperature t and the pressure p. It fails
// with [units.ErrZeroDivision] if p is zero.
func IdealGasVolume(n units.Amount, t units.Temperature, p units.Pressure) (units.Volume, error) {
	pa := si(p)
	if pa == 0 {
		return units.Volume{}, units.ErrZeroDivision
	}

	return units.NewVolume(si(n)*constant.MolarGasConstant*t.Kelvin()/pa, volume.BaseUnit), nil
}

// IdealGasAmount returns the amount p·V/(R·T) (mol) in [amount.BaseUnit] of an
// ideal gas at the pressure p in the volume v at the temperature t. It fails
// with [units.ErrZeroDivision] at absolute zero.
func IdealGasAmount(p units.Pressure, v units.Volume, t units.Temperature) (units.Amount, error) {
	k := t.Kelvin()
	if k == 0 {
		return units.Amount{}, units.ErrZeroDivision
	}

	return units.NewAmount(si(p)*si(v)/(constant.MolarGasConstant*k), amount.BaseUnit), nil
}

// BlackBodyPower returns the power σ·A·T⁴ (W) in [power.BaseUnit] that a black
// body of the surface area a radiates at the temperature t, by the
// Stefan–Boltzmann law. Multiply it by the emissivity for a grey body.
func BlackBodyPower(a units.Area, t units.Temperature) units.Power {
	k := t.Kelvin()

	return units.NewPower(constant.StefanBoltzmannConstant*si(a)*k*k*k*k, power.BaseUnit)
}

// wienConstant is b, the Wien wavelength displacement law constant (CODATA
// 2022), in m·K. It is exact, as it follows from h, c and k.
const wienConstant = 2.897771955e-3

// WienPeakWavelength returns the wavelength b/T (m) in [distance.BaseUnit] at
// which a black body at the temperature t radiates most strongly, by Wien's
// displacement law. It fails with [units.ErrZeroDivision] at absolute zero.
func WienPeakWavelength(t units.Temperature) (units.Distance, error) {
	k := t.Kelvin()
	if k == 0 {
		return units.Distance{}, units.ErrZeroDivision
	}

	return units.NewDistance(wienConstant/k, distance.BaseUnit), nil
}

// Circuits.

// CapacitorEnergy returns the energy C·V²/2 (J) in [energy.BaseUnit] stored in
// the capacitance c charged to the voltage v.
func CapacitorEnergy(c units.Capacitance, v units.Voltage) units.Energy {
	u := si(v)

	return units.NewEnergy(si(c)*u*u/2, energy.BaseUnit)
}

// InductorEnergy returns the energy L·I²/2 (J) in [energy.BaseUnit] stored in the
// inductance l carrying the current i.
func InductorEnergy(l units.Inductance, i units.Current) units.Energy {
	a := si(i)

	return units.NewEnergy(si(l)*a*a/2, energy.BaseUnit)
}

// ResonantFrequency returns the frequency 1/(2π√(L·C)) (Hz) in
// [frequency.BaseUnit] of an LC circuit. It fails with [units.ErrZeroDivision] if
// l or c is zero, and with [units.ErrOutOfRange] if they have opposite signs.
func ResonantFrequency(l units.Inductance, c units.Capacitance) (units.Frequency, error) {
	lc := si(l) * si(c)
	if lc == 0 {
		return units.Frequency{}, units.ErrZeroDivision
	}

	if lc < 0 {
		return units.Frequency{}, errOppositeLC
	}

	return units.NewFrequency(1/(2*math.Pi*math.Sqrt(lc)), frequency.BaseUnit), nil
}
