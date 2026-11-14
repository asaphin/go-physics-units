// Package temperature defines the temperature units of the units package:
// the kelvin, the degree Celsius, the degree Fahrenheit and the degree Rankine.
//
// Temperature scales have offsets (0 °C is 273.15 K), so unlike the other unit
// packages this one has no factors: Convert converts a temperature, a point on
// a scale, and DifferenceRate converts a temperature difference, which has no
// offset. A rise of 10 °C is a rise of 18 °F, while 10 °C is 50 °F.
package temperature

import (
	"iter"
	"math/big"
	"strconv"

	"github.com/asaphin/go-physics-units/internal/rates"
)

// Unit is a temperature unit: one of the constants below. String returns its
// symbol and ParseUnit finds the unit of a symbol.
//
// The zero Unit is the base unit, Kelvin. The numbers behind the constants may
// change between versions: store and exchange units by their symbols.
type Unit uint8

// BaseUnit is the SI unit of temperature.
const BaseUnit = Kelvin

// Temperature units, with their symbols.
const (
	Kelvin     Unit = iota // K
	Celsius                // °C, with U+00B0 DEGREE SIGN; "℃" (U+2103 DEGREE CELSIUS) is accepted too
	Fahrenheit             // °F, with U+00B0 DEGREE SIGN; "℉" (U+2109 DEGREE FAHRENHEIT) is accepted too
	Rankine                // °R, with U+00B0 DEGREE SIGN

	unitCount // the number of units; keep it last
)

// definition is a temperature scale: T in kelvins is (value + offset) · scale,
// with exact rational offset and scale.
type definition struct {
	symbol string
	offset *big.Rat // absolute zero is at -offset
	scale  *big.Rat // kelvins per degree
}

// definitions holds the scales, from their definitions in the SI Brochure (9th
// edition) and NIST SP 811, appendix B.
var definitions = [unitCount]definition{
	Kelvin:     {"K", big.NewRat(0, 1), big.NewRat(1, 1)},
	Celsius:    {"\u00b0C", big.NewRat(27315, 100), big.NewRat(1, 1)},
	Fahrenheit: {"\u00b0F", big.NewRat(45967, 100), big.NewRat(5, 9)},
	Rankine:    {"\u00b0R", big.NewRat(0, 1), big.NewRat(5, 9)},
}

// coefficients converts between two units: y = x·a + b for temperatures and
// y = x·a for differences.
type coefficients struct {
	a, b float64
}

// table holds the coefficients of every pair of units and the absolute zero
// of every unit. It is computed at package init from the exact definitions,
// so each coefficient is rounded only once, and never written afterwards.
var table, absoluteZero = build()

func build() ([unitCount][unitCount]coefficients, [unitCount]float64) {
	var (
		t    [unitCount][unitCount]coefficients
		zero [unitCount]float64
	)

	for from, f := range definitions {
		zero[from], _ = new(big.Rat).Neg(f.offset).Float64()

		for to, d := range definitions {
			// K = (x + f.offset)·f.scale and y = K/d.scale - d.offset, so
			// y = x·(f.scale/d.scale) + (f.offset·f.scale/d.scale - d.offset).
			a := new(big.Rat).Quo(f.scale, d.scale)
			b := new(big.Rat).Mul(f.offset, a)
			b.Sub(b, d.offset)

			t[from][to].a, _ = a.Float64()
			t[from][to].b, _ = b.Float64()
		}
	}

	return t, zero
}

// String returns the unit's symbol. A Unit that isn't one of the constants of
// this package returns "temperature.Unit(<number>)".
func (u Unit) String() string {
	if u >= unitCount {
		return "temperature.Unit(" + strconv.Itoa(int(u)) + ")"
	}

	return definitions[u].symbol
}

// Valid reports whether u is one of the constants of this package. Values in
// any other unit panic when converted.
func (u Unit) Valid() bool {
	return u < unitCount
}

// AbsoluteZero returns absolute zero in unit u: 0 K, -273.15 °C, -459.67 °F
// or 0 °R. It panics if u isn't valid.
func (u Unit) AbsoluteZero() float64 {
	return absoluteZero[u]
}

// ParseUnit returns the unit with the given symbol, and whether there is one.
// U+2103 DEGREE CELSIUS and U+2109 DEGREE FAHRENHEIT are accepted for °C and
// °F, and U+212A KELVIN SIGN for K.
func ParseUnit(symbol string) (Unit, bool) {
	symbol = rates.CanonicalUnit(symbol)

	switch symbol {
	case "\u2103":
		return Celsius, true
	case "\u2109":
		return Fahrenheit, true
	}

	for u, d := range definitions {
		if d.symbol == symbol {
			return Unit(u), true
		}
	}

	return 0, false
}

// Units returns the units in the order of their constants, starting with the
// base unit.
func Units() iter.Seq[Unit] {
	return func(yield func(Unit) bool) {
		for u := range unitCount {
			if !yield(u) {
				return
			}
		}
	}
}

// Convert converts the temperature value, a point on the scale of unit from,
// into unit to: Convert(100, Celsius, Fahrenheit) is 212. It panics if either
// unit isn't valid.
//
// The coefficients of each pair are computed exactly and rounded once, so
// conversions between the defined points are exact: 100 °C is exactly 212 °F
// and 0 °C exactly 32 °F. Other values round in the last bit, as with every
// conversion: 37 °C is 98.60000000000001 °F.
func Convert(value float64, from, to Unit) float64 {
	c := table[from][to]

	// The explicit conversion rounds the product, so that no platform fuses
	// the multiplication and the addition into one, differently rounded,
	// operation.
	return float64(value*c.a) + c.b
}

// DifferenceRate returns the multiplier that converts a temperature
// difference in unit from into unit to: DifferenceRate(Celsius, Fahrenheit)
// is 1.8, so a rise of 10 °C is a rise of 18 °F. It panics if either unit
// isn't valid.
func DifferenceRate(from, to Unit) float64 {
	return table[from][to].a
}
