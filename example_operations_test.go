package units_test

import (
	"fmt"

	units "github.com/asaphin/go-physics-units"
	"github.com/asaphin/go-physics-units/amount"
	"github.com/asaphin/go-physics-units/angle"
	"github.com/asaphin/go-physics-units/angularvelocity"
	"github.com/asaphin/go-physics-units/area"
	"github.com/asaphin/go-physics-units/capacitance"
	"github.com/asaphin/go-physics-units/charge"
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
	"github.com/asaphin/go-physics-units/mass"
	"github.com/asaphin/go-physics-units/massflowrate"
	"github.com/asaphin/go-physics-units/molarity"
	"github.com/asaphin/go-physics-units/molarmass"
	"github.com/asaphin/go-physics-units/momentum"
	"github.com/asaphin/go-physics-units/pace"
	"github.com/asaphin/go-physics-units/power"
	"github.com/asaphin/go-physics-units/pressure"
	"github.com/asaphin/go-physics-units/resistance"
	"github.com/asaphin/go-physics-units/specificenergy"
	"github.com/asaphin/go-physics-units/time"
	"github.com/asaphin/go-physics-units/torque"
	"github.com/asaphin/go-physics-units/velocity"
	"github.com/asaphin/go-physics-units/voltage"
	"github.com/asaphin/go-physics-units/volume"
)

func ExampleDistance_DivideByVelocity() {
	// A marathon at 12 km/h.
	marathon := units.NewDistance(42.195, distance.Kilometer)
	speed := units.NewVelocity(12, velocity.KilometerPerHour)

	duration, err := marathon.DivideByVelocity(speed)
	if err != nil {
		panic(err) // the velocity is zero
	}

	fmt.Printf("%.2f\n", duration.ConvertTo(time.Hour))
	// Output: 3.52 h
}

func ExampleEnergy_DivideByPower() {
	// How long a 1 kWh battery runs a 100 W load.
	battery := units.NewEnergy(1, energy.KilowattHour)
	load := units.NewPower(100, power.Watt)

	runtime, err := battery.DivideByPower(load)
	if err != nil {
		panic(err) // the power is zero
	}

	fmt.Printf("%.1f\n", runtime.ConvertTo(time.Hour))
	// Output: 10.0 h
}

func ExampleMass_MultiplyByVelocity() {
	car := units.NewMass(1000, mass.Kilogram)
	speed := units.NewVelocity(72, velocity.KilometerPerHour)

	p := car.MultiplyByVelocity(speed)

	fmt.Printf("%.0f\n", p.ConvertTo(momentum.NewtonSecond))
	// Output: 20000 N·s
}

func ExampleVelocity_MultiplyByMass() {
	// A baseball pitched at 40 m/s.
	pitch := units.NewVelocity(40, velocity.MeterPerSecond)
	ball := units.NewMass(145, mass.Gram)

	p := pitch.MultiplyByMass(ball)

	fmt.Printf("%.1f\n", p.ConvertTo(momentum.KilogramMeterPerSecond))
	// Output: 5.8 kg·m/s
}

func ExampleMomentum_DivideByMass() {
	p := units.NewMomentum(20000, momentum.NewtonSecond)
	car := units.NewMass(1000, mass.Kilogram)

	speed, err := p.DivideByMass(car)
	if err != nil {
		panic(err) // the mass is zero
	}

	fmt.Printf("%.0f\n", speed.ConvertTo(velocity.KilometerPerHour))
	// Output: 72 km/h
}

func ExampleMomentum_DivideByVelocity() {
	p := units.NewMomentum(5.8, momentum.KilogramMeterPerSecond)
	speed := units.NewVelocity(40, velocity.MeterPerSecond)

	ball, err := p.DivideByVelocity(speed)
	if err != nil {
		panic(err) // the velocity is zero
	}

	fmt.Printf("%.0f\n", ball.ConvertTo(mass.Gram))
	// Output: 145 g
}

