package mass

import "github.com/asaphin/go-physics-units/conversion-factors"

const BaseUnit = Kilogram

const (
	Gram     = "g"
	Kilogram = "kg"
	Pound    = "lb"
	Ounce    = "oz"
)

var massConversionFactors = conversion.Factors{
	Gram:     0.001,
	Kilogram: 1,
	Pound:    0.453592,
	Ounce:    0.0283495,
}

// ConversionFactors shows how many specified units in the base unit of mass - kilogram.
func ConversionFactors() conversion.Factors {
	return conversion.CopyConversionFactors(massConversionFactors)
}
