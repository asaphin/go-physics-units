package units //nolint:testpackage // inspects the unexported unit sets

import (
	"slices"
	"strings"
	"testing"

	"github.com/asaphin/go-physics-units/internal/rates"
	"github.com/asaphin/go-physics-units/temperature"
)

func TestMeasuresSelfCheck(t *testing.T) {
	seen := make(map[string]MeasureType)

	var duplicates []string

	for _, s := range predefinedSets {
		for _, symbol := range s.symbols {
			canonical := rates.CanonicalUnit(symbol)
			if canonical != symbol {
				t.Errorf("%s unit %q isn't in canonical form %q; ParseUnit looks units up by canonical symbol",
					s.name, symbol, canonical)
			}

			if other, ok := seen[canonical]; ok {
				duplicates = append(duplicates, canonical+" ("+string(other)+", "+string(s.name)+")")

				continue
			}

			seen[canonical] = s.name
		}
	}

	if len(duplicates) != 0 {
		slices.Sort(duplicates)
		t.Errorf("non unique units found: %s", strings.Join(duplicates, "; "))
	}

	for u := range temperature.Units() {
		if other, ok := seen[u.String()]; ok {
			t.Errorf("temperature symbol %s is also a %s unit", u, other)
		}
	}

	if len(predefinedUnits) != len(seen) {
		t.Errorf("predefinedUnits has %d symbols, the sets %d", len(predefinedUnits), len(seen))
	}
}

func TestPredefinedSetsAreConsistent(t *testing.T) {
	names := make(map[MeasureType]bool)

	for _, s := range predefinedSets {
		if names[s.name] {
			t.Errorf("two sets are named %s", s.name)
		}

		names[s.name] = true

		if !isReservedMeasureType(s.name) {
			t.Errorf("%s isn't reserved", s.name)
		}

		n := len(s.symbols)
		if len(s.factors) != n || len(s.rates) != n*n || len(s.index) != n || s.box == nil {
			t.Errorf("%s: %d symbols, %d factors, %d rates, %d index entries", s.name,
				n, len(s.factors), len(s.rates), len(s.index))
		}

		for i, symbol := range s.symbols {
			m := s.box(1, uint16(i))
			if m.Type() != s.name || m.Symbol() != symbol || m.generic().set != s {
				t.Errorf("%s: unit %d boxes as %v of type %s", s.name, i, m, m.Type())
			}

			if pu := predefinedUnits[symbol]; pu.set != s || int(pu.unit) != i {
				t.Errorf("predefinedUnits[%s] = unit %d of %v", symbol, pu.unit, pu.set.name)
			}
		}
	}

	if len(names) != len(predefinedSets) {
		t.Errorf("%d names for %d sets", len(names), len(predefinedSets))
	}
}
