package dotfile

import (
	// "fmt"
	// "os"
	"testing"

	// "golang.org/x/tools/godoc/vfs"

	// "github.com/absfs/osfs"
	// "github.com/absfs/basefs"
)

// type mockedFS struct {
// 	osFS
// 	reportErr bool
// 	reportSize int64
// }

// type mockedFileInfo struct {
// 	os.FileInfo
// 	size int64
// }

// func (m mockedFileInfo) Size() int64 { return m.size }

// func (m mockedFS) Create(name string) (os.FileInfo, error) {
// 	return mockedFileInfo{}, nil
// }

// // golang.org/x/tools/godoc/vfs/mapfs
// // func (m mockedFS) Open(name string) (os.FileInfo, error) {
// // 	return mockedFileInfo{}, nil
// // }

// func (m mockedFS) Stat(name string) (os.FileInfo, error) {
// 	if m.reportErr {
// 		return nil, os.ErrNotExist
// 	}
// 	return mockedFileInfo{size: m.reportSize}, nil
// }

var copyTests = []struct{
	// in *CopyTask
	in string
	out string
}{
	// {&CopyTask{Source:"a", Target:"b"}, "c"},
	{"a", "c"},
}

// TestCopy ...
func TestCopy(t *testing.T) {
	// oldFs := fs
	// // Create and "install" mocked fs:
	// ofs := fs
	// ofs, err := osfs.NewFS()
	// if err != nil {
	// 	t.Errorf("osfs failed: %s", err)
	// }
	// bfs, err := basefs.NewFS(ofs, "/tmp")
	// if err != nil {
	// 	t.Errorf("osfs failed: %s", err)
	// }
	// fs = bfs
	// // Make sure fs is restored after this test:
	// defer func() {
	// 	fs = ofs
	// }()

	/*
	// Test when filesystem.Stat() reports error:
	mfs.reportErr = true
	if _, err := getSize("hello.go"); err == nil {
		t.Error("Expected error, but err is nil!")
	}

	// Test when no error and size is returned:
	mfs.reportErr = false
	mfs.reportSize = 123
	if size, err := getSize("hello.go"); err != nil {
		t.Errorf("Expected no error, got: %v", err)
	} else if size != 123 {
		t.Errorf("Expected size %d, got: %d", 123, size)
	}
	*/

	// for _, tt := range copyTests {
	// 	task := &CopyTask{Source: tt.in}
	// 	task.Do("install")
	// 	task.Install()
	// }
}
