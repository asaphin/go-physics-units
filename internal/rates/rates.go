// Package rates validates the units of a measure and builds their conversion
// rates. The unit packages and the root package use it at package init, and
// NewUnitSet once per custom measure; nothing here runs on a hot path.
package rates

import (
	"errors"
	"fmt"
	"math"
	"strings"
)

// Validate checks the units of one measure, given as parallel slices indexed
// by unit: every symbol must be non-empty and unique in canonical form (see
// CanonicalUnit), and every factor positive and finite.
func Validate(symbols []string, factors []float64) error {
	if len(symbols) == 0 {
		return errors.New("conversion factors are empty")
	}

	if len(symbols) != len(factors) {
		return fmt.Errorf("%d symbols but %d factors", len(symbols), len(factors))
	}

	seen := make(map[string]string, len(symbols))

	for i, symbol := range symbols {
		if symbol == "" {
			return errors.New("unit symbol must not be empty")
		}

		if f := factors[i]; !(f > 0) || math.IsInf(f, 1) {
			return fmt.Errorf("conversion factor of unit %s must be positive and finite, got %v", symbol, f)
		}

		canonical := CanonicalUnit(symbol)
		if prev, ok := seen[canonical]; ok {
			if prev == symbol {
				return fmt.Errorf("unit symbol %s is defined twice", symbol)
			}

			return fmt.Errorf("unit symbols %s and %s differ only in the micro sign", prev, symbol)
		}

		seen[canonical] = symbol
	}

	return nil
}

// MustValidate is like Validate but panics if the units of the named measure
// are invalid. It is intended for the predefined measures, built from
// constant data.
func MustValidate(measure string, symbols []string, factors []float64) {
	if err := Validate(symbols, factors); err != nil {
		panic("rates: invalid " + measure + " units: " + err.Error())
	}
}

// Matrix returns the conversion rates between every pair of units, flattened
// in row-major order: the rate from unit i to unit j, factors[i] / factors[j],
// is at index i·n + j. A value in unit i times that rate is the value in unit
// j. The rate from a unit to itself is exactly 1.
func Matrix(factors []float64) []float64 {
	n := len(factors)
	m := make([]float64, n*n)

	for i, from := range factors {
		for j, to := range factors {
			m[i*n+j] = from / to
		}
	}

	return m
}

// CanonicalUnit returns unit with U+00B5 MICRO SIGN replaced by U+03BC GREEK
// SMALL LETTER MU, the same mapping Unicode NFKC applies. Both characters are
// in common use for the micro prefix, so "µm" and "μm" denote the same unit.
// It doesn't allocate unless unit contains a micro sign.
func CanonicalUnit(unit string) string {
	return strings.ReplaceAll(unit, "\u00b5", "\u03bc")
}
