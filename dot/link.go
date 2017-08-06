package dot

import (
	"errors"
	// "fmt"
	"os"
)

var (
	// ErrLinkExist error
	ErrLinkExist = errors.New("Link exists")
	// ErrFileExist error
	ErrFileExist = errors.New("File exists")
)

// LinkTask struct
type LinkTask struct {
	Source, Target string
}

// Install link
// func (l *LinkTask) Install() error {
// 	changed, err := Link(l.Source, l.Target)
// 	if err != nil {
// 		return err
// 	}
// 	prefix := "# "
// 	if changed {
// 		prefix = ""
// 	}
// 	fmt.Printf("%sln -s %s %s\n", prefix, l.Source, l.Target)
// 	return nil
// }

// Link task
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
		return false, ErrLinkExist // fmt.Errorf("%s is a link to %s, not to %s", dst, real, src)
	}
	fi, err := os.Stat(dst)
	if err != nil && os.IsExist(err) {
		return false, err
	}
	if fi != nil {
		// fmt.Fprintf(os.Stderr, "# %s is already a file", dst)
		// os.Exit(1)
		return false, ErrFileExist // fmt.Errorf("%s already exists, could not link %s", dst, src)
	}
	err = os.Symlink(src, dst)
	if err != nil {
		return false, err
	}
	return true, nil
}

// ReadLink path
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

// IsSymlink check
func IsSymlink(fi os.FileInfo) bool {
	return fi != nil && fi.Mode()&os.ModeSymlink != 0
}
