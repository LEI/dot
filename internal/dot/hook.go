package dot

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

var (
	defaultExecShell = "sh"
)

// Hook command to execute
type Hook struct {
	Task    // Action, If, OS
	Command string
	Shell   string
	ExecDir string
}

func (h *Hook) String() string {
	s := strings.TrimRight(h.Command, "\n")
	return fmt.Sprintf("%s", s)
}

// DoString string
func (h *Hook) DoString() string {
	return h.String()
}

// UndoString string
func (h *Hook) UndoString() string {
	return "" // h.String()
}

// Status check task
func (h *Hook) Status() error {
	return nil
}

// Do task
func (h *Hook) Do() error {
	if err := h.Status(); err != nil {
		switch err {
		case ErrAlreadyExist, ErrSkip:
			return nil
		default:
			return err
		}
	}
	if h.Shell == "" {
		h.Shell = defaultExecShell
	}
	cmd := exec.Command(h.Shell, []string{"-c", h.Command}...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Dir = h.ExecDir
	return nil
}

// Undo task (non applicable)
func (h *Hook) Undo() error {
	if err := h.Status(); err != nil {
		switch err {
		case ErrSkip:
			return nil
		case ErrAlreadyExist:
			// continue
		default:
			return err
		}
	}
	return fmt.Errorf("not implemented")
}
