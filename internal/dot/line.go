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

// Status check task
func (l *Line) Status() error {
	ok, err := lineExists(l.Target, l.Data)
	if err != nil {
		return err
	}
	if ok {
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
		if err != ErrAlreadyExist {
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
