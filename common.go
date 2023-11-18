package units

func copyConversionFactors(cf ConversionFactors) ConversionFactors {
	copiedCf := make(ConversionFactors)
	for u, f := range cf {
		copiedCf[u] = f
	}

	return copiedCf
}
