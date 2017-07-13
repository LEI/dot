package helpers

import (
	"fmt"
	"os"
)

var (
// ErrSample = errors.New("Sample")
)

func Link(src, dst string) (bool, error) {
	real, err := ReadLink(dst)
	if err != nil && os.IsExist(err) {
		return false, err
	}
	if real == src { // Symlink already exists
		return false, nil
	}
	if real != "" {
		// fmt.Fprintf(os.Stderr, "# %s is a link to %s, not %s", dst, real, src)
		// os.Exit(1)
		return false, fmt.Errorf("%s is a link to %s, not to %s", dst, real, src)
	}
	fi, err := os.Stat(dst)
	if err != nil && os.IsExist(err) {
		return false, err
	}
	if fi != nil {
		// fmt.Fprintf(os.Stderr, "# %s is already a file", dst)
		// os.Exit(1)
		return false, fmt.Errorf("%s already exists, could not link %s", dst, src)
	}
	err = os.Symlink(src, dst)
	if err != nil {
		return false, err
	}
	return true, nil
}

func ReadLink(path string) (string, error) {
	fi, err := os.Lstat(path)
	if err != nil { // os.IsExist(err)
		// if os.IsNotExist(err) {
		// return path, nil
		// }
		return "", err
	}
	if !IsSymlink(fi) {
		return "", nil
	}
	real, err := os.Readlink(path)
	return real, err
}

func IsSymlink(fi os.FileInfo) bool {
	return fi != nil && fi.Mode()&os.ModeSymlink != 0
}
