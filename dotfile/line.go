package dotfile

import (
	"io/ioutil"
	"os"
	"strings"
)

func LineInFile(path string, line string) (changed bool, err error) {
	fi, err := os.Stat(path)
	if err != nil && os.IsExist(err) {
		return false, err
	}
	if fi != nil {
		contains, err := hasLineInFile(path, line)
		if err != nil || contains {
			return false, err
		}
	}
	// fi, err := os.Create(path)
	// if err != nil {
	// 	return false, err
	// }
	// defer fi.Close()
	return appendStringToFile(path, line+"\n")
}

func hasLineInFile(path string, line string) (bool, error) {
	b, err := ioutil.ReadFile(path)
	if err != nil || len(b) == 0 {
		return false, err
	}
	// defer b.Close()?
	fileContent := string(b)
	if fileContent != "" {
		for _, str := range strings.Split(fileContent, "\n") {
			if strings.Contains(str, line) {
				return true, nil
			}
		}
	}
	return false, err
}

func appendStringToFile(path string, str string) (changed bool, err error) {
	// fi, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, os.ModeAppend)
	fi, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0611)
	defer fi.Close()
	if os.IsNotExist(err) {
		fi, err = os.Create(path)
	}
	if err != nil {
		return false, err
	}
	defer fi.Close()
	n, err := fi.WriteString(str)
	if err != nil {
		return false, err
	}
	if n > 0 {
		return true, nil
	}
	return false, nil
}
