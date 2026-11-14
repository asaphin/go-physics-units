package units_test

import (
	"fmt"

	units "github.com/asaphin/go-physics-units"
	"github.com/asaphin/go-physics-units/absorbeddose"
	"github.com/asaphin/go-physics-units/acceleration"
	"github.com/asaphin/go-physics-units/amount"
	"github.com/asaphin/go-physics-units/angle"
	"github.com/asaphin/go-physics-units/angularvelocity"
	"github.com/asaphin/go-physics-units/area"
	"github.com/asaphin/go-physics-units/capacitance"
	"github.com/asaphin/go-physics-units/charge"
	"github.com/asaphin/go-physics-units/conversion"
	"github.com/asaphin/go-physics-units/current"
	"github.com/asaphin/go-physics-units/datarate"
	"github.com/asaphin/go-physics-units/datasize"
	"github.com/asaphin/go-physics-units/density"
	"github.com/asaphin/go-physics-units/distance"
	"github.com/asaphin/go-physics-units/doserate"
	"github.com/asaphin/go-physics-units/energy"
	"github.com/asaphin/go-physics-units/energydensity"
	"github.com/asaphin/go-physics-units/equivalentdose"
	"github.com/asaphin/go-physics-units/flowrate"
	"github.com/asaphin/go-physics-units/force"
	"github.com/asaphin/go-physics-units/frequency"
	"github.com/asaphin/go-physics-units/fuelconsumption"
	"github.com/asaphin/go-physics-units/fueleconomy"
	"github.com/asaphin/go-physics-units/illuminance"
	"github.com/asaphin/go-physics-units/inductance"
	"github.com/asaphin/go-physics-units/luminousflux"
	"github.com/asaphin/go-physics-units/magneticfield"
	"github.com/asaphin/go-physics-units/mass"
	"github.com/asaphin/go-physics-units/massflowrate"
	"github.com/asaphin/go-physics-units/molarity"
	"github.com/asaphin/go-physics-units/molarmass"
	"github.com/asaphin/go-physics-units/momentum"
	"github.com/asaphin/go-physics-units/pace"
	"github.com/asaphin/go-physics-units/power"
	"github.com/asaphin/go-physics-units/pressure"
	"github.com/asaphin/go-physics-units/radioactivity"
	"github.com/asaphin/go-physics-units/ratio"
	"github.com/asaphin/go-physics-units/resistance"
	"github.com/asaphin/go-physics-units/specificenergy"
	"github.com/asaphin/go-physics-units/temperature"
	"github.com/asaphin/go-physics-units/time"
	"github.com/asaphin/go-physics-units/torque"
	"github.com/asaphin/go-physics-units/velocity"
	"github.com/asaphin/go-physics-units/voltage"
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

	speed, err := marathon.DivideByTime(finish)
	if err != nil {
		panic(err) // the time is zero
	}

	fmt.Printf("%.4f\n", speed) // in the base unit, m/s
	fmt.Printf("%.4f\n", speed.ConvertTo(velocity.KilometerPerHour))
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

func ExampleNewCurrent() {
	led := units.NewCurrent(20, current.Milliampere)

	fmt.Println(led.ConvertTo(current.Ampere))
	// Output: 0.02 A
}

func ExampleNewVoltage() {
	mains := units.NewVoltage(230, voltage.Volt)

	fmt.Println(mains.ConvertTo(voltage.Kilovolt))
	// Output: 0.23 kV
}

func ExampleNewCharge() {
	battery := units.NewCharge(5000, charge.MilliampereHour)

	fmt.Println(battery.ConvertTo(charge.Coulomb))
	// Output: 18000 C
}

func ExampleNewResistance() {
	pullUp, err := units.StringToResistance("4.7 kΩ")
	if err != nil {
		panic(err)
	}

	fmt.Println(pullUp.ConvertToBaseUnits())
	// Output: 4700 Ω
}

func ExampleNewAngularVelocity() {
	record := units.NewAngularVelocity(33.3, angularvelocity.RevolutionPerMinute)

	fmt.Printf("%.3f\n", record.ConvertToBaseUnits())
	// Output: 3.487 rad/s
}

func ExampleVoltage_MultiplyByCurrent() {
	kettle := units.NewVoltage(230, voltage.Volt).MultiplyByCurrent(units.NewCurrent(10, current.Ampere))

	fmt.Printf("%.1f\n", kettle.ConvertTo(power.Kilowatt))
	// Output: 2.3 kW
}

func ExampleCurrent_MultiplyByVoltage() {
	p := units.NewCurrent(500, current.Milliampere).MultiplyByVoltage(units.NewVoltage(5, voltage.Volt))

	fmt.Println(p)
	// Output: 2.5 W
}

func ExampleCurrent_MultiplyByResistance() {
	drop := units.NewCurrent(20, current.Milliampere).MultiplyByResistance(units.NewResistance(150, resistance.Ohm))

	fmt.Println(drop)
	// Output: 3 V
}

func ExampleResistance_MultiplyByCurrent() {
	drop := units.NewResistance(1, resistance.Kiloohm).MultiplyByCurrent(units.NewCurrent(5, current.Milliampere))

	fmt.Println(drop)
	// Output: 5 V
}

func ExampleCurrent_MultiplyByTime() {
	q := units.NewCurrent(2, current.Ampere).MultiplyByTime(units.NewTime(3, time.Hour))

	fmt.Println(q.ConvertTo(charge.AmpereHour))
	// Output: 6 Ah
}

func ExampleVoltage_DivideByCurrent() {
	r, err := units.NewVoltage(12, voltage.Volt).DivideByCurrent(units.NewCurrent(3, current.Ampere))
	if err != nil {
		panic(err) // the current is zero
	}

	fmt.Println(r)
	// Output: 4 Ω
}

func ExampleVoltage_DivideByResistance() {
	i, err := units.NewVoltage(5, voltage.Volt).DivideByResistance(units.NewResistance(250, resistance.Ohm))
	if err != nil {
		panic(err) // the resistance is zero
	}

	fmt.Println(i.ConvertTo(current.Milliampere))
	// Output: 20 mA
}

func ExampleVoltage_MultiplyByCharge() {
	e := units.NewVoltage(3.7, voltage.Volt).MultiplyByCharge(units.NewCharge(5000, charge.MilliampereHour))

	fmt.Printf("%.1f\n", e.ConvertTo(energy.WattHour))
	// Output: 18.5 Wh
}

func ExampleCharge_MultiplyByVoltage() {
	e := units.NewCharge(1, charge.Coulomb).MultiplyByVoltage(units.NewVoltage(9, voltage.Volt))

	fmt.Println(e)
	// Output: 9 J
}

func ExampleCharge_DivideByTime() {
	i, err := units.NewCharge(3000, charge.MilliampereHour).DivideByTime(units.NewTime(2, time.Hour))
	if err != nil {
		panic(err) // the time is zero
	}

	fmt.Println(i)
	// Output: 1.5 A
}

// How long a battery lasts: its charge divided by the current drawn.
func ExampleCharge_DivideByCurrent() {
	battery := units.NewCharge(3000, charge.MilliampereHour)

	t, err := battery.DivideByCurrent(units.NewCurrent(250, current.Milliampere))
	if err != nil {
		panic(err) // the current is zero
	}

	fmt.Println(t.ConvertTo(time.Hour))
	// Output: 12 h
}

func ExamplePower_DivideByVoltage() {
	i, err := units.NewPower(2.3, power.Kilowatt).DivideByVoltage(units.NewVoltage(230, voltage.Volt))
	if err != nil {
		panic(err) // the voltage is zero
	}

	fmt.Println(i)
	// Output: 10 A
}

func ExamplePower_DivideByCurrent() {
	v, err := units.NewPower(60, power.Watt).DivideByCurrent(units.NewCurrent(5, current.Ampere))
	if err != nil {
		panic(err) // the current is zero
	}

	fmt.Println(v)
	// Output: 12 V
}

func ExampleEnergy_DivideByCharge() {
	v, err := units.NewEnergy(18, energy.Joule).DivideByCharge(units.NewCharge(2, charge.Coulomb))
	if err != nil {
		panic(err) // the charge is zero
	}

	fmt.Println(v)
	// Output: 9 V
}

func ExampleAngle_DivideByTime() {
	w, err := units.NewAngle(90, angle.Degree).DivideByTime(units.NewTime(3, time.Second))
	if err != nil {
		panic(err) // the time is zero
	}

	fmt.Printf("%.0f\n", w.ConvertTo(angularvelocity.DegreePerSecond))
	// Output: 30 °/s
}

func ExampleAngularVelocity_MultiplyByTime() {
	turned := units.NewAngularVelocity(60, angularvelocity.RevolutionPerMinute).
		MultiplyByTime(units.NewTime(30, time.Second))

	fmt.Printf("%.0f\n", turned.ConvertTo(angle.Revolution))
	// Output: 30 rev
}

func ExampleNewTorque() {
	wheelNut := units.NewTorque(100, torque.PoundForceFoot)

	fmt.Printf("%.1f\n", wheelNut.ConvertTo(torque.NewtonMeter))
	// Output: 135.6 N·m
}

func ExampleNewFlowRate() {
	pump := units.NewFlowRate(10, flowrate.GallonPerMinute)

	fmt.Printf("%.1f\n", pump.ConvertTo(flowrate.LiterPerMinute))
	// Output: 37.9 L/min
}

func ExampleNewMassFlowRate() {
	// An airliner burns about 2.4 t of fuel an hour in cruise.
	burn := units.NewMassFlowRate(2400, massflowrate.KilogramPerHour)

	fmt.Printf("%.0f\n", burn.ConvertTo(massflowrate.PoundPerHour))
	// Output: 5291 lb/h
}

func ExampleNewMomentum() {
	// The Mars Climate Orbiter was lost in 1999 because software read
	// impulses in lbf·s as if they were in N·s.
	impulse := units.NewMomentum(1, momentum.PoundForceSecond)

	fmt.Printf("%.2f\n", impulse.ConvertTo(momentum.NewtonSecond))
	// Output: 4.45 N·s
}

func ExampleNewPace() {
	p := units.NewPace(5, pace.MinutePerKilometer)

	fmt.Printf("%.2f\n", p.ConvertTo(pace.MinutePerMile))
	// Output: 8.05 min/mi
}

func ExampleNewFuelEconomy() {
	car := units.NewFuelEconomy(40, fueleconomy.MilePerGallon)

	fmt.Printf("%.1f\n", car.ConvertTo(fueleconomy.KilometerPerLiter))
	// Output: 17.0 km/L
}

func ExampleNewFuelConsumption() {
	car := units.NewFuelConsumption(6.5, fuelconsumption.LiterPer100Kilometers)

	fmt.Printf("%.2f\n", car.ConvertTo(fuelconsumption.GallonPer100Miles))
	// Output: 2.76 gal/100mi
}

func ExampleNewAmount() {
	n := units.NewAmount(250, amount.Millimole)

	fmt.Println(n.ConvertTo(amount.Mole))
	// Output: 0.25 mol
}

func ExampleNewMolarMass() {
	water := units.NewMolarMass(18.015, molarmass.GramPerMole)

	fmt.Println(water.ConvertToBaseUnits())
	// Output: 0.018015 kg/mol
}

func ExampleNewMolarity() {
	saline := units.NewMolarity(0.154, molarity.MolePerLiter)

	fmt.Println(saline.ConvertTo(molarity.Millimolar))
	// Output: 154 mM
}

func ExampleNewRadioactivity() {
	// A household smoke detector holds about 1 μCi of americium-241.
	source := units.NewRadioactivity(1, radioactivity.Microcurie)

	fmt.Println(source.ConvertTo(radioactivity.Kilobecquerel))
	// Output: 37 kBq
}

func ExampleNewAbsorbedDose() {
	// A typical radiotherapy fraction.
	fraction := units.NewAbsorbedDose(2, absorbeddose.Gray)

	fmt.Println(fraction.ConvertTo(absorbeddose.Centigray))
	// Output: 200 cGy
}

func ExampleNewEquivalentDose() {
	// The worldwide average yearly dose from natural sources.
	background := units.NewEquivalentDose(2.4, equivalentdose.Millisievert)

	fmt.Println(background.ConvertTo(equivalentdose.Millirem))
	// Output: 240 mrem
}

func ExampleNewDoseRate() {
	// At cruising altitude, cosmic rays give about 5 μSv/h.
	cruise := units.NewDoseRate(5, doserate.MicrosievertPerHour)

	fmt.Printf("%.1f\n", cruise.ConvertTo(doserate.MillisievertPerYear))
	// Output: 43.8 mSv/yr
}

func ExampleNewCapacitance() {
	c := units.NewCapacitance(4700, capacitance.Nanofarad)

	fmt.Println(c.ConvertTo(capacitance.Microfarad))
	// Output: 4.7 μF
}

func ExampleNewInductance() {
	coil := units.NewInductance(220, inductance.Microhenry)

	fmt.Println(coil.ConvertTo(inductance.Millihenry))
	// Output: 0.22 mH
}

func ExampleNewMagneticField() {
	// The Earth's magnetic field at the surface is about half a gauss.
	earth := units.NewMagneticField(0.5, magneticfield.Gauss)

	fmt.Printf("%.0f\n", earth.ConvertTo(magneticfield.Microtesla))
	// Output: 50 μT
}

func ExampleNewIlluminance() {
	office := units.NewIlluminance(50, illuminance.FootCandle)

	fmt.Printf("%.0f\n", office.ConvertTo(illuminance.Lux))
	// Output: 538 lx
}

func ExampleNewLuminousFlux() {
	lamp := units.NewLuminousFlux(1.6, luminousflux.Kilolumen)

	fmt.Println(lamp.ConvertTo(luminousflux.Lumen))
	// Output: 1600 lm
}

func ExampleNewSpecificEnergy() {
	// A lithium-ion cell stores about 250 Wh per kilogram.
	cell := units.NewSpecificEnergy(250, specificenergy.WattHourPerKilogram)

	fmt.Println(cell.ConvertTo(specificenergy.MegajoulePerKilogram))
	// Output: 0.9 MJ/kg
}

func ExampleNewEnergyDensity() {
	// Petrol holds about 34.2 MJ per litre.
	petrol := units.NewEnergyDensity(34.2, energydensity.MegajoulePerLiter)

	fmt.Printf("%.1f\n", petrol.ConvertTo(energydensity.KilowattHourPerLiter))
	// Output: 9.5 kWh/L
}

func ExampleNewRatio() {
	// The CO₂ in the air.
	co2 := units.NewRatio(420, ratio.PartsPerMillion)

	fmt.Printf("%.3f\n", co2.ConvertTo(ratio.Percent))
	// Output: 0.042 %
}
