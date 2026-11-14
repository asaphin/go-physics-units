package distance_test

import (
	"fmt"

	"github.com/asaphin/go-physics-units/distance"
)

func ExampleRate() {
	fmt.Println(distance.Rate(distance.Mile, distance.Kilometer))
	fmt.Println(distance.Rate(distance.Micrometer, distance.Millimeter))
	// Output:
	// 1.609344
	// 0.001
}

func ExampleParseUnit() {
	u, ok := distance.ParseUnit("km")
	fmt.Println(u, ok, u == distance.Kilometer)

	// Both forms of the micro sign are accepted.
	u, ok = distance.ParseUnit("\u00b5m")
	fmt.Println(u == distance.Micrometer, ok)

	_, ok = distance.ParseUnit("furlong")
	fmt.Println(ok)
	// Output:
	// km true true
	// true true
	// false
}

func ExampleUnits() {
	for u := range distance.Units() {
		if u.Factor() >= 1000 && u.Factor() < 1e4 {
			fmt.Println(u, u.Factor())
		}
	}
	// Output:
	// km 1000
	// mi 1609.344
	// nmi 1852
}

func ExampleConversionFactors() {
	factors := distance.ConversionFactors() // a copy: changing it is safe

	fmt.Println(factors[distance.Kilometer.String()], factors[distance.BaseUnit.String()])
	// Output: 1000 1
}
