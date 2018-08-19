package pkg

import (
	"os/exec"
)

var yum = &Pm{
	Sudo:    true,
	Bin:     "yum",
	Install: "install",
	Remove:  "remove",
	Opts: []string{
		"--assumeyes",
		// "--error=0",
		"--quiet",
	},
	Has: func(name string) (bool, error) {
		// sudo yum info
		// yum -C list installed
		cmd := exec.Command("rpm", "-q", name) // --quiet
		// cmd.Stdout = os.Stdout
		// cmd.Stderr = os.Stderr
		// cmd.Stdin = os.Stdin
		err := cmd.Run()
		return err == nil, nil
	},
}
