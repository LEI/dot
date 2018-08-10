package system

import (
	"fmt"
	"os"
)

// CheckSymlink ... (verify/validate)
func CheckSymlink(src, dst string) error {
	// fmt.Println("CheckSymlink", src, dst)
	// if src == "" || dst == "" {
	// 	return fmt.Errorf("missing symlink arg: [src:%s dst:%s]", src, dst)
	// }
	if !Exists(src) {
		// return ErrIsNotExist
		return fmt.Errorf("%s: no such file or directory (to link %s)", src, dst)
	}
	if !Exists(dst) {
		// Stop here if the target does not exist
		return nil
	}
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
	if real == "" {
		return fmt.Errorf("%s: unable to read symlink", dst)
	}
	if real != src {
		return fmt.Errorf("%s: already a symlink to %s, want %s", dst, real, src)
	}
	return ErrLinkExist // exact same link already exists
}

// Symlink ...
func Symlink(src, dst string) error {
	// if src == "" || dst == "" {
	// 	return fmt.Errorf("missing symlink arg! [src:%s dst:%s]", src, dst)
	// }
	if DryRun {
		return nil
	}
	return os.Symlink(src, dst)
}

// Unlink ...
func Unlink(dst string) error {
	// if src == "" || dst == "" {
	// 	return fmt.Errorf("missing symlink arg! [src:%s dst:%s]", src, dst)
	// }
	if DryRun {
		return nil
	}
	return os.Remove(dst)
}

// IsSymlink checks a given file info corresponds to a symbolic link
func IsSymlink(fi os.FileInfo) bool {
	return fi != nil && fi.Mode()&os.ModeSymlink != 0
}
