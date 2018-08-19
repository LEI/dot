package dot

import "errors"

var (
	// ErrAlreadyExist ...
	ErrAlreadyExist = errors.New("already exists")

	// ErrNotEmpty ...
	ErrNotEmpty = errors.New("not empty")

	// ErrSkip ...
	ErrSkip = errors.New("skip task")
)

// TaskError ...
type TaskError struct {
	Op   string
	Task Tasker
	Err  error
	// Code    string // ErrorCode
	// Detail  interface{}
	// Format  string
	// Message string
}

func (e *TaskError) Error() string {
	return e.Op + " " + e.Task.String() + ": " + e.Err.Error()
	// return fmt.Sprintf("%s %s: %s", e.Op, e.Task.String(), e.Err.Error())
}

// // Code is used as a prefix
// // If Format and Detail are given, use it as a template for Message format
// // If only Format is given, apply it to Message
// // Otherwise just use Message
// func (e *TaskError) Error() string {
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
