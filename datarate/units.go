// Package datarate defines the data rate units of the units package
// and their conversion factors to the base unit, the byte per second.
package datarate

// Unit is a data rate unit: one of the constants below. String returns its
// symbol and ParseUnit finds the unit of a symbol.
//
// The zero Unit is the base unit, BytePerSecond. The numbers behind the constants may
// change between versions: store and exchange units by their symbols.
type Unit uint8

// BaseUnit is the unit of data rate that all data rate factors are relative to.
const BaseUnit = BytePerSecond

// Data rate units, with their symbols.
const (
	BytePerSecond     Unit = iota // B/s
	BitPerSecond                  // bit/s
	KilobitPerSecond              // kbit/s
	KilobytePerSecond             // kB/s
	KibibytePerSecond             // KiB/s
	MegabitPerSecond              // Mbit/s
	MegabytePerSecond             // MB/s
	MebibytePerSecond             // MiB/s
	GigabitPerSecond              // Gbit/s
	GigabytePerSecond             // GB/s
	GibibytePerSecond             // GiB/s

	unitCount // the number of units; keep it last
)

// definitions holds the symbol of each unit and its factor: how many base
// units (B/s) there are in one unit.
//
// The factors are exact: defined values, written as constant expressions that Go
// evaluates exactly before rounding once.
//
// Source: IEC 80000-13, as for data sizes.
var definitions = [unitCount]definition{
	BytePerSecond:     {"B/s", 1},
	BitPerSecond:      {"bit/s", 0.125},
	KilobitPerSecond:  {"kbit/s", 125},
	KilobytePerSecond: {"kB/s", 1e3},
	KibibytePerSecond: {"KiB/s", 1 << 10},
	MegabitPerSecond:  {"Mbit/s", 125e3},
	MegabytePerSecond: {"MB/s", 1e6},
	MebibytePerSecond: {"MiB/s", 1 << 20},
	GigabitPerSecond:  {"Gbit/s", 125e6},
	GigabytePerSecond: {"GB/s", 1e9},
	GibibytePerSecond: {"GiB/s", 1 << 30},
}
