// Package datasize defines the data size unit symbols supported by the units package
// and their conversion factors to the base unit, the byte.
package datasize

import (
	"github.com/asaphin/go-physics-units/conversion"
	"github.com/asaphin/go-physics-units/internal/rates"
)

// BaseUnit is the unit of data size that all data size factors are relative to.
const BaseUnit = Byte

// Data size unit symbols.
const (
	Bit      = "bit"
	Byte     = "B" // an octet
	Kilobit  = "kbit"
	Kilobyte = "kB"  // decimal (SI) prefix: 1000 bytes
	Kibibyte = "KiB" // binary (IEC) prefix: 1024 bytes
	Megabit  = "Mbit"
	Megabyte = "MB"
	Mebibyte = "MiB"
	Gigabit  = "Gbit"
	Gigabyte = "GB"
	Gibibyte = "GiB"
	Terabyte = "TB"
	Tebibyte = "TiB"
	Petabyte = "PB"
	Pebibyte = "PiB"
)

// The factors are exact: defined values, written as constant expressions that Go
// evaluates exactly before rounding once.
//
// Source: IEC 80000-13: SI prefixes are powers of 1000, binary prefixes (Ki,
// Mi, …) powers of 1024.
var conversionFactors = conversion.Factors{
	Bit:      0.125,
	Byte:     1,
	Kilobit:  125,
	Kilobyte: 1e3,
	Kibibyte: 1 << 10,
	Megabit:  125e3,
	Megabyte: 1e6,
	Mebibyte: 1 << 20,
	Gigabit:  125e6,
	Gigabyte: 1e9,
	Gibibyte: 1 << 30,
	Terabyte: 1e12,
	Tebibyte: 1 << 40,
	Petabyte: 1e15,
	Pebibyte: 1 << 50,
}

// ConversionFactors returns a copy of the data size conversion factors: how many
// base units (B) there are in one unit of each symbol.
func ConversionFactors() conversion.Factors {
	return conversion.CopyConversionFactors(conversionFactors)
}

// Rate returns the multiplier that converts a value in unit from into unit to,
// and whether both are data size units. Both forms of the micro sign are accepted.
func Rate(from, to string) (float64, bool) {
	return rates.FactorRate(conversionFactors, from, to)
}
