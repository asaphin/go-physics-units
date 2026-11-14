package units

import "github.com/asaphin/go-physics-units/magneticfield"

// MagneticField is a magnetic flux density: a value in a magneticfield.Unit,
// such as 50 μT. Create one with NewMagneticField or StringToMagneticField.
//
// A MagneticField is a small immutable value: every method returns a new value,
// so copy and share it freely. The zero MagneticField is 0 T.
type MagneticField struct {
	value float64
	unit  magneticfield.Unit
}
