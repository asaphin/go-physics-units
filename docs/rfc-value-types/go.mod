// Prototype for docs/rfc-value-types.md (ROADMAP 2.1). A separate module, so
// it never becomes part of the library's build, API or coverage. Benchmark it
// with: cd docs/rfc-value-types && go test -bench . -benchmem
module github.com/asaphin/go-physics-units/docs/rfc-value-types

go 1.27.0

require github.com/asaphin/go-physics-units v0.0.0

replace github.com/asaphin/go-physics-units => ../..
