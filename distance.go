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

const DistanceBaseUnit = Meter

// distanceConversionFactors shows how many specified units in base unit of distance - m (meter)
var distanceConversionFactors = ConversionFactors{
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

type Distance interface {
	Measurement
}

func NewDistance(value float64, unit string) (Distance, error) {
	if _, ok := distanceConversionFactors[unit]; !ok {
		return nil, fmt.Errorf("unknown Distance unit %s", unit)
	}

	return NewBaseMeasurement(value, unit, distanceConversionFactors)
}
