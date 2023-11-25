package units

import "github.com/asaphin/go-physics-units/conversion"

func getUnits(factors conversion.Factors) []string {
	units := make([]string, len(factors))

	i := 0
	for k := range factors {
		units[i] = k
		i++
	}

	return units
}
