package dot

import (
	"fmt"
	"os"

	"github.com/LEI/dot/internal/prompt"
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

// Type task name
func (l *Link) Type() string {
	return "link" // symbolic
}

// DoString string
func (l *Link) DoString() string {
	return fmt.Sprintf("ln -s %s %s", tildify(l.Source), tildify(l.Target))
}

// UndoString string
func (l *Link) UndoString() string {
	return fmt.Sprintf("rm %s", tildify(l.Target))
}

// Status check task
func (l *Link) Status() error {
	exists, err := linkExists(l.Source, l.Target)
	if err != nil {
		perr, ok := err.(*os.PathError)
		// if ok {
		// 	return perr
		// }
		// return err
		if !ok {
			return err
		}
		switch perr.Err {
		case ErrFileExist, ErrLinkExist:
			if l.current != "install" {
				fmt.Println("Skip", l.current, l.Target, "("+perr.Err.Error()+")")
				return ErrSkip
			}
			// Confirm override
			if prompt.AskConfirmation("Remove existing " + l.Target + "?") {
				if err := os.Remove(l.Target); err != nil {
					return err
				}
				return nil
			}
			// if err := os.Remove(); e
			return perr // .Err
		}
		return perr
	}
	if exists {
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
		// case ErrFileExist, ErrLinkExist:
		// 	// Confirm override
		// 	if !prompt.AskConfirmation("Remove existing " + l.Target + "?") {
		// 		return ErrSkip
		// 	}
		// 	if rmerr := os.Remove(l.Target); rmerr != nil {
		// 		return rmerr
		// 	}
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
		// return false, &os.PathError{Op: "source link", Path: src, Err: ErrNotExist}
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
		// return false, fmt.Errorf("%s: not a symlink", dst)
		return false, &os.PathError{Op: "target link", Path: src, Err: ErrFileExist}
	}
	real, err := os.Readlink(dst)
	if err != nil {
		return false, err
	}
	if real == "" {
		return false, fmt.Errorf("%s: unable to read symlink real target", dst)
	}
	if real != src {
		// return false, fmt.Errorf("%s: already a symlink to %s, want %s", dst, real, src)
		return false, &os.PathError{Op: "target link (real: " + real + ")", Path: src, Err: ErrLinkExist}
	}
	return true, nil
}

// isSymlink checks a given file info corresponds to a symbolic link
func isSymlink(fi os.FileInfo) bool {
	return fi != nil && fi.Mode()&os.ModeSymlink != 0
}
