package system

import (
	"fmt"
	"os"
)

var (
	// DirMode default
	DirMode os.FileMode = 0755
)

// CheckDirectory ... (verify/validate)
func CheckDirectory(dir string) error {
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

// CreateDirectory ...
func CreateDirectory(dir string) error {
	// if dir == "" {
	// 	return fmt.Errorf("missing dir arg!")
	// }
	fmt.Printf("$ mkdir -p %s\n", dir)
	if DryRun {
		return nil
	}
	return os.MkdirAll(dir, DirMode)
}

// RemoveDirectory ...
func RemoveDirectory(dir string) error {
	// if dir == "" {
	// 	return fmt.Errorf("missing dir arg!")
	// }
	fmt.Printf("$ rmdir %s\n", dir)
	if DryRun {
		return nil
	}
	return os.Remove(dir)
}
