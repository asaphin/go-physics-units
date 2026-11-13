// Package rates builds immutable tables of conversion rates between the units
// of a single measure.
package rates

import (
	"errors"
	"fmt"
	"math"
	"strings"

	"github.com/asaphin/go-physics-units/conversion"
)

// Table holds the conversion rate for every ordered pair of units of one measure,
// along with the name of that measure. A Table is immutable once built and safe
// for concurrent use.
//
// Table is a handle passed by value: its only field points to unexported data,
// so a copy handed to a caller can't be used to change the table it came from.
// (A *Table could be: *p = other overwrites the shared table.) The zero Table
// has no units.
type Table struct {
	t *table
}

type table struct {
	measure string
	rates   map[pair]float64
}

// pair is the lookup key of a rate. A struct key keeps (from, to) pairs distinct;
// concatenating the symbols does not ("m"+"mm" == "mm"+"m").
type pair struct {
	from, to string
}

type entry struct {
	unit, canonical string
	factor          float64
}

// NewTable validates the factors and builds the rate table of the named measure.
// Every factor must be positive and finite, and every unit symbol non-empty.
// Symbols are stored in canonical form (see CanonicalUnit). The table doesn't
// keep a reference to f.
func NewTable(measure string, f conversion.Factors) (Table, error) {
	entries, err := validEntries(f)
	if err != nil {
		return Table{}, err
	}

	rates := make(map[pair]float64, len(entries)*len(entries))

	for i, from := range entries {
		if _, ok := rates[pair{from.canonical, from.canonical}]; ok {
			return Table{}, microSignCollision(entries[:i], from)
		}

		rates[pair{from.canonical, from.canonical}] = 1
	}

	for _, from := range entries {
		for _, to := range entries {
			if from.canonical != to.canonical {
				rates[pair{from.canonical, to.canonical}] = from.factor / to.factor
			}
		}
	}

	return Table{t: &table{measure: measure, rates: rates}}, nil
}

// validEntries checks the factors and returns them with canonical symbols.
func validEntries(f conversion.Factors) ([]entry, error) {
	if len(f) == 0 {
		return nil, errors.New("conversion factors are empty")
	}

	entries := make([]entry, 0, len(f))

	for unit, factor := range f {
		if unit == "" {
			return nil, errors.New("unit symbol must not be empty")
		}

		if !(factor > 0) || math.IsInf(factor, 1) {
			return nil, fmt.Errorf("conversion factor of unit %s must be positive and finite, got %v", unit, factor)
		}

		entries = append(entries, entry{unit: unit, canonical: CanonicalUnit(unit), factor: factor})
	}

	return entries, nil
}

// microSignCollision reports the earlier entry whose symbol differs from e's
// only in the micro sign.
func microSignCollision(earlier []entry, e entry) error {
	prev := e

	for _, p := range earlier {
		if p.canonical == e.canonical {
			prev = p

			break
		}
	}

	return fmt.Errorf("unit symbols %s and %s differ only in the micro sign", prev.unit, e.unit)
}

// MustNewTable is like NewTable but panics if the factors are invalid.
// It is intended for package-level tables built from constant data.
func MustNewTable(measure string, f conversion.Factors) Table {
	t, err := NewTable(measure, f)
	if err != nil {
		panic(fmt.Sprintf("rates: invalid %s conversion factors: %v", measure, err))
	}

	return t
}

// Measure returns the name of the measure the table belongs to, or "" for the
// zero Table.
func (t Table) Measure() string {
	if t.t == nil {
		return ""
	}

	return t.t.measure
}

// Rate returns the multiplier that converts a value in unit from into unit to,
// and whether both units belong to the table. Both symbols must be canonical.
func (t Table) Rate(from, to string) (float64, bool) {
	if t.t == nil {
		return 0, false
	}

	r, ok := t.t.rates[pair{from, to}]

	return r, ok
}

// Has reports whether the canonical unit symbol belongs to the table.
func (t Table) Has(unit string) bool {
	_, ok := t.Rate(unit, unit)

	return ok
}

// CanonicalUnit returns unit with U+00B5 MICRO SIGN replaced by U+03BC GREEK
// SMALL LETTER MU, the same mapping Unicode NFKC applies. Both characters are
// in common use for the micro prefix, so "µm" and "μm" denote the same unit.
func CanonicalUnit(unit string) string {
	return strings.ReplaceAll(unit, "\u00b5", "\u03bc")
}
