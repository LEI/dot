package dot

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

var (
	testDir string
)

func setupTestCase(t *testing.T) func(t *testing.T) {
	// t.Log("setup test case")
	var err error
	testDir, err = ioutil.TempDir("", "test")
	if err != nil {
		t.Fatal(err)
	}
	return func(t *testing.T) {
		// t.Log("teardown test case")
		if err := os.RemoveAll(testDir); err != nil {
			t.Fatal(err)
		}
	}
}

func TestDoDir(t *testing.T) {
	teardownTestCase := setupTestCase(t)
	defer teardownTestCase(t)

	cases := []struct {
		target string
		in     string
		out    string
	}{
		{testDir, "a", "a"},
	}

	for _, tc := range cases {
		// tc.in = filepath.Join(tc.target, tc.in)
		tc.out = filepath.Join(tc.target, tc.out)
		r := &Role{
			Dirs: []*Dir{
				&Dir{Path: tc.in},
			},
		}
		if err := r.ParseDirs(tc.target); err != nil {
			t.Fatal(err)
		}
		// d := r.Dirs[0]
		for _, d := range r.Dirs {
			d.SetAction("install")
			err := d.Status()
			ok := IsExist(err)
			if !ok && err != nil {
				t.Fatal(err)
			}
			if ok { // t.Log("Already present")
				continue
			}
			// if err := d.Check(); err != nil {
			// 	t.Fatal(err)
			// }
			fmt.Println(d.DoString())
			if err := d.Do(); err != nil {
				t.Fatal(err)
			}
		}
		if !exists(tc.out) {
			t.Fatalf("%s: directory should exist", tc.out)
		}
	}
}

func TestUndoDir(t *testing.T) {
	teardownTestCase := setupTestCase(t)
	defer teardownTestCase(t)

	cases := []struct {
		target string
		in     string
		out    string
	}{
		{testDir, "a", "a"},
	}

	for _, tc := range cases {
		// tc.in = filepath.Join(testDir, tc.in)
		tc.out = filepath.Join(testDir, tc.out)
		r := &Role{
			Dirs: []*Dir{
				&Dir{Path: tc.in},
			},
		}
		if err := r.ParseDirs(tc.target); err != nil {
			t.Fatal(err)
		}
		for _, d := range r.Dirs {
			d.SetAction("remove")
			err := d.Status()
			ok := IsExist(err)
			if !ok && err != nil {
				t.Fatal(err)
			}
			if !ok { // t.Log("Already absent")
				continue
			}
			// if err := d.Check(); err != nil {
			// 	t.Fatal(err)
			// }
			fmt.Println(d.UndoString())
			if err := d.Undo(); err != nil {
				t.Fatal(err)
			}
		}
		if exists(tc.out) {
			t.Fatalf("%s: directory should not exist", tc.out)
		}
	}
}
