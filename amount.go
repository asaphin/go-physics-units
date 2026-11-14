package units

import "github.com/asaphin/go-physics-units/amount"

// Amount is an amount of substance: a value in a amount.Unit, such as 2.5 mmol.
// Create one with NewAmount or StringToAmount.
//
// An Amount is a small immutable value: every method returns a new value, so
// copy and share it freely. The zero Amount is 0 mol.
type Amount struct {
	value float64
	unit  amount.Unit
}
