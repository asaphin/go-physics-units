package units

import "fmt"

const (
	Millimeter   = "mm"
	Centimeter   = "cm"
	Decimeter    = "dm"
	Meter        = "m"
	Kilometer    = "km"
	Inch         = "in"
	Feet         = "ft"
	Yard         = "yd"
	Mile         = "mi"
	NauticalMile = "nmi"
)

const LengthBaseUnit = Meter

// lengthConversionFactors shows how many specified units in base unit of length - m (meter)
var lengthConversionFactors = ConversionFactors{
	Millimeter:   0.0001,
	Centimeter:   0.01,
	Decimeter:    0.1,
	Meter:        1,
	Kilometer:    1000,
	Inch:         0.0254,
	Feet:         0.3048,
	Yard:         0.9144,
	Mile:         1609.344,
	NauticalMile: 1852,
}

type Length interface {
	Measurement
}

func NewLength(value float64, unit string) (Length, error) {
	if _, ok := lengthConversionFactors[unit]; !ok {
		return nil, fmt.Errorf("unknown length unit %s", unit)
	}

	return NewBaseMeasurement(value, unit, lengthConversionFactors)
}
