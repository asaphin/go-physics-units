package units_test

import (
	"bytes"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"slices"
	"strings"
	"testing"
)

// TestUseCasesDoc checks that USECASES.md shows the scenarios of
// example_usecases_test.go verbatim: each declaration, and the body and
// output of each example. go test verifies the outputs, so the document can't
// show code that doesn't compile or output that the code doesn't print.
func TestUseCasesDoc(t *testing.T) {
	const docFile, srcFile = "USECASES.md", "example_usecases_test.go"

	doc, err := os.ReadFile(docFile)
	if err != nil {
		t.Fatal(err)
	}

	src, err := os.ReadFile(srcFile)
	if err != nil {
		t.Fatal(err)
	}

	fset := token.NewFileSet()

	file, err := parser.ParseFile(fset, srcFile, src, parser.ParseComments)
	if err != nil {
		t.Fatal(err)
	}

	goBlocks, textBlocks := fencedBlocks(string(doc), "go"), fencedBlocks(string(doc), "text")

	var decls, codes []string

	for _, decl := range file.Decls {
		if fn, ok := decl.(*ast.FuncDecl); ok && strings.HasPrefix(fn.Name.Name, "Example") {
			code, output := exampleParts(
				string(src[fset.Position(fn.Body.Lbrace).Offset+1 : fset.Position(fn.Body.Rbrace).Offset]),
			)
			codes = append(codes, code)

			if !slices.Contains(goBlocks, code) {
				t.Errorf("%s doesn't show the code of %s verbatim:\n%s", docFile, fn.Name.Name, code)
			}

			if !slices.Contains(textBlocks, output) {
				t.Errorf("%s doesn't show the output of %s verbatim:\n%s", docFile, fn.Name.Name, output)
			}

			continue
		}

		if gen, ok := decl.(*ast.GenDecl); ok && gen.Tok == token.IMPORT {
			continue
		}

		start := decl.Pos()
		if doc := declDoc(decl); doc != nil {
			start = doc.Pos()
		}

		// To the end of the line, with any trailing comment.
		end := fset.Position(decl.End()).Offset
		if i := bytes.IndexByte(src[end:], '\n'); i >= 0 {
			end += i
		}

		decls = append(decls, string(src[fset.Position(start).Offset:end]))
	}

	if len(codes) == 0 {
		t.Fatalf("%s has no examples", srcFile)
	}

	// Every declaration appears in a Go block, and every Go block is either an
	// example's code or declarations only, so the document has no stale code.
	for _, d := range decls {
		if !slices.ContainsFunc(goBlocks, func(b string) bool { return strings.Contains(b, d) }) {
			t.Errorf("%s doesn't show this declaration verbatim:\n%s", docFile, d)
		}
	}

	for _, b := range goBlocks {
		rest := b
		for _, d := range decls {
			rest = strings.ReplaceAll(rest, d, "")
		}

		if strings.TrimSpace(rest) != "" && !slices.Contains(codes, b) {
			t.Errorf("%s has a Go block that isn't in %s:\n%s", docFile, srcFile, b)
		}
	}
}

// fencedBlocks returns the contents of the fenced code blocks of lang in a
// Markdown document.
func fencedBlocks(doc, lang string) []string {
	var blocks []string

	var block strings.Builder

	in := false

	for line := range strings.Lines(doc) {
		switch {
		case !in && strings.TrimSuffix(line, "\n") == "```"+lang:
			in = true

			block.Reset()
		case in && strings.TrimSuffix(line, "\n") == "```":
			in = false
			blocks = append(blocks, strings.TrimSuffix(block.String(), "\n"))
		case in:
			block.WriteString(line)
		}
	}

	return blocks
}

// exampleParts splits the body of an example function into its code and its
// output, and returns them in that order, both without the function's
// indentation.
func exampleParts(body string) (string, string) {
	var lines []string
	for line := range strings.Lines(strings.Trim(body, "\n")) {
		lines = append(lines, strings.TrimPrefix(strings.TrimSuffix(line, "\n"), "\t"))
	}

	i := slices.Index(lines, "// Output:")
	if i < 0 {
		return strings.Join(lines, "\n"), ""
	}

	out := make([]string, 0, len(lines)-i-1)
	for _, line := range lines[i+1:] {
		out = append(out, strings.TrimPrefix(strings.TrimPrefix(line, "//"), " "))
	}

	return strings.TrimRight(strings.Join(lines[:i], "\n"), "\n"), strings.Join(out, "\n")
}

// declDoc returns the doc comment of a declaration, or nil.
func declDoc(decl ast.Decl) *ast.CommentGroup {
	switch d := decl.(type) {
	case *ast.FuncDecl:
		return d.Doc
	case *ast.GenDecl:
		return d.Doc
	default:
		return nil
	}
}
