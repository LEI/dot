package fileutil

import (
	"fmt"
	"github.com/LEI/dot/cli"
	"os"
)

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
		if ok := cli.Confirm(msg); ok {
			err := os.Remove(target)
			if err != nil {
				return err
			}
		}
	} else if fi != nil {
		backup := target+".backup"
		msg := fmt.Sprintf("%s exists, move to %s and replace with %s?", target, backup, source)
		if ok := cli.Confirm(msg); ok {
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
