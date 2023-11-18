package units

import "fmt"

const (
	MetersPerSecond   = "m/s"
	KilometersPerHour = "km/h"
	FeetPerSecond     = "ft/s"
	MilesPerHour      = "mph"
)

const VelocityBaseUnit = MetersPerSecond

var velocityConversionFactors = ConversionFactors{
	MetersPerSecond:   1,
	KilometersPerHour: 3.6,
	FeetPerSecond:     3.28084,
	MilesPerHour:      2.23694,
}

type Velocity interface {
	Measurement
}

func NewVelocity(value float64, unit string) (Velocity, error) {
	if _, ok := velocityConversionFactors[unit]; !ok {
		return nil, fmt.Errorf("unknown velocity unit %s", unit)
	}

	return NewBaseMeasurement(value, unit, velocityConversionFactors)
}

func NewVelocityFromDistanceAntTime(d Distance, t Time) Velocity {
	d, _ = d.ConvertTo(DistanceBaseUnit)
	t, _ = t.ConvertTo(TimeBaseUnit)

	v, _ := NewVelocity(d.Value()/t.Value(), VelocityBaseUnit)

	return v
}
