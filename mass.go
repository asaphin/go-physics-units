package units

import "fmt"

const (
	Gram      = "g"
	Kilogram  = "kg"
	Milligram = "mg"
	Microgram = "μg"
	Pound     = "lb"
	Ounce     = "oz"
	Ton       = "ton"
)

const MassBaseUnit = Kilogram

var massConversionFactors = ConversionFactors{
	Gram:      0.001,
	Kilogram:  1,
	Milligram: 1e-6,
	Microgram: 1e-9,
	Pound:     0.453592,
	Ounce:     0.0283495,
	Ton:       1000,
}

type Mass interface {
	Measurement
}

func NewMass(value float64, unit string) (Mass, error) {
	if _, ok := massConversionFactors[unit]; !ok {
		return nil, fmt.Errorf("unknown mass unit %s", unit)
	}

	return NewBaseMeasurement(value, unit, massConversionFactors)
}
