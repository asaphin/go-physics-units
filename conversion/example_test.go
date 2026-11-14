package conversion_test

import (
	"fmt"

	"github.com/asaphin/go-physics-units/conversion"
)

func ExampleCopyConversionFactors() {
	factors := conversion.Factors{"USD": 1, "EUR": 1.1}

	copied := conversion.CopyConversionFactors(factors)
	copied["EUR"] = 1.2

	fmt.Println(factors["EUR"], copied["EUR"])
	// Output: 1.1 1.2
}
