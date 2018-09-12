package dot

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/LEI/dot/internal/host"
	"github.com/LEI/dot/internal/shell"
)

var (
	// Action contains the running command name
	Action string
)

// Tasker interface
type Tasker interface {
	String() string
	Init() error
	Status() error
	// Sync() error
	Do() error
	Undo() error

	// Already implemented
	CheckAction() error
	CheckOS() error
	CheckIf() error
	Check() error
}

// Task struct
type Task struct {
	Tasker
	// Action specifies a single action for which the task should be run
	Action string   `mapstructure:",omitempty"`
	OS     []string `mapstructure:",omitempty"`
	If     []string `mapstructure:",omitempty"`
}

var (
	homeDir = shell.HomeDir
)

func (t *Task) String() string {
	// FIXME invalid memory address or nil pointer dereference
	return "<task interface>"
}

// Check conditions
func (t *Task) Check() error {
	// TODO Debug Verbose
	if err := t.CheckAction(); err != nil {
		// fmt.Printf("> If %q != %q -> %s\n", t.Action, Action, err)
		return err
	}
	if err := t.CheckOS(); err != nil {
		// fmt.Printf("> OS %s -> %s\n", t.OS, err)
		return err
	}
	if err := t.CheckIf(); err != nil {
		// fmt.Printf("> If %q -> %s\n", t.If, err)
		return err
	}
	return nil
}

// CheckAction task
func (t *Task) CheckAction() error {
	if len(t.Action) == 0 {
		// FIXME: detect if Task.Action is ignored
		// e.g. private Task.state or just omitted
		return nil
	}
	if t.Action != Action {
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
	// https://golang.org/pkg/text/template/#hdr-Functions
	for i, cond := range t.If {
		// str, err := TemplateData("", cond, varsMap, funcMap)
		// if err != nil {
		// 	fmt.Fprintf(os.Stderr, "err tpl: %s\n", err)
		// 	continue
		// }
		name := fmt.Sprintf("if(%d)", i+1)
		c, err := buildTplEnv(name, cond)
		if err != nil {
			return err
		}
		// _, stdErr, status := executils.ExecuteBuf(shell.Get(), "-c", str)
		// // out := strings.TrimRight(string(stdOut), "\n")
		// strErr := strings.TrimRight(string(stdErr), "\n")
		// // if out != "" {
		// // 	fmt.Printf("stdout: %s\n", out)
		// // }
		// if strErr != "" {
		// 	fmt.Fprintf(os.Stderr, "'%s' stderr: %s\n", str, strErr)
		// }
		// if status == 0 {
		// 	return true
		// }
		// fmt.Printf("EXEC COND: %q\n", c)
		cmd := exec.Command(shell.Get(), "-c", c)

		// cmd.Stdout = os.Stdout
		// cmd.Stderr = os.Stderr
		// cmd.Stdin = os.Stdin

		// cmd.Env = expand
		if err := cmd.Run(); err != nil {
			// if Verbose > 1 {
			// 	fmt.Fprintf(os.Stderr, "skip task because %s -> %s\n", c, err)
			// }
			return ErrSkip
		}
	}
	return nil
}

func tildify(path string) string {
	prefix := filepath.Join(homeDir)
	// +string(os.PathSeparator)
	if !strings.HasPrefix(path, prefix) {
		return path
	}
	s := shell.GetHomeShortcutString()
	return s + strings.TrimPrefix(path, prefix)
}

// exists checks if a file is present
func exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil || os.IsExist(err)
}
