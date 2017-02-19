package fileutil

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

func LineInFile(path string, line string) error {
	fi, err := os.Stat(path)
	if err != nil && os.IsExist(err) {
		return err
	}
	if fi != nil {
		contains, err := IsLineInFile(path, line)
		if err != nil {
			return err
		}
		if contains {
			// Already exists
			return nil
		}
	} else {
		fi, err := os.Create(path)
		if err != nil {
			return err
		}
		defer fi.Close()
	}
	fmt.Printf("$ echo '%s' >> %s\n", line, path)
	err = AppendStringToFile(path, line+"\n")
	if err != nil {
		return err
	}
	return nil
}

func IsLineInFile(path string, line string) (bool, error) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return false, err
	}
	// defer b.Close()?
	fileContent := string(b)
	if fileContent != "" {
		for _, str := range strings.Split(fileContent, "\n") {
			if strings.Contains(str, line) {
				fmt.Printf("%s: already contains the line '%s'\n", path, line)
				return true, nil
			}
		}
	}
	return false, err
}

func AppendStringToFile(path string, str string) error {
	// fi, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, os.ModeAppend)
	fi, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0611)
	defer fi.Close()
	if os.IsNotExist(err) {
		fi, err = os.Create(path)
	}
	if err != nil {
		return err
	}
	defer fi.Close()
	_, err = fi.WriteString(str)
	if err != nil {
		return err
	}
	return nil
}
