package units

import "github.com/asaphin/go-physics-units/conversion-factors"

type immutableConversionFactors struct {
	conversionFactors conversion.Factors
}

func newImmutableConversionFactors(cf conversion.Factors) *immutableConversionFactors {
	cfCopy := make(conversion.Factors)

	for u, f := range cf {
		cfCopy[u] = f
	}

	return &immutableConversionFactors{conversionFactors: cfCopy}
}

func (f *immutableConversionFactors) GetFactor(unit string) float64 {
	return f.conversionFactors[unit]
}

func (f *immutableConversionFactors) HasFactor(unit string) (float64, bool) {
	factor, ok := f.conversionFactors[unit]

	return factor, ok
}
