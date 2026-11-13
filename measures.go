package units

import (
	"github.com/asaphin/go-physics-units/area"
	"github.com/asaphin/go-physics-units/conversion"
	"github.com/asaphin/go-physics-units/distance"
	"github.com/asaphin/go-physics-units/internal/rates"
	"github.com/asaphin/go-physics-units/mass"
	"github.com/asaphin/go-physics-units/time"
	"github.com/asaphin/go-physics-units/velocity"
)

// MeasureType identifies a measure (a kind of physical quantity), such as
// distance or time. Measurements can only be combined with measurements of the
// same MeasureType.
type MeasureType string

// Predefined measure types. MeasureCustom is the type of measurements created by
// NewBaseMeasurement and the result of DetectMeasureType for unknown units;
// measurements from a MeasurementCreator have the type the creator was given.
const (
	MeasureDistance MeasureType = "distance"
	MeasureTime     MeasureType = "time"
	MeasureMass     MeasureType = "mass"
	MeasureVelocity MeasureType = "velocity"
	MeasureArea     MeasureType = "area"
	MeasureCustom   MeasureType = "custom"
)

var measureToConversionFactorsMapping = map[MeasureType]conversion.Factors{
	MeasureDistance: distance.ConversionFactors(),
	MeasureTime:     time.ConversionFactors(),
	MeasureMass:     mass.ConversionFactors(),
	MeasureVelocity: velocity.ConversionFactors(),
	MeasureArea:     area.ConversionFactors(),
}

var measureToConversionRatesMapping = map[MeasureType]rates.Table{
	MeasureDistance: distance.ConversionRates(),
	MeasureTime:     time.ConversionRates(),
	MeasureMass:     mass.ConversionRates(),
	MeasureVelocity: velocity.ConversionRates(),
	MeasureArea:     area.ConversionRates(),
}

// unitToMeasureType indexes every predefined unit symbol by its measure.
var unitToMeasureType = buildUnitIndex()

func buildUnitIndex() map[string]MeasureType {
	index := make(map[string]MeasureType)

	for mt, factors := range measureToConversionFactorsMapping {
		for unit := range factors {
			index[rates.CanonicalUnit(unit)] = mt
		}
	}

	return index
}

// DetectMeasureType returns the predefined measure that the unit symbol belongs
// to, or MeasureCustom if it isn't a predefined symbol.
func DetectMeasureType(unit string) MeasureType {
	if mt, ok := unitToMeasureType[rates.CanonicalUnit(unit)]; ok {
		return mt
	}

	return MeasureCustom
}

// isReservedMeasureType reports whether mt is one of the predefined measure types.
func isReservedMeasureType(mt MeasureType) bool {
	_, ok := measureToConversionRatesMapping[mt]

	return ok || mt == MeasureCustom
}
