package units_test

import (
	units "github.com/asaphin/go-physics-units"
	"github.com/asaphin/go-physics-units/conversion"
	"testing"
)

func TestMeasurementCreator(t *testing.T) {
	type Currency interface {
		units.Measurement
	}

	currencyFactors := conversion.Factors{
		"EUR": 1.1,
		"USD": 1,
		"UAH": 0.027,
	}

	mc, err := units.NewMeasurementCreator[Currency]("currency", "USD", currencyFactors)
	if err != nil {
		t.Error(err)
	}

	cur, err := mc.New(1, "USD")
	if err != nil {
		t.Error(err)
	}

	cur2, err := cur.ConvertToMeasurement("UAH")
	if err != nil {
		t.Error(err)
	}

	expectedValue := 37.0

	if !almostEqual(37.0, cur2.Value(), 1e-2) {
		t.Errorf("actual value %v, but expected %v", cur2.Value(), expectedValue)
	}
}
