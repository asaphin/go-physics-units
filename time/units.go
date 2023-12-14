package time

import (
	"github.com/asaphin/go-physics-units/conversion"
	"github.com/asaphin/go-physics-units/internal/immutable"
	"github.com/asaphin/go-physics-units/internal/rates"
)

const BaseUnit = Second

const (
	PlankTime   = "tp"
	Nanosecond  = "ns"
	Microsecond = "μs"
	Millisecond = "ms"
	Second      = "s"
	Minute      = "min"
	Hour        = "h"
	Day         = "d"
	Week        = "wk"
	Year        = "yr"
)

var conversionFactors = conversion.Factors{
	PlankTime:   5.39e-44,
	Nanosecond:  1e-9,
	Microsecond: 1e-6,
	Millisecond: 0.001,
	Second:      1,
	Minute:      60,
	Hour:        3600,
	Day:         86400,
	Week:        604800,
	Year:        31557600,
}

// ConversionFactors shows how many base units (s) in specified unit.
func ConversionFactors() conversion.Factors {
	return conversion.CopyConversionFactors(conversionFactors)
}

var conversionRates immutable.Float64Map

// ConversionRates returns pointer to conversion rates storage.
// Rates stored by composite keys unitFrom + unitTo.
func ConversionRates() immutable.Float64Map {
	if conversionRates == nil {
		conversionRates = immutable.MakeImmutable(rates.FactorsToRates(conversionFactors))
	}

	return conversionRates
}
