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
	// s := ""
	s := strings.TrimRight(h.Command, "\n")
	if strings.Contains(s, "\n") && !strings.HasPrefix(s, "(") {
		s = fmt.Sprintf("(%s)", s)
	}
	switch h.GetAction() {
	case "install":
	case "remove":
		s = "# noop"
	}
	return s
}

// Status check task
func (h *Hook) Status() error {
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
	// fmt.Println("EXEC:", h.Command)
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
	return fmt.Errorf("not implemented")
}
