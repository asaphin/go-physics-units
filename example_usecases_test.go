package units_test

// The scenarios of USECASES.md. TestUseCasesDoc checks that the document
// shows each declaration and each example body below verbatim, with its
// output, so edit them here and copy the change into USECASES.md.

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"maps"
	"slices"
	builtintime "time"

	units "github.com/asaphin/go-physics-units"
	"github.com/asaphin/go-physics-units/acceleration"
	"github.com/asaphin/go-physics-units/amount"
	"github.com/asaphin/go-physics-units/angularvelocity"
	"github.com/asaphin/go-physics-units/area"
	"github.com/asaphin/go-physics-units/charge"
	"github.com/asaphin/go-physics-units/current"
	"github.com/asaphin/go-physics-units/datasize"
	"github.com/asaphin/go-physics-units/density"
	"github.com/asaphin/go-physics-units/distance"
	"github.com/asaphin/go-physics-units/doserate"
	"github.com/asaphin/go-physics-units/energy"
	"github.com/asaphin/go-physics-units/equivalentdose"
	"github.com/asaphin/go-physics-units/fuelconsumption"
	"github.com/asaphin/go-physics-units/fueleconomy"
	"github.com/asaphin/go-physics-units/mass"
	"github.com/asaphin/go-physics-units/molarity"
	"github.com/asaphin/go-physics-units/molarmass"
	"github.com/asaphin/go-physics-units/physics"
	"github.com/asaphin/go-physics-units/power"
	"github.com/asaphin/go-physics-units/pressure"
	"github.com/asaphin/go-physics-units/ratio"
	"github.com/asaphin/go-physics-units/temperature"
	"github.com/asaphin/go-physics-units/time"
	"github.com/asaphin/go-physics-units/torque"
	"github.com/asaphin/go-physics-units/velocity"
	"github.com/asaphin/go-physics-units/voltage"
	"github.com/asaphin/go-physics-units/volume"
)

func Example_fitnessTracker() {
	// Runners type distances in many ways: "6.2 Miles", "10km", "5000 metres".
	p, err := units.NewParser(
		units.WithAliases(map[string]string{"mile": "mi", "miles": "mi", "metres": "m"}),
		units.CaseInsensitive(),
		units.AllowNoSpace(),
	)
	if err != nil {
		panic(err)
	}

	m, err := p.Parse("6.2 Miles")
	if err != nil {
		panic(err) // not a number and a unit the parser knows
	}

	run, err := units.ToDistance(m)
	if err != nil {
		panic(err) // a unit of another measure, such as "6.2 kg"
	}

	// The watch reports the elapsed time as a time.Duration.
	elapsed := units.ParseDuration(52*builtintime.Minute + 30*builtintime.Second)

	speed, err := run.DivideByTime(elapsed)
	if err != nil {
		panic(err) // the time is zero
	}

	// The pace is the time per kilometre.
	pace, err := elapsed.Div(run.ConvertTo(distance.Kilometer).Value())
	if err != nil {
		panic(err) // the distance is zero
	}

	paceDuration, err := pace.ToDuration()
	if err != nil {
		panic(err) // beyond the range of time.Duration
	}

	fmt.Printf("distance: %.2f\n", run.ConvertTo(distance.Kilometer))
	fmt.Printf("speed:    %.1f\n", speed.ConvertTo(velocity.KilometerPerHour))
	fmt.Printf("pace:     %v per km\n", paceDuration.Round(builtintime.Second))
	// Output:
	// distance: 9.98 km
	// speed:    11.4 km/h
	// pace:     5m16s per km
}

