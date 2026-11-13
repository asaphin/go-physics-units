package units

import (
	"reflect"

	"github.com/asaphin/go-physics-units/conversion"
	"github.com/asaphin/go-physics-units/internal/rates"
)

// MeasurementCreator creates measurements of a user-defined measure, such as a
// currency.
type MeasurementCreator[T Measurement] interface {
	// New creates a measurement with the given value and unit. It fails if the
	// unit doesn't belong to the measure.
	New(value float64, unit string) (T, error)
	// ParseString parses a measurement from a string of the form
	// "<value> <unit>", e.g. "10 USD".
	ParseString(s string) (T, error)
}

type measurementCreator[T Measurement] struct {
	conversionRates rates.Table
}

// NewMeasurementCreator returns a creator for the user-defined measure named
// measureType, whose units are described by factors (see conversion.Factors).
// The factors are copied, so later changes to the map have no effect.
//
// T must be satisfied by Measurement alone: Measurement itself or an interface
// that only embeds it, e.g. type Currency interface{ units.Measurement }.
// The predefined measure interfaces such as Distance are rejected; use their own
// constructors instead.
//
// measureType must be non-empty and differ from the predefined MeasureType
// values, baseUnit must be one of the factors' units, and every factor must be
// positive and finite.
func NewMeasurementCreator[T Measurement](
	measureType, baseUnit string, factors conversion.Factors,
) (MeasurementCreator[T], error) {
	if _, ok := Measurement(&baseMeasurement{}).(T); !ok {
		return nil, newError(ErrInvalidMeasure, "unsupported measurement type "+reflect.TypeFor[T]().String()+
			": it must be satisfied by units.Measurement alone", nil)
	}

	if measureType == "" {
		return nil, newError(ErrInvalidMeasure, "measure type must not be empty", nil)
	}

	if isReservedMeasureType(MeasureType(measureType)) {
		return nil, newError(
			ErrInvalidMeasure,
			"measure type "+measureType+" is reserved for a predefined measure",
			nil,
		)
	}

	table, err := rates.NewTable(measureType, factors)
	if err != nil {
		return nil, newError(ErrInvalidFactors, err.Error(), nil)
	}

	if !table.Has(rates.CanonicalUnit(baseUnit)) {
		return nil, newError(ErrInvalidFactors, "provided factors doesn't contain unit "+baseUnit, nil)
	}

	return &measurementCreator[T]{conversionRates: table}, nil
}

// New creates a measurement with the given value and unit.
func (c *measurementCreator[T]) New(value float64, unit string) (T, error) {
	bm, err := newBaseMeasurement(value, unit, c.conversionRates)
	if err != nil {
		var zero T

		return zero, err
	}

	return Measurement(&bm).(T), nil //nolint:errcheck // NewMeasurementCreator checked that *baseMeasurement satisfies T
}

// ParseString parses a measurement from a string of the form "<value> <unit>".
func (c *measurementCreator[T]) ParseString(s string) (T, error) {
	value, unit, err := parseValueAndUnit(s)
	if err != nil {
		var zero T

		return zero, err
	}

	return c.New(value, unit)
}
