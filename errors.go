package units

// Error is the type of the package's sentinel errors, such as ErrUnknownUnit.
// Test for them with errors.Is:
//
//	if errors.Is(err, units.ErrUnknownUnit) { … }
//
// Every error the package returns matches exactly one sentinel, and keeps its
// own, more detailed message. To branch on the kind, errors.AsType also finds
// the sentinel:
//
//	if kind, ok := errors.AsType[units.Error](err); ok {
//		switch kind { … }
//	}
//
// The sentinels are constants, not variables, so no importer can reassign them.
type Error string

// Error returns the sentinel's message.
func (e Error) Error() string {
	return string(e)
}

// The sentinel errors. Use errors.Is to test for them.
const (
	// ErrUnknownUnit: the unit symbol doesn't belong to the measure, as in
	// StringToDistance("1 kg") or UnitOf, or, for NewMeasurement and
	// ParseString, to any predefined measure. These errors are *UnitError
	// values, which carry the measure and the unit.
	ErrUnknownUnit Error = "unknown unit"
	// ErrIncompatibleMeasures: the measurements belong to different measures, as
	// in ToDistance on a mass or AddMeasurement of a distance and a mass, or one
	// of them is nil. Two user-defined measures with the same name but different
	// unit sets are different measures.
	ErrIncompatibleMeasures Error = "incompatible measures"
	// ErrZeroDivision: the divisor is zero, as in Div(0) or DivideByTime with a
	// zero time.
	ErrZeroDivision Error = "division by zero"
	// ErrInvalidFormat: the string isn't of the form "<value> <unit>" with a
	// finite value. If the value doesn't parse, the error also wraps the
	// *strconv.NumError.
	ErrInvalidFormat Error = "invalid format"
	// ErrInvalidFactors: the conversion factors are empty, have an empty unit
	// symbol or a factor that isn't positive and finite, have two symbols that
	// differ only in the micro sign, or lack the base unit.
	ErrInvalidFactors Error = "invalid conversion factors"
	// ErrInvalidMeasure: NewUnitSet got an empty or predefined measure name, or
	// a user-defined measure's Units method returns the zero UnitSet.
	ErrInvalidMeasure Error = "invalid measure"
	// ErrOutOfRange: the value can't be represented in the target type, as in
	// Time.ToDuration of NaN or of more than about 292 years, or MarshalText of
	// a value that isn't finite.
	ErrOutOfRange Error = "value out of range"
	// ErrAmbiguousUnit: a case-insensitive Parser got a unit that matches
	// several symbols differing only in case, such as "KN" for kN and kn.
	ErrAmbiguousUnit Error = "ambiguous unit"
	// ErrInvalidOption: NewParser got an invalid option, such as an alias that
	// doesn't denote a unit, or two measures that share a symbol.
	ErrInvalidOption Error = "invalid option"
)

// UnitError reports a unit that doesn't belong to a measure. It matches
// ErrUnknownUnit:
//
//	if ue, ok := errors.AsType[*units.UnitError](err); ok {
//		fmt.Println(ue.Measure, ue.Unit) // distance kg
//	}
type UnitError struct {
	// Measure is the measure the unit was expected to belong to. It is empty
	// when no predefined measure has the unit (NewMeasurement, ParseString).
	Measure MeasureType
	// Unit is the unit symbol.
	Unit string
}

// Error returns "unknown <measure> unit <unit>", or "unknown unit <unit>" when
// Measure is empty.
func (e *UnitError) Error() string {
	// Built here rather than when the error is created: returning a UnitError
	// then costs one allocation.
	if e.Measure == "" {
		return "unknown unit " + e.Unit
	}

	return "unknown " + string(e.Measure) + " unit " + e.Unit
}

// Unwrap returns ErrUnknownUnit.
func (*UnitError) Unwrap() error {
	return ErrUnknownUnit
}

// incompatibleError is an ErrIncompatibleMeasures with a fixed message of its
// own. Its values are constants, so returning one doesn't allocate.
type incompatibleError string

const (
	errNilMeasurement incompatibleError = "measurement is nil"
	errNoUnits        incompatibleError = "the measurement's measure has no units: its Units method returns the zero UnitSet"
)

func (e incompatibleError) Error() string {
	return string(e)
}

func (incompatibleError) Unwrap() error {
	return ErrIncompatibleMeasures
}

// kindError is an error with a message of its own, of the kind given by a
// sentinel. It may also wrap the error that caused it.
type kindError struct {
	msg  string
	errs [2]error // the sentinel, then the cause or nil
}

// newError returns an error of the given kind with message msg, wrapping cause
// unless it is nil.
func newError(kind Error, msg string, cause error) error {
	return &kindError{msg: msg, errs: [2]error{kind, cause}}
}

func (e *kindError) Error() string {
	return e.msg
}

// Unwrap returns the sentinel and, if there is one, the cause. The slice
// refers to the error's own array, so errors.Is and errors.As don't allocate.
func (e *kindError) Unwrap() []error {
	if e.errs[1] == nil {
		return e.errs[:1]
	}

	return e.errs[:]
}

// unknownUnitError returns the error for a unit that doesn't belong to measure
// mt, or to any predefined measure if mt is empty.
func unknownUnitError(mt MeasureType, unit string) error {
	return &UnitError{Measure: mt, Unit: unit}
}
