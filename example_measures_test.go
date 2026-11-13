package units_test

import (
	"fmt"

	units "github.com/asaphin/go-physics-units"
	"github.com/asaphin/go-physics-units/area"
	"github.com/asaphin/go-physics-units/conversion"
	"github.com/asaphin/go-physics-units/distance"
	"github.com/asaphin/go-physics-units/mass"
	"github.com/asaphin/go-physics-units/time"
	"github.com/asaphin/go-physics-units/velocity"
)

// Examples of each measure, the cross-measure operations and custom measures,
// for pkg.go.dev. The README's examples are in example_test.go.

func ExampleNewDistance() {
	marathon, err := units.NewDistance(26.2, distance.Mile)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%.3f\n", marathon.MustConvertTo(distance.Kilometer))
	// Output: 42.165 km
}

func ExampleNewArea() {
	plot, err := units.NewArea(1, area.SquareMile)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%.4f\n", plot.MustConvertTo(area.SquareKilometer))
	// Output: 2.5900 km²
}

func ExampleNewMass() {
	m, err := units.NewMass(150, mass.Pound)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%.1f\n", m.MustConvertTo(mass.Kilogram))
	// Output: 68.0 kg
}

func ExampleNewTime() {
	film, err := units.NewTime(135, time.Minute)
	if err != nil {
		panic(err)
	}

	fmt.Println(film.MustConvertTo(time.Hour))
	// Output: 2.25 h
}

func ExampleNewVelocity() {
	limit, err := units.NewVelocity(100, velocity.KilometerPerHour)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%.2f\n", limit.ConvertToBaseUnits())
	fmt.Printf("%.1f\n", limit.MustConvertTo(velocity.MilePerHour))
	// Output:
	// 27.78 m/s
	// 62.1 mph
}

func ExampleDistance_DivideByTime() {
	marathon, _ := units.NewDistance(42.195, distance.Kilometer)
	finish, _ := units.NewTime(2, time.Hour)

	pace, err := marathon.DivideByTime(finish)
	if err != nil {
		panic(err) // the time is zero
	}

	fmt.Printf("%.4f\n", pace) // in the base unit, m/s
	fmt.Printf("%.4f\n", pace.MustConvertTo(velocity.KilometerPerHour))
	// Output:
	// 5.8604 m/s
	// 21.0975 km/h
}

func ExampleDistance_Square() {
	side, _ := units.NewDistance(20, distance.Meter)

	fmt.Println(side.Square())
	// Output: 400 m²
}

func ExampleArea_DivideByDistance() {
	field, _ := units.NewArea(7140, area.SquareMeter)
	width, _ := units.NewDistance(68, distance.Meter)

	length, err := field.DivideByDistance(width)
	if err != nil {
		panic(err) // the distance is zero
	}

	fmt.Println(length)
	// Output: 105 m
}

func ExampleVelocity_MultiplyByTime() {
	speed, _ := units.NewVelocity(60, velocity.KilometerPerHour)
	trip, _ := units.NewTime(30, time.Minute)

	fmt.Printf("%.2f\n", speed.MultiplyByTime(trip).MustConvertTo(distance.Kilometer))
	// Output: 30.00 km
}

// Measurements implement fmt.Formatter: the float verbs format the value and
// keep the unit.
func Example_formatting() {
	d, _ := units.NewDistance(1234.5678, distance.Meter)

	fmt.Printf("%v\n", d)
	fmt.Printf("%.1f\n", d)
	fmt.Printf("[%12.2f]\n", d)
	fmt.Printf("%.3e\n", d)
	fmt.Printf("%q\n", d)
	fmt.Printf("%#v\n", d)
	// Output:
	// 1234.5678 m
	// 1234.6 m
	// [     1234.57 m]
	// 1.235e+03 m
	// "1234.5678 m"
	// units.Measurement{Value:1234.5678, Unit:"m", Type:"distance"}
}

// NewBaseMeasurement suits one-off values. It builds a conversion table on
// every call; for many values, use NewMeasurementCreator.
func ExampleNewBaseMeasurement() {
	data := conversion.Factors{"B": 1, "KiB": 1024, "MiB": 1024 * 1024}

	size, err := units.NewBaseMeasurement(3.5, "MiB", data)
	if err != nil {
		panic(err)
	}

	kib, _ := size.ConvertToMeasurement("KiB")

	fmt.Println(kib, kib.Type())
	// Output: 3584 KiB custom
}
