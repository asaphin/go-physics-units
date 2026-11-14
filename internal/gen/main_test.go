package main

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

// TestGenerated fails when a generated file differs from what the generator
// produces now. Run go generate ./... to update them.
func TestGenerated(t *testing.T) {
	files, err := generate()
	if err != nil {
		t.Fatal(err)
	}

	for path, want := range files {
		got, err := os.ReadFile(filepath.Join("..", "..", path))
		if err != nil {
			t.Errorf("%s: %v", path, err)

			continue
		}

		if !bytes.Equal(got, want) {
			t.Errorf("%s is out of date: run go generate ./...", path)
		}
	}
}
