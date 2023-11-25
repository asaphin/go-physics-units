package distance

import (
	"github.com/asaphin/go-physics-units/conversion"
	"github.com/asaphin/go-physics-units/internal/immutable"
	"github.com/asaphin/go-physics-units/internal/rates"
	"sync"
)

const BaseUnit = Meter

const (
	PlankLength      = "lp"
	Femtometer       = "fm"
	Picometer        = "pm"
	Angstrom         = "Å"
	Nanometer        = "nm"
	Micrometer       = "μm"
	Millimeter       = "mm"
	Centimeter       = "cm"
	Decimeter        = "dm"
	Meter            = "m"
	Kilometer        = "km"
	Inch             = "in"
	Feet             = "ft"
	Yard             = "yd"
	Mile             = "mi"
	NauticalMile     = "nmi"
	AstronomicalUnit = "AU"
	LightYear        = "ly"
	Parsec           = "pc"
)

var conversionFactors = conversion.Factors{
	PlankLength:      1.616255e-35,
	Femtometer:       1e-15,
	Picometer:        1e-12,
	Angstrom:         1e-10,
	Nanometer:        1e-9,
	Micrometer:       1e-6,
	Millimeter:       0.001,
	Centimeter:       0.01,
	Decimeter:        0.1,
	Meter:            1,
	Kilometer:        1000,
	Inch:             0.0254,
	Feet:             0.3048,
	Yard:             0.9144,
	Mile:             1609.344,
	NauticalMile:     1852,
	AstronomicalUnit: 1.495978707e11,
	LightYear:        9.4607e15,
	Parsec:           3.0857e16,
}

// ConversionFactors shows how many base units (m) in specified unit.
func ConversionFactors() conversion.Factors {
	return conversion.CopyConversionFactors(conversionFactors)
}

var conversionRates immutable.Float64Map
var conversionRatesSync sync.Once

// ConversionRates returns pointer to conversion rates storage.
// Rates stored by composite keys unitFrom + unitTo
func ConversionRates() immutable.Float64Map {
	conversionRatesSync.Do(func() {
		conversionRates = immutable.MakeImmutable(rates.FactorsToRates(conversionFactors))
	})

	return conversionRates
}
