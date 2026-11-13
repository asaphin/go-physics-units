package units

import (
	"errors"
)

var errZeroDivision = errors.New("division by zero")

// Arithmetics is the arithmetic implemented by every predefined measure type T.
// Add and Sub accept an operand in any unit of the same measure and return the
// result in the receiver's unit.
type Arithmetics[T any] interface {
	// Add returns the sum of the receiver and measurement.
	Add(measurement T) T
	// Sub returns the difference of the receiver and measurement.
	Sub(measurement T) T
	// Mul returns the receiver multiplied by multiplier.
	Mul(multiplier float64) T
	// Div returns the receiver divided by divisor. It fails if divisor is zero.
	Div(divisor float64) (T, error)
}

// UnitConverter is the unit conversion implemented by every predefined measure
// type T.
type UnitConverter[T any] interface {
	// ConvertTo returns the receiver expressed in unit. It fails if unit doesn't
	// belong to the measure.
	ConvertTo(unit string) (T, error)
	// MustConvertTo is like ConvertTo but panics if unit doesn't belong to the
	// measure.
	MustConvertTo(unit string) T
	// ConvertToBaseUnits returns the receiver expressed in the measure's SI base
	// unit.
	ConvertToBaseUnits() T

	valueInBaseUnits() float64
	valueInUnits(unit string) float64
}
