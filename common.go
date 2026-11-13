package units

import (
	"cmp"
	"math"
)

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

// Comparer is the comparison implemented by every predefined measure type T.
//
// Values in the same unit are compared exactly. Values in different units are
// compared in the measure's SI base unit, so the comparison doesn't depend on
// the order of the operands: a.Compare(b) == -b.Compare(a), and a.Equal(b) ==
// b.Equal(a).
//
// Conversions round, so values that are equal in exact arithmetic can differ in
// the last bit once converted: 12 in is 0.30479999999999996 m, while 1 ft is
// 0.3048 m. Use EqualWithinRel to compare with a tolerance.
type Comparer[T any] interface {
	// Compare returns -1, 0 or +1 as the receiver is less than, equal to or
	// greater than other. It follows cmp.Compare: a NaN value is less than any
	// other value and equal to NaN, and -0 equals +0.
	Compare(other T) int
	// Equal reports whether Compare(other) is 0.
	Equal(other T) bool
	// EqualWithinRel reports whether the receiver and other are equal within
	// the relative tolerance tol: |a-b| <= tol·max(|a|, |b|), with a and b
	// compared as in Compare. EqualWithinRel(other, 0) is Equal(other), and an
	// infinite value only equals the same infinity.
	EqualWithinRel(other T, tol float64) bool
}

// equalWithinRel implements Comparer.EqualWithinRel for two values in the same
// unit.
func equalWithinRel(a, b, tol float64) bool {
	if cmp.Compare(a, b) == 0 {
		return true
	}

	if math.IsInf(a, 0) || math.IsInf(b, 0) {
		return false // otherwise Inf ≤ tol·Inf would hold
	}

	return math.Abs(a-b) <= tol*max(math.Abs(a), math.Abs(b))
}
