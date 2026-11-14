package units_test

import (
	"fmt"

	units "github.com/asaphin/go-physics-units"
	"github.com/asaphin/go-physics-units/absorbeddose"
	"github.com/asaphin/go-physics-units/area"
	"github.com/asaphin/go-physics-units/density"
	"github.com/asaphin/go-physics-units/distance"
	"github.com/asaphin/go-physics-units/energy"
	"github.com/asaphin/go-physics-units/equivalentdose"
	"github.com/asaphin/go-physics-units/frequency"
	"github.com/asaphin/go-physics-units/fuelconsumption"
	"github.com/asaphin/go-physics-units/fueleconomy"
	"github.com/asaphin/go-physics-units/mass"
	"github.com/asaphin/go-physics-units/molarmass"
	"github.com/asaphin/go-physics-units/pace"
	"github.com/asaphin/go-physics-units/radioactivity"
	"github.com/asaphin/go-physics-units/ratio"
	"github.com/asaphin/go-physics-units/time"
	"github.com/asaphin/go-physics-units/velocity"
	"github.com/asaphin/go-physics-units/volume"
)

func ExampleDistance_CircleArea() {
	// A round table with a radius of 60 cm.
	table := units.NewDistance(60, distance.Centimeter)

	fmt.Printf("%.2f\n", table.CircleArea())
	// Output: 1.13 m²
}

func ExampleDistance_SphereArea() {
	// The surface of the Earth, from its mean radius.
	earth := units.NewDistance(6371, distance.Kilometer)

	fmt.Printf("%.3e\n", earth.SphereArea().ConvertTo(area.SquareKilometer))
	// Output: 5.101e+08 km²
}

func ExampleDistance_SphereVolume() {
	// A football with a radius of 11 cm.
	ball := units.NewDistance(11, distance.Centimeter)

	fmt.Printf("%.2f\n", ball.SphereVolume().ConvertTo(volume.Liter))
	// Output: 5.58 L
}

func ExampleArea_SquareSide() {
	field := units.NewArea(1, area.Hectare)

	side, err := field.SquareSide()
	if err != nil {
		panic(err) // the area is negative
	}

	fmt.Printf("%.0f\n", side)
	// Output: 100 m
}

func ExampleArea_CircleRadius() {
	// A pizza of 1000 cm².
	pizza := units.NewArea(1000, area.SquareCentimeter)

	r, err := pizza.CircleRadius()
	if err != nil {
		panic(err) // the area is negative
	}

	fmt.Printf("%.1f\n", r.ConvertTo(distance.Centimeter))
	// Output: 17.8 cm
}

func ExampleVolume_CubeEdge() {
	// All the gold ever mined, about 216 000 t, would fit in one cube.
	gold, err := units.NewMass(216000, mass.Tonne).
		DivideByDensity(units.NewDensity(19.32, density.GramPerCubicCentimeter))
	if err != nil {
		panic(err) // the density is zero
	}

	fmt.Printf("%.1f\n", gold.CubeEdge())
	// Output: 22.4 m
}

func ExampleVolume_SphereRadius() {
	// The equivalent diameter of a 0.05 mL raindrop.
	drop := units.NewVolume(0.05, volume.Milliliter)

	fmt.Printf("%.2f\n", drop.SphereRadius().Mul(2).ConvertTo(distance.Millimeter))
	// Output: 4.57 mm
}

func ExampleMass_KineticEnergy() {
	car := units.NewMass(1500, mass.Kilogram)
	speed := units.NewVelocity(100, velocity.KilometerPerHour)

	fmt.Printf("%.0f\n", car.KineticEnergy(speed).ConvertTo(energy.Kilojoule))
	// Output: 579 kJ
}

func ExampleMass_Weight() {
	climber := units.NewMass(75, mass.Kilogram)

	fmt.Printf("%.1f\n", climber.Weight())
	// Output: 735.5 N
}

func ExampleMass_RestEnergy() {
	// The energy equivalent of one gram of matter, as TNT.
	gram := units.NewMass(1, mass.Gram)

	fmt.Printf("%.1f\n", gram.RestEnergy().ConvertTo(energy.KilotonOfTNT))
	// Output: 21.5 ktTNT
}

func ExampleEnergy_MassEquivalent() {
	// The mass converted into energy by a 15 kt explosion.
	blast := units.NewEnergy(15, energy.KilotonOfTNT)

	fmt.Printf("%.2f\n", blast.MassEquivalent().ConvertTo(mass.Gram))
	// Output: 0.70 g
}

