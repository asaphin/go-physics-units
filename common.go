package units

import (
	"errors"
)

var errZeroDivision = errors.New("division by zero")

type Arithmetics[T any] interface {
	Add(measurement T) T
	Sub(measurement T) T
	Mul(multiplier float64) T
	Div(divisor float64) (T, error)
}

type UnitConverter[T any] interface {
	ConvertTo(unit string) (T, error)
	MustConvertTo(unit string) T
	ConvertToBaseUnits() T

	valueInBaseUnits() float64
	valueInUnits(unit string) float64
}
