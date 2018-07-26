package dotfile

import (
	"errors"
	"fmt"
	"os"
	// "path"
	// "path/filepath"
	// "strings"
)

var (
	// ErrLinkExist error
	ErrLinkExist = errors.New("link exists")
	// ErrFileExist error
	ErrFileExist = errors.New("file exists")
)

// LinkTask struct
type LinkTask struct {
	Source, Target string
	Task
}

// Status link
func (t *LinkTask) Status() bool {
	return true
}

// Do ...
func (t *LinkTask) Do(a string) (string, error) {
	return do(t, a)
}

// List link
func (t *LinkTask) List() (string, error) {
	str := fmt.Sprintf("Link: %s -> %s", t.Source, t.Target)
	return str, nil
}

// Install link
func (t *LinkTask) Install() (string, error) {
	if err := createBaseDir(t.Target); err != nil && err != ErrDirShouldExist {
		return "", err
	}
	changed, err := Link(t.Source, t.Target)
	if err != nil {
		return "", err
	}
	prefix := ""
	if !changed {
		prefix = "# "
	}
	return fmt.Sprintf("%sln -s %s %s", prefix, t.Source, t.Target), nil
}

// Remove link
func (t *LinkTask) Remove() (string, error) {
	changed, err := Unlink(t.Source, t.Target)
	if err != nil {
		return "", err
	}
	prefix := ""
	if !changed {
		prefix = "# "
	}
	if RemoveEmptyDirs {
		if err := removeBaseDir(t.Target); err != nil {
			return "", err
		}
	}
	return fmt.Sprintf("%srm %s", prefix, t.Target), nil
}

// IsSymlink check
func IsSymlink(fi os.FileInfo) bool {
	return fi != nil && fi.Mode()&os.ModeSymlink != 0
}

// ReadLink path
func ReadLink(path string) (string, error) {
	// fi, err := os.Lstat(path)
	// fmt.Println("FI", IsSymlink(fi), err)
	// if err != nil { // os.IsExist(err)
	// 	// if os.IsNotExist(err) {
	// 	// return path, nil
	// 	// }
	// 	return "", err
	// }
	// real, err := os.Readlink(path)
	// fmt.Println("--->", real, err)
	// if !IsSymlink(fi) {
	// 	// Quickfix: directories seem to be ignored
	// 	real, err = filepath.EvalSymlinks(path)
	// 	fmt.Println("===>",real, err)
	// } else if !IsSymlink(fi) {
	// 	return "", nil
	// }
	// return real, err
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

// Link task
func Link(src, dst string) (bool, error) {
	real, err := ReadLink(dst)
	// if real == "" {
	// 	// Quickfix: directories seem to be ignored,
	// 	// so try harder to find the link target
	// 	real, err = filepath.EvalSymlinks(dst)
	// 	if real == dst {
	// 		real = ""
	// 	}
	// }
	if err != nil && os.IsExist(err) {
		// ErrFileExist
		if real == src && err == ErrLinkExist {
			return false, nil
		}
		q := fmt.Sprintf("Replace %s with a link to %s", dst, src)
		if !AskConfirmation(q) {
			fmt.Fprintf(os.Stderr, "Skipping symlink %s because its target is an existing link: %s", src, dst)
			return false, nil
		}
		// TODO remove?
		// fmt.Fprintf(os.Stderr, "# %s is a file? at least not a link to %s\n", dst, src)
		// return false, err
	}
	if real == src { // Symlink already exists
		return false, nil
	}
	fi, err := os.Stat(dst)
	if err != nil && os.IsExist(err) {
		return false, err
	}
	if fi != nil && real == "" {
		// return false, ErrLinkExist // fmt.Errorf("%s is a link to %s, not to %s", dst, real, src)
		q := fmt.Sprintf("Replace %s link to %s with %s", dst, real, src)
		if !AskConfirmation(q) {
			fmt.Fprintf(os.Stderr, "Skipping symlink %s because its target is an existing file: %s", src, dst)
			return false, nil
		}
		// if err := Backup(dst); err != nil {
		// 	return false, err
		// }
	}
	if DryRun {
		return true, nil
	}
	err = os.Symlink(src, dst)
	if err != nil {
		return false, err
	}
	// TODO: cache[dst] = src?
	// if err := dotCache.Put(dst, content); err != nil {
	// 	return true, err
	// }
	return true, nil
}

// Unlink task
func Unlink(src, dst string) (bool, error) {
	real, err := ReadLink(dst)
	if err != nil && os.IsExist(err) {
		return false, err
	}
	if real != "" && real != src {
		return false, ErrLinkExist // fmt.Errorf("%s is a link to %s, not to %s", dst, real, src)
	}
	fi, err := os.Stat(dst)
	if err != nil && os.IsExist(err) {
		return false, err
	}
	if fi == nil { // Target does no exist
		return false, nil
	}
	if DryRun {
		return true, nil
	}
	err = os.Remove(dst)
	if err != nil {
		return false, err
	}
	return true, nil
}
