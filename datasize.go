package units

import (
	"cmp"

	"github.com/asaphin/go-physics-units/datarate"
	"github.com/asaphin/go-physics-units/datasize"
	"github.com/asaphin/go-physics-units/time"
)

// DataSize is a data size measurement. Its unit symbols are defined in the datasize
// package.
type DataSize interface {
	Measurement
	UnitConverter[DataSize]
	Arithmetics[DataSize]
	Comparer[DataSize]

	// DivideByTime returns the data rate s/t in datarate.BaseUnit (B/s).
	// It fails if t is zero.
	DivideByTime(t Time) (DataRate, error)
	// DivideByDataRate returns the time s/r, for example to transfer s at
	// rate r, in time.BaseUnit (s). It fails if r is zero.
	DivideByDataRate(r DataRate) (Time, error)
}

// dataSizeImplementation is the concrete implementation of DataSize.
type dataSizeImplementation struct {
	baseMeasurement
}

// ConvertTo returns the data size expressed in unit.
func (s *dataSizeImplementation) ConvertTo(unit string) (DataSize, error) {
	msm, err := s.convertTo(unit)
	if err != nil {
		return nil, err
	}

	return &dataSizeImplementation{msm}, nil
}

// MustConvertTo is like ConvertTo but panics if unit isn't a data size unit.
func (s *dataSizeImplementation) MustConvertTo(unit string) DataSize {
	msm, err := s.convertTo(unit)
	if err != nil {
		panic(err)
	}

	return &dataSizeImplementation{msm}
}

// ConvertToBaseUnits returns the data size expressed in bytes.
func (s *dataSizeImplementation) ConvertToBaseUnits() DataSize {
	return &dataSizeImplementation{s.unsafeConvertTo(datasize.BaseUnit)}
}

func (s *dataSizeImplementation) valueInBaseUnits() float64 {
	return s.unsafeGetValueIn(datasize.BaseUnit)
}

func (s *dataSizeImplementation) valueInUnits(unit string) float64 {
	return s.unsafeGetValueIn(unit)
}

// Add returns the sum of the data sizes, expressed in the receiver's unit.
func (s *dataSizeImplementation) Add(other DataSize) DataSize {
	return &dataSizeImplementation{s.withValue(s.value + other.valueInUnits(s.unit))}
}

// Sub returns the difference of the data sizes, expressed in the receiver's unit.
func (s *dataSizeImplementation) Sub(other DataSize) DataSize {
	return &dataSizeImplementation{s.withValue(s.value - other.valueInUnits(s.unit))}
}

// Mul returns the data size multiplied by multiplier.
func (s *dataSizeImplementation) Mul(multiplier float64) DataSize {
	return &dataSizeImplementation{s.withValue(s.value * multiplier)}
}

// Div returns the data size divided by divisor. It fails if divisor is zero.
func (s *dataSizeImplementation) Div(divisor float64) (DataSize, error) {
	if divisor == 0 {
		return nil, ErrZeroDivision
	}

	return &dataSizeImplementation{s.withValue(s.value / divisor)}, nil
}

// DivideByTime returns the data rate s/t in B/s, or an error if t is zero.
func (s *dataSizeImplementation) DivideByTime(t Time) (DataRate, error) {
	if t.Value() == 0 {
		return nil, ErrZeroDivision
	}

	return newDataRate(s.valueInBaseUnits()/t.valueInBaseUnits(), datarate.BaseUnit), nil
}

// DivideByDataRate returns the time s/r in s, or an error if r is zero.
func (s *dataSizeImplementation) DivideByDataRate(r DataRate) (Time, error) {
	if r.Value() == 0 {
		return nil, ErrZeroDivision
	}

	return newTime(s.valueInBaseUnits()/r.valueInBaseUnits(), time.BaseUnit), nil
}

// Compare compares the data sizes; see Comparer.
func (s *dataSizeImplementation) Compare(other DataSize) int {
	return cmp.Compare(s.commonValues(other.base()))
}

// Equal reports whether the data sizes are equal; see Comparer.
func (s *dataSizeImplementation) Equal(other DataSize) bool {
	return s.Compare(other) == 0
}

// EqualWithinRel reports whether the data sizes are equal within the relative
// tolerance tol; see Comparer.
func (s *dataSizeImplementation) EqualWithinRel(other DataSize, tol float64) bool {
	x, y := s.commonValues(other.base())

	return equalWithinRel(x, y, tol)
}

// NewDataSize creates a data size of value in unit, one of the symbols defined in
// the datasize package.
func NewDataSize(value float64, unit string) (DataSize, error) {
	msm, err := newBaseMeasurement(value, unit, dataSizeRates)
	if err != nil {
		return nil, err
	}

	return &dataSizeImplementation{msm}, nil
}

// newDataSize is an unchecked constructor for internal use with valid, canonical
// units.
func newDataSize(value float64, unit string) DataSize {
	return &dataSizeImplementation{baseMeasurement{value: value, unit: unit, conversionRates: dataSizeRates}}
}

const errDataSizeConversion incompatibleError = "not a data size measure"

// ToDataSize converts m, for example a result of ParseString, into a DataSize.
// It fails if m isn't a data size.
func ToDataSize(m Measurement) (DataSize, error) {
	if m == nil || m.Type() != MeasureDataSize {
		return nil, errDataSizeConversion
	}

	return &dataSizeImplementation{*m.base()}, nil
}

// StringToDataSize parses a data size from a string such as "1.5 GB".
// It fails if the unit isn't a data size unit.
func StringToDataSize(s string) (DataSize, error) {
	// Parse and construct directly: one allocation, where ParseString
	// followed by ToDataSize would allocate twice.
	value, unit, err := parseValueAndUnit(s)
	if err != nil {
		return nil, err
	}

	return NewDataSize(value, unit)
}
