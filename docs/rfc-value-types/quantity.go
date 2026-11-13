package prototype

import (
	"cmp"
	"errors"
	"math"
)

// Quantity[M]: the typed API for user-defined measures. The user declares a
// marker type M whose Units method returns an immutable UnitSet built once;
// no global registry is involved (P2), and the library needs no code per
// measure (P3).

// Measure is implemented by marker types. Units must return the same UnitSet
// on every call.
type Measure interface {
	Units() *UnitSet
}

// UnitSet is the immutable set of units of one measure. Its fields are
// unexported, so a *UnitSet is a read-only handle.
type UnitSet struct {
	name    string
	symbols []string
	factors []float64
	index   map[string]uint16
}

// NewUnitSet validates and copies factors.
func NewUnitSet(name string, f map[string]float64) (*UnitSet, error) {
	s := &UnitSet{name: name, index: make(map[string]uint16, len(f))}

	for sym, factor := range f {
		if sym == "" || !(factor > 0) || math.IsInf(factor, 1) {
			return nil, errors.New("invalid factor of unit " + sym)
		}

		s.index[sym] = uint16(len(s.symbols))
		s.symbols = append(s.symbols, sym)
		s.factors = append(s.factors, factor)
	}

	return s, nil
}

// MustNewUnitSet is NewUnitSet for constant data.
func MustNewUnitSet(name string, f map[string]float64) *UnitSet {
	s, err := NewUnitSet(name, f)
	if err != nil {
		panic(err)
	}

	return s
}

// U is a unit of measure M: an index into M's UnitSet.
type U[M Measure] struct{ i uint16 }

// UnitOf looks up a unit of M by symbol, once; values then use the U.
func UnitOf[M Measure](symbol string) (U[M], error) {
	var m M

	i, ok := m.Units().index[symbol]
	if !ok {
		return U[M]{}, errors.New("unknown " + m.Units().name + " unit " + symbol)
	}

	return U[M]{i}, nil
}

// Quantity is a value of the user-defined measure M.
type Quantity[M Measure] struct {
	v float64
	u U[M]
}

// New returns a quantity of M.
func New[M Measure](v float64, u U[M]) Quantity[M] { return Quantity[M]{v, u} }

// Value returns the value in the quantity's unit.
func (q Quantity[M]) Value() float64 { return q.v }

// Unit returns the unit symbol.
func (q Quantity[M]) Unit() string {
	var m M

	return m.Units().symbols[q.u.i]
}

// ConvertTo returns q in unit u.
func (q Quantity[M]) ConvertTo(u U[M]) Quantity[M] {
	var m M

	f := m.Units().factors

	return Quantity[M]{q.v * f[q.u.i] / f[u.i], u}
}

// Add returns q+o in q's unit.
func (q Quantity[M]) Add(o Quantity[M]) Quantity[M] {
	if o.u == q.u {
		return Quantity[M]{q.v + o.v, q.u}
	}

	var m M

	f := m.Units().factors

	return Quantity[M]{q.v + o.v*f[o.u.i]/f[q.u.i], q.u}
}

// Compare compares q and o in the reference unit unless they share a unit.
func (q Quantity[M]) Compare(o Quantity[M]) int {
	if o.u == q.u {
		return cmp.Compare(q.v, o.v)
	}

	var m M

	f := m.Units().factors

	return cmp.Compare(q.v*f[q.u.i], o.v*f[o.u.i])
}
