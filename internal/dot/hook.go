package dot

import (
	"fmt"
)

// Hook command to execute
type Hook struct {
	Task
	Command string
	Shell   string
	Action  string // install, remove
	OS      []string
}

func (h *Hook) String() string {
	return fmt.Sprintf("%s", h.Command)
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
	// if hookExists(h.Target) {
	// 	return ErrAlreadyExist
	// }
	return nil
}

// Do task
func (h *Hook) Do() error {
	if err := h.Status(); err != nil {
		if err == ErrAlreadyExist {
			return nil
		}
		return err
	}
	fmt.Println("todo", h)
	return nil
}

// Undo task
func (h *Hook) Undo() error {
	return fmt.Errorf("not implemented")
}