func ExampleEnergy_PhotonFrequency() {
	photon := units.NewEnergy(2, energy.Electronvolt)

	fmt.Printf("%.0f\n", photon.PhotonFrequency().ConvertTo(frequency.Terahertz))
	// Output: 484 THz
}

func ExampleFrequency_PhotonEnergy() {
	// Green light.
	green := units.NewFrequency(540, frequency.Terahertz)

	fmt.Printf("%.2f\n", green.PhotonEnergy().ConvertTo(energy.Electronvolt))
	// Output: 2.23 eV
}

func ExampleVelocity_Pace() {
	speed := units.NewVelocity(12, velocity.KilometerPerHour)

	p, err := speed.Pace()
	if err != nil {
		panic(err) // the speed is zero
	}

	fmt.Printf("%.1f\n", p.ConvertTo(pace.MinutePerKilometer))
	// Output: 5.0 min/km
}

func ExampleVelocity_LorentzFactor() {
	// Cosmic-ray muons reach the ground because their clocks run slow.
	muon := units.NewVelocity(0.998, velocity.SpeedOfLight)

	gamma, err := muon.LorentzFactor()
	if err != nil {
		panic(err) // not below the speed of light
	}

	fmt.Printf("%.1f\n", gamma)
	// Output: 15.8
}

func ExampleTime_Frequency() {
	// A heartbeat every 0.8 s.
	beat := units.NewTime(0.8, time.Second)

	rate, err := beat.Frequency()
	if err != nil {
		panic(err) // the time is zero
	}

	fmt.Printf("%.0f\n", rate.ConvertTo(frequency.BeatPerMinute))
	// Output: 75 bpm
}

func ExamplePace_Velocity() {
	p := units.NewPace(4, pace.MinutePerKilometer)

	speed, err := p.Velocity()
	if err != nil {
		panic(err) // the pace is zero
	}

	fmt.Printf("%.1f\n", speed.ConvertTo(velocity.KilometerPerHour))
	// Output: 15.0 km/h
}

func ExampleFuelEconomy_Consumption() {
	car := units.NewFuelEconomy(30, fueleconomy.MilePerGallon)

	c, err := car.Consumption()
	if err != nil {
		panic(err) // the fuel economy is zero
	}

	fmt.Printf("%.1f\n", c.ConvertTo(fuelconsumption.LiterPer100Kilometers))
	// Output: 7.8 L/100km
}

func ExampleFuelConsumption_Economy() {
	car := units.NewFuelConsumption(5, fuelconsumption.LiterPer100Kilometers)

	e, err := car.Economy()
	if err != nil {
		panic(err) // the fuel consumption is zero
	}

	fmt.Printf("%.1f\n", e.ConvertTo(fueleconomy.MilePerGallon))
	// Output: 47.0 mpg
}

func ExampleAmount_Entities() {
	// A mole of water weighs its molar mass in grams.
	water := units.NewMolarMass(18.015, molarmass.GramPerMole)

	n, err := units.NewMass(18.015, mass.Gram).DivideByMolarMass(water)
	if err != nil {
		panic(err) // the molar mass is zero
	}

	fmt.Printf("%.3e molecules\n", n.Entities())
	// Output: 6.022e+23 molecules
}

func ExampleRadioactivity_DecayAfter() {
	// Iodine-131, used in medicine, has a half-life of 8.02 days.
	dose := units.NewRadioactivity(370, radioactivity.Megabecquerel)

	left, err := dose.DecayAfter(units.NewTime(30, time.Day), units.NewTime(8.02, time.Day))
	if err != nil {
		panic(err) // the half-life is zero
	}

	fmt.Printf("%.1f\n", left)
	// Output: 27.7 MBq
}

func ExampleAbsorbedDose_Equivalent() {
	// Alpha particles have a radiation weighting factor of 20.
	alpha := units.NewAbsorbedDose(0.1, absorbeddose.Milligray)

	fmt.Printf("%.0f\n", alpha.Equivalent(20).ConvertTo(equivalentdose.Millisievert))
	// Output: 2 mSv
}

func ExampleRatio_Fraction() {
	// A 7.5 % share of a year's consumption.
	share := units.NewRatio(7.5, ratio.Percent)
	year := units.NewEnergy(1200, energy.KilowattHour)

	fmt.Println(year.Mul(share.Fraction()))
	// Output: 90 kWh
}

func ExampleDistance_Ratio() {
	done := units.NewDistance(12, distance.Kilometer)
	marathon := units.NewDistance(42.195, distance.Kilometer)

	progress, err := done.Ratio(marathon)
	if err != nil {
		panic(err) // the marathon is zero
	}

	fmt.Printf("%.1f\n", progress.ConvertTo(ratio.Percent))
	// Output: 28.4 %
}
