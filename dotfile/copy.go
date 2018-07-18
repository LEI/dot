package dotfile

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	// "path"
	// "strings"

	"github.com/LEI/dot/utils"
)

// CopyTask struct
type CopyTask struct {
	Source, Target string
	Task
}

// Do ...
func (t *CopyTask) Do(a string) error {
	return do(t, a)
}

// Install copy
func (t *CopyTask) Install() error {
	if err := createBaseDir(t.Target); err != nil && err != ErrDirShouldExist {
		return err
	}
	changed, err := Copy(t.Source, t.Target)
	if err != nil {
		return err
	}
	prefix := ""
	if !changed {
		prefix = "# "
	}
	fmt.Printf("%scp %s %s\n", prefix, t.Source, t.Target)
	return nil
}

// Remove copy
func (t *CopyTask) Remove() error {
	changed, err := Uncopy(t.Source, t.Target)
	if err != nil {
		return err
	}
	prefix := ""
	if !changed {
		prefix = "# "
	}
	fmt.Printf("%srm %s\n", prefix, t.Target)
	if RemoveEmptyDirs {
		if err := removeBaseDir(t.Target); err != nil {
			return err
		}
	}
	return nil
}

// Copy task
// https://stackoverflow.com/a/21067803/7796750
func Copy(src, dst string) (bool, error) {
	// TODO: check cache (see tplCache)
	if utils.Exist(dst) {
		ok, err := checkCopy(src, dst)
		if err != nil {
			return false, err
		}
		if ok {
			return false, nil
		}
		if !ok {
			return false, fmt.Errorf("Different copy target: %s", dst)
		}
	}
	/*fi, err := os.Stat(src)
	if err != nil {
		return false, err
	}
	if !fi.Mode().IsRegular() {
		// cannot copy non-regular files (e.g., directories, symlinks, devices, etc.)
		return false, fmt.Errorf("Copy: non-regular source file %s (%q)", fi.Name(), fi.Mode().String())
	}*/
	in, err := os.Open(src)
	if err != nil {
		return false, err
	}
	defer in.Close()
	out, err := os.Create(dst)
	if err != nil {
		return false, err
	}
	// defer func() {
	// 	cerr := out.Close()
	// 	if err == nil {
	// 		err = cerr
	// 	}
	// }()
	if _, err = io.Copy(out, in); err != nil {
		return false, err
	}
	err = out.Sync()
	if err := out.Close(); err != nil {
		return false, err
	}
	b, err := ioutil.ReadFile(dst)
	if err != nil && os.IsExist(err) {
		return false, err
	}
	c := string(b)
	if err := dotCache.Put(dst, c); err != nil {
		return false, err
	}
	return true, nil

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
}

// Uncopy task
func Uncopy(src, dst string) (bool, error) {
	if !utils.Exist(dst) {
		return false, nil
	}
	ok, err := checkCopy(src, dst)
	if err != nil {
		return false, err
	}
	if !ok {
		return false, fmt.Errorf("Different uncopy target: %s", dst)
	}
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
	if DryRun {
		return true, nil
	}
	if err := os.Remove(dst); err != nil {
		return false, err
	}
	if err := dotCache.Del(dst); err != nil {
		return false, err
	}
	return true, nil
}

func checkCopy(src, dst string) (bool, error) {
	ok, err := utils.Compare(src, dst)
	if err != nil {
		return false, err
	}
	if ok {
		return true, nil
	}
	b, err := ioutil.ReadFile(dst)
	if err != nil && os.IsExist(err) {
		return false, err
	}
	c := string(b)
	ok, err = dotCache.Validate(dst, c)
	if err != nil {
		return false, err
	}
	if !ok {
		q := fmt.Sprintf("Overwrite existing copy target: %s", dst)
		if !AskConfirmation(q) {
			fmt.Fprintf(os.Stderr, "Skipping copy %s because its target exists: %s", src, dst)
			return false, nil
		}
		if err := Backup(dst); err != nil {
			return false, err
		}
	}
	return true, nil
}
