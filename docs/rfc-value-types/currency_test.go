package prototype_test

import (
	"fmt"

	prototype "github.com/asaphin/go-physics-units/docs/rfc-value-types"
)

// A user-defined measure with the full typed API: the marker type, its unit
// set and an optional alias. This is the "about 10 lines" criterion.

type currency struct{}

var currencyUnits = prototype.MustNewUnitSet("currency", map[string]float64{"USD": 1, "EUR": 1.1, "UAH": 0.027})

func (currency) Units() *prototype.UnitSet { return currencyUnits }

// Money is the typed currency value.
type Money = prototype.Quantity[currency]

func Example_customMeasure() {
	usd, _ := prototype.UnitOf[currency]("USD")
	eur, _ := prototype.UnitOf[currency]("EUR")

	var total Money = prototype.New(10, usd).Add(prototype.New(10, eur))

	fmt.Println(total.Value(), total.Unit())
	fmt.Println(total.ConvertTo(eur).Value())
	// Output:
	// 21 USD
	// 19.09090909090909
}
