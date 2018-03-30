package helpers

import (
	// "fmt"
	"io/ioutil"
	"os"
	"strings"
)

var ()

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
