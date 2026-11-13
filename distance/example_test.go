package distance_test

import (
	"fmt"

	"github.com/asaphin/go-physics-units/distance"
)

func ExampleRate() {
	r, ok := distance.Rate(distance.Mile, distance.Kilometer)
	fmt.Println(r, ok)

	// Both forms of the micro sign are accepted.
	r, ok = distance.Rate("\u00b5m", distance.Millimeter)
	fmt.Println(r, ok)

	_, ok = distance.Rate(distance.Mile, "furlong")
	fmt.Println(ok)
	// Output:
	// 1.609344 true
	// 0.001 true
	// false
}

func ExampleConversionFactors() {
	factors := distance.ConversionFactors() // a copy: changing it is safe

	fmt.Println(factors[distance.Kilometer], factors[distance.BaseUnit])
	// Output: 1000 1
}
