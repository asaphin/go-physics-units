// Package prototype compares value-type designs for go-physics-units with the
// current interface-based one. See docs/rfc-value-types.md.
//
//   - Option S: a Distance struct with typed string units (strunits.go).
//   - Option I: a Distance struct with typed integer units that index immutable
//     arrays (intunits.go).
//   - Quantity[M]: a generic value type for user-defined measures, whose units
//     are an immutable UnitSet supplied by a marker type (quantity.go).
package prototype
