package system

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

// CheckLine ... (verify/validate)
func CheckLine(file, line string) error {
	if !Exists(file) {
		return nil
	}
	lines, err := getLines(file)
	if err != nil {
		return err
	}
	for _, l := range lines {
		if strings.Contains(l, line) {
			return ErrLineAlreadyExist
		}
	}
	return nil
}

// CreateLine ...
func CreateLine(file, line string) (err error) {
	lines := []string{}
	if Exists(file) {
		lines, err = getLines(file)
		if err != nil {
			return err
		}
	}
	// Add line
	lines = append(lines, line) // +"\n"
	output := strings.Join(lines, "\n")
	if DryRun {
		return nil
	}
	if err := ioutil.WriteFile(file, []byte(output), FileMode); err != nil {
		return err
	}
	return nil
}

// RemoveLine ...
func RemoveLine(file, line string) error {
	// if !Exists(file) {
	// 	return nil
	// }
	lines, err := getLines(file)
	if err != nil {
		return err
	}
	// if len(lines) == 0 {
	// 	return nil
	// }
	index := -1 // First match
	for i, l := range lines {
		if strings.Contains(l, line) {
			index = i
			break
		}
	}
	if index < 0 {
		return fmt.Errorf("%s: line not found", file)
	}
	// Remove line
	lines = append(lines[:index], lines[index+1:]...)
	output := strings.Join(lines, "\n")
	if DryRun {
		return nil
	}
	if err := ioutil.WriteFile(file, []byte(output), FileMode); err != nil {
		return err
	}
	return nil
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
