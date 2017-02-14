package main

import (
	// "fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func linesInFiles(src string, dest string, lines map[string]string) error {
	for file, line := range lines {
		dst := filepath.Join(dest, file)

		fi, err := os.Stat(dst)
		if err != nil && os.IsExist(err) {
			return err
		}
		if fi != nil { // err != nil && os.IsExist(err)
			contains, err := hasLineInFile(dst, line)
			if err != nil {
				return err
			}
			if contains {
				logSuccess.Printf("'%s' => %s\n", line, dst)
				return nil
			}
		} else { // os.IsNotExist(err)
			fi, err := os.Create(dst)
			if err != nil {
				return err
			}
			defer fi.Close()
		}

		err = appendStringToFile(dst, line+"\n")
		if err != nil {
			return err
		}

		logSuccess.Printf("'%s' -> %s\n", line, dst)
	}
	return nil
}

func linesOutFiles(src string, dest string, lines map[string]string) error {
	for file, line := range lines {
		dst := filepath.Join(dest, file)

		fi, err := os.Stat(dst)
		if err != nil && os.IsExist(err) {
			return err
		}

		if err != nil && os.IsNotExist(err) {
			// logSuccess.Printf("'%s' <= %s\n", line, dst)
			return nil
		}

		if fi != nil {
			contains, err := hasLineInFile(dst, line)
			if err != nil {
				return err
			}

			if contains {
				logInfo.Printf("%s: %s <- %s\n", "TODO LineOutFile", line, dst)
				// logSuccess.Printf("'%s' <- %s\n", line, dst)
			}
			// } else {
			// 	logSuccess.Printf("'%s' <= %s\n", line, dst)
		}
	}
	return nil
}

func hasLineInFile(path string, line string) (bool, error) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return false, err
	}
	content := string(b)
	if content != "" {
		for _, str := range strings.Split(content, "\n") {
			if strings.Contains(str, line) {
				// fmt.Printf("%s: already contains the line '%s'\n", path, line)
				return true, nil
			}
		}
	}
	return false, err
}

func appendStringToFile(path string, text string) error {
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

	// fmt.Fprintf(fi, line+"\n")
	_, err = fi.WriteString(text)
	if err != nil {
		return err
	}
	return nil
}
