package system

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/LEI/dot/pkg/comp"
)

var copyTests = []struct {
	in  []byte
	out bool
}{
	{[]byte("abc"), true},
}

// func TestCheckCopy(t *testing.T) {
// 	// Create test dir
// 	if !Exists(testDir) {
// 		if err := os.Mkdir(testDir, DirMode); err != nil {
// 			t.Fatal(err)
// 		}
// 	}
// 	// Check dir present
// 	if err := CheckDir(testDir); err != ErrDirAlreadyExist {
// 		t.Fatalf("CheckDir (DirAlreadyExist) %s: %s", testDir, err)
// 	}
// }

func TestCopy(t *testing.T) {
	// Create test dir
	if !Exists(testDir) {
		if err := os.Mkdir(testDir, DirMode); err != nil {
			t.Fatal(err)
		}
	}
	for _, tt := range copyTests {
		// Create test file
		tmpFile, err := ioutil.TempFile(testDir, "src")
		if err != nil {
			t.Fatal(err)
		}
		defer os.Remove(tmpFile.Name()) // clean up
		if _, err := tmpFile.Write(tt.in); err != nil {
			t.Fatal(err)
		}
		if err := tmpFile.Close(); err != nil {
			t.Fatal(err)
		}
		// tmpDst, err := ioutil.TempFile(testDir, "dst")
		// if err != nil {
		// 	t.Fatal(err)
		// }
		// tmpDir, err := filepath.Abs(filepath.Dir(tmpFile.Name()))
		// if err != nil {
		// 	t.Fatal(err)
		// }
		src := tmpFile.Name()
		dst := filepath.Join(filepath.Dir(tmpFile.Name()), "dst")
		if err := Copy(src, dst); err != nil {
			t.Fatal(err)
		}
		defer os.Remove(dst)
		ok, err := comp.FileEquals(src, dst)
		if err != nil {
			t.Fatal(err)
		}
		if ok != tt.out { // !ok
			t.Fatalf("%s: not equal to %s", src, dst)
		}
	}
}

func TestCopyExist(t *testing.T) {
	// Create test dir
	if !Exists(testDir) {
		if err := os.Mkdir(testDir, DirMode); err != nil {
			t.Fatal(err)
		}
	}
	// Create test file
	content := []byte("abc")
	tmpFile, err := ioutil.TempFile(testDir, "src")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name()) // clean up
	if _, err := tmpFile.Write(content); err != nil {
		t.Fatal(err)
	}
	if err := tmpFile.Close(); err != nil {
		t.Fatal(err)
	}
	tmpDst, err := ioutil.TempFile(filepath.Dir(tmpFile.Name()), "dst")
	if err != nil {
		t.Fatal(err)
	}
	src := tmpFile.Name()
	dst := tmpDst.Name()

	ok, err := comp.FileEquals(src, dst)
	if err != nil {
		t.Fatal(err)
	}
	if ok && string(content) != "" {
		t.Fatalf("%s: content should not be equal to %s", src, dst)
	}

	if err := Copy(src, dst); err != nil {
		t.Fatal(err)
	}
	defer os.Remove(dst)
	ok, err = comp.FileEquals(src, dst)
	if err != nil {
		t.Fatal(err)
	}
	if !ok {
		t.Fatalf("%s: content should be equal to %s", src, dst)
	}
}

// func TestRemoveDir(t *testing.T) {
// 	if !Exists(testDir) {
// 		if err := os.Mkdir(testDir, DirMode); err != nil {
// 			t.Fatal(err)
// 		}
// 	}
// 	if err := RemoveDir(testDir); err != nil {
// 		t.Fatal(err)
// 	}
// 	if Exists(testDir) {
// 		t.Fatalf("TestRemoveDir %s: failed", testDir)
// 	}
// }

// func TestIsEmptyDir(t *testing.T) {
// 	if !Exists(testDir) {
// 		if err := os.Mkdir(testDir, DirMode); err != nil {
// 			t.Fatal(err)
// 		}
// 	}
// 	empty, err := IsEmptyDir(testDir)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	if !empty {
// 		t.Fatalf("TestIsEmptyDir %s: failed", testDir)
// 	}
// }
