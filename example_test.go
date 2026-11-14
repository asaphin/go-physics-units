package units_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"slices"

	units "github.com/asaphin/go-physics-units"
	"github.com/asaphin/go-physics-units/conversion"
	"github.com/asaphin/go-physics-units/distance"
	"github.com/asaphin/go-physics-units/time"
	"github.com/asaphin/go-physics-units/velocity"
)

func Example() {
	// Creating a distance of 5 kilometers. Units are typed constants, so
	// a unit of another measure doesn't compile.
	dist := units.NewDistance(5, distance.Kilometer)

	fmt.Println(dist)

	// Parsing a distance from a string.
	dist, err := units.StringToDistance("10 km")
	if err != nil {
		panic(err)
	}

	// Converting the distance to nautical miles.
	nmi := dist.ConvertTo(distance.NauticalMile)

	fmt.Println(nmi)

	// fmt verbs format the value and keep the unit.
	fmt.Printf("%.2f\n", nmi)

	// Accessing the value and unit.
	fmt.Println(dist.Value(), dist.Unit())

	// Dividing a distance by a time gives a velocity in m/s.
	tm := units.NewTime(10, time.Minute)

	vel, err := dist.DivideByTime(tm)
	if err != nil {
		panic(err) // the time is zero
	}

	fmt.Println(vel)
	fmt.Println(vel.ConvertTo(velocity.KilometerPerHour))
	// Output:
	// 5 km
	// 5.399568034557235 nmi
	// 5.40 nmi
	// 10 km
	// 16.666666666666668 m/s
	// 60 km/h
}

func ExampleParseString() {
	m, err := units.ParseString("1500 m")
	if err != nil {
		panic(err)
	}

	fmt.Println(m.Type())

	// The result has the measure's concrete type.
	if d, ok := m.(units.Distance); ok {
		fmt.Println(d.ConvertTo(distance.Kilometer))
	}
	// Output:
	// distance
	// 1.5 km
}

func ExampleDistance_compare() {
	mile := units.NewDistance(1, distance.Mile)
	foot := units.NewDistance(1, distance.Feet)
	inches := units.NewDistance(12, distance.Inch)

	fmt.Println(mile.Compare(foot))

	// Conversions round: in metres, 12 in and 1 ft differ in the last bit.
	fmt.Println(foot.Equal(inches), foot.EqualWithinRel(inches, 1e-12))

	ds := []units.Distance{mile, foot, inches}
	slices.SortFunc(ds, units.Distance.Compare)
	fmt.Println(ds)
	// Output:
	// 1
	// false true
	// [12 in 1 ft 1 mi]
}

func ExampleTime_ToDuration() {
	tm := units.NewTime(1.5, time.Hour)

	d, err := tm.ToDuration()
	if err != nil {
		panic(err)
	}

	fmt.Println(d)
	fmt.Println(units.ParseDuration(d))
	// Output:
	// 1h30m0s
	// 5400 s
}

func ExampleUnitError() {
	_, err := units.StringToDistance("3 kg")

	// Every error matches one of the sentinels, such as ErrUnknownUnit.
	fmt.Println(errors.Is(err, units.ErrUnknownUnit))

	// Unknown-unit errors carry the measure and the unit.
	if ue, ok := errors.AsType[*units.UnitError](err); ok {
		fmt.Println(ue.Measure, ue.Unit)
	}

	fmt.Println(err)
	// Output:
	// true
	// distance kg
	// unknown distance unit kg
}

func ExampleDistance_add() {
	km := units.NewDistance(1, distance.Kilometer)
	m := units.NewDistance(250, distance.Meter)

	// The result keeps the receiver's unit.
	fmt.Println(km.Add(m))
	fmt.Println(m.Add(km))
	// Output:
	// 1.25 km
	// 1250 m
}

func ExampleNewParser() {
	p, err := units.NewParser(
		units.WithAliases(map[string]string{"metres": "m", "kph": "km/h"}),
		units.CaseInsensitive(),
		units.AllowNoSpace(),
	)
	if err != nil {
		panic(err)
	}

	for _, s := range []string{"5 metres", "100kph", "10 KM", "1 KN"} {
		m, err := p.Parse(s)
		if err != nil {
			fmt.Println(err)

			continue
		}

		fmt.Println(m, m.Type())
	}
	// Output:
	// 5 m distance
	// 100 km/h velocity
	// 10 km distance
	// ambiguous unit KN: it could be any of kN, kn
}

// A user-defined measure: a marker type, its units, and a name for the value
// type. Money then has the same typed API as the predefined measures.
type exampleCurrency struct{}

var exampleCurrencyUnits = units.MustNewUnitSet("currency", "USD", conversion.Factors{
	"USD": 1,
	"EUR": 1.1,
	"UAH": 0.027,
})

func (exampleCurrency) Units() units.UnitSet { return exampleCurrencyUnits }

type ExampleMoney = units.Quantity[exampleCurrency]

func ExampleQuantity() {
	usd := units.MustUnitOf[exampleCurrency]("USD")
	uah := units.MustUnitOf[exampleCurrency]("UAH")

	price := units.NewQuantity(1, usd)
	tip, err := units.ParseQuantity[exampleCurrency]("0.5 EUR")
	if err != nil {
		panic(err)
	}

	total := price.Add(tip) // an ExampleMoney

	fmt.Println(total)
	fmt.Println(total.ConvertTo(uah))
	fmt.Println(total.Type())
	// Output:
	// 1.55 USD
	// 57.40740740740741 UAH
	// currency
}

// surveying extends the predefined distance units: copy their factors, add
// units, and define a measure of your own.
type surveying struct{}

var surveyingUnits = func() units.UnitSet {
	factors := distance.ConversionFactors() // a copy: changing it is safe
	factors["fur"] = 201.168                // furlong, 660 international feet
	factors["ch"] = 20.1168                 // chain, 66 international feet

	return units.MustNewUnitSet("surveying", distance.BaseUnit.String(), factors)
}()

func (surveying) Units() units.UnitSet { return surveyingUnits }

func ExampleNewUnitSet_extendPredefined() {
	furlong, err := units.ParseQuantity[surveying]("1 fur")
	if err != nil {
		panic(err)
	}

	fmt.Println(furlong.ConvertTo(units.MustUnitOf[surveying]("ch")), furlong.ConvertToBaseUnits(), furlong.Type())

	// The predefined distance measure is unchanged.
	_, err = units.StringToDistance("1 fur")
	fmt.Println(err)
	// Output:
	// 10 ch 201.168 m surveying
	// unknown distance unit fur
}

func ExampleDistance_MarshalText() {
	type Route struct {
		Name   string         `json:"name"`
		Length units.Distance `json:"length"`
	}

	b, err := json.Marshal(Route{"loop", units.NewDistance(12.5, distance.Kilometer)})
	if err != nil {
		panic(err)
	}

	fmt.Println(string(b))

	var r Route
	if err := json.Unmarshal([]byte(`{"name": "ridge", "length": "3.2 mi"}`), &r); err != nil {
		panic(err)
	}

	fmt.Printf("%.2f\n", r.Length.ConvertTo(distance.Kilometer))
	// Output:
	// {"name":"loop","length":"12.5 km"}
	// 5.15 km
}