func ExampleForce_MultiplyByTime() {
	// The impulse of a 4 s engine burn.
	thrust := units.NewForce(500, force.Newton)
	burn := units.NewTime(4, time.Second)

	impulse := thrust.MultiplyByTime(burn)

	fmt.Printf("%.1f\n", impulse.ConvertTo(momentum.PoundForceSecond))
	// Output: 449.6 lbf·s
}

func ExampleMomentum_DivideByTime() {
	impulse := units.NewMomentum(2000, momentum.NewtonSecond)
	burn := units.NewTime(4, time.Second)

	thrust, err := impulse.DivideByTime(burn)
	if err != nil {
		panic(err) // the time is zero
	}

	fmt.Printf("%.0f\n", thrust.ConvertTo(force.Newton))
	// Output: 500 N
}

func ExampleForce_TorqueAt() {
	// Pushing on a 40 cm wrench.
	push := units.NewForce(200, force.Newton)
	wrench := units.NewDistance(40, distance.Centimeter)

	tq := push.TorqueAt(wrench)

	fmt.Printf("%.1f\n", tq.ConvertTo(torque.PoundForceFoot))
	// Output: 59.0 lbf·ft
}

func ExampleTorque_DivideByDistance() {
	tq := units.NewTorque(80, torque.NewtonMeter)
	wrench := units.NewDistance(40, distance.Centimeter)

	push, err := tq.DivideByDistance(wrench)
	if err != nil {
		panic(err) // the distance is zero
	}

	fmt.Printf("%.0f\n", push.ConvertTo(force.Newton))
	// Output: 200 N
}

func ExampleTorque_MultiplyByAngularVelocity() {
	// An engine's power from its torque and speed.
	tq := units.NewTorque(400, torque.NewtonMeter)
	speed := units.NewAngularVelocity(3000, angularvelocity.RevolutionPerMinute)

	p := tq.MultiplyByAngularVelocity(speed)

	fmt.Printf("%.0f\n", p.ConvertTo(power.Horsepower))
	// Output: 169 hp
}

func ExampleAngularVelocity_MultiplyByTorque() {
	speed := units.NewAngularVelocity(3000, angularvelocity.RevolutionPerMinute)
	tq := units.NewTorque(400, torque.NewtonMeter)

	p := speed.MultiplyByTorque(tq)

	fmt.Printf("%.1f\n", p.ConvertTo(power.Kilowatt))
	// Output: 125.7 kW
}

func ExamplePower_DivideByAngularVelocity() {
	// The torque of a 300 hp engine at 6000 rev/min.
	p := units.NewPower(300, power.Horsepower)
	speed := units.NewAngularVelocity(6000, angularvelocity.RevolutionPerMinute)

	tq, err := p.DivideByAngularVelocity(speed)
	if err != nil {
		panic(err) // the angular velocity is zero
	}

	fmt.Printf("%.0f\n", tq.ConvertTo(torque.PoundForceFoot))
	// Output: 263 lbf·ft
}

func ExamplePower_DivideByTorque() {
	p := units.NewPower(100, power.Kilowatt)
	tq := units.NewTorque(250, torque.NewtonMeter)

	speed, err := p.DivideByTorque(tq)
	if err != nil {
		panic(err) // the torque is zero
	}

	fmt.Printf("%.0f\n", speed.ConvertTo(angularvelocity.RevolutionPerMinute))
	// Output: 3820 rev/min
}

func ExampleTorque_MultiplyByAngle() {
	tq := units.NewTorque(50, torque.NewtonMeter)
	turns := units.NewAngle(10, angle.Revolution)

	work := tq.MultiplyByAngle(turns)

	fmt.Printf("%.2f\n", work.ConvertTo(energy.Kilojoule))
	// Output: 3.14 kJ
}

func ExampleVolume_DivideByTime() {
	// Filling a bath in 8 minutes.
	bath := units.NewVolume(200, volume.Liter)
	fill := units.NewTime(8, time.Minute)

	flow, err := bath.DivideByTime(fill)
	if err != nil {
		panic(err) // the time is zero
	}

	fmt.Printf("%.0f\n", flow.ConvertTo(flowrate.LiterPerMinute))
	// Output: 25 L/min
}

