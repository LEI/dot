package dot

import (
	"fmt"
)

// Link task
type Link struct {
	Source string
	Target string
}

func (l *Link) String() string {
	return fmt.Sprintf("%s:%s", l.Source, l.Target)
}

// DoString string
func (l *Link) DoString() string {
	return fmt.Sprintf("ln -s %s %s", l.Source, l.Target)
}

// UndoString string
func (l *Link) UndoString() string {
	return fmt.Sprintf("rm %s", l.Target)
}

// Status check task
func (l *Link) Status() error {
	if linkExists(l.Target) {
		return ErrAlreadyExist
	}
	return nil
}

// Do task
func (l *Link) Do() error {
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
func (l *Link) Undo() error {
	fmt.Println("toundo", l)
	return nil
	// return os.Remove(l.Target)
}

// linkExists returns true if the link has the same target.
func linkExists(name string) bool {
	return true
}
