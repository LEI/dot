package comp

// https://github.com/sergi/go-diff
// https://github.com/google/go-cmp

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
)

// Equals compare data
func Equals(a, b []byte) (bool, error) {
	ok := bytes.Compare(a, b) == 0
	return ok, nil
}

// PathEquals compares the content of a give path and data
func PathEquals(path string, a []byte) (bool, error) {
	b, err := ioutil.ReadFile(path)
	if err != nil { // && os.IsExist(err) {
		return false, err
	}
	return Equals(a, b)
}

// FileEquals compares the content of two given paths
func FileEquals(path1, path2 string) (bool, error) {
	b, err := ioutil.ReadFile(path1)
	if err != nil { // && os.IsExist(err) {
		return false, err
	}
	return PathEquals(path2, b)
}

// RegularFile compares content
func RegularFile(path string, a []byte) (bool, error) {
	fi, err := os.Open(path)
	if err != nil {
		return false, err
	}
	defer fi.Close()
	stat, err := fi.Stat()
	if err != nil && os.IsExist(err) {
		return false, err
	}
	if stat != nil && !stat.Mode().IsRegular() {
		// stat.Name(), stat.Mode().String()
		return false, fmt.Errorf("%s: not a regular file", path)
	}
	b, err := ioutil.ReadAll(fi)
	// content := string(b)
	if err != nil {
		return false, err
	}
	return Equals(a, b)
}