func Example_iotTelemetry() {
	// Sensors from different vendors report in their own units.
	type Reading struct {
		Device      string            `json:"device"`
		Temperature units.Temperature `json:"temperature"`
		Pressure    units.Pressure    `json:"pressure"`
		Battery     units.Voltage     `json:"battery"`
	}

	payload := `[
		{"device": "greenhouse", "temperature": "77 °F", "pressure": "1013.25 hPa", "battery": "3.7 V"},
		{"device": "warehouse", "temperature": "18.5 °C", "pressure": "14.5 psi", "battery": "3100 mV"},
		{"device": "cellar", "temperature": "285.15 K", "pressure": "29.92 inHg", "battery": "2.8 V"}
	]`

	var readings []Reading
	if err := json.Unmarshal([]byte(payload), &readings); err != nil {
		panic(err) // a malformed value, or a unit of the wrong measure
	}

	lowBattery := units.NewVoltage(3, voltage.Volt)

	for _, r := range readings {
		status := "ok"
		if r.Battery.Compare(lowBattery) < 0 {
			status = "replace battery"
		}

		fmt.Printf("%-10s %6.1f %9.1f  %-7v %s\n", r.Device,
			r.Temperature.ConvertTo(temperature.Celsius), r.Pressure.ConvertTo(pressure.Hectopascal),
			r.Battery, status)
	}

	// A unit of the wrong measure is rejected when the payload is decoded.
	var bad Reading

	err := json.Unmarshal([]byte(`{"device": "attic", "battery": "3.7 kg"}`), &bad)
	fmt.Println(errors.Is(err, units.ErrUnknownUnit), err)
	// Output:
	// greenhouse   25.0 °C    1013.2 hPa  3.7 V   ok
	// warehouse    18.5 °C     999.7 hPa  3100 mV ok
	// cellar       12.0 °C    1013.2 hPa  2.8 V   replace battery
	// true unknown voltage unit kg
}

func Example_serviceConfig() {
	// Limits in a configuration file are readable and carry their units.
	type Config struct {
		MaxUpload units.DataSize `json:"max_upload"`
		Uplink    units.DataRate `json:"uplink"`
	}

	var cfg Config
	if err := json.Unmarshal([]byte(`{"max_upload": "25 MiB", "uplink": "8 Mbit/s"}`), &cfg); err != nil {
		panic(err)
	}

	// The same types work as command-line flags, which override the file.
	fs := flag.NewFlagSet("server", flag.ContinueOnError)
	fs.TextVar(&cfg.MaxUpload, "max-upload", cfg.MaxUpload, "the largest accepted upload")

	if err := fs.Parse([]string{"-max-upload", "100 MB"}); err != nil {
		panic(err)
	}

	// Reject an upload by its Content-Length before reading it.
	upload := units.NewDataSize(120_000_000, datasize.Byte)
	if upload.Compare(cfg.MaxUpload) > 0 {
		fmt.Printf("rejected: %.0f is larger than %v\n", upload.ConvertTo(datasize.Megabyte), cfg.MaxUpload)
	}

	// The time the largest upload takes over the uplink, with 50% headroom,
	// as a timeout for context.WithTimeout or http.Server.
	transfer, err := cfg.MaxUpload.DivideByDataRate(cfg.Uplink)
	if err != nil {
		panic(err) // the uplink is zero
	}

	timeout, err := transfer.Mul(1.5).ToDuration()
	if err != nil {
		panic(err) // beyond the range of time.Duration
	}

	fmt.Println("timeout:", timeout)
	// Output:
	// rejected: 120 MB is larger than 100 MB
	// timeout: 2m30s
}

func Example_dataImport() {
	// An instrument exports one reading per line, in whatever unit it was set to.
	lines := []string{"12.5 kg", "350 g", "1.2 m", "20 s", "3 mi", "40 min", "9 parsecs"}

	totals := make(map[units.MeasureType]units.Measurement)

	for _, line := range lines {
		m, err := units.ParseString(line)
		if err != nil {
			fmt.Println("skipped:", err)

			continue
		}

		if total, ok := totals[m.Type()]; ok {
			// The sum keeps the unit of the first reading of each measure.
			if m, err = total.AddMeasurement(m); err != nil {
				panic(err) // never: both have the same measure
			}
		}

		totals[m.Type()] = m
	}

	for _, mt := range slices.Sorted(maps.Keys(totals)) {
		fmt.Printf("%-8s %v\n", mt, totals[mt])
	}
	// Output:
	// skipped: unknown unit parsecs
	// distance 4829.232 m
	// mass     12.85 kg
	// time     2420 s
}

