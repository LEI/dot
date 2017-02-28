package dotfile

import (
	"bytes"
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
		b, err := ioutil.ReadFile(path)
		if err != nil || len(b) == 0 {
			return false, err
		}
		if Contains(string(b), line) {
			return false, err
		}
	}
	// fi, err := os.Create(path)
	// if err != nil {
	// 	return false, err
	// }
	// defer fi.Close()
	if DryRun {
		return true, nil
	}
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
		b, err := ioutil.ReadFile(path)
		if err != nil || len(b) == 0 {
			return false, err
		}
		if len(b) == 0 {
			if RemoveEmptyFile || prompt.Confirm("> Remove empty file %s?", path) {
				if DryRun {
					return true, nil
				}
				err := os.Remove(path)
				if err != nil {
					return false, err
				}
				return true, nil
			}
			return false, nil
		}
		if !Contains(string(b), line) {
			return false, err
		}
	}
	if DryRun {
		return true, nil
	}
	return removeStringInFile(path, line+"\n")
}

func Contains(file string, line string) bool {
	if file == "" {
		return false
	}
	for _, s := range strings.Split(file, "\n") {
		if strings.Contains(s, line) {
			return true
		}
	}
	return false
}

func appendStringInFile(path string, str string) (changed bool, err error) {
	// fi, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, os.ModeAppend)
	fi, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_RDWR, FileMode)
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
	err = ioutil.WriteFile(path, []byte(new), FileMode)
	if err != nil {
		return false, err
	}
	if len(new) == 0 {
		if RemoveEmptyFile || prompt.Confirm("> Remove empty file %s?", path) {
			err := os.Remove(path)
			if err != nil {
				return false, err
			}
		}
		return true, nil
	}
	return true, nil
}
