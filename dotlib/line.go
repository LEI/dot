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
	fmt.Printf("%secho '%s' >> \"%s\"\n", prefix, l.Line, l.File)
	return nil
}

// Remove line
func (l *LineTask) Remove() error {
	prefix := "TODO: "
	c := fmt.Sprintf("sed %s\n", l.File)
	fmt.Printf("%s%s\n", prefix, c)
	return nil
}

// Line task
func LineInFile(file string, line string) (bool, error) {
	input, err := ioutil.ReadFile(file)
	if err != nil && os.IsExist(err) {
		return false, err
	}
	lines := strings.Split(string(input), "\n")
	for _, l := range lines {
		if strings.Contains(l, line) {
			return false, nil
		}
	}
	lines = append(lines, line+"\n")
	output := strings.Join(lines, "\n")
	err = ioutil.WriteFile(file, []byte(output), FileMode)
	if err != nil {
		return false, err
	}
	return true, nil
}
