package units_test

import (
	units "github.com/asaphin/go-physics-units"
	"testing"
)

func TestDistanceConversionFactors_Safety(t *testing.T) {
	dcf := units.DistanceConversionFactors()

	dcf["km"] = 2000

	otherDcf := units.DistanceConversionFactors()

	if otherDcf["km"] == 2000 {
		t.Error("DistanceConversionFactors function isn't safe")
	}
}
