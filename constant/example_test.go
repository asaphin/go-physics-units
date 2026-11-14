package constant_test

import (
	"fmt"

	units "github.com/asaphin/go-physics-units"
	"github.com/asaphin/go-physics-units/constant"
	"github.com/asaphin/go-physics-units/velocity"
)

// The constants are in SI base units, so they combine with the measures
// through their base units.
func Example() {
	c := units.NewVelocity(constant.SpeedOfLight, velocity.MeterPerSecond)

	fmt.Println(c.ConvertTo(velocity.KilometerPerSecond))
	fmt.Printf("R = %.4f J/(mol·K)\n", constant.MolarGasConstant)
	// Output:
	// 299792.458 km/s
	// R = 8.3145 J/(mol·K)
}
