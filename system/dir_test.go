package system

import (
	"io/ioutil"
	"log"
	"os"
	"testing"
)

var (
	testDir       string
	testTmpDir    = "" // "/tmp"
	testTmpPrefix = "test"
)

func TestMain(m *testing.M) {
	dir, err := ioutil.TempDir(testTmpDir, testTmpPrefix)
	if err != nil {
		log.Fatal(err)
	}
	testDir = dir
	// fmt.Printf("Using test dir: %s\n", dir)
	exitCode := m.Run()
	// defer os.RemoveAll(dir)
	if err := os.RemoveAll(dir); err != nil {
		log.Fatal(err)
	}
	os.Exit(exitCode)
}

func TestCheckDirExist(t *testing.T) {
	// Create test dir
	if !Exists(testDir) {
		if err := os.Mkdir(testDir, DirMode); err != nil {
			t.Fatal(err)
		}
	}
	// Check dir present
	if err := CheckDir(testDir); err != ErrDirAlreadyExist {
		t.Fatalf("CheckDir (DirAlreadyExist) %s: %s", testDir, err)
	}
}

func TestCheckDirNotExist(t *testing.T) {
	// Remove test dir
	if Exists(testDir) {
		if err := os.Remove(testDir); err != nil {
			t.Fatal(err)
		}
	}
	// Check dir absent
	if err := CheckDir(testDir); err != nil {
		t.Fatalf("CheckDir %s: %s", testDir, err)
	}
}

func TestCreateDir(t *testing.T) {
	if Exists(testDir) {
		if err := os.Remove(testDir); err != nil {
			t.Fatal(err)
		}
	}
	if err := CreateDir(testDir); err != nil {
		t.Fatal(err)
	}
	if !Exists(testDir) {
		t.Fatalf("TestCreateDir %s: failed", testDir)
	}
}

func TestRemoveDir(t *testing.T) {
	if !Exists(testDir) {
		if err := os.Mkdir(testDir, DirMode); err != nil {
			t.Fatal(err)
		}
	}
	if err := RemoveDir(testDir); err != nil {
		t.Fatal(err)
	}
	if Exists(testDir) {
		t.Fatalf("TestRemoveDir %s: failed", testDir)
	}
}

func TestIsEmptyDir(t *testing.T) {
	if !Exists(testDir) {
		if err := os.Mkdir(testDir, DirMode); err != nil {
			t.Fatal(err)
		}
	}
	empty, err := IsEmptyDir(testDir)
	if err != nil {
		t.Fatal(err)
	}
	if !empty {
		t.Fatalf("TestIsEmptyDir %s: failed", testDir)
	}
}