package dot

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/LEI/dot/internal/env"
	"github.com/LEI/dot/internal/shell"
)

var (
	defaultShell string

	hookEnvPrefix = "$"
)

func init() {
	defaultShell = shell.Get()
}

// Hook command to execute
// FIXME: fmt.Println(t) -> stack exceeds limit
type Hook struct {
	Task    `mapstructure:",squash"` // Action, If, OS
	Command string
	Shell   string
	Env     *Env
	ExecDir string
}

// NewHook task
func NewHook(s string) *Hook {
	return &Hook{Command: s}
}

// Init hook: set default shell for next commands
// and return arguments to be executed
func (t *Hook) buildCmd() (*exec.Cmd, error) {
	bin := t.Shell
	if bin == "" {
		bin = defaultShell
	}
	c := t.Command
	args := []string{"-c", "set -e; " + c}
	// fmt.Printf("EXEC HOOK: %q\n", t.Command)
	cmd := exec.Command(bin, args...)
	cmd.Stdout = Stdout
	cmd.Stderr = os.Stderr
	cmd.Dir = t.ExecDir
	for k, v := range *t.Env {
		v = env.ExpandEnvVar(k, v, *t.Env)
		// fmt.Printf("%s %s=%q\n", hookEnvPrefix, k, v)
		cmd.Env = append(cmd.Env, k+"="+v)
	}
	return cmd, nil
}

func (t *Hook) String() string {
	// TODO: verbosity >= 2?
	// bin, args := t.build()
	// s := fmt.Sprintf("%s %s", bin, shell.FormatArgs(args))
	s := strings.TrimRight(t.Command, "\n")
	if strings.Contains(s, "\n") && !strings.HasPrefix(s, "(") {
		s = fmt.Sprintf("(%s)", s)
	}
	return s
}

// Init task
func (t *Hook) Init() error {
	// ...
	return nil
}

// Status check task
func (t *Hook) Status() error {
	// t.Command == "" &&
	// if t.URL != "" && t.Dest != "" {
	// 	if exists(t.Dest) {
	// 		return ErrExist
	// 	}
	// }
	// Always run hooks
	switch Action {
	case "install":
		// return nil
	case "remove":
		return ErrExist
	}
	return nil
}

// Do task
func (t *Hook) Do() error {
	if err := t.Status(); err != nil {
		switch err {
		case ErrExist, ErrSkip:
			return nil
		default:
			return err
		}
	}
	cmd, err := t.buildCmd()
	if err != nil {
		return err
	}
	return cmd.Run()
}

// Undo task (non applicable)
func (t *Hook) Undo() error {
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
	// if t.URL != "" && t.Dest != "" {
	// 	// TODO: check remote file?
	// 	return os.Remove(t.Dest)
	// }
	// cmd, err := t.buildCmd()
	// if err != nil {
	// 	return err
	// }
	// return cmd.Run()
	return fmt.Errorf("not implemented")
}