func Example_electronics() {
	// The series resistor for an LED: R = (supply - forward voltage) / current.
	supply := units.NewVoltage(5, voltage.Volt)
	forward := units.NewVoltage(2, voltage.Volt)
	led := units.NewCurrent(20, current.Milliampere)

	drop := supply.Sub(forward)

	r, err := drop.DivideByCurrent(led)
	if err != nil {
		panic(err) // the current is zero
	}

	heat := drop.MultiplyByCurrent(led) // the power the resistor dissipates

	fmt.Printf("resistor: %.0f, dissipating %.0f\n", r, heat.ConvertTo(power.Milliwatt))

	// How long a 2000 mAh battery runs the LED and a 5 mA microcontroller.
	battery := units.NewCharge(2000, charge.MilliampereHour)
	draw := led.Add(units.NewCurrent(5, current.Milliampere))

	life, err := battery.DivideByCurrent(draw)
	if err != nil {
		panic(err) // the current is zero
	}

	fmt.Printf("battery life: %.0f\n", life.ConvertTo(time.Hour))
	// Output:
	// resistor: 150 Ω, dissipating 60 mW
	// battery life: 80 h
}

func Example_homeEnergy() {
	// The energy to heat 1.5 L of water from 15 °C to boiling.
	water := units.NewDensity(1, density.KilogramPerLiter).MultiplyByVolume(units.NewVolume(1.5, volume.Liter))

	tap, err := units.NewTemperature(15, temperature.Celsius)
	if err != nil {
		panic(err) // below absolute zero
	}

	boiling, err := units.NewTemperature(100, temperature.Celsius)
	if err != nil {
		panic(err)
	}

	// A temperature difference converts without an offset: 85 °C is 85 K.
	rise := boiling.Sub(tap).ConvertTo(temperature.Kelvin)

	const waterSpecificHeat = 4186 // J/(kg·K)

	needed := units.NewEnergy(water.ConvertToBaseUnits().Value()*waterSpecificHeat*rise.Value(), energy.Joule)

	// The kettle draws 2 kW for 5 minutes.
	used := units.NewPower(2, power.Kilowatt).MultiplyByTime(units.NewTime(5, time.Minute))

	fmt.Printf("needed %.0f, used %.0f: %.0f%% efficient\n",
		needed.ConvertTo(energy.Kilojoule), used.ConvertTo(energy.Kilojoule), 100*needed.Value()/used.Value())

	// One boil a day for a year, at 0.30 per kWh.
	year := used.Mul(365).ConvertTo(energy.KilowattHour)

	fmt.Printf("a year: %.1f, costing %.2f\n", year, year.Value()*0.30)
	// Output:
	// needed 534 kJ, used 600 kJ: 89% efficient
	// a year: 60.8 kWh, costing 18.25
}

func Example_aviationFuel() {
	// Jet A-1 weighs about 0.804 kg per litre at 15 °C.
	jetA1 := units.NewDensity(0.804, density.KilogramPerLiter)

	// The flight plan asks for 12 t of fuel; the truck meters litres.
	uplift, err := units.NewMass(12, mass.Tonne).DivideByDensity(jetA1)
	if err != nil {
		panic(err) // the density is zero
	}

	fmt.Printf("uplift: %.0f\n", uplift.ConvertTo(volume.Liter))

	// The truck delivered 4000 US gallons. How much fuel is that, in pounds?
	delivered := jetA1.MultiplyByVolume(units.NewVolume(4000, volume.Gallon))

	fmt.Printf("delivered: %.0f = %.0f\n", delivered.ConvertTo(mass.Kilogram), delivered.ConvertTo(mass.Pound))
	// Output:
	// uplift: 14925 L
	// delivered: 12174 kg = 26839 lb
}

func Example_pressure() {
	// Why a stiletto heel dents a floor that an elephant doesn't.
	g := units.NewAcceleration(1, acceleration.StandardGravity)

	heel, err := units.NewMass(60, mass.Kilogram).MultiplyByAcceleration(g).
		DivideByArea(units.NewArea(1, area.SquareCentimeter))
	if err != nil {
		panic(err) // the area is zero
	}

	// A 5 t elephant standing on four feet of about 0.1 m² each.
	elephant, err := units.NewMass(5, mass.Tonne).MultiplyByAcceleration(g).
		DivideByArea(units.NewArea(4*0.1, area.SquareMeter))
	if err != nil {
		panic(err)
	}

	show := func(name string, p units.Pressure) {
		fmt.Printf("%-9s %5.2f = %4.2f = %3.0f\n", name, p.ConvertTo(pressure.KilogramForcePerSquareCentimeter),
			p.ConvertTo(pressure.Megapascal), p.ConvertTo(pressure.PoundPerSquareInch))
	}

	// 60 kg on 1 cm² is 60 kgf/cm²: a kilogram-force per square centimetre is
	// the weight of 1 kg on 1 cm².
	show("heel:", heel)
	show("elephant:", elephant)
	// Output:
	// heel:     60.00 kgf/cm² = 5.88 MPa = 853 psi
	// elephant:  1.25 kgf/cm² = 0.12 MPa =  18 psi
}

