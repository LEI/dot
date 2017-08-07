package dotlib

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"path"
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
	Destination string
}

func (t *LinkTask) String() string {
	return fmt.Sprintf("%s -> %s", t.Source, t.Target)
}

// Register link
func (t *LinkTask) Register(baseDir string, str string) error {
	parts := strings.Split(str, ":")
	if len(parts) == 1 {
		parts = append(parts, t.Destination)
	} else if len(parts) != 2 {
		return fmt.Errorf("Invalid arg: %s", str)
	}
	src := os.ExpandEnv(parts[0])
	if !path.IsAbs(src) {
		src = path.Join(baseDir, src)
	}
	src = path.Clean(src)
	dst := os.ExpandEnv(parts[1])
	if !path.IsAbs(dst) {
		src = path.Join(t.Destination, dst)
	}
	dst = path.Clean(dst)
	t.Source = src
	t.Target = dst
	return nil
}

// Install link
func (t *LinkTask) Install() error {
	changed, err := Link(t.Source, t.Target)
	if err != nil {
		return err
	}
	prefix := "# "
	if changed {
		prefix = ""
	}
	c := fmt.Sprintf("ln -s %s %s\n", t.Source, t.Target)
	fmt.Printf("%s%s\n", prefix, c)
	return nil
}

// Remove link
func (t *LinkTask) Remove() error {
	prefix := "TODO: "
	c := fmt.Sprintf("rm %s\n", t.Target)
	fmt.Printf("%s%s\n", prefix, c)
	return nil
}

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
