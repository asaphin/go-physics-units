package units //nolint: testpackage

import (
	"strings"
	"testing"
)

func TestMeasuresSelfCheck(t *testing.T) {
	uniqueUnitsMap := make(map[string]struct{})
	nonUniqueUnitsMap := make(map[string]struct{})

	for measureType := range measureToConversionFactorsMapping {
		for unit := range measureToConversionFactorsMapping[measureType] {
			if _, ok := uniqueUnitsMap[unit]; !ok {
				uniqueUnitsMap[unit] = struct{}{}

				continue
			}

			nonUniqueUnitsMap[unit] = struct{}{}
		}
	}

	nonUniqueUnits := make([]string, 0, len(nonUniqueUnitsMap))
	for k := range nonUniqueUnitsMap {
		nonUniqueUnits = append(nonUniqueUnits, k)
	}

	if len(nonUniqueUnits) != 0 {
		t.Errorf("non unique units found: %s", strings.Join(nonUniqueUnits, "; "))
	}
}
