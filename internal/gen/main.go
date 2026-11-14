// Command gen generates the code that every predefined measure shares: the
// methods of each unit package's Unit type (<pkg>/unit_gen.go) and the
// API of each measure type in the root package (measures_gen.go).
//
// Run it from the repository root, after adding a measure to the measures
// list below:
//
//	go generate ./...
//
// TestGenerated fails when a generated file is out of date.
package main

import (
	"bytes"
	"fmt"
	"go/format"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

// measure describes one predefined measure.
type measure struct {
	Type       string // root type, e.g. Distance
	Pkg        string // unit package, e.g. distance
	Recv       string // receiver name, e.g. d
	Noun       string // e.g. distance
	Plural     string // e.g. distances
	BaseName   string // the base unit in words, plural: metres
	BaseSymbol string // the base unit's symbol: m
	Example    string // a string that parses: 10 km
}

// measures lists the predefined measures, in the order of their MeasureType
// constants.
var measures = []measure{
	{"Distance", "distance", "d", "distance", "distances", "metres", "m", "10 km"},
	{"Time", "time", "t", "time", "times", "seconds", "s", "10 min"},
	{"Mass", "mass", "m", "mass", "masses", "kilograms", "kg", "10 kg"},
	{"Velocity", "velocity", "v", "velocity", "velocities", "metres per second", "m/s", "10 km/h"},
	{"Area", "area", "a", "area", "areas", "square metres", "m²", "10 m²"},
	{"Volume", "volume", "v", "volume", "volumes", "cubic metres", "m³", "10 L"},
	{
		"Acceleration", "acceleration", "a", "acceleration", "accelerations",
		"metres per second squared", "m/s²", "9.8 m/s²",
	},
	{"Force", "force", "f", "force", "forces", "newtons", "N", "10 N"},
	{"Energy", "energy", "e", "energy", "energies", "joules", "J", "10 kWh"},
	{"Power", "power", "p", "power", "powers", "watts", "W", "100 W"},
	{"Pressure", "pressure", "p", "pressure", "pressures", "pascals", "Pa", "101.3 kPa"},
	{"Frequency", "frequency", "f", "frequency", "frequencies", "hertz", "Hz", "50 Hz"},
	{"Angle", "angle", "a", "angle", "angles", "radians", "rad", "90 °"},
	{"DataSize", "datasize", "s", "data size", "data sizes", "bytes", "B", "1.5 GB"},
	{"Density", "density", "d", "density", "densities", "kilograms per cubic metre", "kg/m³", "1.2 kg/m³"},
	{"DataRate", "datarate", "r", "data rate", "data rates", "bytes per second", "B/s", "100 Mbit/s"},
	{"Current", "current", "i", "electric current", "electric currents", "amperes", "A", "2.5 mA"},
	{"Voltage", "voltage", "v", "voltage", "voltages", "volts", "V", "230 V"},
	{"Charge", "charge", "q", "electric charge", "electric charges", "coulombs", "C", "3000 mAh"},
	{"Resistance", "resistance", "r", "resistance", "resistances", "ohms", "Ω", "4.7 kΩ"},
	{
		"AngularVelocity", "angularvelocity", "w", "angular velocity", "angular velocities",
		"radians per second", "rad/s", "33.3 rev/min",
	},
	{"Torque", "torque", "tq", "torque", "torques", "newton metres", "N·m", "40 N·m"},
	{"FlowRate", "flowrate", "q", "flow rate", "flow rates", "cubic metres per second", "m³/s", "12 L/min"},
	{
		"MassFlowRate",
		"massflowrate",
		"mr",
		"mass flow rate",
		"mass flow rates",
		"kilograms per second",
		"kg/s",
		"250 kg/h",
	},
	{"Momentum", "momentum", "p", "momentum", "momenta", "kilogram metres per second", "kg·m/s", "10 N·s"},
	{"Pace", "pace", "p", "pace", "paces", "seconds per metre", "s/m", "5 min/km"},
	{"FuelEconomy", "fueleconomy", "fe", "fuel economy", "fuel economies", "metres per cubic metre", "m/m³", "15 km/L"},
	{
		"FuelConsumption",
		"fuelconsumption",
		"fc",
		"fuel consumption",
		"fuel consumptions",
		"cubic metres per metre",
		"m³/m",
		"6.5 L/100km",
	},
	{"Amount", "amount", "n", "amount of substance", "amounts of substance", "moles", "mol", "2.5 mmol"},
	{"MolarMass", "molarmass", "mm", "molar mass", "molar masses", "kilograms per mole", "kg/mol", "18.015 g/mol"},
	{
		"Molarity",
		"molarity",
		"c",
		"molar concentration",
		"molar concentrations",
		"moles per cubic metre",
		"mol/m³",
		"150 mM",
	},
	{"Radioactivity", "radioactivity", "a", "radioactivity", "radioactivities", "becquerels", "Bq", "3.7 GBq"},
	{"AbsorbedDose", "absorbeddose", "d", "absorbed dose", "absorbed doses", "grays", "Gy", "2 Gy"},
	{"EquivalentDose", "equivalentdose", "h", "equivalent dose", "equivalent doses", "sieverts", "Sv", "2.4 mSv"},
	{"DoseRate", "doserate", "r", "dose rate", "dose rates", "sieverts per second", "Sv/s", "0.12 μSv/h"},
	{"Capacitance", "capacitance", "c", "capacitance", "capacitances", "farads", "F", "100 μF"},
	{"Inductance", "inductance", "l", "inductance", "inductances", "henries", "H", "10 mH"},
	{
		"MagneticField",
		"magneticfield",
		"mf",
		"magnetic flux density",
		"magnetic flux densities",
		"teslas",
		"T",
		"50 μT",
	},
	{"Illuminance", "illuminance", "e", "illuminance", "illuminances", "lux", "lx", "500 lx"},
	{"LuminousFlux", "luminousflux", "f", "luminous flux", "luminous fluxes", "lumens", "lm", "800 lm"},
	{
		"SpecificEnergy",
		"specificenergy",
		"e",
		"specific energy",
		"specific energies",
		"joules per kilogram",
		"J/kg",
		"250 Wh/kg",
	},
	{
		"EnergyDensity",
		"energydensity",
		"u",
		"energy density",
		"energy densities",
		"joules per cubic metre",
		"J/m³",
		"700 Wh/L",
	},
	{"Ratio", "ratio", "r", "ratio", "ratios", "the unit 1, as a plain fraction", "1", "12.5 %"},
}

func main() {
	files, err := generate()
	if err != nil {
		fmt.Fprintln(os.Stderr, "gen:", err)
		os.Exit(1)
	}

	for path, content := range files {
		if err := os.WriteFile(path, content, 0o600); err != nil { // the mode applies only to new files
			fmt.Fprintln(os.Stderr, "gen:", err)
			os.Exit(1)
		}
	}
}

// generate returns the generated files by path, relative to the repository
// root.
func generate() (map[string][]byte, error) {
	files := make(map[string][]byte, len(measures)+1)

	for i := range measures {
		m := &measures[i]

		src, err := execute(unitTemplate, m)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", m.Pkg, err)
		}

		files[filepath.Join(m.Pkg, "unit_gen.go")] = src
	}

	src, err := execute(rootTemplate, measures)
	if err != nil {
		return nil, fmt.Errorf("measures_gen.go: %w", err)
	}

	files["measures_gen.go"] = src

	return files, nil
}

