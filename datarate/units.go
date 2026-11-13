// Package datarate defines the data rate unit symbols supported by the units package
// and their conversion factors to the base unit, the byte per second.
package datarate

import (
	"github.com/asaphin/go-physics-units/conversion"
	"github.com/asaphin/go-physics-units/internal/rates"
)

// BaseUnit is the unit of data rate that all data rate factors are relative to.
const BaseUnit = BytePerSecond

// Data rate unit symbols.
const (
	BitPerSecond      = "bit/s"
	BytePerSecond     = "B/s"
	KilobitPerSecond  = "kbit/s"
	KilobytePerSecond = "kB/s"
	KibibytePerSecond = "KiB/s"
	MegabitPerSecond  = "Mbit/s"
	MegabytePerSecond = "MB/s"
	MebibytePerSecond = "MiB/s"
	GigabitPerSecond  = "Gbit/s"
	GigabytePerSecond = "GB/s"
	GibibytePerSecond = "GiB/s"
)

// The factors are exact: defined values, written as constant expressions that Go
// evaluates exactly before rounding once.
//
// Source: IEC 80000-13, as for data sizes.
var conversionFactors = conversion.Factors{
	BitPerSecond:      0.125,
	BytePerSecond:     1,
	KilobitPerSecond:  125,
	KilobytePerSecond: 1e3,
	KibibytePerSecond: 1 << 10,
	MegabitPerSecond:  125e3,
	MegabytePerSecond: 1e6,
	MebibytePerSecond: 1 << 20,
	GigabitPerSecond:  125e6,
	GigabytePerSecond: 1e9,
	GibibytePerSecond: 1 << 30,
}

// ConversionFactors returns a copy of the data rate conversion factors: how many
// base units (B/s) there are in one unit of each symbol.
func ConversionFactors() conversion.Factors {
	return conversion.CopyConversionFactors(conversionFactors)
}

// Rate returns the multiplier that converts a value in unit from into unit to,
// and whether both are data rate units. Both forms of the micro sign are accepted.
func Rate(from, to string) (float64, bool) {
	return rates.FactorRate(conversionFactors, from, to)
}
