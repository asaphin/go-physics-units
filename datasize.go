package units

import (
	"github.com/asaphin/go-physics-units/datarate"
	"github.com/asaphin/go-physics-units/datasize"
	"github.com/asaphin/go-physics-units/time"
)

// DataSize is an amount of data: a value in a datasize.Unit, such as 1.5 GB. Create one with
// NewDataSize or StringToDataSize.
//
// A DataSize is a small immutable value: every method returns a new value, so
// copy and share it freely. The zero DataSize is 0 B.
type DataSize struct {
	value float64
	unit  datasize.Unit
}

// DivideByTime returns the data rate s/t (B/s) in [datarate.BaseUnit]. It fails with
// [ErrZeroDivision] if t is zero.
func (s DataSize) DivideByTime(t Time) (DataRate, error) {
	if t.value == 0 {
		return DataRate{}, ErrZeroDivision
	}

	return DataRate{value: s.baseValue() / t.baseValue(), unit: datarate.BaseUnit}, nil
}

// DivideByDataRate returns the time s/r (s), for example to transfer s at rate
// r, in [time.BaseUnit]. It fails with [ErrZeroDivision] if r is zero.
func (s DataSize) DivideByDataRate(r DataRate) (Time, error) {
	if r.value == 0 {
		return Time{}, ErrZeroDivision
	}

	return Time{value: s.baseValue() / r.baseValue(), unit: time.BaseUnit}, nil
}
