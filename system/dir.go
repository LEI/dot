package system

import (
	// "fmt"
	"io"
	"os"
)

var (
	// DirMode default
	DirMode os.FileMode = 0755
)

// CheckDir ... (verify/validate)
func CheckDir(dir string) error {
	fi, err := os.Stat(dir)
	if err != nil && os.IsExist(err) {
	    return err
	}
	if fi != nil && fi.IsDir() {
	    return ErrDirExist
	} else if fi != nil {
	    return ErrFileExist
	}
	return nil
}

// CreateDir ...
func CreateDir(dir string) error {
	// if dir == "" {
	// 	return fmt.Errorf("missing dir arg!")
	// }
	if DryRun {
		return nil
	}
	return os.MkdirAll(dir, DirMode)
}

// RemoveDir ...
func RemoveDir(dir string) error {
	// if dir == "" {
	// 	return fmt.Errorf("missing dir arg!")
	// }
	if DryRun {
		return nil
	}
	return os.Remove(dir)
}

// IsEmptyDir ...
func IsEmptyDir(name string) (bool, error) {
	f, err := os.Open(name)
	if err != nil {
		return false, err
	}
	defer f.Close()
	_, err = f.Readdirnames(1) // Or f.Readdir(1)
	if err == io.EOF {
		return true, nil
	}
	return false, err // Either not empty or error, suits both cases
}
