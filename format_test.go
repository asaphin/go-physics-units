package units_test

import (
	"fmt"
	"testing"

	units "github.com/asaphin/go-physics-units"
	"github.com/asaphin/go-physics-units/distance"
)

func TestFormat(t *testing.T) {
	d := units.NewDistance(5, distance.Kilometer)

	m, err := units.ParseString("1234.5678 m")
	if err != nil {
		t.Fatal(err)
	}

	cases := []struct {
		format string
		arg    any
		want   string
	}{
		{"%v", d, "5 km"},
		{"%+v", d, "5 km"},
		{"%s", d, "5 km"},
		{"%10v|", d, "      5 km|"},
		{"%-10s|", d, "5 km      |"},
		{"%q", d, `"5 km"`},
		{"%f", d, "5.000000 km"},
		{"%.2f", d, "5.00 km"},
		{"%8.2f|", d, "    5.00 km|"},
		{"%+.1f", d, "+5.0 km"},
		{"%e", d, "5.000000e+00 km"},
		{"%x", d, "0x1.4p+02 km"},
		{"%.3g", m, "1.23e+03 m"},
		{"%.2f", m, "1234.57 m"},
		{"%#v", d, `units.Distance{Value:5, Unit:"km"}`},
		{"%#v", m, `units.Distance{Value:1234.5678, Unit:"m"}`},
		{"%d", d, "%!d(units.Distance=5 km)"},
		{"%v", []units.Distance{d, d}, "[5 km 5 km]"},
	}

	for _, c := range cases {
		if got := fmt.Sprintf(c.format, c.arg); got != c.want {
			t.Errorf("Sprintf(%q, %v) = %q, want %q", c.format, c.arg, got, c.want)
		}
	}

	if got := fmt.Sprint(d, m); got != "5 km 1234.5678 m" {
		t.Errorf("Sprint = %q, want %q", got, "5 km 1234.5678 m")
	}
}
