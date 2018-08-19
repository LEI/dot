package dot

import (
	"fmt"
	"os"
)

// // LinkError type
// type LinkError struct {
// 	// taskError
// 	Action string
// 	Path   string
// 	Err    error
// 	// skip   bool
// }

// func (e *LinkError) Error() string {
// 	return e.Action + " " + e.Path + ": " + e.Err.Error()
// }

// Link task
type Link struct {
	Task   `mapstructure:",squash"` // Action, If, OS
	Source string
	Target string
}

func (l *Link) String() string {
	return fmt.Sprintf("%s:%s", l.Source, l.Target)
}

// DoString string
func (l *Link) DoString() string {
	return fmt.Sprintf("ln -s %s %s", l.Source, l.Target)
}

// UndoString string
func (l *Link) UndoString() string {
	return fmt.Sprintf("rm %s", l.Target)
}

// Status check task
func (l *Link) Status() error {
	ok, err := linkExists(l.Source, l.Target)
	if err != nil {
		return nil
	}
	if ok {
		return ErrAlreadyExist
	}
	return nil
}

// Do task
func (l *Link) Do() error {
	if err := l.Status(); err != nil {
		switch err {
		case ErrAlreadyExist, ErrSkip:
			return nil
		default:
			return err
		}
	}
	return os.Symlink(l.Source, l.Target)
}

// Undo task
func (l *Link) Undo() error {
	if err := l.Status(); err != nil {
		switch err {
		case ErrSkip:
			return nil
		case ErrAlreadyExist:
			// continue
		default:
			return err
		}
	}
	return os.Remove(l.Target)
}

// linkExists returns true if the link has the same target.
func linkExists(src, dst string) (bool, error) {
	if !exists(src) {
		return false, fmt.Errorf("%s: no such file or directory (to link %s)", src, dst)
	}
	if !exists(dst) {
		// Stop here if the target does not exist
		return false, nil
	}
	fi, err := os.Lstat(dst)
	if err != nil {
		return false, err
	}
	if !isSymlink(fi) {
		return false, fmt.Errorf("%s: not a symlink", dst)
	}
	real, err := os.Readlink(dst)
	if err != nil {
		return false, err
	}
	if real == "" {
		return false, fmt.Errorf("%s: unable to read symlink", dst)
	}
	if real != src {
		return false, fmt.Errorf("%s: already a symlink to %s, want %s", dst, real, src)
	}
	return true, nil
}

// isSymlink checks a given file info corresponds to a symbolic link
func isSymlink(fi os.FileInfo) bool {
	return fi != nil && fi.Mode()&os.ModeSymlink != 0
}