func execute(t *template.Template, data any) ([]byte, error) {
	var buf bytes.Buffer
	if err := t.Execute(&buf, data); err != nil {
		return nil, fmt.Errorf("executing the template: %w", err)
	}

	src, err := format.Source(buf.Bytes())
	if err != nil {
		return nil, fmt.Errorf("formatting: %w\n%s", err, buf.Bytes())
	}

	return src, nil
}

// funcs are the template functions: a adds the indefinite article.
var funcs = template.FuncMap{
	"a": func(noun string) string {
		if strings.ContainsRune("aeiou", rune(noun[0])) {
			return "an " + noun
		}

		return "a " + noun
	},
}

const header = `// Code generated by "go run ./internal/gen"; DO NOT EDIT.

`

var unitTemplate = template.Must(template.New("unit").Funcs(funcs).Parse(header + `package {{.Pkg}}

import (
	"iter"
	"strconv"

	"github.com/asaphin/go-physics-units/conversion"
	"github.com/asaphin/go-physics-units/internal/rates"
)

// definition is the symbol and the factor of a unit.
type definition struct {
	symbol string
	factor float64
}

// rateTable and bySymbol are built from definitions at package init and never
// written afterwards.
var rateTable, bySymbol = build()

// build validates definitions and returns the rate between every pair of
// units and the units by canonical symbol.
func build() ([unitCount][unitCount]float64, map[string]Unit) {
	symbols := make([]string, unitCount)
	factors := make([]float64, unitCount)

	for u, d := range definitions {
		symbols[u], factors[u] = d.symbol, d.factor
	}

	rates.MustValidate("{{.Noun}}", symbols, factors)

	var table [unitCount][unitCount]float64

	matrix := rates.Matrix(factors)
	for i := range table {
		copy(table[i][:], matrix[i*int(unitCount):])
	}

	index := make(map[string]Unit, unitCount)
	for u, symbol := range symbols {
		index[rates.CanonicalUnit(symbol)] = Unit(u)
	}

	return table, index
}

// String returns the unit's symbol. A Unit that isn't one of the constants of
// this package returns "{{.Pkg}}.Unit(<number>)".
func (u Unit) String() string {
	if u >= unitCount {
		return "{{.Pkg}}.Unit(" + strconv.Itoa(int(u)) + ")"
	}

	return definitions[u].symbol
}

// Factor returns how many base units ({{.BaseSymbol}}) there are in one u. It panics if u
// isn't valid.
func (u Unit) Factor() float64 {
	return definitions[u].factor
}

// Valid reports whether u is one of the constants of this package. Values in
// any other unit panic when converted.
func (u Unit) Valid() bool {
	return u < unitCount
}

// ParseUnit returns the unit with the given symbol, and whether there is one.
// Both forms of the micro sign are accepted.
func ParseUnit(symbol string) (Unit, bool) {
	u, ok := bySymbol[rates.CanonicalUnit(symbol)]

	return u, ok
}

// Units returns the units in the order of their constants, starting with the
// base unit.
func Units() iter.Seq[Unit] {
	return func(yield func(Unit) bool) {
		for u := range Unit(unitCount) {
			if !yield(u) {
				return
			}
		}
	}
}

// ConversionFactors returns the factors of the {{.Noun}} units by symbol: how many
// base units ({{.BaseSymbol}}) there are in one unit. The map is a new copy on every
// call, so it can be changed freely, for example to define a custom measure
// with more units (see units.NewUnitSet).
func ConversionFactors() conversion.Factors {
	factors := make(conversion.Factors, unitCount)
	for _, d := range definitions {
		factors[d.symbol] = d.factor
	}

	return factors
}

// Rate returns the multiplier that converts a value in unit from into unit to:
// from.Factor() / to.Factor(), computed once. The rate from a unit to itself
// is exactly 1. Rate panics if either unit isn't valid.
func Rate(from, to Unit) float64 {
	return rateTable[from][to]
}
`))

