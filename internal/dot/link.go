package dot

import (
	"fmt"
	"os"

	"github.com/LEI/dot/internal/shell"
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

// NewLink task
func NewLink(s string) *Link {
	return &Link{Source: s}
}

func (t *Link) String() string {
	s := fmt.Sprintf("%s:%s", t.Source, t.Target)
	switch Action {
	case "install":
		s = fmt.Sprintf("ln -s %s %s", tildify(t.Source), tildify(t.Target))
	case "remove":
		s = fmt.Sprintf("rm %s", tildify(t.Target))
	}
	return s
}

// Init task
func (t *Link) Init() error {
	// ...
	return nil
}

// Status check task
func (t *Link) Status() error {
	exists, err := linkExists(t.Source, t.Target)
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
		// TODO os.LinkError Err: ErrExist
		case ErrFileExist, ErrLinkExist:
			if Action != "install" {
				fmt.Println("Skip", Action, t.Target, "("+perr.Err.Error()+")")
				return ErrSkip
			}
			// Confirm override
			if shell.AskConfirmation("Remove existing " + t.Target + "?") {
				if err := os.Remove(t.Target); err != nil {
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
		return ErrExist
	}
	return nil
}

// Do task
func (t *Link) Do() error {
	if err := t.Status(); err != nil {
		switch err {
		case ErrExist, ErrSkip:
			return nil
		// case ErrFileExist, ErrLinkExist:
		// 	// Confirm override
		// 	if !shell.AskConfirmation("Remove existing " + t.Target + "?") {
		// 		return ErrSkip
		// 	}
		// 	if rmerr := os.Remove(t.Target); rmerr != nil {
		// 		return rmerr
		// 	}
		default:
			return err
		}
	}
	return os.Symlink(t.Source, t.Target)
}

// Undo task
func (t *Link) Undo() error {
	if err := t.Status(); err != nil {
		switch err {
		case ErrExist:
			// continue
		case ErrSkip:
			return nil
		default:
			return err
		}
	}
	return os.Remove(t.Target)
}

// linkExists returns true if the link has the same target.
func linkExists(src, dst string) (bool, error) {
	if !exists(src) {
		// fmt.Fprintf(os.Stderr, "%s: no such file or directory (to link to %s)\n", src, dst)
		// return false, ErrSkip // fmt.Errorf("%s: no such file or directory (to link to %s)", src, dst)
		return false, &os.PathError{Op: "link source", Path: src, Err: ErrNotExist}
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
		return false, &os.PathError{
			Op:   "target link",
			Path: src,
			Err:  ErrFileExist,
		}
	}
	real, err := os.Readlink(dst)
	if err != nil {
		return false, err
	}
	if real == "" {
		return false, fmt.Errorf("%s: unable to read symlink real target", dst)
		// return false, &os.LinkError{
		// 	Op: "readllink"
		// 	Path: dst,
		// 	Err: os.ErrInvalid,
		// }
	}
	if real != src {
		// return false, fmt.Errorf("%s: already a symlink to %s, want %s", dst, real, src)
		return false, &os.PathError{
			Op:   "target link (real: " + real + ")",
			Path: src,
			Err:  ErrLinkExist,
		}
		// return false, &os.LinkError{
		// 	Op: "link "+src+" to "+dst,
		// 	Path: real,
		// 	Err: os.ErrExist,
		// }
	}
	return true, nil
}

// isSymlink checks a given file info corresponds to a symbolic link
func isSymlink(fi os.FileInfo) bool {
	return fi != nil && fi.Mode()&os.ModeSymlink != 0
}
