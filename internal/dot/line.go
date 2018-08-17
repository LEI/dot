package dot

import (
	"fmt"
	"path/filepath"
)

// // Lines task map
// type Lines map[string]string

// Lines task list
type Lines []*Line

// Line task
type Line struct {
	Task
	Target string // Target file
	Data   string // Line content
}

func (l *Line) String() string {
	return fmt.Sprintf("%s:%s", l.Target, l.Data)
}

// DoString string
func (l *Line) DoString() string {
	return fmt.Sprintf("echo '%s' >> %s", l.Data, l.Target)
}

// UndoString string
func (l *Line) UndoString() string {
	return fmt.Sprintf("sed -i '#^%s$#d' %s", l.Data, l.Target)
}

// Prepare task
func (l *Line) Prepare(target string) error {
	if !filepath.IsAbs(l.Target) {
		l.Target = filepath.Join(target, l.Target)
	}
	return nil
}

// Status check task
func (l *Line) Status() error {
	if lineExists(l.Target) {
		return ErrAlreadyExist
	}
	return nil
}

// Do task
func (l *Line) Do() error {
	if err := l.Status(); err != nil {
		if err == ErrAlreadyExist {
			return nil
		}
		return err
	}
	fmt.Println("todo", l)
	return nil
}

// Undo task
func (l *Line) Undo() error {
	fmt.Println("toundo", l)
	return nil
	// return os.Remove(l.Target)
}

// lineExists returns true if the target file contains the line.
func lineExists(name string) bool {
	return true
}
