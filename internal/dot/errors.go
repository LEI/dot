package dot

// https://golang.org/src/os/error.go
// https://github.com/golang/go/blob/master/src/os/error.go

import (
	"errors"
	"fmt"
	"os"
	"strings"
)

var (
	// ErrExist used when task is installed
	ErrExist = errors.New("already exists")

	// ErrNotExist ...
	ErrNotExist = errors.New("does not exists")

	// ErrFileExist ...
	ErrFileExist = errors.New("file exists")

	// ErrLinkExist ...
	ErrLinkExist = errors.New("link exists")

	// ErrNotEmpty ...
	ErrNotEmpty = errors.New("not empty")

	// ErrDirNotEmpty ...
	ErrDirNotEmpty = errors.New("directory not empty")

	// ErrSkip task
	ErrSkip = errors.New("skip")

	// ErrInvalid ...
	ErrInvalid = errors.New("invalid file")
)

// OpError ...
type OpError struct {
	Op   string
	Task Tasker
	Err  error
	// Code    string // ErrorCode
	// Detail  interface{}
	// Format  string
	// Message string
}

func (e *OpError) Error() string {
	// return e.Op + " " + fmt.Sprintf("[%s]", e.Task) + ": " + e.Err.Error()
	return fmt.Sprintf("%s %s: %s", e.Op, e.Task.String(), e.Err.Error())
}

// // Code is used as a prefix
// // If Format and Detail are given, use it as a template for Message format
// // If only Format is given, apply it to Message
// // Otherwise just use Message
// func (e *FmtError) Error() string {
// 	msg := e.Message
// 	if e.Format != "" && e.Detail == nil {
// 		msg = fmt.Sprintf(e.Format, e.Message)
// 	} else if e.Format != "" { // e.Detail != nil
// 		t, err := template.New("err" + e.Code).Parse(e.Format)
// 		if err == nil {
// 			var tpl bytes.Buffer
// 			if err := t.Execute(&tpl, e.Detail); err == nil {
// 				msg = fmt.Sprintf(tpl.String(), e.Message)
// 			}
// 		}
// 	}
// 	if e.Code == "" {
// 		e.Code = "task error"
// 	}
// 	return fmt.Sprintf("%s: %s", e.Code, msg)
// }

// DiffError ...
type DiffError struct {
	Src, Dst string
	Full     string
}

func (e *DiffError) Error() string {
	// return e.Src + " != " + e.Dst
	return fmt.Sprintf("%s differs from source %s", tildify(e.Dst), tildify(e.Src))
}

// Full diff ouput
func (e *DiffError) String() string {
	fullDiff := strings.TrimSuffix(e.Full, "\n")
	return fmt.Sprintf(
		"--- %s\n+++ %s\n%s",
		tildify(e.Src),
		tildify(e.Dst),
		fullDiff,
	)
}

// IsExist error
func IsExist(err error) bool {
	if err == nil {
		return false
	}
	// if terr, ok := err.(*OpError); ok {
	// 	err = terr
	// 	// if terr.Task err == ErrNotEmpty {}
	// }
	switch err {
	case ErrExist:
		return true
	// case ErrFileExist, ErrLinkExist:
	// 	return true
	default:
		return false
	}
}

// IsNotExist error
func IsNotExist(err error) bool {
	return !IsExist(err)
}

// IsSkip error
func IsSkip(err error) bool {
	if err == nil {
		return false
	}
	err = unwrapError(err)
	switch err {
	case ErrDirNotEmpty, ErrSkip:
		return true
	default:
		return false
	}
}

func unwrapError(err error) error {
	switch e := err.(type) {
	case *os.PathError: // , *OpError:
		if e.Err != nil {
			err = e.Err
		}
	}
	// if terr, ok := err.(*OpError); ok && terr.Err != nil {
	// 	if terr.Op == "rmdir" {
	// 		// ...
	// 	} else {
	// 		err = terr.Err
	// 	}
	// }
	// // Check for custom os.PathError (os.Remove Op: "remove")
	// if perr, ok := err.(*os.PathError); ok && perr.Err != nil {
	// 	if perr.Op == "rmdir" && perr.Err == ErrDirNotEmpty {
	// 		// fmt.Fprintf(os.Stderr, "# rmdir: %s: directory not empty\n", perr.Path)
	// 		// TODO Skip ErrDirNotEmpty
	// 	}
	// 	err = perr.Err
	// }
	return err
}
