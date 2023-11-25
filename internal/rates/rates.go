package rates

import "github.com/asaphin/go-physics-units/conversion"

// Rates stores conversion rates by complex keys from_unit + to_unit.
type Rates = map[string]float64

func FactorsToRates(f conversion.Factors) Rates {
	rates := make(Rates)

	u := getUnits(f)

	for i := range u {
		for j := range u {
			if u[i] == u[j] {
				rates[u[i]+u[j]] = 1

				continue
			}

			rates[u[i]+u[j]] = f[u[i]] / f[u[j]]
		}
	}

	return rates
}

func getUnits(f conversion.Factors) []string {
	units := make([]string, len(f))

	i := 0
	for k := range f {
		units[i] = k
		i++
	}

	return units
}
