package dot

import (
	"bytes"
	"fmt"
	"os"
	"text/template"

	"github.com/LEI/dot/internal/ostype"
)

// taskError type
type taskError struct {
	Code    string // ErrorCode
	Detail  interface{}
	Format  string
	Message string
	// Err error
}

// Code is used as a prefix
// If Format and Detail are given, use it as a template for Message format
// If only Format is given, apply it to Message
// Otherwise just use Message
func (e *taskError) Error() string {
	msg := e.Message
	if e.Format != "" && e.Detail == nil {
		msg = fmt.Sprintf(e.Format, e.Message)
	} else if e.Format != "" { // e.Detail != nil
		t, err := template.New("err" + e.Code).Parse(e.Format)
		if err == nil {
			var tpl bytes.Buffer
			if err := t.Execute(&tpl, e.Detail); err == nil {
				msg = fmt.Sprintf(tpl.String(), e.Message)
			}
		}
	}
	if e.Code == "" {
		e.Code = "task error"
	}
	return fmt.Sprintf("%s: %s", e.Code, msg)
}

// Tasker interface
type Tasker interface {
	IsAction(string) bool
	String() string
	DoString() string
	UndoString() string
	Status() error
	// Sync() error
	Do() error
	Undo() error
	CheckIf() error
	CheckOS() error
	// GetOS() []string
}

// Task struct
type Task struct {
	Tasker
	State string   `mapstructure:"action,omitempty"` // install, remove
	If    []string `mapstructure:",omitempty"`
	OS    []string `mapstructure:",omitempty"`
}

// IsAction task
func (t *Task) IsAction(state string) bool {
	if len(t.State) == 0 {
		// FIXME: detect if Task.State is ignored
		// e.g. private Task.state or just omitted
		return true
	}
	return t.State == state
}

// IsOk status
func IsOk(err error) bool {
	return err == ErrAlreadyExist
}

// exists checks if a file is present
func exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil || os.IsExist(err)
}

// // GetOS ...
// func (t *Task) GetOS() []string {
// 	return t.OS
// }

// CheckOS task
func (t *Task) CheckOS() error {
	if len(t.OS) == 0 {
		return nil
	}
	ok := ostype.Has(t.OS...)
	if !ok {
		return ErrSkip
	}
	return nil
}

// CheckIf task
func (t *Task) CheckIf() error {
	if len(t.If) == 0 {
		return nil
	}
	// varsMap := map[string]interface{}{
	// 	// "DryRun": system.DryRun,
	// 	// "Verbose": tasks.Verbose,
	// 	// "OS":      runtime.GOOS,
	// }
	// funcMap := template.FuncMap{
	// 	"hasOS": ostype.Has,
	// }
	if len(t.If) > 0 {
		fmt.Println("TODO (skip) If:", t.If)
		return ErrSkip
	}
	// https://golang.org/pkg/text/template/#hdr-Functions
	// for _, cond := range t.If {
	// 	str, err := TemplateData("", cond, varsMap, funcMap)
	// 	if err != nil {
	// 		fmt.Fprintf(os.Stderr, "err tpl: %s\n", err)
	// 		continue
	// 	}
	// 	_, stdErr, status := executils.ExecuteBuf("sh", "-c", str)
	// 	// out := strings.TrimRight(string(stdOut), "\n")
	// 	strErr := strings.TrimRight(string(stdErr), "\n")
	// 	// if out != "" {
	// 	// 	fmt.Printf("stdout: %s\n", out)
	// 	// }
	// 	if strErr != "" {
	// 		fmt.Fprintf(os.Stderr, "'%s' stderr: %s\n", str, strErr)
	// 	}
	// 	if status == 0 {
	// 		return true
	// 	}
	// }
	return nil
}
