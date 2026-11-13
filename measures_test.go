package units //nolint:testpackage // inspects the unexported measure registries

import (
	"sort"
	"strings"
	"testing"

	"github.com/asaphin/go-physics-units/internal/rates"
)

func TestMeasuresSelfCheck(t *testing.T) {
	seen := make(map[string]MeasureType)

	var duplicates []string

	for measureType, factors := range measureToConversionFactorsMapping {
		for unit := range factors {
			canonical := rates.CanonicalUnit(unit)
			if other, ok := seen[canonical]; ok {
				duplicates = append(duplicates, canonical+" ("+string(other)+", "+string(measureType)+")")

				continue
			}

			seen[canonical] = measureType
		}
	}

	if len(duplicates) != 0 {
		sort.Strings(duplicates)
		t.Errorf("non unique units found: %s", strings.Join(duplicates, "; "))
	}
}

func TestMeasureRegistriesAreConsistent(t *testing.T) {
	if len(measureToConversionFactorsMapping) != len(measureToConversionRatesMapping) {
		t.Fatal("the factors and rates registries list different measures")
	}

	for mt, table := range measureToConversionRatesMapping {
		if table.Measure() != string(mt) {
			t.Errorf("rates table registered as %s is named %s", mt, table.Measure())
		}

		factors, ok := measureToConversionFactorsMapping[mt]
		if !ok {
			t.Errorf("%s has rates but no factors", mt)

			continue
		}

		for unit := range factors {
			if !table.Has(rates.CanonicalUnit(unit)) {
				t.Errorf("%s: unit %s is missing from the rates table", mt, unit)
			}
		}

		if !isReservedMeasureType(mt) {
			t.Errorf("%s isn't reserved", mt)
		}
	}
}
