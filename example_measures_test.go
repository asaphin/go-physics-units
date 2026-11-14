package units_test

import (
	"fmt"

	units "github.com/asaphin/go-physics-units"
	"github.com/asaphin/go-physics-units/acceleration"
	"github.com/asaphin/go-physics-units/angle"
	"github.com/asaphin/go-physics-units/area"
	"github.com/asaphin/go-physics-units/conversion"
	"github.com/asaphin/go-physics-units/datarate"
	"github.com/asaphin/go-physics-units/datasize"
	"github.com/asaphin/go-physics-units/density"
	"github.com/asaphin/go-physics-units/distance"
	"github.com/asaphin/go-physics-units/energy"
	"github.com/asaphin/go-physics-units/force"
	"github.com/asaphin/go-physics-units/frequency"
	"github.com/asaphin/go-physics-units/mass"
	"github.com/asaphin/go-physics-units/power"
	"github.com/asaphin/go-physics-units/pressure"
	"github.com/asaphin/go-physics-units/temperature"
	"github.com/asaphin/go-physics-units/time"
	"github.com/asaphin/go-physics-units/velocity"
	"github.com/asaphin/go-physics-units/volume"
)

// Examples of each measure, the cross-measure operations and custom measures,
// for pkg.go.dev. The README's examples are in example_test.go.

func ExampleNewDistance() {
	marathon := units.NewDistance(26.2, distance.Mile)

	fmt.Printf("%.3f\n", marathon.ConvertTo(distance.Kilometer))
	// Output: 42.165 km
}

func ExampleNewArea() {
	plot := units.NewArea(1, area.SquareMile)

	fmt.Printf("%.4f\n", plot.ConvertTo(area.SquareKilometer))
	// Output: 2.5900 km²
}

func ExampleNewMass() {
	m := units.NewMass(150, mass.Pound)

	fmt.Printf("%.1f\n", m.ConvertTo(mass.Kilogram))
	// Output: 68.0 kg
}

func ExampleNewTime() {
	film := units.NewTime(135, time.Minute)

	fmt.Println(film.ConvertTo(time.Hour))
	// Output: 2.25 h
}

func ExampleNewVelocity() {
	limit := units.NewVelocity(100, velocity.KilometerPerHour)

	fmt.Printf("%.2f\n", limit.ConvertToBaseUnits())
	fmt.Printf("%.1f\n", limit.ConvertTo(velocity.MilePerHour))
	// Output:
	// 27.78 m/s
	// 62.1 mph
}

func ExampleDistance_DivideByTime() {
	marathon := units.NewDistance(42.195, distance.Kilometer)
	finish := units.NewTime(2, time.Hour)

	pace, err := marathon.DivideByTime(finish)
	if err != nil {
		panic(err) // the time is zero
	}

	fmt.Printf("%.4f\n", pace) // in the base unit, m/s
	fmt.Printf("%.4f\n", pace.ConvertTo(velocity.KilometerPerHour))
	// Output:
	// 5.8604 m/s
	// 21.0975 km/h
}

func ExampleDistance_Square() {
	side := units.NewDistance(20, distance.Meter)

	fmt.Println(side.Square())
	// Output: 400 m²
}

func ExampleArea_DivideByDistance() {
	field := units.NewArea(7140, area.SquareMeter)
	width := units.NewDistance(68, distance.Meter)

	length, err := field.DivideByDistance(width)
	if err != nil {
		panic(err) // the distance is zero
	}

	fmt.Println(length)
	// Output: 105 m
}

func ExampleVelocity_MultiplyByTime() {
	speed := units.NewVelocity(60, velocity.KilometerPerHour)
	trip := units.NewTime(30, time.Minute)

	fmt.Printf("%.2f\n", speed.MultiplyByTime(trip).ConvertTo(distance.Kilometer))
	// Output: 30.00 km
}

// Measurements implement fmt.Formatter: the float verbs format the value and
// keep the unit.
func Example_formatting() {
	d := units.NewDistance(1234.5678, distance.Meter)

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
	// units.Distance{Value:1234.5678, Unit:"m"}
}

// dataUnits is a user-defined measure: a marker type whose Units method
// returns a UnitSet built once.
type dataUnits struct{}

var dataUnitSet = units.MustNewUnitSet("storage", "B", conversion.Factors{"B": 1, "KiB": 1024, "MiB": 1024 * 1024})

func (dataUnits) Units() units.UnitSet { return dataUnitSet }

func ExampleNewUnitSet() {
	kib := units.MustUnitOf[dataUnits]("KiB")

	size, err := units.ParseQuantity[dataUnits]("3.5 MiB")
	if err != nil {
		panic(err)
	}

	fmt.Println(size.ConvertTo(kib), size.Type())
	// Output: 3584 KiB storage
}

func ExampleNewVolume() {
	bottle := units.NewVolume(2, volume.Liter)

	fmt.Printf("%.4f\n", bottle.ConvertTo(volume.Gallon))
	// Output: 0.5283 gal
}

func ExampleNewAcceleration() {
	g := units.NewAcceleration(1, acceleration.StandardGravity)

	fmt.Printf("%.3f\n", g.ConvertTo(acceleration.FootPerSecondSquared))
	// Output: 32.174 ft/s²
}

