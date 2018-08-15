package system

import (
	"os"
	"testing"
)

func TestExists(t *testing.T) {
	// Create test dir
	if !Exists(testDir) {
		if err := os.Mkdir(testDir, DirMode); err != nil {
			t.Fatal(err)
		}
	}
	// Check dir present
	if !Exists(testDir) {
		t.Fatalf("%s: should exist", testDir)
	}
}

func TestIsDir(t *testing.T) {
	// Create test dir
	if !Exists(testDir) {
		if err := os.Mkdir(testDir, DirMode); err != nil {
			t.Fatal(err)
		}
	}
	// Check is dir
	ok, err := IsDir(testDir)
	if err != nil {
		t.Fatal(err)
	}
	if !ok {
		t.Fatalf("%s: should be a directory", testDir)
	}
}

// func TestRemove(t *testing.T) {
// 	// TODO no store.Delete
// }