func ExampleFlowRate_MultiplyByTime() {
	shower := units.NewFlowRate(12, flowrate.LiterPerMinute)
	minutes := units.NewTime(10, time.Minute)

	water := shower.MultiplyByTime(minutes)

	fmt.Printf("%.0f\n", water.ConvertTo(volume.Liter))
	// Output: 120 L
}

func ExampleVolume_DivideByFlowRate() {
	pool := units.NewVolume(50, volume.CubicMeter)
	pump := units.NewFlowRate(10, flowrate.CubicMeterPerHour)

	fill, err := pool.DivideByFlowRate(pump)
	if err != nil {
		panic(err) // the flow rate is zero
	}

	fmt.Printf("%.1f\n", fill.ConvertTo(time.Hour))
	// Output: 5.0 h
}

func ExampleFlowRate_DivideByArea() {
	flow := units.NewFlowRate(30, flowrate.LiterPerMinute)
	pipe := units.NewArea(5, area.SquareCentimeter)

	speed, err := flow.DivideByArea(pipe)
	if err != nil {
		panic(err) // the area is zero
	}

	fmt.Printf("%.1f\n", speed.ConvertTo(velocity.MeterPerSecond))
	// Output: 1.0 m/s
}

func ExampleVelocity_MultiplyByArea() {
	wind := units.NewVelocity(10, velocity.MeterPerSecond)
	opening := units.NewArea(2, area.SquareMeter)

	flow := wind.MultiplyByArea(opening)

	fmt.Printf("%.0f\n", flow.ConvertTo(flowrate.CubicMeterPerHour))
	// Output: 72000 m³/h
}

func ExampleFlowRate_MultiplyByPressure() {
	flow := units.NewFlowRate(60, flowrate.LiterPerMinute)
	head := units.NewPressure(3, pressure.Bar)

	p := flow.MultiplyByPressure(head)

	fmt.Printf("%.0f\n", p.ConvertTo(power.Watt))
	// Output: 300 W
}

func ExamplePressure_MultiplyByFlowRate() {
	system := units.NewPressure(200, pressure.Bar)
	flow := units.NewFlowRate(20, flowrate.LiterPerMinute)

	p := system.MultiplyByFlowRate(flow)

	fmt.Printf("%.2f\n", p.ConvertTo(power.Kilowatt))
	// Output: 6.67 kW
}

func ExampleFlowRate_MultiplyByDensity() {
	flow := units.NewFlowRate(100, flowrate.LiterPerMinute)
	water := units.NewDensity(1, density.KilogramPerLiter)

	massFlow := flow.MultiplyByDensity(water)

	fmt.Printf("%.0f\n", massFlow.ConvertTo(massflowrate.KilogramPerHour))
	// Output: 6000 kg/h
}

func ExampleDensity_MultiplyByFlowRate() {
	air := units.NewDensity(1.2, density.KilogramPerCubicMeter)
	fan := units.NewFlowRate(1000, flowrate.CubicMeterPerHour)

	massFlow := air.MultiplyByFlowRate(fan)

	fmt.Printf("%.0f\n", massFlow.ConvertTo(massflowrate.KilogramPerHour))
	// Output: 1200 kg/h
}

func ExampleMassFlowRate_DivideByDensity() {
	burn := units.NewMassFlowRate(2400, massflowrate.KilogramPerHour)
	fuel := units.NewDensity(0.804, density.KilogramPerLiter)

	flow, err := burn.DivideByDensity(fuel)
	if err != nil {
		panic(err) // the density is zero
	}

	fmt.Printf("%.0f\n", flow.ConvertTo(flowrate.LiterPerHour))
	// Output: 2985 L/h
}

func ExampleMass_DivideByTime() {
	grain := units.NewMass(3, mass.Tonne)
	load := units.NewTime(15, time.Minute)

	rate, err := grain.DivideByTime(load)
	if err != nil {
		panic(err) // the time is zero
	}

	fmt.Printf("%.0f\n", rate.ConvertTo(massflowrate.TonnePerHour))
	// Output: 12 t/h
}

