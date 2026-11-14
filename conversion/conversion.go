// Package conversion defines conversion factor tables, the data that describes
// the units of a measure.
package conversion

import "maps"

// Factors maps unit symbols to how many base units of the measure there are in
// one unit of that symbol. The base unit itself has the factor 1.
type Factors = map[string]float64

// CopyConversionFactors returns an independent copy of cf.
func CopyConversionFactors(cf Factors) Factors {
	copiedCf := make(Factors, len(cf))
	maps.Copy(copiedCf, cf)

	return copiedCf
}
