package conversion

// Factors shows how many base measure units in specified unit.
type Factors = map[string]float64

func CopyConversionFactors(cf Factors) Factors {
	copiedCf := make(Factors)
	for u, f := range cf {
		copiedCf[u] = f
	}

	return copiedCf
}