func ExampleMassFlowRate_MultiplyByTime() {
	burn := units.NewMassFlowRate(2400, massflowrate.KilogramPerHour)
	flight := units.NewTime(6, time.Hour)

	fuel := burn.MultiplyByTime(flight)

	fmt.Printf("%.1f\n", fuel.ConvertTo(mass.Tonne))
	// Output: 14.4 t
}

func ExampleMass_DivideByMassFlowRate() {
	fuel := units.NewMass(12, mass.Tonne)
	burn := units.NewMassFlowRate(2400, massflowrate.KilogramPerHour)

	endurance, err := fuel.DivideByMassFlowRate(burn)
	if err != nil {
		panic(err) // the mass flow rate is zero
	}

	fmt.Printf("%.1f\n", endurance.ConvertTo(time.Hour))
	// Output: 5.0 h
}

func ExamplePace_MultiplyByDistance() {
	p := units.NewPace(5, pace.MinutePerKilometer)
	marathon := units.NewDistance(42.195, distance.Kilometer)

	finish := p.MultiplyByDistance(marathon)

	fmt.Printf("%.2f\n", finish.ConvertTo(time.Hour))
	// Output: 3.52 h
}

func ExampleDistance_MultiplyByPace() {
	race := units.NewDistance(10, distance.Kilometer)
	p := units.NewPace(4.5, pace.MinutePerKilometer)

	finish := race.MultiplyByPace(p)

	fmt.Printf("%.0f\n", finish.ConvertTo(time.Minute))
	// Output: 45 min
}

func ExampleTime_DivideByDistance() {
	run := units.NewTime(25, time.Minute)
	route := units.NewDistance(5, distance.Kilometer)

	p, err := run.DivideByDistance(route)
	if err != nil {
		panic(err) // the distance is zero
	}

	fmt.Printf("%.1f\n", p.ConvertTo(pace.MinutePerKilometer))
	// Output: 5.0 min/km
}

func ExampleDistance_DivideByVolume() {
	trip := units.NewDistance(650, distance.Kilometer)
	fuel := units.NewVolume(40, volume.Liter)

	economy, err := trip.DivideByVolume(fuel)
	if err != nil {
		panic(err) // the volume is zero
	}

	fmt.Printf("%.2f\n", economy.ConvertTo(fueleconomy.KilometerPerLiter))
	// Output: 16.25 km/L
}

func ExampleFuelEconomy_MultiplyByVolume() {
	car := units.NewFuelEconomy(35, fueleconomy.MilePerGallon)
	tank := units.NewVolume(15, volume.Gallon)

	reach := car.MultiplyByVolume(tank)

	fmt.Printf("%.0f\n", reach.ConvertTo(distance.Mile))
	// Output: 525 mi
}

func ExampleDistance_DivideByFuelEconomy() {
	trip := units.NewDistance(300, distance.Kilometer)
	car := units.NewFuelEconomy(15, fueleconomy.KilometerPerLiter)

	fuel, err := trip.DivideByFuelEconomy(car)
	if err != nil {
		panic(err) // the fuel economy is zero
	}

	fmt.Printf("%.0f\n", fuel.ConvertTo(volume.Liter))
	// Output: 20 L
}

func ExampleFuelConsumption_MultiplyByDistance() {
	car := units.NewFuelConsumption(6.5, fuelconsumption.LiterPer100Kilometers)
	trip := units.NewDistance(500, distance.Kilometer)

	fuel := car.MultiplyByDistance(trip)

	fmt.Printf("%.1f\n", fuel.ConvertTo(volume.Liter))
	// Output: 32.5 L
}

func ExampleDistance_MultiplyByFuelConsumption() {
	trip := units.NewDistance(1000, distance.Kilometer)
	car := units.NewFuelConsumption(7, fuelconsumption.LiterPer100Kilometers)

	fuel := trip.MultiplyByFuelConsumption(car)

	fmt.Printf("%.0f\n", fuel.ConvertTo(volume.Liter))
	// Output: 70 L
}

func ExampleVolume_DivideByFuelConsumption() {
	tank := units.NewVolume(50, volume.Liter)
	car := units.NewFuelConsumption(5, fuelconsumption.LiterPer100Kilometers)

	reach, err := tank.DivideByFuelConsumption(car)
	if err != nil {
		panic(err) // the fuel consumption is zero
	}

	fmt.Printf("%.0f\n", reach.ConvertTo(distance.Kilometer))
	// Output: 1000 km
}

