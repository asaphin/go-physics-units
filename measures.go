package units

import (
	"github.com/asaphin/go-physics-units/acceleration"
	"github.com/asaphin/go-physics-units/angle"
	"github.com/asaphin/go-physics-units/area"
	"github.com/asaphin/go-physics-units/conversion"
	"github.com/asaphin/go-physics-units/datasize"
	"github.com/asaphin/go-physics-units/distance"
	"github.com/asaphin/go-physics-units/energy"
	"github.com/asaphin/go-physics-units/force"
	"github.com/asaphin/go-physics-units/frequency"
	"github.com/asaphin/go-physics-units/internal/rates"
	"github.com/asaphin/go-physics-units/mass"
	"github.com/asaphin/go-physics-units/power"
	"github.com/asaphin/go-physics-units/pressure"
	"github.com/asaphin/go-physics-units/time"
	"github.com/asaphin/go-physics-units/velocity"
	"github.com/asaphin/go-physics-units/volume"
)

// MeasureType identifies a measure (a kind of physical quantity), such as
// distance or time. Measurements can only be combined with measurements of the
// same MeasureType.
type MeasureType string

// Predefined measure types. MeasureCustom is the type of measurements created by
// NewBaseMeasurement and the result of DetectMeasureType for unknown units;
// measurements from a MeasurementCreator have the type the creator was given.
const (
	MeasureDistance     MeasureType = "distance"
	MeasureTime         MeasureType = "time"
	MeasureMass         MeasureType = "mass"
	MeasureVelocity     MeasureType = "velocity"
	MeasureArea         MeasureType = "area"
	MeasureVolume       MeasureType = "volume"
	MeasureAcceleration MeasureType = "acceleration"
	MeasureForce        MeasureType = "force"
	MeasureEnergy       MeasureType = "energy"
	MeasurePower        MeasureType = "power"
	MeasurePressure     MeasureType = "pressure"
	MeasureFrequency    MeasureType = "frequency"
	MeasureAngle        MeasureType = "angle"
	MeasureDataSize     MeasureType = "data size"
	MeasureCustom       MeasureType = "custom"
)

var measureToConversionFactorsMapping = map[MeasureType]conversion.Factors{
	MeasureDistance:     distance.ConversionFactors(),
	MeasureTime:         time.ConversionFactors(),
	MeasureMass:         mass.ConversionFactors(),
	MeasureVelocity:     velocity.ConversionFactors(),
	MeasureArea:         area.ConversionFactors(),
	MeasureVolume:       volume.ConversionFactors(),
	MeasureAcceleration: acceleration.ConversionFactors(),
	MeasureForce:        force.ConversionFactors(),
	MeasureEnergy:       energy.ConversionFactors(),
	MeasurePower:        power.ConversionFactors(),
	MeasurePressure:     pressure.ConversionFactors(),
	MeasureFrequency:    frequency.ConversionFactors(),
	MeasureAngle:        angle.ConversionFactors(),
	MeasureDataSize:     datasize.ConversionFactors(),
}

// Rate tables of the predefined measures, built once at package init and never
// written again. Each is named after its MeasureType, which is how a value
// knows its measure.
var (
	distanceRates     = rates.MustNewTable(string(MeasureDistance), distance.ConversionFactors())
	timeRates         = rates.MustNewTable(string(MeasureTime), time.ConversionFactors())
	massRates         = rates.MustNewTable(string(MeasureMass), mass.ConversionFactors())
	velocityRates     = rates.MustNewTable(string(MeasureVelocity), velocity.ConversionFactors())
	areaRates         = rates.MustNewTable(string(MeasureArea), area.ConversionFactors())
	volumeRates       = rates.MustNewTable(string(MeasureVolume), volume.ConversionFactors())
	accelerationRates = rates.MustNewTable(string(MeasureAcceleration), acceleration.ConversionFactors())
	forceRates        = rates.MustNewTable(string(MeasureForce), force.ConversionFactors())
	energyRates       = rates.MustNewTable(string(MeasureEnergy), energy.ConversionFactors())
	powerRates        = rates.MustNewTable(string(MeasurePower), power.ConversionFactors())
	pressureRates     = rates.MustNewTable(string(MeasurePressure), pressure.ConversionFactors())
	frequencyRates    = rates.MustNewTable(string(MeasureFrequency), frequency.ConversionFactors())
	angleRates        = rates.MustNewTable(string(MeasureAngle), angle.ConversionFactors())
	dataSizeRates     = rates.MustNewTable(string(MeasureDataSize), datasize.ConversionFactors())
)

var measureToConversionRatesMapping = map[MeasureType]rates.Table{
	MeasureDistance:     distanceRates,
	MeasureTime:         timeRates,
	MeasureMass:         massRates,
	MeasureVelocity:     velocityRates,
	MeasureArea:         areaRates,
	MeasureVolume:       volumeRates,
	MeasureAcceleration: accelerationRates,
	MeasureForce:        forceRates,
	MeasureEnergy:       energyRates,
	MeasurePower:        powerRates,
	MeasurePressure:     pressureRates,
	MeasureFrequency:    frequencyRates,
	MeasureAngle:        angleRates,
	MeasureDataSize:     dataSizeRates,
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
