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

// Fraction returns the ratio as a plain number, 0.25 for 25 %. Multiply a value
// by it to take that share: d.Mul(r.Fraction()).
func (r Ratio) Fraction() float64 {
	return r.baseValue()
}
