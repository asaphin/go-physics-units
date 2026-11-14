package constant_test

import (
	"math"
	"testing"

	"github.com/asaphin/go-physics-units/acceleration"
	"github.com/asaphin/go-physics-units/constant"
	"github.com/asaphin/go-physics-units/distance"
	"github.com/asaphin/go-physics-units/energy"
	"github.com/asaphin/go-physics-units/mass"
	"github.com/asaphin/go-physics-units/power"
	"github.com/asaphin/go-physics-units/pressure"
	"github.com/asaphin/go-physics-units/velocity"
)

func near(a, b, tol float64) bool {
	return math.Abs(a-b) <= tol*math.Max(math.Abs(a), math.Abs(b))
}

// TestDerivedConstants checks the derived constants against the values CODATA
// 2022 publishes for them.
func TestDerivedConstants(t *testing.T) {
	cases := []struct {
		name      string
		got, want float64
	}{
		{"ħ", constant.ReducedPlanckConstant, 1.054571817e-34},
		{"R", constant.MolarGasConstant, 8.314462618},
		{"F", constant.FaradayConstant, 96485.33212},
		{"σ", constant.StefanBoltzmannConstant, 5.670374419e-8},
		{"ly", constant.LightYear, 9460730472580800},
		{"pc", constant.Parsec, 3.0856775814913673e16},
	}

	for _, c := range cases {
		if !near(c.got, c.want, 1e-9) {
			t.Errorf("%s = %v, want %v", c.name, c.got, c.want)
		}
	}

	// ε₀·μ₀·c² = 1, to the precision of the measured ε₀ and μ₀.
	if got := constant.VacuumPermittivity * constant.VacuumPermeability * constant.SpeedOfLight *
		constant.SpeedOfLight; !near(got, 1, 1e-9) {
		t.Errorf("ε₀·μ₀·c² = %v, want 1", got)
	}
}

// TestConstantsMatchUnits checks that the unit packages and this package agree.
func TestConstantsMatchUnits(t *testing.T) {
	cases := []struct {
		name         string
		factor, want float64
	}{
		{"velocity.SpeedOfLight", velocity.SpeedOfLight.Factor(), constant.SpeedOfLight},
		{"acceleration.StandardGravity", acceleration.StandardGravity.Factor(), constant.StandardGravity},
		{"pressure.Atmosphere", pressure.Atmosphere.Factor(), constant.StandardAtmosphere},
		{"energy.Electronvolt", energy.Electronvolt.Factor(), constant.ElementaryCharge},
		{"mass.Dalton", mass.Dalton.Factor(), constant.AtomicMassConstant},
		{"distance.AstronomicalUnit", distance.AstronomicalUnit.Factor(), constant.AstronomicalUnit},
		{"distance.LightYear", distance.LightYear.Factor(), constant.LightYear},
		{"distance.Parsec", distance.Parsec.Factor(), constant.Parsec},
		{"power.SolarLuminosity", power.SolarLuminosity.Factor(), constant.NominalSolarLuminosity},
	}

	for _, c := range cases {
		if c.factor != c.want {
			t.Errorf("%s = %v, constant = %v", c.name, c.factor, c.want)
		}
	}
}
