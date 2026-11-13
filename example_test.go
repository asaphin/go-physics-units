package units_test

import (
	"fmt"

	units "github.com/asaphin/go-physics-units"
	"github.com/asaphin/go-physics-units/conversion"
	"github.com/asaphin/go-physics-units/distance"
	"github.com/asaphin/go-physics-units/time"
	"github.com/asaphin/go-physics-units/velocity"
)

func Example() {
	// Creating a distance of 5 kilometers.
	dist, err := units.NewDistance(5, distance.Kilometer)
	if err != nil {
		panic(err)
	}

	fmt.Println(dist)

	// Parsing a distance from a string.
	dist, err = units.StringToDistance("10 km")
	if err != nil {
		panic(err)
	}

	// Converting the distance to nautical miles.
	nmi, err := dist.ConvertTo(distance.NauticalMile)
	if err != nil {
		panic(err)
	}

	fmt.Println(nmi)

	// Accessing the value and unit.
	fmt.Println(dist.Value(), dist.Unit())

	// Dividing a distance by a time gives a velocity in m/s.
	tm, err := units.NewTime(10, time.Minute)
	if err != nil {
		panic(err)
	}

	vel, err := dist.DivideByTime(tm)
	if err != nil {
		panic(err)
	}

	fmt.Println(vel)
	fmt.Println(vel.MustConvertTo(velocity.KilometerPerHour))
	// Output:
	// 5 km
	// 5.399568034557235 nmi
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

	d, err := units.ToDistance(m)
	if err != nil {
		panic(err)
	}

	fmt.Println(d.MustConvertTo(distance.Kilometer))
	// Output:
	// distance
	// 1.5 km
}

func ExampleDistance_add() {
	km, _ := units.NewDistance(1, distance.Kilometer)
	m, _ := units.NewDistance(250, distance.Meter)

	// The result keeps the receiver's unit.
	fmt.Println(km.Add(m))
	fmt.Println(m.Add(km))
	// Output:
	// 1.25 km
	// 1250 m
}

func ExampleNewMeasurementCreator() {
	type Currency interface {
		units.Measurement
	}

	currencyFactors := conversion.Factors{
		"EUR": 1.1,
		"USD": 1,
		"UAH": 0.027,
	}

	mc, err := units.NewMeasurementCreator[Currency]("currency", "USD", currencyFactors)
	if err != nil {
		panic(err)
	}

	usd, err := mc.New(1, "USD")
	if err != nil {
		panic(err)
	}

	uah, err := usd.ConvertToMeasurement("UAH")
	if err != nil {
		panic(err)
	}

	fmt.Println(uah)
	fmt.Println(uah.Type())
	// Output:
	// 37.03703703703704 UAH
	// currency
}
