package units

import "fmt"

type MeasureType string

const (
	measureDistance MeasureType = "distance"
	measureTime     MeasureType = "time"
	measureMass     MeasureType = "mass"
	measureVelocity MeasureType = "velocity"
)

var measureToConversionFactorsMapping = map[MeasureType]ConversionFactors{
	measureDistance: distanceConversionFactors,
	measureTime:     timeConversionFactors,
	measureMass:     massConversionFactors,
	measureVelocity: velocityConversionFactors,
}

func DetectMeasureType(unit string) (MeasureType, error) {
	for measureType, factors := range measureToConversionFactorsMapping {
		for k := range factors {
			if k == unit {
				return measureType, nil
			}
		}
	}

	return "", fmt.Errorf("unknown unit %s", unit)
}
