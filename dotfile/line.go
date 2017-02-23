package dotfile

import (
	"bytes"
	// "fmt"
	"github.com/LEI/dot/prompt"
	"io/ioutil"
	"os"
	"strings"
)

var RemoveEmptyFile bool

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
	return appendStringInFile(path, line+"\n")
}

func LineOutFile(path string, line string) (changed bool, err error) {
	fi, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	if fi != nil {
		contains, err := hasLineInFile(path, line)
		if err != nil || !contains {
			return false, err
		}
	}
	return removeStringInFile(path, line+"\n")
}

func hasLineInFile(path string, line string) (bool, error) {
	b, err := ioutil.ReadFile(path)
	if err != nil || len(b) == 0 {
		return false, err
	}
	// defer b.Close()?
	content := string(b)
	if content != "" {
		for _, s := range strings.Split(content, "\n") {
			if strings.Contains(s, line) {
				return true, nil
			}
		}
	}
	return false, err
}

func appendStringInFile(path string, str string) (changed bool, err error) {
	// fi, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, os.ModeAppend)
	fi, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0611)
	defer fi.Close()
	// if err != nil && os.IsNotExist(err) {
	// 	fi, err = os.Create(path)
	// 	defer fi.Close()
	// }
	if err != nil {
		return false, err
	}
	n, err := fi.WriteString(str)
	if err != nil {
		return false, err
	}
	if n == 0 {
		return false, nil
	}
	return true, nil
}

func removeStringInFile(path string, str string) (changed bool, err error) {
	b, err := ioutil.ReadFile(path)
	if err != nil || len(b) == 0 {
		return false, err
	}
	// defer b.Close()?
	new := bytes.Replace(b, []byte(str), []byte(""), 1)
	// content := string(b)
	// new := strings.Replace(content, str, "", 1)
	if len(new) == len(b) {
		return false, nil
	}
	err = ioutil.WriteFile(path, []byte(new), 0611)
	if err != nil {
		return false, err
	}
	if len(new) == 0 {
		if RemoveEmptyFile || prompt.Confirm("> Remove empty file %s?", path) {
			err := os.Remove(path)
			if err != nil {
				return true, err
			}
		}
		return true, nil
	}
	return true, nil
}
