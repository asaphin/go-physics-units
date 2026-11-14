package units

//go:generate go run ./internal/gen

import (
	"github.com/asaphin/go-physics-units/internal/rates"
	"github.com/asaphin/go-physics-units/temperature"
)

// MeasureType identifies a measure (a kind of physical quantity), such as
// distance or time. Measurements can only be combined with measurements of the
// same measure.
type MeasureType string

// Predefined measure types. Values of a user-defined measure have the type
// given to NewUnitSet. MeasureCustom is the result of DetectMeasureType for
// symbols that no predefined measure has.
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

	// MeasureTemperature names the temperature measure in errors and in
	// DetectMeasureType. Temperatures aren't Measurements: see Temperature.
	MeasureTemperature MeasureType = "temperature"
)

// parsedUnit is what a symbol denotes: a unit of a unit set, and how to make
// a Measurement of it. If several symbols of a case-insensitive Parser fold
// to the same key, candidates lists them instead.
type parsedUnit struct {
	set        *unitSet
	unit       uint16
	box        func(float64, uint16) Measurement
	candidates string
}

// predefinedUnits indexes every predefined unit symbol, in canonical form: one
// lookup takes NewMeasurement and ParseString from a symbol to a value. It is
// built at init and never written again, and it is the starting point of
// every Parser.
var predefinedUnits = buildUnitIndex()

func buildUnitIndex() map[string]parsedUnit {
	index := make(map[string]parsedUnit)

	for _, s := range predefinedSets {
		for symbol, u := range s.index {
			index[symbol] = parsedUnit{set: s, unit: u, box: s.box}
		}
	}

	return index
}

// DetectMeasureType returns the predefined measure that the unit symbol belongs
// to, MeasureTemperature for a temperature symbol, or MeasureCustom if it isn't
// a predefined symbol.
func DetectMeasureType(symbol string) MeasureType {
	if pu, ok := predefinedUnits[rates.CanonicalUnit(symbol)]; ok {
		return pu.set.name
	}

	if _, ok := temperature.ParseUnit(symbol); ok {
		return MeasureTemperature
	}

	return MeasureCustom
}

// unknownSymbolError returns the error for a symbol that no Measurement has.
// Temperature symbols get an error of their own, which says where to parse
// them.
func unknownSymbolError(symbol string) error {
	if _, ok := temperature.ParseUnit(symbol); ok {
		return newError(ErrIncompatibleMeasures, "unit "+symbol+
			" is a temperature unit, and temperatures aren't Measurements: use StringToTemperature", nil)
	}

	return unknownUnitError("", symbol)
}

// isReservedMeasureType reports whether mt is the name of a predefined
// measure, MeasureTemperature or MeasureCustom.
func isReservedMeasureType(mt MeasureType) bool {
	if mt == MeasureCustom || mt == MeasureTemperature {
		return true
	}

	for _, s := range predefinedSets {
		if s.name == mt {
			return true
		}
	}

	return false
}
