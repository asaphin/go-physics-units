package units

import (
	"github.com/asaphin/go-physics-units/datarate"
	"github.com/asaphin/go-physics-units/datasize"
)

// DataRate is a data rate: a value in a datarate.Unit, such as 100 Mbit/s. Create one with
// NewDataRate or StringToDataRate.
//
// A DataRate is a small immutable value: every method returns a new value, so
// copy and share it freely. The zero DataRate is 0 B/s.
type DataRate struct {
	value float64
	unit  datarate.Unit
}

// MultiplyByTime returns the data size r·t (B) in [datasize.BaseUnit].
func (r DataRate) MultiplyByTime(t Time) DataSize {
	return DataSize{value: r.baseValue() * t.baseValue(), unit: datasize.BaseUnit}
}
