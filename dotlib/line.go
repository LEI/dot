package dotlib

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

// LineTask struct
type LineTask struct {
	File, Line string
	Task
}

// Install line
func (l *LineTask) Install() error {
	changed, err := LineInFile(l.File, l.Line)
	if err != nil {
		return err
	}
	prefix := "# "
	if changed {
		prefix = ""
	}
	fmt.Printf("%secho '%s' >> %s\n", prefix, l.Line, l.File)
	return nil
}

// Remove line
func (l *LineTask) Remove() error {
	changed, err := LineOutFile(l.File, l.Line)
	if err != nil {
		return err
	}
	prefix := "# "
	if changed {
		prefix = ""
	}
	fmt.Printf("%ssed -i '#^%s$#d' %s\n", prefix, l.Line, l.File)
	return nil
}

// Line task
func LineInFile(file string, line string) (bool, error) {
	input, err := ioutil.ReadFile(file)
	if err != nil && os.IsExist(err) {
		return false, err
	}
	str := string(input)
	lines := []string{}
	if str != "" {
		lines = strings.Split(str, "\n")
	}
	// fmt.Printf("IN: [%+v] (%d)\n", string(input), len(lines))
	for _, l := range lines {
		if strings.Contains(l, line) {
			return false, nil
		}
	}
	lines = append(lines, line) // +"\n"
	output := strings.Join(lines, "\n")
	// fmt.Printf("OUT: [%+v] (%d)\n", output, len(lines))
	err = ioutil.WriteFile(file, []byte(output), FileMode)
	if err != nil {
		return false, err
	}
	return true, nil
}

func LineOutFile(file string, line string) (bool, error) {
	input, err := ioutil.ReadFile(file)
	if err != nil && os.IsExist(err) {
		return false, err
	}
	str := string(input)
	if str == "" {
		return false, nil
	}
	lines := strings.Split(str, "\n")
	// fmt.Printf("IN: [%+v] (%d)\n", string(input), len(lines))
	index := -1 // First match
	for i, l := range lines {
		if strings.Contains(l, line) {
			index = i
			break;
		}
	}
	if index < 0 {
		return false, nil
	}
	// Remove line
	lines = append(lines[:index], lines[index+1:]...)
	output := strings.Join(lines, "\n")
	// fmt.Printf("OUT: [%+v] (%d)\n", output, len(lines))
	err = ioutil.WriteFile(file, []byte(output), FileMode)
	if err != nil {
		return false, err
	}
	return true, nil
}
