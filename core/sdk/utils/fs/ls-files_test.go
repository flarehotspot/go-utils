package fs

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

const d = "/tmp/lsdir"

func TestLsFiles(t *testing.T) {
  // testing none-recursive option

	f := filepath.Join(d, "sample.txt")
	if err := os.MkdirAll(d, os.ModePerm); err != nil {
		t.Errorf("Unable to create directory for testing: %s", d)
	}
	_, err := os.Create(f)
	if err != nil {
		t.Errorf("Unable to create test file: %s", f)
	}

	res, err := LsFiles(d, false)
	if err != nil {
		t.Errorf("LsDir (non-recursive): %s", err.Error())
	}

	expected := []string{f}
	if !reflect.DeepEqual(res, expected) {
		t.Error("Result is not right:\n", "Expected:", expected, "\nResult:", res)
	}

  // testing recursive option

	d2 := filepath.Join(d, "dir2")
	f2 := filepath.Join(d2, "sample2.txt")

	if err := os.MkdirAll(d2, os.ModePerm); err != nil {
		t.Errorf("Unable to create directory for testing: %s", d)
	}
	_, err = os.Create(f2)
	if err != nil {
		t.Errorf("Unable to create test file: %s", f2)
	}

	res, err = LsFiles(d, true)
	if err != nil {
		t.Errorf("Error in LsDir: %s", err.Error())
	}

	expected = []string{f2, f}
	if !reflect.DeepEqual(res, expected) {
		t.Error("LsDir (recursive):\n", "Expected:", expected, "\nResult:", res)
	}
}