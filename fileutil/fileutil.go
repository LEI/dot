package fileutil

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

func Exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func MakeDir(path string) error {
	err := os.MkdirAll(path, 0755)
	if err != nil {
		return err
	}
	return nil
}

func RemoveDir(path string) error {
	err := os.Remove(path)
	if err != nil {
		return err
	}
	return nil
}

func Link(source string, target string) error {
	fi, err := os.Lstat(target)
	if err != nil && os.IsExist(err) {
		return err
	}
	if fi != nil && (fi.Mode() & os.ModeSymlink != 0) {
		link, err := os.Readlink(target)
		if err != nil {
			return err
		}
		if link == source {
			fmt.Printf("%s already linked to %s\n", target, source)
			return nil
		}
		// TODO check broken symlink?
		msg := fmt.Sprintf("%s exists, linked to %s, replace with %s?", target, link, source)
		if ok := confirm(msg); ok {
			err := os.Remove(target)
			if err != nil {
				return err
			}
		}
	} else if fi != nil {
		backup := target+".backup"
		msg := fmt.Sprintf("%s exists, move to %s and replace with %s?", target, backup, source)
		if ok := confirm(msg); ok {
			err := os.Rename(target, target+".backup")
			if err != nil {
				return err
			}
		}
	}
	err = os.Symlink(source, target)
	if err != nil {
		return err
	}
	return nil
}

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
			fmt.Println("LineInFile: OK")
			return nil
		}
	} else {
		fi, err := os.Create(path)
		if err != nil {
			return err
		}
		defer fi.Close()
	}
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

func confirm(str string) bool {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Printf("%s [y/n]: ", str)
		// if ForceYes {
		// 	fmt.Printf("%s", "Forced")
		// 	return true
		// }
		res, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		res = strings.ToLower(strings.TrimSpace(res))
		switch res {
		case "y", "yes":
			return true
		case "n", "no":
			return false
		}
		// TODO limit retries
	}
}
