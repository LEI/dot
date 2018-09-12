package pkg

import (
	"os/exec"
)

var yum *Pm

func init() {
	yum = &Pm{
		Sudo:    true,
		Bin:     "yum",
		Install: "install",
		Remove:  "remove",
		Opts: []string{
			"--assumeyes",
			// "--error=0",
			"--quiet",
		},
		Has: func(pkgs []string) (bool, error) {
			// sudo yum info
			// yum -C list installed
			cmd := exec.Command("rpm", append([]string{"-q"}, pkgs...)...) // --quiet
			// cmd.Stdout = Stdout
			// cmd.Stderr = Stderr
			// cmd.Stdin = Stdin
			err := cmd.Run()
			return err == nil, nil
		},
	}
}
