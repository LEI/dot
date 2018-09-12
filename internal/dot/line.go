package dot

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

// Line task
type Line struct {
	Task   `mapstructure:",squash"` // Action, If, OS
	Target string                   // Target file
	Data   string                   // Line content
}

// NewLine task
func NewLine(s, d string) *Line {
	return &Line{Target: s, Data: d}
}

func (t *Line) String() string {
	s := fmt.Sprintf("%s:%s", t.Target, t.Data)
	switch Action {
	case "install":
		s = fmt.Sprintf("echo '%s' >> %s", t.Data, tildify(t.Target))
	case "remove":
		s = fmt.Sprintf("sed -i '#^%s$#d' %s", t.Data, tildify(t.Target))
	}
	return s
}

// Init task
func (t *Line) Init() error {
	// ...
	return nil
}

// Status check task
func (t *Line) Status() error {
	exists, err := lineExists(t.Target, t.Data)
	if err != nil {
		return err
	}
	if exists {
		return ErrExist
	}
	return nil
}

// Do task
func (t *Line) Do() error {
	if err := t.Status(); err != nil {
		switch err {
		case ErrExist, ErrSkip:
			return nil
		default:
			return err
		}
	}
	var err error
	lines := []string{}
	if exists(t.Target) {
		lines, err = getLines(t.Target)
		if err != nil {
			return err
		}
	}
	// Add line
	lines = append(lines, t.Data) // +"\n"
	output := strings.Join(lines, "\n")
	// Write target file
	if err := ioutil.WriteFile(t.Target, []byte(output), defaultFileMode); err != nil {
		return err
	}
	return nil
}

// Undo task
func (t *Line) Undo() error {
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
	lines, err := getLines(t.Target)
	if err != nil {
		return err
	}
	// if len(lines) == 0 {
	// 	return nil
	// }
	index := -1 // First match
	for i, s := range lines {
		if strings.Contains(s, t.Data) {
			index = i
			break
		}
	}
	if index < 0 {
		return fmt.Errorf("%s: line '%v' not found", t.Target, t.Data)
		// return &os.PathError{Op: "line", Path: t.Target, Err: ErrNotExist}
	}
	// Remove line
	lines = append(lines[:index], lines[index+1:]...)
	output := strings.Join(lines, "\n")
	// Write target file
	if err := ioutil.WriteFile(t.Target, []byte(output), defaultFileMode); err != nil {
		return err
	}
	return nil
}

// lineExists returns true if the target file contains the line.
func lineExists(target, data string) (bool, error) {
	if !exists(target) {
		return false, nil
	}
	lines, err := getLines(target)
	if err != nil {
		return false, err
	}
	for _, t := range lines {
		if strings.Contains(t, data) {
			return true, nil
		}
	}
	return false, nil
}

func getLines(file string) (lines []string, err error) {
	input, err := ioutil.ReadFile(file)
	str := string(input)
	// if err != nil && os.IsExist(err) {
	if err != nil && os.IsExist(err) {
		return
	}
	if str != "" {
		lines = strings.Split(str, "\n")
	}
	return
}
