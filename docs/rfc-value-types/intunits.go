package prototype

import (
	"cmp"
	"strconv"
)

// Option I: typed integer units that index immutable arrays. Conversions are
// two array loads; the zero Unit is the base unit, so the zero Distance is 0 m.
// The set of units is closed: extending it means a custom measure (Quantity).

// Unit is a distance unit. Only the constants below are valid.
type Unit uint8

// Distance units of option I. Meter is the zero value.
const (
	Meter Unit = iota
	Kilometer
	Foot
	Inch
	Mile
	NauticalMile
	unitCount
)

// The arrays are unexported and never written, so they are immutable (P2):
// an importer can't reach them, unlike an exported variable.
var (
	factors = [unitCount]float64{1, 1000, 0.3048, 0.0254, 1609.344, 1852}
	symbols = [unitCount]string{"m", "km", "ft", "in", "mi", "nmi"}
)

// String returns the unit symbol.
func (u Unit) String() string {
	if u >= unitCount {
		return "Unit(" + strconv.Itoa(int(u)) + ")"
	}

	return symbols[u]
}

// ParseUnit returns the unit with the given symbol: the one place a symbol
// can be wrong, so the only conversion step that returns an error.
func ParseUnit(symbol string) (Unit, bool) {
	for u, s := range symbols {
		if s == symbol {
			return Unit(u), true
		}
	}

	return 0, false
}

// Distance is a distance value of option I.
type Distance struct {
	v float64
	u Unit
}

// NewDistance can't fail: u is a valid constant unless the caller converted an
// arbitrary integer, which panics on the first conversion.
func NewDistance(v float64, u Unit) Distance { return Distance{v, u} }

// Value returns the value in Unit.
func (d Distance) Value() float64 { return d.v }

// Unit returns the unit.
func (d Distance) Unit() Unit { return d.u }

// ConvertTo returns d in unit u.
func (d Distance) ConvertTo(u Unit) Distance {
	return Distance{d.v * factors[d.u] / factors[u], u}
}

// Add returns d+o in d's unit.
func (d Distance) Add(o Distance) Distance {
	if o.u == d.u {
		return Distance{d.v + o.v, d.u}
	}

	return Distance{d.v + o.v*factors[o.u]/factors[d.u], d.u}
}

// Compare compares d and o in the base unit unless they share a unit.
func (d Distance) Compare(o Distance) int {
	if o.u == d.u {
		return cmp.Compare(d.v, o.v)
	}

	return cmp.Compare(d.v*factors[d.u], o.v*factors[o.u])
}

// String formats d as "<value> <unit>", like the current String.
func (d Distance) String() string {
	var buf [64]byte

	s := strconv.AppendFloat(buf[:0], d.v, 'g', -1, 64)
	s = append(s, ' ')
	s = append(s, d.u.String()...)

	return string(s)
}
