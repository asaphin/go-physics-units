// Package datasize defines the data size units of the units package
// and their conversion factors to the base unit, the byte.
package datasize

// Unit is a data size unit: one of the constants below. String returns its
// symbol and ParseUnit finds the unit of a symbol.
//
// The zero Unit is the base unit, Byte. The numbers behind the constants may
// change between versions: store and exchange units by their symbols.
type Unit uint8

// BaseUnit is the unit of data size that all data size factors are relative to.
const BaseUnit = Byte

// Data size units, with their symbols.
const (
	Byte     Unit = iota // B: an octet
	Bit                  // bit
	Kilobit              // kbit
	Kilobyte             // kB: decimal (SI) prefix: 1000 bytes
	Kibibyte             // KiB: binary (IEC) prefix: 1024 bytes
	Megabit              // Mbit
	Megabyte             // MB
	Mebibyte             // MiB
	Gigabit              // Gbit
	Gigabyte             // GB
	Gibibyte             // GiB
	Terabyte             // TB
	Tebibyte             // TiB
	Petabyte             // PB
	Pebibyte             // PiB

	unitCount // the number of units; keep it last
)

// definitions holds the symbol of each unit and its factor: how many base
// units (B) there are in one unit.
//
// The factors are exact: defined values, written as constant expressions that Go
// evaluates exactly before rounding once.
//
// Source: IEC 80000-13: SI prefixes are powers of 1000, binary prefixes (Ki,
// Mi, …) powers of 1024.
var definitions = [unitCount]definition{
	Byte:     {"B", 1},
	Bit:      {"bit", 0.125},
	Kilobit:  {"kbit", 125},
	Kilobyte: {"kB", 1e3},
	Kibibyte: {"KiB", 1 << 10},
	Megabit:  {"Mbit", 125e3},
	Megabyte: {"MB", 1e6},
	Mebibyte: {"MiB", 1 << 20},
	Gigabit:  {"Gbit", 125e6},
	Gigabyte: {"GB", 1e9},
	Gibibyte: {"GiB", 1 << 30},
	Terabyte: {"TB", 1e12},
	Tebibyte: {"TiB", 1 << 40},
	Petabyte: {"PB", 1e15},
	Pebibyte: {"PiB", 1 << 50},
}
