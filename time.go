package units

import "fmt"

const (
	Second = "s"
	Minute = "min"
	Hour   = "h"
	Day    = "d"
	Week   = "wk"
)

const TimeBaseUnit = Second

var timeConversionFactors = ConversionFactors{
	Second: 1,
	Minute: 60,
	Hour:   3600,
	Day:    86400,
	Week:   604800,
}

func TimeConversionFactors() ConversionFactors {
	return copyConversionFactors(timeConversionFactors)
}

type Time interface {
	Measurement
}

func NewTime(value float64, unit string) (Time, error) {
	if _, ok := timeConversionFactors[unit]; !ok {
		return nil, fmt.Errorf("unknown time unit %s", unit)
	}

	return NewBaseMeasurement(value, unit, timeConversionFactors)
}