var rootTemplate = template.Must(template.New("root").Funcs(funcs).Parse(header + `package units

import (
	"cmp"
	"fmt"
{{range .}}
	"github.com/asaphin/go-physics-units/{{.Pkg}}"{{end}}
)

// The unit sets of the predefined measures, for the generic Measurement
// methods, NewMeasurement, ParseString and Parser. They are built at package
// init and never written afterwards.
var ({{range .}}
	{{.Pkg}}Set = newPredefinedSet(Measure{{.Type}}, {{.Pkg}}.Units(), func(v float64, u {{.Pkg}}.Unit) Measurement {
		return {{.Type}}{value: v, unit: u}
	}){{end}}
)

// predefinedSets lists the unit sets of the predefined measures.
var predefinedSets = [...]*unitSet{ {{range .}}
	{{.Pkg}}Set,{{end}}
}
{{range .}}
// New{{.Type}} returns {{a .Noun}} of value in unit, such as
// New{{.Type}}(10, [{{.Pkg}}.BaseUnit]).
func New{{.Type}}(value float64, unit {{.Pkg}}.Unit) {{.Type}} {
	return {{.Type}}{value: value, unit: unit}
}

// Value returns the numeric value, expressed in Unit.
func ({{.Recv}} {{.Type}}) Value() float64 {
	return {{.Recv}}.value
}

// Unit returns the unit of the value.
func ({{.Recv}} {{.Type}}) Unit() {{.Pkg}}.Unit {
	return {{.Recv}}.unit
}

// Symbol returns the symbol of the unit, as Unit().String() does.
func ({{.Recv}} {{.Type}}) Symbol() string {
	return {{.Recv}}.unit.String()
}

// Type returns Measure{{.Type}}.
func ({{.Type}}) Type() MeasureType {
	return Measure{{.Type}}
}

// String formats the {{.Noun}} as "<value> <symbol>", such as "{{.Example}}", the form
// StringTo{{.Type}} and ParseString accept. The value is formatted like fmt's %v,
// in the shortest form that parses back to the same float64.
func ({{.Recv}} {{.Type}}) String() string {
	return formatString({{.Recv}}.value, {{.Recv}}.unit.String())
}

// Format implements fmt.Formatter: %.2f formats the value and appends the
// symbol. See the package documentation for every verb.
func ({{.Recv}} {{.Type}}) Format(state fmt.State, verb rune) {
	formatMeasurement(state, verb, "{{.Type}}", {{.Recv}}.value, {{.Recv}}.unit.String())
}

// ConvertTo returns the {{.Noun}} expressed in unit.
func ({{.Recv}} {{.Type}}) ConvertTo(unit {{.Pkg}}.Unit) {{.Type}} {
	return {{.Type}}{value: {{.Recv}}.value * {{.Pkg}}.Rate({{.Recv}}.unit, unit), unit: unit}
}

// MustConvertTo is ConvertTo: like ConvertTo, it panics only if unit isn't
// valid.
//
// Deprecated: Use ConvertTo, which no longer fails.
//
//go:fix inline
func ({{.Recv}} {{.Type}}) MustConvertTo(unit {{.Pkg}}.Unit) {{.Type}} {
	return {{.Recv}}.ConvertTo(unit)
}

// ConvertToBaseUnits returns the {{.Noun}} expressed in {{.BaseName}}.
func ({{.Recv}} {{.Type}}) ConvertToBaseUnits() {{.Type}} {
	return {{.Type}}{value: {{.Recv}}.baseValue(), unit: {{.Pkg}}.BaseUnit}
}

// baseValue returns the value in {{.BaseName}}.
func ({{.Recv}} {{.Type}}) baseValue() float64 {
	return {{.Recv}}.value * {{.Recv}}.unit.Factor()
}

// Add returns the sum of the {{.Plural}}, expressed in the receiver's unit.
func ({{.Recv}} {{.Type}}) Add(other {{.Type}}) {{.Type}} {
	return {{.Type}}{value: {{.Recv}}.value + other.value*{{.Pkg}}.Rate(other.unit, {{.Recv}}.unit), unit: {{.Recv}}.unit}
}

// Sub returns the difference of the {{.Plural}}, expressed in the receiver's unit.
func ({{.Recv}} {{.Type}}) Sub(other {{.Type}}) {{.Type}} {
	return {{.Type}}{value: {{.Recv}}.value - other.value*{{.Pkg}}.Rate(other.unit, {{.Recv}}.unit), unit: {{.Recv}}.unit}
}

// Mul returns the {{.Noun}} multiplied by multiplier.
func ({{.Recv}} {{.Type}}) Mul(multiplier float64) {{.Type}} {
	return {{.Type}}{value: {{.Recv}}.value * multiplier, unit: {{.Recv}}.unit}
}

// Div returns the {{.Noun}} divided by divisor. It fails with [ErrZeroDivision] if
// divisor is zero.
func ({{.Recv}} {{.Type}}) Div(divisor float64) ({{.Type}}, error) {
	if divisor == 0 {
		return {{.Type}}{}, ErrZeroDivision
	}

	return {{.Type}}{value: {{.Recv}}.value / divisor, unit: {{.Recv}}.unit}, nil
}

// Compare compares the {{.Plural}}; see Comparer.
func ({{.Recv}} {{.Type}}) Compare(other {{.Type}}) int {
	return cmp.Compare({{.Recv}}.common(other))
}

// Equal reports whether the {{.Plural}} are equal; see Comparer.
func ({{.Recv}} {{.Type}}) Equal(other {{.Type}}) bool {
	return {{.Recv}}.Compare(other) == 0
}

// EqualWithinRel reports whether the {{.Plural}} are equal within the relative
// tolerance tol; see Comparer.
func ({{.Recv}} {{.Type}}) EqualWithinRel(other {{.Type}}, tol float64) bool {
	x, y := {{.Recv}}.common(other)

	return equalWithinRel(x, y, tol)
}

// Ratio returns {{.Recv}}/other, the ratio of the {{.Plural}}, in [ratio.BaseUnit]:
// convert it with ConvertTo(ratio.Percent) for a percentage. It fails with
// [ErrZeroDivision] if other is zero.
func ({{.Recv}} {{.Type}}) Ratio(other {{.Type}}) (Ratio, error) {
	x, y := {{.Recv}}.common(other)
	if y == 0 {
		return Ratio{}, ErrZeroDivision
	}

	return Ratio{value: x / y, unit: ratio.BaseUnit}, nil
}

// common returns the values of the {{.Plural}} in a common unit: their own if
// they share it, otherwise the base unit. Either way the result doesn't depend
// on the order of the operands, so comparisons are symmetric.
func ({{.Recv}} {{.Type}}) common(other {{.Type}}) (float64, float64) {
	if {{.Recv}}.unit == other.unit {
		return {{.Recv}}.value, other.value
	}

	return {{.Recv}}.baseValue(), other.baseValue()
}

// ConvertToMeasurement is ConvertTo for generic code: it takes a unit symbol.
// It fails with [ErrUnknownUnit] if symbol isn't {{a .Noun}} unit.
func ({{.Recv}} {{.Type}}) ConvertToMeasurement(symbol string) (Measurement, error) {
	unit, ok := {{.Pkg}}.ParseUnit(symbol)
	if !ok {
		return nil, unknownUnitError(Measure{{.Type}}, symbol)
	}

	return {{.Recv}}.ConvertTo(unit), nil
}

// AddMeasurement is Add for generic code. It fails with
// [ErrIncompatibleMeasures] if other isn't {{a .Noun}}.
func ({{.Recv}} {{.Type}}) AddMeasurement(other Measurement) (Measurement, error) {
	value, err := {{.Pkg}}Set.valueIn(other, uint16({{.Recv}}.unit))
	if err != nil {
		return nil, err
	}

	return {{.Type}}{value: {{.Recv}}.value + value, unit: {{.Recv}}.unit}, nil
}

// SubMeasurement is Sub for generic code. It fails with
// [ErrIncompatibleMeasures] if other isn't {{a .Noun}}.
func ({{.Recv}} {{.Type}}) SubMeasurement(other Measurement) (Measurement, error) {
	value, err := {{.Pkg}}Set.valueIn(other, uint16({{.Recv}}.unit))
	if err != nil {
		return nil, err
	}

	return {{.Type}}{value: {{.Recv}}.value - value, unit: {{.Recv}}.unit}, nil
}

// MulMeasurement is Mul for generic code.
func ({{.Recv}} {{.Type}}) MulMeasurement(multiplier float64) Measurement {
	return {{.Recv}}.Mul(multiplier)
}

// DivMeasurement is Div for generic code. It fails with [ErrZeroDivision] if
// divisor is zero.
func ({{.Recv}} {{.Type}}) DivMeasurement(divisor float64) (Measurement, error) {
	if divisor == 0 {
		return nil, ErrZeroDivision
	}

	return {{.Type}}{value: {{.Recv}}.value / divisor, unit: {{.Recv}}.unit}, nil
}

// CompareMeasurement is Compare for generic code. It fails with
// [ErrIncompatibleMeasures] if other isn't {{a .Noun}}.
func ({{.Recv}} {{.Type}}) CompareMeasurement(other Measurement) (int, error) {
	return {{.Pkg}}Set.compare({{.Recv}}.value, uint16({{.Recv}}.unit), other)
}

func ({{.Recv}} {{.Type}}) generic() genericValue {
	return genericValue{value: {{.Recv}}.value, unit: uint16({{.Recv}}.unit), set: {{.Pkg}}Set}
}

const err{{.Type}}Conversion incompatibleError = "not {{a .Noun}} measure"

// To{{.Type}} returns m as a {{.Type}}, for example a result of ParseString. It fails
// with [ErrIncompatibleMeasures] if m isn't {{a .Noun}}.
func To{{.Type}}(m Measurement) ({{.Type}}, error) {
	value, ok := m.({{.Type}})
	if !ok {
		return {{.Type}}{}, err{{.Type}}Conversion
	}

	return value, nil
}

// StringTo{{.Type}} parses {{a .Noun}} from a string such as "{{.Example}}". It fails with
// [ErrInvalidFormat] unless the string has the form "<value> <symbol>" with a
// finite value, and with [ErrUnknownUnit] if the symbol isn't {{a .Noun}} unit.
func StringTo{{.Type}}(s string) ({{.Type}}, error) {
	value, symbol, err := parseValueAndUnit(s)
	if err != nil {
		return {{.Type}}{}, err
	}

	unit, ok := {{.Pkg}}.ParseUnit(symbol)
	if !ok {
		return {{.Type}}{}, unknownUnitError(Measure{{.Type}}, symbol)
	}

	return {{.Type}}{value: value, unit: unit}, nil
}

// AppendText implements encoding.TextAppender: it appends the String form,
// such as "{{.Example}}". It fails with [ErrOutOfRange] if the value isn't finite, and
// with [ErrUnknownUnit] if the unit isn't valid, so that the text always parses
// back.
func ({{.Recv}} {{.Type}}) AppendText(b []byte) ([]byte, error) {
	if !{{.Recv}}.unit.Valid() {
		return b, unknownUnitError(Measure{{.Type}}, {{.Recv}}.unit.String())
	}

	return appendText(b, {{.Recv}}.value, {{.Recv}}.unit.String())
}

// MarshalText implements encoding.TextMarshaler, so {{.Noun}} values encode in
// JSON as strings such as "{{.Example}}". It fails like AppendText.
func ({{.Recv}} {{.Type}}) MarshalText() ([]byte, error) {
	return {{.Recv}}.AppendText(nil)
}

// UnmarshalText implements encoding.TextUnmarshaler: it parses text as
// StringTo{{.Type}} does and stores the result in {{.Recv}}. Unmarshalling needs a
// pointer receiver; it changes only the value it is called on. It fails like
// [StringTo{{.Type}}], and leaves {{.Recv}} unchanged when it does.
func ({{.Recv}} *{{.Type}}) UnmarshalText(text []byte) error {
	value, err := StringTo{{.Type}}(string(text))
	if err != nil {
		return err
	}

	*{{.Recv}} = value

	return nil
}
{{end}}`))
