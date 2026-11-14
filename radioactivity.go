package units

import "github.com/asaphin/go-physics-units/radioactivity"

// Radioactivity is a radioactivity: a value in a radioactivity.Unit, such as
// 3.7 GBq. Create one with NewRadioactivity or StringToRadioactivity.
//
// A Radioactivity is a small immutable value: every method returns a new value,
// so copy and share it freely. The zero Radioactivity is 0 Bq.
type Radioactivity struct {
	value float64
	unit  radioactivity.Unit
}
