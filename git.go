package main

import (
	"fmt"
	"os"
	"os/exec"
)

func gitCloneOrPull(name string, repo string, dir string) error {
	_, err := os.Stat(dir)
	if err != nil && os.IsNotExist(err) {
		err := os.MkdirAll(dir, 0755)
		if err != nil {
			return err
		}
		err := gitExec("clone", repo, dir)
		if err != nil {
			return err
		}
	} else {
		// err := gitExec("-C", dir, "status")
		err := gitExec("-C", dir, "pull")
		// "--git-dir", dir+"/.git",
		// "--work-tree", dir,
		if err != nil {
			return err
		}
	}
	return nil
}

func gitExec(args ...string) error {
	cmd := exec.Command("git", args...)
	out, err := cmd.CombinedOutput()
	if len(out) > 0 {
		fmt.Printf("%s: %s", name, out)
	}
	if err != nil {
		return err
	}
	return nil
}
