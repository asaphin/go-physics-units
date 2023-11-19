package units

import (
	"github.com/asaphin/go-physics-units/conversion-factors"
	"github.com/asaphin/go-physics-units/distance"
	"github.com/asaphin/go-physics-units/mass"
	"github.com/asaphin/go-physics-units/time"
	"github.com/asaphin/go-physics-units/velocity"
)

type MeasureType string

const (
	MeasureDistance MeasureType = "distance"
	MeasureTime     MeasureType = "time"
	MeasureMass     MeasureType = "mass"
	MeasureVelocity MeasureType = "velocity"
	MeasureCustom   MeasureType = "custom"
)

var measureToConversionFactorsMapping = map[MeasureType]conversion.Factors{
	MeasureDistance: distance.ConversionFactors(),
	MeasureTime:     time.ConversionFactors(),
	MeasureMass:     mass.ConversionFactors(),
	MeasureVelocity: velocity.ConversionFactors(),
}

func DetectMeasureType(unit string) MeasureType {
	for measureType, factors := range measureToConversionFactorsMapping {
		for k := range factors {
			if k == unit {
				return measureType
			}
		}
	}

	return MeasureCustom
}
