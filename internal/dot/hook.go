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
// FIXME: fmt.Println(h) -> stack exceeds limit
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
func (h *Hook) buildCmd() (*exec.Cmd, error) {
	bin := h.Shell
	if bin == "" {
		bin = defaultShell
	}
	c := h.Command
	args := []string{"-c", "set -e; " + c}
	// fmt.Printf("EXEC HOOK: %q\n", h.Command)
	cmd := exec.Command(bin, args...)
	cmd.Stdout = Stdout
	cmd.Stderr = os.Stderr
	cmd.Dir = h.ExecDir
	for k, v := range *h.Env {
		v = env.ExpandEnvVar(k, v, *h.Env)
		// fmt.Printf("%s %s=%q\n", hookEnvPrefix, k, v)
		cmd.Env = append(cmd.Env, k+"="+v)
	}
	return cmd, nil
}

func (h *Hook) String() string {
	// TODO: verbosity >= 2?
	// bin, args := h.build()
	// s := fmt.Sprintf("%s %s", bin, shell.FormatArgs(args))
	s := strings.TrimRight(h.Command, "\n")
	if strings.Contains(s, "\n") && !strings.HasPrefix(s, "(") {
		s = fmt.Sprintf("(%s)", s)
	}
	return s
}

// Status check task
func (h *Hook) Status() error {
	// h.Command == "" &&
	// if h.URL != "" && h.Dest != "" {
	// 	if exists(h.Dest) {
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
func (h *Hook) Do() error {
	if err := h.Status(); err != nil {
		switch err {
		case ErrExist, ErrSkip:
			return nil
		default:
			return err
		}
	}
	cmd, err := h.buildCmd()
	if err != nil {
		return err
	}
	return cmd.Run()
}

// Undo task (non applicable)
func (h *Hook) Undo() error {
	if err := h.Status(); err != nil {
		switch err {
		case ErrExist:
			// continue
		case ErrSkip:
			return nil
		default:
			return err
		}
	}
	// if h.URL != "" && h.Dest != "" {
	// 	// TODO: check remote file?
	// 	return os.Remove(h.Dest)
	// }
	// cmd, err := h.buildCmd()
	// if err != nil {
	// 	return err
	// }
	// return cmd.Run()
	return fmt.Errorf("not implemented")
}
