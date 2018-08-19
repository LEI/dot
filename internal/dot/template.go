package dot

import (
	"fmt"
)

// Template task
type Template struct {
	Task        `mapstructure:",squash"` // Action, If, OS
	Source      string
	Target      string
	Env         map[string]string
	Vars        map[string]interface{}
	IncludeVars string `mapstructure:"include_vars"`
}

func (t *Template) String() string {
	return fmt.Sprintf("%s:%s", t.Source, t.Target)
}

// DoString string
func (t *Template) DoString() string {
	return fmt.Sprintf("gotpl %s %s", t.Source, t.Target)
}

// UndoString string
func (t *Template) UndoString() string {
	return fmt.Sprintf("rm %s", t.Target)
}

// Status check task
func (t *Template) Status() error {
	if templateExists(t.Target) {
		return ErrAlreadyExist
	}
	return nil
}

// Do task
func (t *Template) Do() error {
	if err := t.Status(); err != nil {
		switch err {
		case ErrAlreadyExist, ErrSkip:
			return nil
		default:
			return err
		}
	}
	fmt.Println("todo", t)
	return nil
}

// Undo task
func (t *Template) Undo() error {
	if err := t.Status(); err != nil {
		switch err {
		case ErrSkip:
			return nil
		case ErrAlreadyExist:
			// continue
		default:
			return err
		}
	}
	return nil
	// return os.Remove(t.Target)
}

// templateExists returns true if the template is the same.
func templateExists(name string) bool {
	return true
}
