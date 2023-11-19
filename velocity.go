package units

import (
	"errors"
	"fmt"
	"github.com/asaphin/go-physics-units/distance"
	"github.com/asaphin/go-physics-units/time"
	"github.com/asaphin/go-physics-units/velocity"
)

var velocityConversionFactors = newImmutableConversionFactors(velocity.ConversionFactors())

type Velocity struct {
	BaseMeasurement
}

func NewVelocity(value float64, unit string) (*Velocity, error) {
	if _, ok := velocityConversionFactors.HasFactor(unit); !ok {
		return nil, fmt.Errorf("unknown Velocity unit %s", unit)
	}

	m, err := newBaseMeasurement(value, unit, velocityConversionFactors)
	if err != nil {
		return nil, err
	}

	return &Velocity{*m}, nil
}

var velocityConversionErr = errors.New("not a velocity measure")

func ToVelocity(m Measurement) (*Velocity, error) {
	if m.Type() == MeasureVelocity {
		return m.(*Velocity), nil
	}

	return nil, velocityConversionErr
}

func NewVelocityFromDistanceAndTime(d *Distance, t *Time) *Velocity {
	newD, _ := d.convertTo(distance.BaseUnit)
	newT, _ := t.convertTo(time.BaseUnit)

	d, _ = ToDistance(newD)
	t, _ = ToTime(newT)

	v, _ := NewVelocity(d.Value()/t.Value(), velocity.BaseUnit)

	return v
}
