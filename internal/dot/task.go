package dot

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/LEI/dot/internal/host"
	"github.com/docker/docker/pkg/homedir"
)

// Tasker interface
type Tasker interface {
	String() string
	Status() error
	// Sync() error
	Do() error
	Undo() error

	// Already implemented
	SetAction(string) *Task
	GetAction() string
	CheckAction() error
	CheckOS() error
	CheckIf() error
	Check() error
}

// Task struct
type Task struct {
	Tasker
	Action string   `mapstructure:",omitempty"` // install, remove
	OS     []string `mapstructure:",omitempty"`
	If     []string `mapstructure:",omitempty"`

	running string // Current action name
}

var (
	homeDir = homedir.Get()
)

// SetAction name
func (t *Task) SetAction(name string) *Task {
	t.running = name
	return t
}

// GetAction name
func (t *Task) GetAction() string {
	return t.running
}

// Check conditions
func (t *Task) Check() error {
	if err := t.CheckAction(); err != nil {
		// fmt.Println("> Skip "+action, t, err)
		return err
	}
	if err := t.CheckOS(); err != nil {
		// fmt.Println("> Skip OS", t, err)
		return err
	}
	if err := t.CheckIf(); err != nil {
		// fmt.Println("> Skip If", t, err)
		return err
	}
	return nil
}

// CheckAction task
func (t *Task) CheckAction() error {
	if len(t.running) == 0 {
		return fmt.Errorf("unable to check empty action")
	}
	if len(t.Action) == 0 {
		// FIXME: detect if Task.Action is ignored
		// e.g. private Task.state or just omitted
		return nil
	}
	if t.Action != t.running {
		return ErrSkip
	}
	return nil
}

// CheckOS task
func (t *Task) CheckOS() error {
	if len(t.OS) == 0 {
		return nil
	}
	ok := host.HasOS(t.OS...)
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
	// 	"hasOS": host.HasOS,
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
	// 	_, stdErr, status := executils.ExecuteBuf(shell.Get(), "-c", str)
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
