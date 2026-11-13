package units

import (
	"cmp"

	"github.com/asaphin/go-physics-units/area"
	"github.com/asaphin/go-physics-units/distance"
	"github.com/asaphin/go-physics-units/volume"
)

// Volume is a volume measurement. Its unit symbols are defined in the volume
// package.
type Volume interface {
	Measurement
	UnitConverter[Volume]
	Arithmetics[Volume]
	Comparer[Volume]

	// DivideByArea returns the distance v/a in distance.BaseUnit (m).
	// It fails if a is zero.
	DivideByArea(a Area) (Distance, error)
	// DivideByDistance returns the area v/d in area.BaseUnit (m²).
	// It fails if d is zero.
	DivideByDistance(d Distance) (Area, error)
}

// volumeImplementation is the concrete implementation of Volume.
type volumeImplementation struct {
	baseMeasurement
}

// ConvertTo returns the volume expressed in unit.
func (v *volumeImplementation) ConvertTo(unit string) (Volume, error) {
	msm, err := v.convertTo(unit)
	if err != nil {
		return nil, err
	}

	return &volumeImplementation{msm}, nil
}

// MustConvertTo is like ConvertTo but panics if unit isn't a volume unit.
func (v *volumeImplementation) MustConvertTo(unit string) Volume {
	msm, err := v.convertTo(unit)
	if err != nil {
		panic(err)
	}

	return &volumeImplementation{msm}
}

// ConvertToBaseUnits returns the volume expressed in cubic metres.
func (v *volumeImplementation) ConvertToBaseUnits() Volume {
	return &volumeImplementation{v.unsafeConvertTo(volume.BaseUnit)}
}

func (v *volumeImplementation) valueInBaseUnits() float64 {
	return v.unsafeGetValueIn(volume.BaseUnit)
}

func (v *volumeImplementation) valueInUnits(unit string) float64 {
	return v.unsafeGetValueIn(unit)
}

// Add returns the sum of the volumes, expressed in the receiver's unit.
func (v *volumeImplementation) Add(other Volume) Volume {
	return &volumeImplementation{v.withValue(v.value + other.valueInUnits(v.unit))}
}

// Sub returns the difference of the volumes, expressed in the receiver's unit.
func (v *volumeImplementation) Sub(other Volume) Volume {
	return &volumeImplementation{v.withValue(v.value - other.valueInUnits(v.unit))}
}

// Mul returns the volume multiplied by multiplier.
func (v *volumeImplementation) Mul(multiplier float64) Volume {
	return &volumeImplementation{v.withValue(v.value * multiplier)}
}

// Div returns the volume divided by divisor. It fails if divisor is zero.
func (v *volumeImplementation) Div(divisor float64) (Volume, error) {
	if divisor == 0 {
		return nil, ErrZeroDivision
	}

	return &volumeImplementation{v.withValue(v.value / divisor)}, nil
}

// DivideByArea returns the distance v/a in m, or an error if a is zero.
func (v *volumeImplementation) DivideByArea(a Area) (Distance, error) {
	if a.Value() == 0 {
		return nil, ErrZeroDivision
	}

	return newDistance(v.valueInBaseUnits()/a.valueInBaseUnits(), distance.BaseUnit), nil
}

// DivideByDistance returns the area v/d in m², or an error if d is zero.
func (v *volumeImplementation) DivideByDistance(d Distance) (Area, error) {
	if d.Value() == 0 {
		return nil, ErrZeroDivision
	}

	return newArea(v.valueInBaseUnits()/d.valueInBaseUnits(), area.BaseUnit), nil
}

// Compare compares the volumes; see Comparer.
func (v *volumeImplementation) Compare(other Volume) int {
	return cmp.Compare(v.commonValues(other.base()))
}

// Equal reports whether the volumes are equal; see Comparer.
func (v *volumeImplementation) Equal(other Volume) bool {
	return v.Compare(other) == 0
}

// EqualWithinRel reports whether the volumes are equal within the relative
// tolerance tol; see Comparer.
func (v *volumeImplementation) EqualWithinRel(other Volume, tol float64) bool {
	x, y := v.commonValues(other.base())

	return equalWithinRel(x, y, tol)
}

// NewVolume creates a volume of value in unit, one of the symbols defined in
// the volume package.
func NewVolume(value float64, unit string) (Volume, error) {
	msm, err := newBaseMeasurement(value, unit, volumeRates)
	if err != nil {
		return nil, err
	}

	return &volumeImplementation{msm}, nil
}

// newVolume is an unchecked constructor for internal use with valid, canonical
// units.
func newVolume(value float64, unit string) Volume {
	return &volumeImplementation{baseMeasurement{value: value, unit: unit, conversionRates: volumeRates}}
}

const errVolumeConversion incompatibleError = "not a volume measure"

// ToVolume converts m, for example a result of ParseString, into a Volume.
// It fails if m isn't a volume.
func ToVolume(m Measurement) (Volume, error) {
	if m == nil || m.Type() != MeasureVolume {
		return nil, errVolumeConversion
	}

	return &volumeImplementation{*m.base()}, nil
}

// StringToVolume parses a volume from a string such as "10 L".
// It fails if the unit isn't a volume unit.
func StringToVolume(s string) (Volume, error) {
	// Parse and construct directly: one allocation, where ParseString
	// followed by ToVolume would allocate twice.
	value, unit, err := parseValueAndUnit(s)
	if err != nil {
		return nil, err
	}

	return NewVolume(value, unit)
}
