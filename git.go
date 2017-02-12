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
		gitClone := exec.Command("git", "clone", repo, dir)
		out, err := gitClone.CombinedOutput()
		if len(out) > 0 {
			fmt.Printf("%s: %s", name, out)
		}
		if err != nil {
			return err
		}
	} else {
		gitPull := exec.Command("git",
			"--git-dir", dir+"/.git",
			"--work-tree", dir,
			"pull")
		out, err := gitPull.CombinedOutput()
		if len(out) > 0 {
			fmt.Printf("%s: %s", name, out)
		}
		if err != nil {
			return err
		}
	}
	return nil
}