func ExampleAmount_MultiplyByMolarMass() {
	// Two moles of water.
	n := units.NewAmount(2, amount.Mole)
	water := units.NewMolarMass(18.015, molarmass.GramPerMole)

	m := n.MultiplyByMolarMass(water)

	fmt.Printf("%.2f\n", m.ConvertTo(mass.Gram))
	// Output: 36.03 g
}

func ExampleMolarMass_MultiplyByAmount() {
	// Half a mole of table salt, NaCl.
	salt := units.NewMolarMass(58.44, molarmass.GramPerMole)
	n := units.NewAmount(0.5, amount.Mole)

	m := salt.MultiplyByAmount(n)

	fmt.Printf("%.2f\n", m.ConvertTo(mass.Gram))
	// Output: 29.22 g
}

func ExampleMass_DivideByMolarMass() {
	// 100 g of glucose.
	sample := units.NewMass(100, mass.Gram)
	glucose := units.NewMolarMass(180.156, molarmass.GramPerMole)

	n, err := sample.DivideByMolarMass(glucose)
	if err != nil {
		panic(err) // the molar mass is zero
	}

	fmt.Printf("%.0f\n", n.ConvertTo(amount.Millimole))
	// Output: 555 mmol
}

func ExampleMass_DivideByAmount() {
	sample := units.NewMass(4.4, mass.Gram)
	n := units.NewAmount(0.1, amount.Mole)

	mm, err := sample.DivideByAmount(n)
	if err != nil {
		panic(err) // the amount of substance is zero
	}

	fmt.Printf("%.0f\n", mm.ConvertTo(molarmass.GramPerMole))
	// Output: 44 g/mol
}

func ExampleAmount_DivideByVolume() {
	n := units.NewAmount(0.5, amount.Mole)
	flask := units.NewVolume(2, volume.Liter)

	c, err := n.DivideByVolume(flask)
	if err != nil {
		panic(err) // the volume is zero
	}

	fmt.Printf("%.0f\n", c.ConvertTo(molarity.Millimolar))
	// Output: 250 mM
}

func ExampleMolarity_MultiplyByVolume() {
	c := units.NewMolarity(0.1, molarity.Molar)
	flask := units.NewVolume(250, volume.Milliliter)

	n := c.MultiplyByVolume(flask)

	fmt.Printf("%.0f\n", n.ConvertTo(amount.Millimole))
	// Output: 25 mmol
}

func ExampleAmount_DivideByMolarity() {
	n := units.NewAmount(0.01, amount.Mole)
	stock := units.NewMolarity(0.5, molarity.Molar)

	v, err := n.DivideByMolarity(stock)
	if err != nil {
		panic(err) // the molar concentration is zero
	}

	fmt.Printf("%.0f\n", v.ConvertTo(volume.Milliliter))
	// Output: 20 mL
}

func ExampleEquivalentDose_DivideByTime() {
	dose := units.NewEquivalentDose(50, equivalentdose.Microsievert)
	flight := units.NewTime(10, time.Hour)

	rate, err := dose.DivideByTime(flight)
	if err != nil {
		panic(err) // the time is zero
	}

	fmt.Printf("%.0f\n", rate.ConvertTo(doserate.MicrosievertPerHour))
	// Output: 5 μSv/h
}

func ExampleDoseRate_MultiplyByTime() {
	// The dose of a 10-hour flight at 5 μSv/h.
	cruise := units.NewDoseRate(5, doserate.MicrosievertPerHour)
	flight := units.NewTime(10, time.Hour)

	dose := cruise.MultiplyByTime(flight)

	fmt.Printf("%.0f\n", dose.ConvertTo(equivalentdose.Microsievert))
	// Output: 50 μSv
}