func ExampleNewForce() {
	f := units.NewForce(1, force.Kilonewton)

	fmt.Printf("%.1f\n", f.ConvertTo(force.PoundForce))
	// Output: 224.8 lbf
}

func ExampleNewEnergy() {
	diet := units.NewEnergy(2000, energy.Kilocalorie)

	fmt.Printf("%.3f\n", diet.ConvertTo(energy.KilowattHour))
	// Output: 2.324 kWh
}

func ExampleNewPower() {
	engine := units.NewPower(100, power.Horsepower)

	fmt.Printf("%.2f\n", engine.ConvertTo(power.Kilowatt))
	// Output: 74.57 kW
}

func ExampleNewPressure() {
	tyre := units.NewPressure(32, pressure.PoundPerSquareInch)

	fmt.Printf("%.3f\n", tyre.ConvertTo(pressure.Bar))
	// Output: 2.206 bar
}

func ExampleNewFrequency() {
	motor := units.NewFrequency(3000, frequency.RevolutionPerMinute)

	fmt.Println(motor.ConvertTo(frequency.Hertz))
	// Output: 50 Hz
}

func ExampleFrequency_Period() {
	mains := units.NewFrequency(50, frequency.Hertz)

	period, err := mains.Period()
	if err != nil {
		panic(err) // the frequency is zero
	}

	d, _ := period.ToDuration()
	fmt.Println(period, d)
	// Output: 0.02 s 20ms
}

func ExampleVolume_DivideByArea() {
	rain := units.NewVolume(1, volume.Liter)
	gauge := units.NewArea(1, area.SquareMeter)

	depth, err := rain.DivideByArea(gauge)
	if err != nil {
		panic(err) // the area is zero
	}

	fmt.Printf("%.1f\n", depth.ConvertTo(distance.Millimeter))
	// Output: 1.0 mm
}

// Climbing 3 m of stairs in 4 s: the work against gravity and the power it
// takes, through Mass, Force, Energy and Power.
func Example_work() {
	climber := units.NewMass(75, mass.Kilogram)
	g := units.NewAcceleration(1, acceleration.StandardGravity)
	height := units.NewDistance(3, distance.Meter)
	climb := units.NewTime(4, time.Second)

	weight := climber.MultiplyByAcceleration(g)
	work := weight.MultiplyByDistance(height)

	p, err := work.DivideByTime(climb)
	if err != nil {
		panic(err) // the time is zero
	}

	fmt.Printf("%.1f, %.1f, %.1f = %.2f\n", weight, work, p, p.ConvertTo(power.Horsepower))
	// Output: 735.5 N, 2206.5 J, 551.6 W = 0.74 hp
}

func ExampleNewAngle() {
	right := units.NewAngle(90, angle.Degree)

	fmt.Printf("%.4f\n", right.ConvertToBaseUnits())
	fmt.Printf("%.2f\n", right.ConvertTo(angle.Gradian)) // π/180 and π/200 both round
	// Output:
	// 1.5708 rad
	// 100.00 gon
}

func ExampleNewDataSize() {
	disk := units.NewDataSize(1, datasize.Terabyte)

	// What a "1 TB" disk shows in binary units.
	fmt.Printf("%.1f\n", disk.ConvertTo(datasize.Gibibyte))
	// Output: 931.3 GiB
}

func ExampleNewDensity() {
	steel := units.NewDensity(7.85, density.GramPerCubicCentimeter)

	beam := units.NewVolume(0.02, volume.CubicMeter)

	fmt.Println(steel.MultiplyByVolume(beam))
	// Output: 157 kg
}

// How long a download takes: a data size divided by a data rate.
func ExampleDataSize_DivideByDataRate() {
	file := units.NewDataSize(4.7, datasize.Gigabyte)
	link := units.NewDataRate(100, datarate.MegabitPerSecond)

	t, err := file.DivideByDataRate(link)
	if err != nil {
		panic(err) // the rate is zero
	}

	d, _ := t.ToDuration()
	fmt.Println(t, d)
	// Output: 376 s 6m16s
}

func ExampleNewTemperature() {
	room, err := units.NewTemperature(20, temperature.Celsius)
	if err != nil {
		panic(err) // below absolute zero
	}

	fmt.Println(room.ConvertTo(temperature.Fahrenheit))
	fmt.Println(room.Kelvin())

	_, err = units.NewTemperature(-300, temperature.Celsius)
	fmt.Println(err)
	// Output:
	// 68 °F
	// 293.15
	// temperature -300 °C is below absolute zero, -273.15 °C
}

// The difference of two temperatures is a TemperatureDifference, which
// converts without an offset.
func ExampleTemperature_Sub() {
	morning, _ := units.StringToTemperature("12 °C")
	noon, _ := units.StringToTemperature("77 °F")

	rise := noon.Sub(morning) // in °F, the receiver's unit
	fmt.Printf("%.1f\n", rise)
	fmt.Printf("%.1f\n", rise.ConvertTo(temperature.Celsius))

	evening, err := noon.SubDifference(units.NewTemperatureDifference(5, temperature.Kelvin))
	if err != nil {
		panic(err) // below absolute zero
	}

	fmt.Printf("%.1f\n", evening)
	// Output:
	// 23.4 °F
	// 13.0 °C
	// 68.0 °F
}
