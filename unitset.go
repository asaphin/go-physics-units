package units

import (
	"cmp"
	"iter"
	"math"
	"slices"
	"strconv"

	"github.com/asaphin/go-physics-units/conversion"
	"github.com/asaphin/go-physics-units/internal/rates"
)

// UnitSet is the set of units of a user-defined measure: their symbols and
// conversion factors. Build one with NewUnitSet, once, and return it from the
// Units method of the measure's marker type (see Measure).
//
// A UnitSet is immutable and safe for concurrent use. It is a small value that
// refers to data no one can change, so copy it freely. The zero UnitSet has no
// units.
type UnitSet struct {
	s *unitSet
}

// unitSet holds the units of one measure, indexed from 0. The predefined
// measures use the numbers of their Unit constants as indexes; user-defined
// measures have their base unit at 0 and the other symbols in sorted order.
type unitSet struct {
	name    MeasureType
	symbols []string
	factors []float64
	rates   []float64                         // n×n, row-major: the rate from unit i to unit j is at i·n + j
	index   map[string]uint16                 // canonical symbol → unit
	box     func(float64, uint16) Measurement // the predefined measure's value type; nil for user-defined measures
}

// maxUnits is the number of units a UnitSet can hold: units are uint16 indexes.
const maxUnits = math.MaxUint16 + 1

// NewUnitSet returns the unit set of the user-defined measure named measure,
// whose units are described by factors: how many base units there are in one
// unit of each symbol (see conversion.Factors). The factors are copied, so
// later changes to the map have no effect.
//
// measure must be non-empty and differ from the predefined MeasureType
// values, baseUnit must be one of the symbols, and every factor must be
// positive and finite. It fails with ErrInvalidMeasure or ErrInvalidFactors
// otherwise.
//
// The base unit is the zero Unit of the measure, and the unit of the zero
// Quantity. To extend a predefined measure, start from its factors:
//
//	f := distance.ConversionFactors()
//	f["fur"] = 201.168 // the furlong
//	set, err := units.NewUnitSet("long distance", distance.BaseUnit.String(), f)
func NewUnitSet(measure, baseUnit string, factors conversion.Factors) (UnitSet, error) {
	if err := checkMeasureName(measure); err != nil {
		return UnitSet{}, err
	}

	if len(factors) > maxUnits {
		return UnitSet{}, newError(
			ErrInvalidFactors,
			"a measure can't have more than "+strconv.Itoa(maxUnits)+" units",
			nil,
		)
	}

	if _, ok := factors[baseUnit]; !ok {
		return UnitSet{}, newError(ErrInvalidFactors, "provided factors don't contain the base unit "+baseUnit, nil)
	}

	symbols := orderedSymbols(baseUnit, factors)

	values := make([]float64, len(symbols))
	for i, symbol := range symbols {
		values[i] = factors[symbol]
	}

	if err := rates.Validate(symbols, values); err != nil {
		return UnitSet{}, newError(ErrInvalidFactors, err.Error(), nil)
	}

	for i, symbol := range symbols {
		symbols[i] = rates.CanonicalUnit(symbol)
	}

	return UnitSet{s: newUnitSet(MeasureType(measure), symbols, values)}, nil
}

// checkMeasureName checks the name of a user-defined measure.
func checkMeasureName(measure string) error {
	if measure == "" {
		return newError(ErrInvalidMeasure, "measure type must not be empty", nil)
	}

	if isReservedMeasureType(MeasureType(measure)) {
		return newError(ErrInvalidMeasure, "measure type "+measure+" is reserved for a predefined measure", nil)
	}

	return nil
}

// orderedSymbols returns the symbols of factors in the order of a unit set:
// the base unit first, the others sorted. The indexes, and so the zero Unit,
// then don't depend on map iteration order.
func orderedSymbols(baseUnit string, factors conversion.Factors) []string {
	symbols := make([]string, 0, len(factors))

	for symbol := range factors {
		if symbol != baseUnit {
			symbols = append(symbols, symbol)
		}
	}

	slices.Sort(symbols)

	return slices.Insert(symbols, 0, baseUnit)
}