// kitchen extends the predefined volume units with US customary cooking
// measures (NIST Handbook 133). Its values are units.Quantity[kitchen].
type kitchen struct{}

const usFluidOunce = 29.5735295625e-6 // m³, exact

var kitchenUnits = func() units.UnitSet {
	factors := volume.ConversionFactors() // a copy of mL, L, pt, gal, …
	factors["cup"] = 8 * usFluidOunce
	factors["tbsp"] = usFluidOunce / 2
	factors["tsp"] = usFluidOunce / 6

	return units.MustNewUnitSet("kitchen", volume.BaseUnit.String(), factors)
}()

func (kitchen) Units() units.UnitSet { return kitchenUnits }

func Example_recipeScaling() {
	// A recipe for 4 servings, scaled to 6, with the weights for kitchen scales.
	// The densities are typical values: flour's, for one, depends on how it's packed.
	recipe := []struct {
		ingredient, amount string
		density            units.Density
	}{
		{"flour", "2 cup", units.NewDensity(0.53, density.GramPerMilliliter)},
		{"milk", "250 mL", units.NewDensity(1.03, density.GramPerMilliliter)},
		{"sugar", "3 tbsp", units.NewDensity(0.85, density.GramPerMilliliter)},
		{"butter", "4 tbsp", units.NewDensity(0.96, density.GramPerMilliliter)},
	}

	ml := units.MustUnitOf[kitchen]("mL")

	for _, line := range recipe {
		q, err := units.ParseQuantity[kitchen](line.amount)
		if err != nil {
			panic(err) // not a kitchen unit
		}

		scaled := q.Mul(6.0 / 4)

		// The kitchen measure has volume's base unit, the cubic metre, so the
		// amount converts to a units.Volume, and density × volume is its mass.
		vol := units.NewVolume(scaled.ConvertToBaseUnits().Value(), volume.BaseUnit)
		weight := line.density.MultiplyByVolume(vol)

		fmt.Printf("%-7s %-9v %4.0f  %4.0f\n", line.ingredient, scaled, scaled.ConvertTo(ml),
			weight.ConvertTo(mass.Gram))
	}
	// Output:
	// flour   3 cup      710 mL   376 g
	// milk    375 mL     375 mL   386 g
	// sugar   4.5 tbsp    67 mL    57 g
	// butter  6 tbsp      89 mL    85 g
}

func Example_carTrip() {
	// The engine makes 350 N·m at 4000 rev/min.
	engine := units.NewTorque(350, torque.NewtonMeter).MultiplyByAngularVelocity(
		units.NewAngularVelocity(4000, angularvelocity.RevolutionPerMinute))

	fmt.Printf("power:   %.0f = %.0f\n", engine.ConvertTo(power.Kilowatt), engine.ConvertTo(power.Horsepower))

	// A 620 km trip at 6.8 L/100km, with fuel at 1.85 a litre.
	trip := units.NewDistance(620, distance.Kilometer)
	consumption := units.NewFuelConsumption(6.8, fuelconsumption.LiterPer100Kilometers)
	fuel := consumption.MultiplyByDistance(trip).ConvertTo(volume.Liter)

	fmt.Printf("fuel:    %.1f, costing %.2f\n", fuel, fuel.Value()*1.85)

	// The same car, as an American driver would rate it.
	economy, err := consumption.Economy()
	if err != nil {
		panic(err) // the consumption is zero
	}

	fmt.Printf("economy: %.1f\n", economy.ConvertTo(fueleconomy.MilePerGallon))

	// The time it takes at an average of 95 km/h.
	duration, err := trip.DivideByVelocity(units.NewVelocity(95, velocity.KilometerPerHour))
	if err != nil {
		panic(err) // the speed is zero
	}

	d, err := duration.ToDuration()
	if err != nil {
		panic(err) // beyond the range of time.Duration
	}

	fmt.Println("time:   ", d.Round(builtintime.Minute))
	// Output:
	// power:   147 kW = 197 hp
	// fuel:    42.2 L, costing 78.00
	// economy: 34.6 mpg
	// time:    6h32m0s
}

