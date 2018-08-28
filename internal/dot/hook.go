package dot

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/LEI/dot/internal/shell"
)

var (
	defaultShell string
)

func init() {
	defaultShell = shell.Get()
}

// Hook command to execute
type Hook struct {
	Task    `mapstructure:",squash"` // Action, If, OS
	Command string
	Shell   string
	ExecDir string
}

func (h *Hook) String() string {
	s := strings.TrimRight(h.Command, "\n")
	if strings.Contains(s, "\n") && !strings.HasPrefix(s, "(") {
		s = fmt.Sprintf("(%s)", s)
	}
	return s
}

// Status check task
func (h *Hook) Status() error {
	// Always run hooks
	switch h.GetAction() {
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
	if h.Shell == "" {
		h.Shell = defaultShell
	}
	// fmt.Printf("EXEC DO HOOK: %q\n", h.Command)
	cmd := exec.Command(h.Shell, []string{"-c", "set -e; " + h.Command}...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Dir = h.ExecDir
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
	if h.Shell == "" {
		h.Shell = defaultShell
	}
	// fmt.Printf("EXEC UNDO HOOK: %q\n", h.Command)
	cmd := exec.Command(h.Shell, []string{"-c", "set -e; " + h.Command}...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Dir = h.ExecDir
	return cmd.Run()
}
