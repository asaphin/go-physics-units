package units_test

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"testing"

	units "github.com/asaphin/go-physics-units"
	"github.com/asaphin/go-physics-units/distance"
)

// libraryFiles parses the module's non-test Go files.
func libraryFiles(t *testing.T) (*token.FileSet, []*ast.File) {
	t.Helper()

	fset := token.NewFileSet()

	var files []*ast.File

	err := filepath.WalkDir(".", func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			if p != "." && (strings.HasPrefix(d.Name(), ".") || d.Name() == "testdata" || d.Name() == "coverage") {
				return filepath.SkipDir
			}

			return nil
		}

		if !strings.HasSuffix(p, ".go") || strings.HasSuffix(p, "_test.go") {
			return nil
		}

		file, err := parser.ParseFile(fset, p, nil, parser.SkipObjectResolution)
		if err != nil {
			return fmt.Errorf("parse %s: %w", p, err)
		}

		files = append(files, file)

		return nil
	})
	if err != nil {
		t.Fatal(err)
	}

	return fset, files
}

// TestNoExportedPackageVariables guards immutability after init: any importer
// can reassign an exported package-level variable, so the library has none.
// Use constants, or functions that return copies.
func TestNoExportedPackageVariables(t *testing.T) {
	fset, files := libraryFiles(t)

	for _, file := range files {
		for _, decl := range file.Decls {
			gen, ok := decl.(*ast.GenDecl)
			if !ok || gen.Tok != token.VAR {
				continue
			}

			for _, spec := range gen.Specs {
				for _, name := range spec.(*ast.ValueSpec).Names {
					if name.IsExported() {
						t.Errorf("%s: exported package variable %s can be reassigned by any importer",
							fset.Position(name.Pos()), name.Name)
					}
				}
			}
		}
	}
}

// TestNoInternalTypesInExportedAPI keeps the shared rate tables out of reach:
// no exported function, method, field or interface method may use a type from
// an internal package. Callers can't name such types, and a handle to shared
// state could let them modify it.
func TestNoInternalTypesInExportedAPI(t *testing.T) {
	fset, files := libraryFiles(t)

	for _, file := range files {
		if strings.Contains(filepath.ToSlash(fset.Position(file.Package).Filename), "internal/") {
			continue
		}

		internal := internalImports(file)
		if len(internal) == 0 {
			continue
		}

		check := func(name string, node ast.Node) {
			ast.Inspect(node, func(n ast.Node) bool {
				sel, ok := n.(*ast.SelectorExpr)
				if !ok {
					return true
				}

				if pkg, ok := sel.X.(*ast.Ident); ok && internal[pkg.Name] {
					t.Errorf("%s: exported %s uses internal type %s.%s",
						fset.Position(sel.Pos()), name, pkg.Name, sel.Sel.Name)
				}

				return true
			})
		}

		for _, decl := range file.Decls {
			switch decl := decl.(type) {
			case *ast.FuncDecl:
				if decl.Name.IsExported() {
					check(decl.Name.Name, decl.Type)
				}
			case *ast.GenDecl:
				for _, spec := range decl.Specs {
					if ts, ok := spec.(*ast.TypeSpec); ok && ts.Name.IsExported() {
						checkExportedType(ts, check)
					}
				}
			}
		}
	}
}

// internalImports returns the local names of the file's imports of internal
// packages.
func internalImports(file *ast.File) map[string]bool {
	names := make(map[string]bool)

	for _, imp := range file.Imports {
		importPath, err := strconv.Unquote(imp.Path.Value)
		if err != nil || !strings.Contains(importPath, "/internal/") {
			continue
		}

		name := path.Base(importPath)
		if imp.Name != nil {
			name = imp.Name.Name
		}

		names[name] = true
	}

	return names
}

// checkExportedType checks the parts of an exported type that callers can see:
// exported fields of a struct, exported methods of an interface, or the whole
// type otherwise.
func checkExportedType(ts *ast.TypeSpec, check func(name string, node ast.Node)) {
	var fields *ast.FieldList

	switch typ := ts.Type.(type) {
	case *ast.StructType:
		fields = typ.Fields
	case *ast.InterfaceType:
		fields = typ.Methods
	default:
		check(ts.Name.Name, ts.Type)

		return
	}

	for _, field := range fields.List {
		if len(field.Names) == 0 {
			check(ts.Name.Name, field.Type) // embedded

			continue
		}

		for _, name := range field.Names {
			if name.IsExported() {
				check(ts.Name.Name+"."+name.Name, field.Type)
			}
		}
	}
}

// TestConcurrentUse runs the public API from many goroutines. Under the race
// detector (make check, CI) it fails if any shared state is written after init.
func TestConcurrentUse(t *testing.T) {
	mc := newCurrencyCreator(t)

	var wg sync.WaitGroup

	for range 8 {
		wg.Add(1)

		go func() {
			defer wg.Done()

			for range 100 {
				useConcurrently(t, mc)
			}
		}()
	}

	wg.Wait()
}

func useConcurrently(t *testing.T, mc units.MeasurementCreator[Currency]) {
	t.Helper()

	m, err := units.ParseString("1 km")
	if err != nil {
		t.Error(err)

		return
	}

	d, err := units.ToDistance(m)
	if err != nil {
		t.Error(err)

		return
	}

	if got := d.MustConvertTo(distance.Millimeter).Value(); got != 1e6 {
		t.Errorf("1 km = %v mm, want 1e6", got)
	}

	if got := d.String(); got != "1 km" {
		t.Errorf("String() = %q, want 1 km", got)
	}

	// Callers own the factors they get back; changing them affects nothing shared.
	f := distance.ConversionFactors()
	f[distance.Millimeter] = 1

	if r, ok := distance.Rate(distance.Kilometer, distance.Meter); !ok || r != 1000 {
		t.Errorf("Rate(km, m) = %v, %v; want 1000, true", r, ok)
	}

	eur, err := mc.New(1, "EUR")
	if err != nil {
		t.Error(err)

		return
	}

	if _, err := eur.AddMeasurement(eur); err != nil {
		t.Error(err)
	}

	if _, err := units.NewBaseMeasurement(1, "EUR", currencyFactors()); err != nil {
		t.Error(err)
	}
}
