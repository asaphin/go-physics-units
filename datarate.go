package units

import (
	"cmp"

	"github.com/asaphin/go-physics-units/datarate"
	"github.com/asaphin/go-physics-units/datasize"
)

// DataRate is a data rate measurement. Its unit symbols are defined in the datarate
// package.
type DataRate interface {
	Measurement
	UnitConverter[DataRate]
	Arithmetics[DataRate]
	Comparer[DataRate]

	// MultiplyByTime returns the data size r·t in datasize.BaseUnit (B).
	MultiplyByTime(t Time) DataSize
}

// dataRateImplementation is the concrete implementation of DataRate.
type dataRateImplementation struct {
	baseMeasurement
}

// ConvertTo returns the data rate expressed in unit.
func (r *dataRateImplementation) ConvertTo(unit string) (DataRate, error) {
	msm, err := r.convertTo(unit)
	if err != nil {
		return nil, err
	}

	return &dataRateImplementation{msm}, nil
}

// MustConvertTo is like ConvertTo but panics if unit isn't a data rate unit.
func (r *dataRateImplementation) MustConvertTo(unit string) DataRate {
	msm, err := r.convertTo(unit)
	if err != nil {
		panic(err)
	}

	return &dataRateImplementation{msm}
}

// ConvertToBaseUnits returns the data rate expressed in bytes per second.
func (r *dataRateImplementation) ConvertToBaseUnits() DataRate {
	return &dataRateImplementation{r.unsafeConvertTo(datarate.BaseUnit)}
}

func (r *dataRateImplementation) valueInBaseUnits() float64 {
	return r.unsafeGetValueIn(datarate.BaseUnit)
}

func (r *dataRateImplementation) valueInUnits(unit string) float64 {
	return r.unsafeGetValueIn(unit)
}

// Add returns the sum of the data rates, expressed in the receiver's unit.
func (r *dataRateImplementation) Add(other DataRate) DataRate {
	return &dataRateImplementation{r.withValue(r.value + other.valueInUnits(r.unit))}
}

// Sub returns the difference of the data rates, expressed in the receiver's unit.
func (r *dataRateImplementation) Sub(other DataRate) DataRate {
	return &dataRateImplementation{r.withValue(r.value - other.valueInUnits(r.unit))}
}

// Mul returns the data rate multiplied by multiplier.
func (r *dataRateImplementation) Mul(multiplier float64) DataRate {
	return &dataRateImplementation{r.withValue(r.value * multiplier)}
}

// Div returns the data rate divided by divisor. It fails if divisor is zero.
func (r *dataRateImplementation) Div(divisor float64) (DataRate, error) {
	if divisor == 0 {
		return nil, ErrZeroDivision
	}

	return &dataRateImplementation{r.withValue(r.value / divisor)}, nil
}

// MultiplyByTime returns the data size r·t in B.
func (r *dataRateImplementation) MultiplyByTime(t Time) DataSize {
	return newDataSize(r.valueInBaseUnits()*t.valueInBaseUnits(), datasize.BaseUnit)
}

// Compare compares the data rates; see Comparer.
func (r *dataRateImplementation) Compare(other DataRate) int {
	return cmp.Compare(r.commonValues(other.base()))
}

// Equal reports whether the data rates are equal; see Comparer.
func (r *dataRateImplementation) Equal(other DataRate) bool {
	return r.Compare(other) == 0
}

// EqualWithinRel reports whether the data rates are equal within the relative
// tolerance tol; see Comparer.
func (r *dataRateImplementation) EqualWithinRel(other DataRate, tol float64) bool {
	x, y := r.commonValues(other.base())

	return equalWithinRel(x, y, tol)
}

// NewDataRate creates a data rate of value in unit, one of the symbols defined in
// the datarate package.
func NewDataRate(value float64, unit string) (DataRate, error) {
	msm, err := newBaseMeasurement(value, unit, dataRateRates)
	if err != nil {
		return nil, err
	}

	return &dataRateImplementation{msm}, nil
}

// newDataRate is an unchecked constructor for internal use with valid, canonical
// units.
func newDataRate(value float64, unit string) DataRate {
	return &dataRateImplementation{baseMeasurement{value: value, unit: unit, conversionRates: dataRateRates}}
}

const errDataRateConversion incompatibleError = "not a data rate measure"

// ToDataRate converts m, for example a result of ParseString, into a DataRate.
// It fails if m isn't a data rate.
func ToDataRate(m Measurement) (DataRate, error) {
	if m == nil || m.Type() != MeasureDataRate {
		return nil, errDataRateConversion
	}

	return &dataRateImplementation{*m.base()}, nil
}

// StringToDataRate parses a data rate from a string such as "100 Mbit/s".
// It fails if the unit isn't a data rate unit.
func StringToDataRate(s string) (DataRate, error) {
	// Parse and construct directly: one allocation, where ParseString
	// followed by ToDataRate would allocate twice.
	value, unit, err := parseValueAndUnit(s)
	if err != nil {
		return nil, err
	}

	return NewDataRate(value, unit)
}
