package units

import "github.com/asaphin/go-physics-units/ratio"

// Ratio is a ratio: a value in a ratio.Unit, such as 12.5 %. Create one with
// NewRatio or StringToRatio.
//
// A Ratio is a small immutable value: every method returns a new value, so copy
// and share it freely. The zero Ratio is 0 1.
type Ratio struct {
	value float64
	unit  ratio.Unit
}
