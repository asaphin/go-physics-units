package time

import "github.com/asaphin/go-physics-units/conversion-factors"

const BaseUnit = Second

const (
	Second = "s"
	Minute = "min"
	Hour   = "h"
	Day    = "d"
	Week   = "wk"
)

var timeConversionFactors = conversion.Factors{
	Second: 1,
	Minute: 60,
	Hour:   3600,
	Day:    86400,
	Week:   604800,
}

// ConversionFactors shows how many base time units (seconds) in specified unit.
func ConversionFactors() conversion.Factors {
	return conversion.CopyConversionFactors(timeConversionFactors)
}
