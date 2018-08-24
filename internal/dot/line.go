package dot

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

var (
	defaultFileMode os.FileMode = 0644
)

// Line task
type Line struct {
	Task   `mapstructure:",squash"` // Action, If, OS
	Target string                   // Target file
	Data   string                   // Line content
}

func (l *Line) String() string {
	s := fmt.Sprintf("%s:%s", l.Target, l.Data)
	switch l.GetAction() {
	case "install":
		s = fmt.Sprintf("echo '%s' >> %s", l.Data, tildify(l.Target))
	case "remove":
		s = fmt.Sprintf("sed -i '#^%s$#d' %s", l.Data, tildify(l.Target))
	}
	return s
}

// Type task name
func (l *Line) Type() string {
	return "line" // in file
}

// Status check task
func (l *Line) Status() error {
	exists, err := lineExists(l.Target, l.Data)
	if err != nil {
		return err
	}
	if exists {
		return ErrExist
	}
	return nil
}

// Do task
func (l *Line) Do() error {
	if err := l.Status(); err != nil {
		switch err {
		case ErrExist, ErrSkip:
			return nil
		default:
			return err
		}
	}
	var err error
	lines := []string{}
	if exists(l.Target) {
		lines, err = getLines(l.Target)
		if err != nil {
			return err
		}
	}
	// Add line
	lines = append(lines, l.Data) // +"\n"
	output := strings.Join(lines, "\n")
	// Write target file
	if err := ioutil.WriteFile(l.Target, []byte(output), defaultFileMode); err != nil {
		return err
	}
	return nil
}

// Undo task
func (l *Line) Undo() error {
	if err := l.Status(); err != nil {
		switch err {
		case ErrExist:
			// continue
		case ErrSkip:
			return nil
		default:
			return err
		}
	}
	if err := l.Status(); err != nil {
		if err != ErrExist {
			return err
		}
	}
	lines, err := getLines(l.Target)
	if err != nil {
		return err
	}
	// if len(lines) == 0 {
	// 	return nil
	// }
	index := -1 // First match
	for i, s := range lines {
		if strings.Contains(s, l.Data) {
			index = i
			break
		}
	}
	if index < 0 {
		return fmt.Errorf("%s: line '%v' not found", l.Target, l.Data)
		// return &os.PathError{Op: "line", Path: l.Target, Err: ErrNotExist}
	}
	// Remove line
	lines = append(lines[:index], lines[index+1:]...)
	output := strings.Join(lines, "\n")
	// Write target file
	if err := ioutil.WriteFile(l.Target, []byte(output), defaultFileMode); err != nil {
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
	for _, l := range lines {
		if strings.Contains(l, data) {
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
