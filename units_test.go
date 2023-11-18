package units

import (
	"testing"
)

func TestBaseMeasurement(t *testing.T) {
	ln, err := NewBaseMeasurement(1, "m", distanceConversionFactors)
	if err != nil {
		t.Error(err)
	}

	ln2, err := ln.ConvertTo("km")
	if err != nil {
		t.Error(err)
	}

	ln2Value := ln2.Value()

	if ln2Value != 1000 {
		t.Errorf("value after konversion should be 1000, but was %v", ln2Value)
	}
}
