package dot

// https://golang.org/src/os/error.go

import (
	"errors"
	"os"
)

var (
	// ErrAlreadyExist ...
	ErrAlreadyExist = errors.New("already exists")

	// ErrNotEmpty ...
	ErrNotEmpty = errors.New("not empty")

	// ErrNotExist ...
	ErrNotExist = errors.New("does not exists")

	// ErrSkip ...
	ErrSkip = errors.New("skip task")
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
	return e.Op + " " + e.Task.String() + ": " + e.Err.Error()
	// return fmt.Sprintf("%s %s: %s", e.Op, e.Task.String(), e.Err.Error())
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

// IsExist unwraps task error
func IsExist(err error) bool {
	if err == nil {
		return false
	}
	// if terr, ok := err.(*OpError); ok {
	// 	err = terr
	// 	// if terr.Task err == ErrNotEmpty {}
	// }
	return err == ErrAlreadyExist
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
	// if ok {
	// 	// terr.Op == "undo dir"
	// 	if terr.Err == ErrNotEmpty {
	// 		// Skip rm empty directory
	// 		err = ErrSkip
	// 	} else {
	// 		err = terr.Err
	// 	}
	// }
	if terr, ok := err.(*OpError); ok && terr.Err != nil {
		err = terr.Err
	}
	if perr, ok := err.(*os.PathError); ok && perr.Err != nil {
		err = perr.Err
	}
	// switch e := err.(type) {
	// case *os.PathError, *OpError:
	// 	if e.Err != nil {
	// 		err = e.Err
	// 	}
	// }
	switch err {
	case ErrNotEmpty, ErrSkip:
		return true
	default:
		return false
	}
}
