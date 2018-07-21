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

// Do ...
func (t *LinkTask) Do(a string) error {
	return do(t, a)
}

// Install link
func (t *LinkTask) Install() error {
	if err := createBaseDir(t.Target); err != nil && err != ErrDirShouldExist {
		return err
	}
	changed, err := Link(t.Source, t.Target)
	if err != nil {
		return err
	}
	prefix := ""
	if !changed {
		prefix = "# "
	}
	fmt.Printf("%sln -s %s %s\n", prefix, t.Source, t.Target)
	return nil
}

// Remove link
func (t *LinkTask) Remove() error {
	changed, err := Unlink(t.Source, t.Target)
	if err != nil {
		return err
	}
	prefix := ""
	if !changed {
		prefix = "# "
	}
	fmt.Printf("%srm %s\n", prefix, t.Target)
	if RemoveEmptyDirs {
		if err := removeBaseDir(t.Target); err != nil {
			return err
		}
	}
	return nil
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
		fmt.Fprintf(os.Stderr, "# %s is a file? at least not a link to %s\n", dst, src)
		return false, err
	}
	if real == src { // Symlink already exists
		return false, nil
	}
	if real != "" {
		// return false, ErrLinkExist // fmt.Errorf("%s is a link to %s, not to %s", dst, real, src)
		q := fmt.Sprintf("Replace %s link to %s with %s", dst, real, src)
		if !AskConfirmation(q) {
			fmt.Fprintf(os.Stderr, "Skipping symlink %s because its target exists: %s", src, dst)
			return false, nil
		}
		// if err := Backup(dst); err != nil {
		// 	return false, err
		// }
	}
	if err != nil && os.IsExist(err) {
		return false, err
	}
	fi, err := os.Stat(dst)
	if err != nil && os.IsExist(err) {
		return false, err
	}
	// fmt.Println("SRC:", src)
	// fmt.Println("REAL:", real)
	// fmt.Println("ERR:", err)
	// fmt.Println("FI:", fi)
	if fi != nil {
		fmt.Fprintf(os.Stderr, "# %s is already a file\n", dst)
		// os.Exit(1)
		return false, ErrFileExist // fmt.Errorf("%s already exists, could not link %s", dst, src)
	}
	if DryRun {
		return true, nil
	}
	err = os.Symlink(src, dst)
	if err != nil {
		return false, err
	}
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
