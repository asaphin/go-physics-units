package units_test

import (
	"errors"
	"math"
	"testing"

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

// TestOperations checks each operation between measures added in v0.2.0
// against a value computed from the definitions of the units, in the result's
// base unit.
func TestOperations(t *testing.T) {
	cases := []struct {
		name  string
		op    func() (units.Measurement, error)
		value float64
		unit  string
	}{
		{
			"Distance.DivideByVelocity",
			func() (units.Measurement, error) {
				return fallible(
					units.NewDistance(42.195, distance.Kilometer).
						DivideByVelocity(units.NewVelocity(12, velocity.KilometerPerHour)),
				)
			},
			42195 / (12000.0 / 3600),
			time.BaseUnit.String(),
		},
		{
			"Energy.DivideByPower",
			func() (units.Measurement, error) {
				return fallible(units.NewEnergy(1, energy.KilowattHour).DivideByPower(units.NewPower(100, power.Watt)))
			},
			3.6e6 / 100,
			time.BaseUnit.String(),
		},
		{
			"Mass.MultiplyByVelocity",
			func() (units.Measurement, error) {
				return ok(
					units.NewMass(1000, mass.Kilogram).
						MultiplyByVelocity(units.NewVelocity(72, velocity.KilometerPerHour)),
				)
			},
			1000 * 20,
			momentum.BaseUnit.String(),
		},
		{
			"Velocity.MultiplyByMass",
			func() (units.Measurement, error) {
				return ok(units.NewVelocity(40, velocity.MeterPerSecond).MultiplyByMass(units.NewMass(145, mass.Gram)))
			},
			40 * 0.145,
			momentum.BaseUnit.String(),
		},
		{
			"Momentum.DivideByMass",
			func() (units.Measurement, error) {
				return fallible(
					units.NewMomentum(20000, momentum.NewtonSecond).DivideByMass(units.NewMass(1000, mass.Kilogram)),
				)
			},
			20000.0 / 1000,
			velocity.BaseUnit.String(),
		},
		{
			"Momentum.DivideByVelocity",
			func() (units.Measurement, error) {
				return fallible(
					units.NewMomentum(5.8, momentum.KilogramMeterPerSecond).
						DivideByVelocity(units.NewVelocity(40, velocity.MeterPerSecond)),
				)
			},
			5.8 / 40,
			mass.BaseUnit.String(),
		},
		{
			"Force.MultiplyByTime",
			func() (units.Measurement, error) {
				return ok(units.NewForce(500, force.Newton).MultiplyByTime(units.NewTime(4, time.Second)))
			},
			500 * 4,
			momentum.BaseUnit.String(),
		},
		{
			"Momentum.DivideByTime",
			func() (units.Measurement, error) {
				return fallible(
					units.NewMomentum(2000, momentum.NewtonSecond).DivideByTime(units.NewTime(4, time.Second)),
				)
			},
			2000.0 / 4,
			force.BaseUnit.String(),
		},
		{
			"Force.TorqueAt",
			func() (units.Measurement, error) {
				return ok(units.NewForce(200, force.Newton).TorqueAt(units.NewDistance(40, distance.Centimeter)))
			},
			200 * 0.4,
			torque.BaseUnit.String(),
		},
		{
			"Torque.DivideByDistance",
			func() (units.Measurement, error) {
				return fallible(
					units.NewTorque(80, torque.NewtonMeter).
						DivideByDistance(units.NewDistance(40, distance.Centimeter)),
				)
			},
			80 / 0.4,
			force.BaseUnit.String(),
		},
		{
			"Torque.MultiplyByAngularVelocity",
			func() (units.Measurement, error) {
				return ok(
					units.NewTorque(400, torque.NewtonMeter).
						MultiplyByAngularVelocity(units.NewAngularVelocity(3000, angularvelocity.RevolutionPerMinute)),
				)
			},
			400 * 3000 * 2 * math.Pi / 60,
			power.BaseUnit.String(),
		},
		{
			"AngularVelocity.MultiplyByTorque",
			func() (units.Measurement, error) {
				return ok(
					units.NewAngularVelocity(3000, angularvelocity.RevolutionPerMinute).
						MultiplyByTorque(units.NewTorque(400, torque.NewtonMeter)),
				)
			},
			3000 * 2 * math.Pi / 60 * 400,
			power.BaseUnit.String(),
		},
		{
			"Power.DivideByAngularVelocity",
			func() (units.Measurement, error) {
				return fallible(
					units.NewPower(300, power.Horsepower).
						DivideByAngularVelocity(units.NewAngularVelocity(6000, angularvelocity.RevolutionPerMinute)),
				)
			},
			300 * 745.69987158227022 / (6000 * 2 * math.Pi / 60),
			torque.BaseUnit.String(),
		},
		{
			"Power.DivideByTorque",
			func() (units.Measurement, error) {
				return fallible(
					units.NewPower(100, power.Kilowatt).DivideByTorque(units.NewTorque(250, torque.NewtonMeter)),
				)
			},
			100e3 / 250,
			angularvelocity.BaseUnit.String(),
		},
		{
			"Torque.MultiplyByAngle",
			func() (units.Measurement, error) {
				return ok(units.NewTorque(50, torque.NewtonMeter).MultiplyByAngle(units.NewAngle(10, angle.Revolution)))
			},
			50 * 10 * 2 * math.Pi,
			energy.BaseUnit.String(),
		},
		{
			"Volume.DivideByTime",
			func() (units.Measurement, error) {
				return fallible(units.NewVolume(200, volume.Liter).DivideByTime(units.NewTime(8, time.Minute)))
			},
			0.2 / 480,
			flowrate.BaseUnit.String(),
		},
		{
			"FlowRate.MultiplyByTime",
			func() (units.Measurement, error) {
				return ok(units.NewFlowRate(12, flowrate.LiterPerMinute).MultiplyByTime(units.NewTime(10, time.Minute)))
			},
			12e-3 / 60 * 600,
			volume.BaseUnit.String(),
		},
		{
			"Volume.DivideByFlowRate",
			func() (units.Measurement, error) {
				return fallible(
					units.NewVolume(50, volume.CubicMeter).
						DivideByFlowRate(units.NewFlowRate(10, flowrate.CubicMeterPerHour)),
				)
			},
			50 / (10.0 / 3600),
			time.BaseUnit.String(),
		},
		{
			"FlowRate.DivideByArea",
			func() (units.Measurement, error) {
				return fallible(
					units.NewFlowRate(30, flowrate.LiterPerMinute).
						DivideByArea(units.NewArea(5, area.SquareCentimeter)),
				)
			},
			30e-3 / 60 / 5e-4,
			velocity.BaseUnit.String(),
		},
		{
			"Velocity.MultiplyByArea",
			func() (units.Measurement, error) {
				return ok(
					units.NewVelocity(10, velocity.MeterPerSecond).MultiplyByArea(units.NewArea(2, area.SquareMeter)),
				)
			},
			10 * 2,
			flowrate.BaseUnit.String(),
		},
		{
			"FlowRate.MultiplyByPressure",
			func() (units.Measurement, error) {
				return ok(
					units.NewFlowRate(60, flowrate.LiterPerMinute).
						MultiplyByPressure(units.NewPressure(3, pressure.Bar)),
				)
			},
			60e-3 / 60 * 3e5,
			power.BaseUnit.String(),
		},
		{
			"Pressure.MultiplyByFlowRate",
			func() (units.Measurement, error) {
				return ok(
					units.NewPressure(200, pressure.Bar).
						MultiplyByFlowRate(units.NewFlowRate(20, flowrate.LiterPerMinute)),
				)
			},
			200e5 * 20e-3 / 60,
			power.BaseUnit.String(),
		},
		{
			"FlowRate.MultiplyByDensity",
			func() (units.Measurement, error) {
				return ok(
					units.NewFlowRate(100, flowrate.LiterPerMinute).
						MultiplyByDensity(units.NewDensity(1, density.KilogramPerLiter)),
				)
			},
			100e-3 / 60 * 1000,
			massflowrate.BaseUnit.String(),
		},
		{
			"Density.MultiplyByFlowRate",
			func() (units.Measurement, error) {
				return ok(
					units.NewDensity(1.2, density.KilogramPerCubicMeter).
						MultiplyByFlowRate(units.NewFlowRate(1000, flowrate.CubicMeterPerHour)),
				)
			},
			1.2 * 1000.0 / 3600,
			massflowrate.BaseUnit.String(),
		},
		{
			"MassFlowRate.DivideByDensity",
			func() (units.Measurement, error) {
				return fallible(
					units.NewMassFlowRate(2400, massflowrate.KilogramPerHour).
						DivideByDensity(units.NewDensity(0.804, density.KilogramPerLiter)),
				)
			},
			2400.0 / 3600 / 804,
			flowrate.BaseUnit.String(),
		},
		{
			"Mass.DivideByTime",
			func() (units.Measurement, error) {
				return fallible(units.NewMass(3, mass.Tonne).DivideByTime(units.NewTime(15, time.Minute)))
			},
			3000.0 / 900,
			massflowrate.BaseUnit.String(),
		},
		{
			"MassFlowRate.MultiplyByTime",
			func() (units.Measurement, error) {
				return ok(
					units.NewMassFlowRate(2400, massflowrate.KilogramPerHour).
						MultiplyByTime(units.NewTime(6, time.Hour)),
				)
			},
			2400.0 / 3600 * 6 * 3600,
			mass.BaseUnit.String(),
		},
		{
			"Mass.DivideByMassFlowRate",
			func() (units.Measurement, error) {
				return fallible(
					units.NewMass(12, mass.Tonne).
						DivideByMassFlowRate(units.NewMassFlowRate(2400, massflowrate.KilogramPerHour)),
				)
			},
			12000 / (2400.0 / 3600),
			time.BaseUnit.String(),
		},
		{
			"Pace.MultiplyByDistance",
			func() (units.Measurement, error) {
				return ok(
					units.NewPace(5, pace.MinutePerKilometer).
						MultiplyByDistance(units.NewDistance(42.195, distance.Kilometer)),
				)
			},
			0.3 * 42195,
			time.BaseUnit.String(),
		},
		{
			"Distance.MultiplyByPace",
			func() (units.Measurement, error) {
				return ok(
					units.NewDistance(10, distance.Kilometer).
						MultiplyByPace(units.NewPace(4.5, pace.MinutePerKilometer)),
				)
			},
			10000 * 0.27,
			time.BaseUnit.String(),
		},
		{
			"Time.DivideByDistance",
			func() (units.Measurement, error) {
				return fallible(
					units.NewTime(25, time.Minute).DivideByDistance(units.NewDistance(5, distance.Kilometer)),
				)
			},
			1500.0 / 5000,
			pace.BaseUnit.String(),
		},
		{
			"Distance.DivideByVolume",
			func() (units.Measurement, error) {
				return fallible(
					units.NewDistance(650, distance.Kilometer).DivideByVolume(units.NewVolume(40, volume.Liter)),
				)
			},
			650000 / 0.04,
			fueleconomy.BaseUnit.String(),
		},
		{
			"FuelEconomy.MultiplyByVolume",
			func() (units.Measurement, error) {
				return ok(
					units.NewFuelEconomy(35, fueleconomy.MilePerGallon).
						MultiplyByVolume(units.NewVolume(15, volume.Gallon)),
				)
			},
			35 * 15 * 1609.344,
			distance.BaseUnit.String(),
		},
		{
			"Distance.DivideByFuelEconomy",
			func() (units.Measurement, error) {
				return fallible(
					units.NewDistance(300, distance.Kilometer).
						DivideByFuelEconomy(units.NewFuelEconomy(15, fueleconomy.KilometerPerLiter)),
				)
			},
			300000 / 1.5e7,
			volume.BaseUnit.String(),
		},
		{
			"FuelConsumption.MultiplyByDistance",
			func() (units.Measurement, error) {
				return ok(
					units.NewFuelConsumption(6.5, fuelconsumption.LiterPer100Kilometers).
						MultiplyByDistance(units.NewDistance(500, distance.Kilometer)),
				)
			},
			6.5e-8 * 500000,
			volume.BaseUnit.String(),
		},
		{
			"Distance.MultiplyByFuelConsumption",
			func() (units.Measurement, error) {
				return ok(
					units.NewDistance(1000, distance.Kilometer).
						MultiplyByFuelConsumption(units.NewFuelConsumption(7, fuelconsumption.LiterPer100Kilometers)),
				)
			},
			1e6 * 7e-8,
			volume.BaseUnit.String(),
		},
		{
			"Volume.DivideByFuelConsumption",
			func() (units.Measurement, error) {
				return fallible(
					units.NewVolume(50, volume.Liter).
						DivideByFuelConsumption(units.NewFuelConsumption(5, fuelconsumption.LiterPer100Kilometers)),
				)
			},
			0.05 / 5e-8,
			distance.BaseUnit.String(),
		},
		{
			"Amount.MultiplyByMolarMass",
			func() (units.Measurement, error) {
				return ok(
					units.NewAmount(2, amount.Mole).
						MultiplyByMolarMass(units.NewMolarMass(18.015, molarmass.GramPerMole)),
				)
			},
			2 * 0.018015,
			mass.BaseUnit.String(),
		},
		{
			"MolarMass.MultiplyByAmount",
			func() (units.Measurement, error) {
				return ok(
					units.NewMolarMass(58.44, molarmass.GramPerMole).
						MultiplyByAmount(units.NewAmount(0.5, amount.Mole)),
				)
			},
			0.05844 * 0.5,
			mass.BaseUnit.String(),
		},
		{
			"Mass.DivideByMolarMass",
			func() (units.Measurement, error) {
				return fallible(
					units.NewMass(100, mass.Gram).DivideByMolarMass(units.NewMolarMass(180.156, molarmass.GramPerMole)),
				)
			},
			0.1 / 0.180156,
			amount.BaseUnit.String(),
		},
		{
			"Mass.DivideByAmount",
			func() (units.Measurement, error) {
				return fallible(units.NewMass(4.4, mass.Gram).DivideByAmount(units.NewAmount(0.1, amount.Mole)))
			},
			0.0044 / 0.1,
			molarmass.BaseUnit.String(),
		},
		{
			"Amount.DivideByVolume",
			func() (units.Measurement, error) {
				return fallible(units.NewAmount(0.5, amount.Mole).DivideByVolume(units.NewVolume(2, volume.Liter)))
			},
			0.5 / 0.002,
			molarity.BaseUnit.String(),
		},
		{
			"Molarity.MultiplyByVolume",
			func() (units.Measurement, error) {
				return ok(
					units.NewMolarity(0.1, molarity.Molar).MultiplyByVolume(units.NewVolume(250, volume.Milliliter)),
				)
			},
			100 * 0.00025,
			amount.BaseUnit.String(),
		},
		{
			"Amount.DivideByMolarity",
			func() (units.Measurement, error) {
				return fallible(
					units.NewAmount(0.01, amount.Mole).DivideByMolarity(units.NewMolarity(0.5, molarity.Molar)),
				)
			},
			0.01 / 500,
			volume.BaseUnit.String(),
		},
		{
			"EquivalentDose.DivideByTime",
			func() (units.Measurement, error) {
				return fallible(
					units.NewEquivalentDose(50, equivalentdose.Microsievert).DivideByTime(units.NewTime(10, time.Hour)),
				)
			},
			50e-6 / 36000,
			doserate.BaseUnit.String(),
		},
		{
			"DoseRate.MultiplyByTime",
			func() (units.Measurement, error) {
				return ok(
					units.NewDoseRate(5, doserate.MicrosievertPerHour).MultiplyByTime(units.NewTime(10, time.Hour)),
				)
			},
			5e-6 / 3600 * 36000,
			equivalentdose.BaseUnit.String(),
		},
		{
			"EquivalentDose.DivideByDoseRate",
			func() (units.Measurement, error) {
				return fallible(
					units.NewEquivalentDose(1, equivalentdose.Millisievert).
						DivideByDoseRate(units.NewDoseRate(0.5, doserate.MicrosievertPerHour)),
				)
			},
			1e-3 / (0.5e-6 / 3600),
			time.BaseUnit.String(),
		},
		{
			"Capacitance.MultiplyByVoltage",
			func() (units.Measurement, error) {
				return ok(
					units.NewCapacitance(100, capacitance.Microfarad).
						MultiplyByVoltage(units.NewVoltage(12, voltage.Volt)),
				)
			},
			100e-6 * 12,
			charge.BaseUnit.String(),
		},
		{
			"Voltage.MultiplyByCapacitance",
			func() (units.Measurement, error) {
				return ok(
					units.NewVoltage(5, voltage.Volt).MultiplyByCapacitance(units.NewCapacitance(1, capacitance.Farad)),
				)
			},
			5 * 1,
			charge.BaseUnit.String(),
		},
		{
			"Charge.DivideByVoltage",
			func() (units.Measurement, error) {
				return fallible(
					units.NewCharge(1, charge.Millicoulomb).DivideByVoltage(units.NewVoltage(10, voltage.Volt)),
				)
			},
			1e-3 / 10,
			capacitance.BaseUnit.String(),
		},
		{
			"Charge.DivideByCapacitance",
			func() (units.Measurement, error) {
				return fallible(
					units.NewCharge(50, charge.Microcoulomb).
						DivideByCapacitance(units.NewCapacitance(10, capacitance.Microfarad)),
				)
			},
			50e-6 / 10e-6,
			voltage.BaseUnit.String(),
		},
		{
			"Resistance.MultiplyByCapacitance",
			func() (units.Measurement, error) {
				return ok(
					units.NewResistance(10, resistance.Kiloohm).
						MultiplyByCapacitance(units.NewCapacitance(100, capacitance.Microfarad)),
				)
			},
			1e4 * 1e-4,
			time.BaseUnit.String(),
		},
		{
			"Capacitance.MultiplyByResistance",
			func() (units.Measurement, error) {
				return ok(
					units.NewCapacitance(47, capacitance.Nanofarad).
						MultiplyByResistance(units.NewResistance(4.7, resistance.Kiloohm)),
				)
			},
			47e-9 * 4700,
			time.BaseUnit.String(),
		},
		{
			"Inductance.DivideByResistance",
			func() (units.Measurement, error) {
				return fallible(
					units.NewInductance(10, inductance.Millihenry).
						DivideByResistance(units.NewResistance(100, resistance.Ohm)),
				)
			},
			0.01 / 100,
			time.BaseUnit.String(),
		},
		{
			"LuminousFlux.DivideByArea",
			func() (units.Measurement, error) {
				return fallible(
					units.NewLuminousFlux(800, luminousflux.Lumen).DivideByArea(units.NewArea(4, area.SquareMeter)),
				)
			},
			800.0 / 4,
			illuminance.BaseUnit.String(),
		},
		{
			"Illuminance.MultiplyByArea",
			func() (units.Measurement, error) {
				return ok(
					units.NewIlluminance(500, illuminance.Lux).MultiplyByArea(units.NewArea(10, area.SquareMeter)),
				)
			},
			500 * 10,
			luminousflux.BaseUnit.String(),
		},
		{
			"Energy.DivideByMass",
			func() (units.Measurement, error) {
				return fallible(
					units.NewEnergy(60, energy.KilowattHour).DivideByMass(units.NewMass(400, mass.Kilogram)),
				)
			},
			60 * 3.6e6 / 400,
			specificenergy.BaseUnit.String(),
		},
		{
			"SpecificEnergy.MultiplyByMass",
			func() (units.Measurement, error) {
				return ok(
					units.NewSpecificEnergy(45.5, specificenergy.MegajoulePerKilogram).
						MultiplyByMass(units.NewMass(50, mass.Kilogram)),
				)
			},
			45.5e6 * 50,
			energy.BaseUnit.String(),
		},
		{
			"Mass.MultiplyBySpecificEnergy",
			func() (units.Measurement, error) {
				return ok(
					units.NewMass(1, mass.Kilogram).
						MultiplyBySpecificEnergy(units.NewSpecificEnergy(120, specificenergy.MegajoulePerKilogram)),
				)
			},
			1 * 120e6,
			energy.BaseUnit.String(),
		},
		{
			"Energy.DivideBySpecificEnergy",
			func() (units.Measurement, error) {
				return fallible(
					units.NewEnergy(10, energy.KilowattHour).
						DivideBySpecificEnergy(units.NewSpecificEnergy(250, specificenergy.WattHourPerKilogram)),
				)
			},
			10 * 3.6e6 / (250 * 3600),
			mass.BaseUnit.String(),
		},
		{
			"Energy.DivideByVolume",
			func() (units.Measurement, error) {
				return fallible(
					units.NewEnergy(75, energy.KilowattHour).DivideByVolume(units.NewVolume(300, volume.Liter)),
				)
			},
			75 * 3.6e6 / 0.3,
			energydensity.BaseUnit.String(),
		},
		{
			"EnergyDensity.MultiplyByVolume",
			func() (units.Measurement, error) {
				return ok(
					units.NewEnergyDensity(34.2, energydensity.MegajoulePerLiter).
						MultiplyByVolume(units.NewVolume(50, volume.Liter)),
				)
			},
			34.2e9 * 0.05,
			energy.BaseUnit.String(),
		},
		{
			"Energy.DivideByEnergyDensity",
			func() (units.Measurement, error) {
				return fallible(
					units.NewEnergy(1, energy.KilowattHour).
						DivideByEnergyDensity(units.NewEnergyDensity(700, energydensity.WattHourPerLiter)),
				)
			},
			3.6e6 / (700 * 3.6e6),
			volume.BaseUnit.String(),
		},
		{
			"Velocity.DivideByFrequency",
			func() (units.Measurement, error) {
				return fallible(
					units.NewVelocity(343, velocity.MeterPerSecond).
						DivideByFrequency(units.NewFrequency(440, frequency.Hertz)),
				)
			},
			343.0 / 440,
			distance.BaseUnit.String(),
		},
		{
			"Velocity.DivideByDistance",
			func() (units.Measurement, error) {
				return fallible(
					units.NewVelocity(1, velocity.SpeedOfLight).
						DivideByDistance(units.NewDistance(21.106, distance.Centimeter)),
				)
			},
			299792458 / 0.21106,
			frequency.BaseUnit.String(),
		},
		{
			"Frequency.MultiplyByDistance",
			func() (units.Measurement, error) {
				return ok(
					units.NewFrequency(2.4, frequency.Gigahertz).
						MultiplyByDistance(units.NewDistance(12.5, distance.Centimeter)),
				)
			},
			2.4e9 * 0.125,
			velocity.BaseUnit.String(),
		},
	}

	for _, c := range cases {
		m, err := c.op()
		if err != nil {
			t.Fatalf("%s: %v", c.name, err)
		}

		assertMeasurement(t, m, c.value, c.unit)
	}

	_ = math.Pi
}

// zeroDivision checks that op fails with ErrZeroDivision and returns the zero
// value of its result type.
func zeroDivision[T comparable](t *testing.T, name string, op func() (T, error)) {
	t.Helper()

	var zero T
	if got, err := op(); !errors.Is(err, units.ErrZeroDivision) || got != zero {
		t.Errorf("%s by zero = %v, %v; want the zero value and ErrZeroDivision", name, got, err)
	}
}

func TestOperations_ZeroDivision(t *testing.T) {
	zeroDivision(t, "Distance.DivideByVelocity", func() (units.Time, error) {
		return units.NewDistance(42.195, distance.Kilometer).DivideByVelocity(units.NewVelocity(0, velocity.BaseUnit))
	})
	zeroDivision(t, "Energy.DivideByPower", func() (units.Time, error) {
		return units.NewEnergy(1, energy.KilowattHour).DivideByPower(units.NewPower(0, power.BaseUnit))
	})
	zeroDivision(t, "Momentum.DivideByMass", func() (units.Velocity, error) {
		return units.NewMomentum(20000, momentum.NewtonSecond).DivideByMass(units.NewMass(0, mass.BaseUnit))
	})
	zeroDivision(t, "Momentum.DivideByVelocity", func() (units.Mass, error) {
		return units.NewMomentum(5.8, momentum.KilogramMeterPerSecond).
			DivideByVelocity(units.NewVelocity(0, velocity.BaseUnit))
	})
	zeroDivision(t, "Momentum.DivideByTime", func() (units.Force, error) {
		return units.NewMomentum(2000, momentum.NewtonSecond).DivideByTime(units.NewTime(0, time.BaseUnit))
	})
	zeroDivision(t, "Torque.DivideByDistance", func() (units.Force, error) {
		return units.NewTorque(80, torque.NewtonMeter).DivideByDistance(units.NewDistance(0, distance.BaseUnit))
	})
	zeroDivision(t, "Power.DivideByAngularVelocity", func() (units.Torque, error) {
		return units.NewPower(300, power.Horsepower).
			DivideByAngularVelocity(units.NewAngularVelocity(0, angularvelocity.BaseUnit))
	})
	zeroDivision(t, "Power.DivideByTorque", func() (units.AngularVelocity, error) {
		return units.NewPower(100, power.Kilowatt).DivideByTorque(units.NewTorque(0, torque.BaseUnit))
	})
	zeroDivision(t, "Volume.DivideByTime", func() (units.FlowRate, error) {
		return units.NewVolume(200, volume.Liter).DivideByTime(units.NewTime(0, time.BaseUnit))
	})
	zeroDivision(t, "Volume.DivideByFlowRate", func() (units.Time, error) {
		return units.NewVolume(50, volume.CubicMeter).DivideByFlowRate(units.NewFlowRate(0, flowrate.BaseUnit))
	})
	zeroDivision(t, "FlowRate.DivideByArea", func() (units.Velocity, error) {
		return units.NewFlowRate(30, flowrate.LiterPerMinute).DivideByArea(units.NewArea(0, area.BaseUnit))
	})
	zeroDivision(t, "MassFlowRate.DivideByDensity", func() (units.FlowRate, error) {
		return units.NewMassFlowRate(2400, massflowrate.KilogramPerHour).
			DivideByDensity(units.NewDensity(0, density.BaseUnit))
	})
	zeroDivision(t, "Mass.DivideByTime", func() (units.MassFlowRate, error) {
		return units.NewMass(3, mass.Tonne).DivideByTime(units.NewTime(0, time.BaseUnit))
	})
	zeroDivision(t, "Mass.DivideByMassFlowRate", func() (units.Time, error) {
		return units.NewMass(12, mass.Tonne).DivideByMassFlowRate(units.NewMassFlowRate(0, massflowrate.BaseUnit))
	})
	zeroDivision(t, "Time.DivideByDistance", func() (units.Pace, error) {
		return units.NewTime(25, time.Minute).DivideByDistance(units.NewDistance(0, distance.BaseUnit))
	})
	zeroDivision(t, "Distance.DivideByVolume", func() (units.FuelEconomy, error) {
		return units.NewDistance(650, distance.Kilometer).DivideByVolume(units.NewVolume(0, volume.BaseUnit))
	})
	zeroDivision(t, "Distance.DivideByFuelEconomy", func() (units.Volume, error) {
		return units.NewDistance(300, distance.Kilometer).
			DivideByFuelEconomy(units.NewFuelEconomy(0, fueleconomy.BaseUnit))
	})
	zeroDivision(t, "Volume.DivideByFuelConsumption", func() (units.Distance, error) {
		return units.NewVolume(50, volume.Liter).
			DivideByFuelConsumption(units.NewFuelConsumption(0, fuelconsumption.BaseUnit))
	})
	zeroDivision(t, "Mass.DivideByMolarMass", func() (units.Amount, error) {
		return units.NewMass(100, mass.Gram).DivideByMolarMass(units.NewMolarMass(0, molarmass.BaseUnit))
	})
	zeroDivision(t, "Mass.DivideByAmount", func() (units.MolarMass, error) {
		return units.NewMass(4.4, mass.Gram).DivideByAmount(units.NewAmount(0, amount.BaseUnit))
	})
	zeroDivision(t, "Amount.DivideByVolume", func() (units.Molarity, error) {
		return units.NewAmount(0.5, amount.Mole).DivideByVolume(units.NewVolume(0, volume.BaseUnit))
	})
	zeroDivision(t, "Amount.DivideByMolarity", func() (units.Volume, error) {
		return units.NewAmount(0.01, amount.Mole).DivideByMolarity(units.NewMolarity(0, molarity.BaseUnit))
	})
	zeroDivision(t, "EquivalentDose.DivideByTime", func() (units.DoseRate, error) {
		return units.NewEquivalentDose(50, equivalentdose.Microsievert).DivideByTime(units.NewTime(0, time.BaseUnit))
	})
	zeroDivision(t, "EquivalentDose.DivideByDoseRate", func() (units.Time, error) {
		return units.NewEquivalentDose(1, equivalentdose.Millisievert).
			DivideByDoseRate(units.NewDoseRate(0, doserate.BaseUnit))
	})
	zeroDivision(t, "Charge.DivideByVoltage", func() (units.Capacitance, error) {
		return units.NewCharge(1, charge.Millicoulomb).DivideByVoltage(units.NewVoltage(0, voltage.BaseUnit))
	})
	zeroDivision(t, "Charge.DivideByCapacitance", func() (units.Voltage, error) {
		return units.NewCharge(50, charge.Microcoulomb).
			DivideByCapacitance(units.NewCapacitance(0, capacitance.BaseUnit))
	})
	zeroDivision(t, "Inductance.DivideByResistance", func() (units.Time, error) {
		return units.NewInductance(10, inductance.Millihenry).
			DivideByResistance(units.NewResistance(0, resistance.BaseUnit))
	})
	zeroDivision(t, "LuminousFlux.DivideByArea", func() (units.Illuminance, error) {
		return units.NewLuminousFlux(800, luminousflux.Lumen).DivideByArea(units.NewArea(0, area.BaseUnit))
	})
	zeroDivision(t, "Energy.DivideByMass", func() (units.SpecificEnergy, error) {
		return units.NewEnergy(60, energy.KilowattHour).DivideByMass(units.NewMass(0, mass.BaseUnit))
	})
	zeroDivision(t, "Energy.DivideBySpecificEnergy", func() (units.Mass, error) {
		return units.NewEnergy(10, energy.KilowattHour).
			DivideBySpecificEnergy(units.NewSpecificEnergy(0, specificenergy.BaseUnit))
	})
	zeroDivision(t, "Energy.DivideByVolume", func() (units.EnergyDensity, error) {
		return units.NewEnergy(75, energy.KilowattHour).DivideByVolume(units.NewVolume(0, volume.BaseUnit))
	})
	zeroDivision(t, "Energy.DivideByEnergyDensity", func() (units.Volume, error) {
		return units.NewEnergy(1, energy.KilowattHour).
			DivideByEnergyDensity(units.NewEnergyDensity(0, energydensity.BaseUnit))
	})
	zeroDivision(t, "Velocity.DivideByFrequency", func() (units.Distance, error) {
		return units.NewVelocity(343, velocity.MeterPerSecond).
			DivideByFrequency(units.NewFrequency(0, frequency.BaseUnit))
	})
	zeroDivision(t, "Velocity.DivideByDistance", func() (units.Frequency, error) {
		return units.NewVelocity(1, velocity.SpeedOfLight).DivideByDistance(units.NewDistance(0, distance.BaseUnit))
	})
}
