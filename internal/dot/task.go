package dot

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/LEI/dot/internal/ostype"
	"github.com/docker/docker/pkg/homedir"
)

// Tasker interface
type Tasker interface {
	String() string
	Type() string
	SetAction(string)
	Check() error
	CheckAct() error
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
	Action string   `mapstructure:",omitempty"` // install, remove
	If     []string `mapstructure:",omitempty"`
	OS     []string `mapstructure:",omitempty"`

	current string // Current action name
}

var (
	homeDir = homedir.Get()
)

// SetAction name
func (t *Task) SetAction(name string) {
	t.current = name
}

// Check conditions
func (t *Task) Check() error {
	if err := t.CheckAction(); err != nil {
		// fmt.Println("> Skip "+action, t, err)
		return err
	}
	if err := t.CheckIf(); err != nil {
		// fmt.Println("> Skip If", t, err)
		return err
	}
	if err := t.CheckOS(); err != nil {
		// fmt.Println("> Skip OS", t, err)
		return err
	}
	return nil
}

// CheckAction task
func (t *Task) CheckAction() error {
	if len(t.current) == 0 {
		return fmt.Errorf("unable to check empty action")
	}
	if len(t.Action) == 0 {
		// FIXME: detect if Task.Action is ignored
		// e.g. private Task.state or just omitted
		return nil
	}
	if t.Action != t.current {
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

func tildify(path string) string {
	prefix := filepath.Join(homeDir)
	// +string(os.PathSeparator)
	if !strings.HasPrefix(path, prefix) {
		return path
	}
	s := homedir.GetShortcutString()
	return s + strings.TrimPrefix(path, prefix)
}

// exists checks if a file is present
func exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil || os.IsExist(err)
}
