package units

import (
	"fmt"
	"testing"
)

func TestMeasuresSelfCheck(t *testing.T) {
	uniqueUnits := make(map[string]struct{})
	nonUniqueUnits := make(map[string]struct{})

	for measureType := range measureToConversionFactorsMapping {
		for unit := range measureToConversionFactorsMapping[measureType] {
			if _, ok := uniqueUnits[unit]; !ok {
				uniqueUnits[unit] = struct{}{}

				continue
			}

			nonUniqueUnits[unit] = struct{}{}
		}
	}

	if len(nonUniqueUnits) != 0 {
		t.Errorf("non unique units found")
		for unit := range nonUniqueUnits {
			fmt.Println(unit)
		}
	}
}
