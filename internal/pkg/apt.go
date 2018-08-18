package pkg

import (
	"os"
	"os/exec"
)

// https://manpages.debian.org/stretch/apt/apt-get.8.en.html
var aptGet = &Pm{
	Sudo: true,
	Bin:  "apt-get",
	Acts: map[string]interface{}{
		"install": "install",
		"remove":  "remove",
	},
	Opts: []*Opt{
		{
			// Args: []string{"-qqy"},
			Args: []string{
				"--assume-yes",
				"--no-install-recommends",
				"--no-install-suggests",
				"--quiet",
				"--quiet",
			},
		},
	},
	Has: func(name string) (bool, error) {
		args := []string{"-Wf'${db:Status-abbrev}'", name}
		cmd := exec.Command("dpkg-query", args...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Stdin = os.Stdin
		err := cmd.Run()
		return err == nil, err
	},
}
