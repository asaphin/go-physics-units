package units_test

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
	"path/filepath"
	"reflect"
	"strings"
	"sync"
	"testing"

	units "github.com/asaphin/go-physics-units"
	"github.com/asaphin/go-physics-units/area"
	"github.com/asaphin/go-physics-units/distance"
	"github.com/asaphin/go-physics-units/mass"
	"github.com/asaphin/go-physics-units/time"
	"github.com/asaphin/go-physics-units/velocity"
)

// TestNoExportedPackageVariables guards immutability after init: any importer
// can reassign an exported package-level variable, so the library has none.
// Use constants, or functions that return copies.
func TestNoExportedPackageVariables(t *testing.T) {
	fset := token.NewFileSet()

	err := filepath.WalkDir(".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			if path != "." && (strings.HasPrefix(d.Name(), ".") || d.Name() == "testdata" || d.Name() == "coverage") {
				return filepath.SkipDir
			}

			return nil
		}

		if !strings.HasSuffix(path, ".go") || strings.HasSuffix(path, "_test.go") {
			return nil
		}

		file, err := parser.ParseFile(fset, path, nil, parser.SkipObjectResolution)
		if err != nil {
			return fmt.Errorf("parse %s: %w", path, err)
		}

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

		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
}

// TestConversionRates_CannotModifySharedTable guards the package-level rate
// tables. ConversionRates must return a value with no exported fields: a
// pointer would let callers overwrite the shared table with *p = other.
func TestConversionRates_CannotModifySharedTable(t *testing.T) {
	tables := map[string]any{
		"area":     area.ConversionRates(),
		"distance": distance.ConversionRates(),
		"mass":     mass.ConversionRates(),
		"time":     time.ConversionRates(),
		"velocity": velocity.ConversionRates(),
	}

	for name, table := range tables {
		typ := reflect.TypeOf(table)
		if typ.Kind() != reflect.Struct {
			t.Errorf("%s.ConversionRates() returns a %s, want a struct value", name, typ.Kind())

			continue
		}

		for i := range typ.NumField() {
			if f := typ.Field(i); f.IsExported() {
				t.Errorf("%s.ConversionRates() result has exported field %s", name, f.Name)
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

	_ = distance.ConversionRates()

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
