package system

import (
	"fmt"
	"os"
)

// CheckSymlink ... (verify/validate)
func CheckSymlink(src, dst string) error {
	fi, err := os.Lstat(dst)
	if err != nil {
		return err
	}
	if !IsSymlink(fi) {
		return fmt.Errorf("%s: not a symlink", dst)
	}
	real, err := os.Readlink(dst)
	if err != nil {
		return err
	}
	if real != "" && real != src {
		return fmt.Errorf("%s: already symlinked to %s", dst, real)
	}
	return nil
}

// Symlink ...
func Symlink(src, dst string) error {
	if DryRun {
		return nil
	}
	fmt.Println("DO SYMLINK", src, dst)
	return nil // os.Symlink(src, dst)
}

// IsSymlink checks a given file info corresponds to a symbolic link
func IsSymlink(fi os.FileInfo) bool {
	return fi != nil && fi.Mode()&os.ModeSymlink != 0
}
