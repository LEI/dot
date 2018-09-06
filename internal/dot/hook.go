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
func (h *Hook) build() (string, []string) {
	if h.Shell == "" {
		h.Shell = defaultShell
	}
	args := []string{"-c", "set -e; " + h.Command}
	return h.Shell, args
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
	// fmt.Printf("EXEC DO HOOK: %q\n", h.Command)
	bin, args := h.build()
	cmd := exec.Command(bin, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Dir = h.ExecDir
	cmd.Env = os.Environ()
	for k, v := range *h.Env {
		v = env.ExpandEnvVar(k, v, *h.Env)
		// fmt.Printf("%s %s=%q\n", hookEnvPrefix, k, v)
		cmd.Env = append(cmd.Env, k+"="+v)
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
	bin, args := h.build()
	// fmt.Printf("EXEC UNDO HOOK: %q\n", h.Command)
	cmd := exec.Command(bin, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Dir = h.ExecDir
	for k, v := range *h.Env {
		v = env.ExpandEnvVar(k, v, *h.Env)
		// fmt.Printf("%s %s=%q\n", hookEnvPrefix, k, v)
		cmd.Env = append(cmd.Env, k+"="+v)
	}
	return cmd.Run()
}
