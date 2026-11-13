package prototype

import (
	"cmp"
	"errors"
)

// Option S: typed string units. Untyped literals ("km") still compile, units
// of another measure don't, and every conversion looks the symbols up in a map.

// StrUnit is a distance unit symbol.
type StrUnit string

// Distance units of option S.
const (
	SMeter        StrUnit = "m"
	SKilometer    StrUnit = "km"
	SFoot         StrUnit = "ft"
	SInch         StrUnit = "in"
	SMile         StrUnit = "mi"
	SNauticalMile StrUnit = "nmi"
)

// strFactors is written only at init (P2).
var strFactors = map[StrUnit]float64{
	SMeter: 1, SKilometer: 1000, SFoot: 0.3048, SInch: 0.0254, SMile: 1609.344, SNauticalMile: 1852,
}

var errUnknownUnit = errors.New("unknown distance unit")

// DistanceS is a distance value of option S. Its zero value has no unit.
type DistanceS struct {
	v float64
	u StrUnit
}

// NewDistanceS validates u: a StrUnit can hold any string.
func NewDistanceS(v float64, u StrUnit) (DistanceS, error) {
	if _, ok := strFactors[u]; !ok {
		return DistanceS{}, errUnknownUnit
	}

	return DistanceS{v, u}, nil
}

// ConvertTo returns d in unit u.
func (d DistanceS) ConvertTo(u StrUnit) (DistanceS, error) {
	to, ok := strFactors[u]
	if !ok {
		return DistanceS{}, errUnknownUnit
	}

	return DistanceS{d.v * strFactors[d.u] / to, u}, nil
}

// Add returns d+o in d's unit.
func (d DistanceS) Add(o DistanceS) DistanceS {
	if o.u == d.u {
		return DistanceS{d.v + o.v, d.u}
	}

	return DistanceS{d.v + o.v*strFactors[o.u]/strFactors[d.u], d.u}
}

// Compare compares d and o in the base unit unless they share a unit.
func (d DistanceS) Compare(o DistanceS) int {
	if o.u == d.u {
		return cmp.Compare(d.v, o.v)
	}

	return cmp.Compare(d.v*strFactors[d.u], o.v*strFactors[o.u])
}
