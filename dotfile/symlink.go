package dotfile

import (
	// "fmt"
	"os"
)

func IsSymlink(fi os.FileInfo) bool {
	return fi != nil && fi.Mode()&os.ModeSymlink != 0
}

func IsLink(path string) (string, error) {
	fi, err := os.Lstat(path)
	if err != nil && os.IsExist(err) {
		return path, err
	}
	if !IsSymlink(fi) {
		return "", nil
	}
	real, err := os.Readlink(path)
	return real, err
}

func InstallSymlink(source, target string, backup func(string, string) (bool, error)) (bool, error) {
	link, err := IsLink(target)
	if err != nil && os.IsExist(err) {
		return false, err
	}
	if link == source { // Already linked to this source
		// logger.Infof("# ln -s %s %s\n", source, target)
		return false, nil
	}
	if link != "" {
		ok, err := backup(target, link)
		if err != nil {
			return false, err
		}
		if !ok {
			return false, nil
		}
	}
	fi, err := os.Stat(target)
	if err != nil && os.IsExist(err) {
		return false, err
	}
	if fi != nil && backup != nil {
		ok, err := backup(target, "")
		if err != nil {
			return false, err
		}
		if !ok {
			return false, nil
		}
	}
	if DryRun {
		return true, nil
	}
	err = os.Symlink(source, target)
	if err != nil {
		return false, err
	}
	return true, nil
}

func RemoveSymlink(source, target string) (bool, error) {
	_, err := IsLink(target)
	if err != nil && os.IsExist(err) {
		return false, err
	}
	if err != nil && os.IsNotExist(err) {
		return false, nil
	}
	if DryRun {
		return true, nil
	}
	return true, os.Remove(target)
}
