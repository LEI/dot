package dotfile

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

// Do ...
func (t *LineTask) Do(a string) error {
	return do(t, a)
}

// Install line
func (t *LineTask) Install() error {
	if err := createBaseDir(t.File); err != nil && err != ErrDirShouldExist {
		return err
	}
	changed, err := LineInFile(t.File, t.Line)
	if err != nil {
		return err
	}
	prefix := ""
	if !changed {
		prefix = "# "
	}
	fmt.Printf("%secho '%s' >> %s\n", prefix, t.Line, t.File)
	return nil
}

// Remove line
func (t *LineTask) Remove() error {
	changed, err := LineOutFile(t.File, t.Line)
	if err != nil {
		return err
	}
	prefix := ""
	if !changed {
		prefix = "# "
	}
	fmt.Printf("%ssed -i '#^%s$#d' %s\n", prefix, t.Line, t.File)
	if RemoveEmptyDirs {
		if err := removeBaseDir(t.File); err != nil {
			return err
		}
	}
	return nil
}

// LineInFile task
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
	if DryRun {
		return true, nil
	}
	// Add line
	lines = append(lines, line) // +"\n"
	output := strings.Join(lines, "\n")
	// fmt.Printf("OUT: [%+v] (%d)\n", output, len(lines))
	err = ioutil.WriteFile(file, []byte(output), FileMode)
	if err != nil {
		return false, err
	}
	return true, nil
}

// LineOutFile task
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
			break
		}
	}
	if index < 0 {
		return false, nil
	}
	if DryRun {
		return true, nil
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
