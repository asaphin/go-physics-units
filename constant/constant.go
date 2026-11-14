// Package constant defines physical constants in SI units, as untyped Go
// constants: the defining constants of the SI, which are exact, constants
// derived from them, the CODATA 2022 recommended values of measured
// constants, and astronomical constants of the IAU.
//
// They are constants, not variables, so no importer can change them. Combine
// them with the measures of the units package through the SI base units:
//
//	c := units.NewVelocity(constant.SpeedOfLight, velocity.MeterPerSecond)
//
// Each comment gives the constant's symbol and SI unit, and whether it is
// exact. A measured constant's comment gives its relative standard
// uncertainty: the precision of the value, not of its float64 representation.
package constant

import "math"

// The defining constants of the SI (the SI Brochure, 9th edition, 2019).
// They are exact by definition.
const (
	SpeedOfLight      = 299792458       // c, m/s
	PlanckConstant    = 6.62607015e-34  // h, J·s
	ElementaryCharge  = 1.602176634e-19 // e, C
	BoltzmannConstant = 1.380649e-23    // k, J/K
	AvogadroConstant  = 6.02214076e23   // N_A, 1/mol
)

// Constants derived from the defining constants. They are exact, but for π,
// which Go's constant arithmetic carries to far more digits than a float64.
const (
	// ReducedPlanckConstant is ħ = h/2π, in J·s.
	ReducedPlanckConstant = PlanckConstant / (2 * math.Pi)
	// MolarGasConstant is R = N_A·k, in J/(mol·K).
	MolarGasConstant = AvogadroConstant * BoltzmannConstant
	// FaradayConstant is F = N_A·e, in C/mol.
	FaradayConstant = AvogadroConstant * ElementaryCharge
	// StefanBoltzmannConstant is σ = 2π⁵k⁴/(15h³c²), in W/(m²·K⁴).
	StefanBoltzmannConstant = 2 * math.Pi * math.Pi * math.Pi * math.Pi * math.Pi *
		BoltzmannConstant * BoltzmannConstant * BoltzmannConstant * BoltzmannConstant /
		(15 * PlanckConstant * PlanckConstant * PlanckConstant * SpeedOfLight * SpeedOfLight)
)

// Conventional constants, exact by definition.
const (
	// StandardGravity is gₙ, in m/s² (3rd CGPM, 1901).
	StandardGravity = 9.80665
	// StandardAtmosphere is the standard atmosphere, in Pa (10th CGPM, 1954).
	StandardAtmosphere = 101325
)

// Measured constants: the CODATA 2022 recommended values, with their relative
// standard uncertainties.
const (
	// GravitationalConstant is G, in m³/(kg·s²); relative uncertainty 2.2e-5.
	GravitationalConstant = 6.67430e-11
	// ElectronMass is mₑ, in kg; relative uncertainty 3.1e-10.
	ElectronMass = 9.1093837139e-31
	// ProtonMass is mₚ, in kg; relative uncertainty 3.1e-10.
	ProtonMass = 1.67262192595e-27
	// AtomicMassConstant is u, the dalton, in kg; relative uncertainty 3.1e-10.
	AtomicMassConstant = 1.66053906892e-27
	// FineStructureConstant is α, dimensionless; relative uncertainty 1.6e-10.
	FineStructureConstant = 7.2973525643e-3
	// VacuumPermittivity is ε₀, in F/m; relative uncertainty 1.6e-10.
	VacuumPermittivity = 8.8541878188e-12
	// VacuumPermeability is μ₀, in N/A²; relative uncertainty 1.6e-10.
	VacuumPermeability = 1.25663706127e-6
)

// Astronomical constants of the IAU. The astronomical unit, the Julian year,
// the light-year and the parsec are exact by definition; the nominal values
// (IAU 2015 Resolution B3) are conversion constants, not measurements.
const (
	// AstronomicalUnit is au, in m (IAU 2012 Resolution B2).
	AstronomicalUnit = 149597870700
	// JulianYear is 365.25 days, in s.
	JulianYear = 31557600
	// LightYear is the distance light travels in a Julian year, in m.
	LightYear = SpeedOfLight * JulianYear
	// Parsec is au·648000/π, in m (IAU 2015 Resolution B2).
	Parsec = AstronomicalUnit * 648000 / math.Pi
	// NominalSolarMassParameter is GM☉, in m³/s². It is known far more
	// precisely than the solar mass, because G isn't.
	NominalSolarMassParameter = 1.3271244e20
	// NominalEarthMassParameter is GM⊕, in m³/s².
	NominalEarthMassParameter = 3.986004e14
	// NominalSolarRadius is R☉, in m.
	NominalSolarRadius = 6.957e8
	// NominalSolarLuminosity is L☉, in W.
	NominalSolarLuminosity = 3.828e26
	// NominalEarthEquatorialRadius is the Earth's equatorial radius, in m.
	NominalEarthEquatorialRadius = 6.3781e6
	// NominalEarthPolarRadius is the Earth's polar radius, in m.
	NominalEarthPolarRadius = 6.3568e6
)