func ExampleEquivalentDose_DivideByDoseRate() {
	limit := units.NewEquivalentDose(1, equivalentdose.Millisievert)
	rate := units.NewDoseRate(0.5, doserate.MicrosievertPerHour)

	duration, err := limit.DivideByDoseRate(rate)
	if err != nil {
		panic(err) // the dose rate is zero
	}

	fmt.Printf("%.0f\n", duration.ConvertTo(time.Hour))
	// Output: 2000 h
}

func ExampleCapacitance_MultiplyByVoltage() {
	c := units.NewCapacitance(100, capacitance.Microfarad)
	v := units.NewVoltage(12, voltage.Volt)

	q := c.MultiplyByVoltage(v)

	fmt.Printf("%.1f\n", q.ConvertTo(charge.Millicoulomb))
	// Output: 1.2 mC
}

func ExampleVoltage_MultiplyByCapacitance() {
	v := units.NewVoltage(5, voltage.Volt)
	supercap := units.NewCapacitance(1, capacitance.Farad)

	q := v.MultiplyByCapacitance(supercap)

	fmt.Printf("%.0f\n", q.ConvertTo(charge.Coulomb))
	// Output: 5 C
}

func ExampleCharge_DivideByVoltage() {
	q := units.NewCharge(1, charge.Millicoulomb)
	v := units.NewVoltage(10, voltage.Volt)

	c, err := q.DivideByVoltage(v)
	if err != nil {
		panic(err) // the voltage is zero
	}

	fmt.Printf("%.0f\n", c.ConvertTo(capacitance.Microfarad))
	// Output: 100 μF
}

func ExampleCharge_DivideByCapacitance() {
	q := units.NewCharge(50, charge.Microcoulomb)
	c := units.NewCapacitance(10, capacitance.Microfarad)

	v, err := q.DivideByCapacitance(c)
	if err != nil {
		panic(err) // the capacitance is zero
	}

	fmt.Printf("%.0f\n", v.ConvertTo(voltage.Volt))
	// Output: 5 V
}

func ExampleResistance_MultiplyByCapacitance() {
	r := units.NewResistance(10, resistance.Kiloohm)
	c := units.NewCapacitance(100, capacitance.Microfarad)

	tau := r.MultiplyByCapacitance(c)

	fmt.Printf("%.1f\n", tau.ConvertTo(time.Second))
	// Output: 1.0 s
}

func ExampleCapacitance_MultiplyByResistance() {
	c := units.NewCapacitance(47, capacitance.Nanofarad)
	r := units.NewResistance(4.7, resistance.Kiloohm)

	tau := c.MultiplyByResistance(r)

	fmt.Printf("%.1f\n", tau.ConvertTo(time.Microsecond))
	// Output: 220.9 μs
}

func ExampleInductance_DivideByResistance() {
	l := units.NewInductance(10, inductance.Millihenry)
	r := units.NewResistance(100, resistance.Ohm)

	tau, err := l.DivideByResistance(r)
	if err != nil {
		panic(err) // the resistance is zero
	}

	fmt.Printf("%.0f\n", tau.ConvertTo(time.Microsecond))
	// Output: 100 μs
}

func ExampleLuminousFlux_DivideByArea() {
	bulb := units.NewLuminousFlux(800, luminousflux.Lumen)
	floor := units.NewArea(4, area.SquareMeter)

	e, err := bulb.DivideByArea(floor)
	if err != nil {
		panic(err) // the area is zero
	}

	fmt.Printf("%.0f\n", e.ConvertTo(illuminance.Lux))
	// Output: 200 lx
}

func ExampleIlluminance_MultiplyByArea() {
	target := units.NewIlluminance(500, illuminance.Lux)
	room := units.NewArea(10, area.SquareMeter)

	f := target.MultiplyByArea(room)

	fmt.Printf("%.1f\n", f.ConvertTo(luminousflux.Kilolumen))
	// Output: 5.0 klm
}

func ExampleEnergy_DivideByMass() {
	// A 60 kWh battery pack weighing 400 kg.
	pack := units.NewEnergy(60, energy.KilowattHour)
	weight := units.NewMass(400, mass.Kilogram)

	stored, err := pack.DivideByMass(weight)
	if err != nil {
		panic(err) // the mass is zero
	}

	fmt.Printf("%.0f\n", stored.ConvertTo(specificenergy.WattHourPerKilogram))
	// Output: 150 Wh/kg
}

