package dot

import (
	"fmt"
	"os"

	"github.com/LEI/dot/internal/ostype"
)

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
	GetOS() []string
}

// Task struct
type Task struct {
	Tasker
	state string   `mapstructure:"action,omitempty"` // install, remove
	If    []string `mapstructure:",omitempty"`
	OS    []string `mapstructure:",omitempty"`
}

// IsAction task
func (t *Task) IsAction(state string) bool {
	return t.state == "" || t.state == state
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

// GetOS task
func (t *Task) GetOS() []string {
	return t.OS
}

// CheckOS task
func (t *Task) CheckOS() error {
	if len(t.OS) == 0 {
		return nil
	}
	if ostype.Has(t.OS...) {
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
