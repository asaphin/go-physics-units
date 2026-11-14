package units_test

import (
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

// sentinelName matches the name of a sentinel error, such as ErrUnknownUnit,
// with or without a package qualifier and doc-link brackets.
var sentinelName = regexp.MustCompile(`\bErr[A-Z][A-Za-z]+\b`)

// TestDocContracts checks that the doc comment of every exported function and
// method in the module says how it fails or panics: those that return an
// error name a sentinel error ("It fails with ErrZeroDivision if …") or point
// to a function whose failures they share ("It fails like StringToDistance"),
// and those whose name starts with Must say that they panic. A user, or an AI
// agent, can then handle every error without reading the code.
func TestDocContracts(t *testing.T) {
	fset := token.NewFileSet()
	checked := 0

	err := filepath.WalkDir(".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() && path != "." &&
			(strings.HasPrefix(d.Name(), ".") || d.Name() == "internal" || d.Name() == "bin") {
			return filepath.SkipDir
		}

		if d.IsDir() || !strings.HasSuffix(path, ".go") || strings.HasSuffix(path, "_test.go") {
			return nil
		}

		file, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
		if err != nil {
			return err
		}

		for _, decl := range file.Decls {
			if fn, ok := decl.(*ast.FuncDecl); ok && fn.Name.IsExported() {
				checkContract(t, fset, fn)
				checked++
			}
		}

		return nil
	})
	if err != nil {
		t.Fatal(err)
	}

	if checked < 1000 {
		t.Errorf("checked only %d exported functions; the walk misses files", checked)
	}
}

func checkContract(t *testing.T, fset *token.FileSet, fn *ast.FuncDecl) {
	t.Helper()

	name := fn.Name.Name
	if fn.Recv != nil && len(fn.Recv.List) == 1 {
		recv := typeName(fn.Recv.List[0].Type)
		if !ast.IsExported(recv) || fn.Name.Name == "Unwrap" {
			return // not part of the API, or it returns the wrapped error rather than failing
		}

		name = recv + "." + name
	}

	doc := strings.Join(strings.Fields(fn.Doc.Text()), " ")
	where := fset.Position(fn.Pos()).String()

	if strings.HasPrefix(fn.Name.Name, "Must") && !strings.Contains(doc, "panic") {
		t.Errorf("%s: %s doesn't say when it panics", where, name)
	}

	if !returnsError(fn) {
		return
	}

	if !strings.Contains(doc, "fail") || (!sentinelName.MatchString(doc) && !strings.Contains(doc, "fails like")) {
		t.Errorf("%s: %s returns an error, but its doc doesn't name the sentinel it fails with", where, name)
	}
}

func returnsError(fn *ast.FuncDecl) bool {
	if fn.Type.Results == nil {
		return false
	}

	for _, r := range fn.Type.Results.List {
		if id, ok := r.Type.(*ast.Ident); ok && id.Name == "error" {
			return true
		}
	}

	return false
}

func typeName(expr ast.Expr) string {
	switch e := expr.(type) {
	case *ast.StarExpr:
		return typeName(e.X)
	case *ast.IndexExpr:
		return typeName(e.X)
	case *ast.Ident:
		return e.Name
	default:
		return "?"
	}
}
