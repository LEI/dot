package dotfile

import (
	"fmt"
	// "io/ioutil"
	// "os"
	// "path"
	// "strings"
)

// CopyTask struct
type CopyTask struct {
	Source, Target string
	Task
}

// Install copy
func (t *CopyTask) Install() error {
	changed, err := Copy(t)
	if err != nil {
		return err
	}
	prefix := "# "
	if changed {
		prefix = ""
	}
	fmt.Printf("%scp %s %s\n", prefix, t.Source, t.Target)
	return nil
}

// Remove copy
func (t *CopyTask) Remove() error {
	changed, err := Uncopy(t)
	if err != nil {
		return err
	}
	prefix := "# "
	if changed {
		prefix = ""
	}
	fmt.Printf("%srm %s\n", prefix, t.Target)
	return nil
}

// Copy task
func Copy(t *CopyTask) (bool, error) {
	// str, err := t.Parse()
	// if err != nil {
	// 	return false, err
	// }
	// b, err := ioutil.ReadFile(t.Target)
	// if err != nil && os.IsExist(err) {
	// 	return false, err
	// }
	// if str == string(b) {
	// 	return false, nil
	// }
	// if DryRun {
	// 	return true, nil
	// }
	// if err := ioutil.WriteFile(t.Target, []byte(str), FileMode); err != nil {
	// 	return false, err
	// }
	return true, nil
}

// Uncopy task
func Uncopy(t *CopyTask) (bool, error) {
	// str, err := t.Parse()
	// if err != nil {
	// 	return false, err
	// }
	// b, err := ioutil.ReadFile(t.Target)
	// if err != nil && os.IsExist(err) {
	// 	return false, err
	// }
	// if len(b) == 0 { // Empty file
	// 	return false, nil
	// }
	// if str != string(b) { // Mismatching content
	// 	fmt.Printf("Warn: mismatching content %s\n", t.Target)
	// 	return false, nil
	// }
	// if DryRun {
	// 	return true, nil
	// }
	// if err := os.Remove(t.Target); err != nil {
	// 	return false, err
	// }
	return true, nil
}
