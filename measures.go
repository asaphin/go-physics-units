package units

import (
	"github.com/asaphin/go-physics-units/acceleration"
	"github.com/asaphin/go-physics-units/angle"
	"github.com/asaphin/go-physics-units/area"
	"github.com/asaphin/go-physics-units/conversion"
	"github.com/asaphin/go-physics-units/datarate"
	"github.com/asaphin/go-physics-units/datasize"
	"github.com/asaphin/go-physics-units/density"
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
	MeasureDensity      MeasureType = "density"
	MeasureDataRate     MeasureType = "data rate"
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
	MeasureDensity:      density.ConversionFactors(),
	MeasureDataRate:     datarate.ConversionFactors(),
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
	densityRates      = rates.MustNewTable(string(MeasureDensity), density.ConversionFactors())
	dataRateRates     = rates.MustNewTable(string(MeasureDataRate), datarate.ConversionFactors())
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
	MeasureDensity:      densityRates,
	MeasureDataRate:     dataRateRates,
}

// predefinedUnits indexes every predefined unit symbol, in canonical form, by
// its unit and table: one lookup takes NewMeasurement and ParseString from a
// symbol to a value. It is built at init and never written again, and it is
// the starting point of every Parser.
var predefinedUnits = buildUnitIndex()

func buildUnitIndex() map[string]parsedUnit {
	index := make(map[string]parsedUnit)

	for _, table := range measureToConversionRatesMapping {
		for unit := range table.Units() {
			index[unit] = parsedUnit{unit: unit, table: table}
		}
	}

	return index
}

// DetectMeasureType returns the predefined measure that the unit symbol belongs
// to, or MeasureCustom if it isn't a predefined symbol.
func DetectMeasureType(unit string) MeasureType {
	if pu, ok := predefinedUnits[rates.CanonicalUnit(unit)]; ok {
		return MeasureType(pu.table.Measure())
	}

	return MeasureCustom
}

// isReservedMeasureType reports whether mt is one of the predefined measure types.
func isReservedMeasureType(mt MeasureType) bool {
	_, ok := measureToConversionRatesMapping[mt]

	return ok || mt == MeasureCustom
}