// MustNewUnitSet is like NewUnitSet but panics if the arguments are invalid.
// It is intended for package-level variables built from constant data:
//
//	var currencyUnits = units.MustNewUnitSet("currency", "USD", conversion.Factors{"USD": 1, "EUR": 1.1})
func MustNewUnitSet(measure, baseUnit string, factors conversion.Factors) UnitSet {
	s, err := NewUnitSet(measure, baseUnit, factors)
	if err != nil {
		panic(err)
	}

	return s
}

// Type returns the name of the measure, or "" for the zero UnitSet.
func (s UnitSet) Type() MeasureType {
	if s.s == nil {
		return ""
	}

	return s.s.name
}

// Symbols returns the unit symbols, in canonical form: the base unit first,
// then the others in sorted order.
func (s UnitSet) Symbols() iter.Seq[string] {
	return func(yield func(string) bool) {
		if s.s == nil {
			return
		}

		for _, symbol := range s.s.symbols {
			if !yield(symbol) {
				return
			}
		}
	}
}

// ConversionFactors returns a copy of the factors the set was built from,
// keyed by canonical symbol.
func (s UnitSet) ConversionFactors() conversion.Factors {
	if s.s == nil {
		return conversion.Factors{}
	}

	factors := make(conversion.Factors, len(s.s.symbols))
	for i, symbol := range s.s.symbols {
		factors[symbol] = s.s.factors[i]
	}

	return factors
}

// newUnitSet builds a set from validated, canonical symbols and their factors.
func newUnitSet(name MeasureType, symbols []string, factors []float64) *unitSet {
	index := make(map[string]uint16, len(symbols))
	for i, symbol := range symbols {
		index[symbol] = uint16(i)
	}

	return &unitSet{
		name:    name,
		symbols: symbols,
		factors: factors,
		rates:   rates.Matrix(factors),
		index:   index,
	}
}

// newPredefinedSet builds the set of a predefined measure from its unit
// package, whose Units iterates over the constants 0, 1, ….
func newPredefinedSet[U interface {
	~uint8
	String() string
	Factor() float64
}](name MeasureType, all iter.Seq[U], box func(float64, U) Measurement,
) *unitSet {
	var (
		symbols []string
		factors []float64
	)

	for u := range all {
		symbols = append(symbols, u.String())
		factors = append(factors, u.Factor())
	}

	s := newUnitSet(name, symbols, factors)
	s.box = func(v float64, u uint16) Measurement { return box(v, U(u)) }

	return s
}

// rate returns the multiplier that converts a value in unit from into unit to.
func (s *unitSet) rate(from, to uint16) float64 {
	return s.rates[int(from)*len(s.symbols)+int(to)]
}

// valueIn returns the value of m in unit to of s. It fails if m is nil or of
// another measure.
func (s *unitSet) valueIn(m Measurement, to uint16) (float64, error) {
	o, err := s.operand(m)
	if err != nil {
		return 0, err
	}

	return o.value * s.rate(o.unit, to), nil
}

// compare compares a value in unit u of s with m, as Comparer.Compare does:
// exactly in a shared unit, otherwise through the factors.
func (s *unitSet) compare(value float64, u uint16, m Measurement) (int, error) {
	o, err := s.operand(m)
	if err != nil {
		return 0, err
	}

	if o.unit == u {
		return cmp.Compare(value, o.value), nil
	}

	return cmp.Compare(value*s.factors[u], o.value*s.factors[o.unit]), nil
}

// operand returns m as a value of s, or an error if m is nil or of another
// measure.
func (s *unitSet) operand(m Measurement) (genericValue, error) {
	if m == nil {
		return genericValue{}, errNilMeasurement
	}

	o := m.generic()
	if o.set == s {
		return o, nil
	}

	if o.set == nil {
		return genericValue{}, errNoUnits
	}

	if o.set.name == s.name {
		return genericValue{}, newError(ErrIncompatibleMeasures,
			"incompatible "+string(s.name)+" measures: they have different unit sets", nil)
	}

	return genericValue{}, newError(ErrIncompatibleMeasures,
		"incompatible types "+string(s.name)+" and "+string(o.set.name), nil)
}

// unitOf returns the index of the unit with the given symbol, or an error
// if s has no such unit.
func (s *unitSet) unitOf(symbol string) (uint16, error) {
	u, ok := s.index[rates.CanonicalUnit(symbol)]
	if !ok {
		return 0, unknownUnitError(s.name, symbol)
	}

	return u, nil
}