func Example_labSolution() {
	// 500 mL of 0.15 M sodium chloride, NaCl (58.44 g/mol).
	stock := units.NewMolarity(0.15, molarity.Molar)
	salt := stock.MultiplyByVolume(units.NewVolume(500, volume.Milliliter))

	fmt.Printf("weigh out %.3f of NaCl (%.0f)\n",
		salt.MultiplyByMolarMass(units.NewMolarMass(58.44, molarmass.GramPerMole)).ConvertTo(mass.Gram),
		salt.ConvertTo(amount.Millimole))

	// Dilute 10 mL of the stock to 250 mL.
	aliquot := stock.MultiplyByVolume(units.NewVolume(10, volume.Milliliter))

	diluted, err := aliquot.DivideByVolume(units.NewVolume(250, volume.Milliliter))
	if err != nil {
		panic(err) // the volume is zero
	}

	fmt.Printf("diluted: %.0f\n", diluted.ConvertTo(molarity.Millimolar))
	fmt.Printf("sodium ions in the aliquot: %.2e\n", aliquot.Entities())
	// Output:
	// weigh out 4.383 g of NaCl (75 mmol)
	// diluted: 6 mM
	// sodium ions in the aliquot: 9.03e+20
}

func Example_flightRadiation() {
	// Cosmic rays give about 5 μSv/h at cruising altitude.
	cruise := units.NewDoseRate(5, doserate.MicrosievertPerHour)
	flight := units.NewTime(11, time.Hour)

	dose := cruise.MultiplyByTime(flight)

	fmt.Printf("one flight: %.0f\n", dose.ConvertTo(equivalentdose.Microsievert))

	// Compared with a year of natural background radiation.
	share, err := dose.Ratio(units.NewEquivalentDose(2.4, equivalentdose.Millisievert))
	if err != nil {
		panic(err) // the background dose is zero
	}

	fmt.Printf("of a year's background: %.1f\n", share.ConvertTo(ratio.Percent))

	// The hours in the air that add up to 1 mSv.
	hours, err := units.NewEquivalentDose(1, equivalentdose.Millisievert).DivideByDoseRate(cruise)
	if err != nil {
		panic(err) // the dose rate is zero
	}

	flights, err := hours.Ratio(flight)
	if err != nil {
		panic(err) // the flight takes no time
	}

	fmt.Printf("1 mSv after %.0f, about %.0f such flights\n", hours.ConvertTo(time.Hour), flights.Fraction())
	// Output:
	// one flight: 55 μSv
	// of a year's background: 2.3 %
	// 1 mSv after 200 h, about 18 such flights
}

func Example_orbits() {
	earth := units.NewMass(1, mass.EarthMass)
	radius := units.NewDistance(6371, distance.Kilometer)

	// The International Space Station flies about 420 km up.
	iss := radius.Add(units.NewDistance(420, distance.Kilometer))

	speed, err := physics.OrbitalVelocity(earth, iss)
	if err != nil {
		panic(err) // the radius is zero
	}

	period, err := physics.OrbitalPeriod(earth, iss)
	if err != nil {
		panic(err) // the mass is zero
	}

	fmt.Printf("ISS: %.2f, once every %.0f\n", speed.ConvertTo(velocity.KilometerPerSecond),
		period.ConvertTo(time.Minute))

	escape, err := physics.EscapeVelocity(earth, radius)
	if err != nil {
		panic(err) // the radius is zero
	}

	fmt.Printf("escape from the surface: %.1f\n", escape.ConvertTo(velocity.KilometerPerSecond))

	// The kinetic energy of 1 kg at the station's speed, as TNT.
	impact := units.NewMass(1, mass.Kilogram).KineticEnergy(speed)

	fmt.Printf("1 kg at that speed: %.1f, like %.1f kg of TNT\n", impact.ConvertTo(energy.Megajoule),
		impact.ConvertTo(energy.GramOfTNT).Value()/1000)
	// Output:
	// ISS: 7.66 km/s, once every 93 min
	// escape from the surface: 11.2 km/s
	// 1 kg at that speed: 29.3 MJ, like 7.0 kg of TNT
}