func ExampleSpecificEnergy_MultiplyByMass() {
	// 50 kg of diesel.
	diesel := units.NewSpecificEnergy(45.5, specificenergy.MegajoulePerKilogram)
	fuel := units.NewMass(50, mass.Kilogram)

	e := diesel.MultiplyByMass(fuel)

	fmt.Printf("%.0f\n", e.ConvertTo(energy.KilowattHour))
	// Output: 632 kWh
}

func ExampleMass_MultiplyBySpecificEnergy() {
	// 1 kg of hydrogen.
	fuel := units.NewMass(1, mass.Kilogram)
	hydrogen := units.NewSpecificEnergy(120, specificenergy.MegajoulePerKilogram)

	e := fuel.MultiplyBySpecificEnergy(hydrogen)

	fmt.Printf("%.1f\n", e.ConvertTo(energy.KilowattHour))
	// Output: 33.3 kWh
}

func ExampleEnergy_DivideBySpecificEnergy() {
	need := units.NewEnergy(10, energy.KilowattHour)
	cells := units.NewSpecificEnergy(250, specificenergy.WattHourPerKilogram)

	m, err := need.DivideBySpecificEnergy(cells)
	if err != nil {
		panic(err) // the specific energy is zero
	}

	fmt.Printf("%.0f\n", m.ConvertTo(mass.Kilogram))
	// Output: 40 kg
}

func ExampleEnergy_DivideByVolume() {
	pack := units.NewEnergy(75, energy.KilowattHour)
	space := units.NewVolume(300, volume.Liter)

	u, err := pack.DivideByVolume(space)
	if err != nil {
		panic(err) // the volume is zero
	}

	fmt.Printf("%.0f\n", u.ConvertTo(energydensity.WattHourPerLiter))
	// Output: 250 Wh/L
}

func ExampleEnergyDensity_MultiplyByVolume() {
	// A 50 L tank of petrol.
	petrol := units.NewEnergyDensity(34.2, energydensity.MegajoulePerLiter)
	tank := units.NewVolume(50, volume.Liter)

	e := petrol.MultiplyByVolume(tank)

	fmt.Printf("%.0f\n", e.ConvertTo(energy.KilowattHour))
	// Output: 475 kWh
}

func ExampleEnergy_DivideByEnergyDensity() {
	need := units.NewEnergy(1, energy.KilowattHour)
	cells := units.NewEnergyDensity(700, energydensity.WattHourPerLiter)

	v, err := need.DivideByEnergyDensity(cells)
	if err != nil {
		panic(err) // the energy density is zero
	}

	fmt.Printf("%.2f\n", v.ConvertTo(volume.Liter))
	// Output: 1.43 L
}

func ExampleVelocity_DivideByFrequency() {
	// The wavelength of the note A in air.
	sound := units.NewVelocity(343, velocity.MeterPerSecond)
	note := units.NewFrequency(440, frequency.Hertz)

	wavelength, err := sound.DivideByFrequency(note)
	if err != nil {
		panic(err) // the frequency is zero
	}

	fmt.Printf("%.0f\n", wavelength.ConvertTo(distance.Centimeter))
	// Output: 78 cm
}

func ExampleVelocity_DivideByDistance() {
	// The 21 cm line of hydrogen, in radio astronomy.
	light := units.NewVelocity(1, velocity.SpeedOfLight)
	line := units.NewDistance(21.106, distance.Centimeter)

	f, err := light.DivideByDistance(line)
	if err != nil {
		panic(err) // the distance is zero
	}

	fmt.Printf("%.0f\n", f.ConvertTo(frequency.Megahertz))
	// Output: 1420 MHz
}

func ExampleFrequency_MultiplyByDistance() {
	wifi := units.NewFrequency(2.4, frequency.Gigahertz)
	wavelength := units.NewDistance(12.5, distance.Centimeter)

	speed := wifi.MultiplyByDistance(wavelength)

	fmt.Printf("%.0f\n", speed.ConvertTo(velocity.KilometerPerSecond))
	// Output: 300000 km/s
}
