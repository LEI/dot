package dot

import (
	"os"

	"github.com/LEI/dot/internal/ostype"
)

// Tasker interface
type Tasker interface {
	String() string
	Type() string
	Check(string) error
	CheckAction(string) error
	CheckIf() error
	CheckOS() error
	// GetOS() []string
	DoString() string
	UndoString() string
	Status() error
	// Sync() error
	Do() error
	Undo() error
}

// Task struct
type Task struct {
	Tasker
	State string   `mapstructure:"action,omitempty"` // install, remove
	If    []string `mapstructure:",omitempty"`
	OS    []string `mapstructure:",omitempty"`
}

// Check conditions
func (t *Task) Check(action string) error {
	if err := t.CheckAction(action); err != nil {
		return err
	}
	if err := t.CheckIf(); err != nil {
		// fmt.Println("> Skip If", err)
		return err
	}
	if err := t.CheckOS(); err != nil {
		// fmt.Println("> Skip OS", err)
		return err
	}
	return nil
}

// CheckAction task
func (t *Task) CheckAction(name string) error {
	if len(t.State) == 0 {
		// FIXME: detect if Task.State is ignored
		// e.g. private Task.state or just omitted
		return nil
	}
	if t.State != name {
		return ErrSkip
	}
	return nil
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
		return ErrSkip // &OpError{"check os", t, ErrSkip}
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
		return ErrSkip // &OpError{"check if", t, ErrSkip}
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

// exists checks if a file is present
func exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil || os.IsExist(err)
}
