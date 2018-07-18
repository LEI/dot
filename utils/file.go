package utils

// https://github.com/sergi/go-diff
// https://kgrz.io/reading-files-in-go-an-overview.html

import (
	"fmt"
	"io/ioutil"
	"os"
)

var (
	// ErrEmpty ...
	ErrEmpty = fmt.Errorf("file empty")
)

// Read ...
func Read(s string) ([]byte, error) {
	return ioutil.ReadFile(s)
}

// Exist ...
func Exist(s string) bool {
	_, err := os.Stat(s)
	return err == nil || os.IsExist(err)
	// return !os.IsNotExist(err)
}

// IsFile ...
func IsFile(s string) bool {
	fi, err := os.Stat(s)
	return !os.IsNotExist(err) && !fi.IsDir()
}

// func writeFile(fileName string, data []byte) error {
// 	os.MkdirAll(path.Dir(fileName), os.ModePerm)
// 	return ioutil.WriteFile(fileName, data, 0655)
// }

// // SameContent ...
// func SameContent(src, dst string) (bool, error) {
// 	b, err := ioutil.ReadFile(src)
// 	if err != nil { // && os.IsNotExist(err) {
// 		return false, err
// 	}
// 	_, ok, err := CompareFileContent(dst, string(b))
// 	return ok, err
// }

// Compare ...
func Compare(src, dst string) (bool, error) {
	// Retrieve source file content
	b, err := ioutil.ReadFile(src)
	if err != nil {
		return false, err
	}
	c := string(b)
	// Compare with target content
	content, ok, err := CompareFileContent(dst, c)
	// fmt.Println("------------------", dst, err)
	if err != nil && os.IsExist(err) {
		return false, err
	}
	if content == "" {
		return ok, ErrEmpty
	}
	return ok, nil
}

// CompareFileContent ...
func CompareFileContent(s, str string) (string, bool, error) {
	fi, err := os.Open(s)
	if err != nil {
		return "", false, err
	}
	// if err != nil && os.IsExist(err) {
	// 	return "", false, err
	// }
	// if fi == nil {
	// 	return "", false, nil
	// }
	defer fi.Close()
	stat, err := fi.Stat()
	if err != nil && os.IsExist(err) {
		return "", false, err
	}
	if stat != nil && !stat.Mode().IsRegular() {
		return "", false, fmt.Errorf("Not a regular file: %s (%q)", stat.Name(), stat.Mode().String())
	}
	// b, err := ioutil.ReadFile(s)
	// if err != nil && os.IsExist(err) {
	// 	return false, err
	// }
	b, err := ioutil.ReadAll(fi)
	content := string(b)
	if err != nil {
		return content, false, err
	}
	// fmt.Println("COMPARED FILE CONTENT", s, len(str), "vs", len(content), "->", content == str)
	return content, content == str, nil
}
