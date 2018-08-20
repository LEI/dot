package pkg

import (
	"fmt"
	"os"
	"os/exec"
)

// https://manpages.debian.org/stretch/apt/apt-get.8.en.html
var aptGet = &Pm{
	Sudo:    true,
	Bin:     "apt-get",
	Install: "install",
	Remove:  "remove",
	Opts: []string{
		// -qqy
		"--assume-yes",
		"--no-install-recommends",
		"--no-install-suggests",
		"--quiet",
		"--quiet",
	},
	Has: func(m *Pm, pkgs []string) (bool, error) {
		// dpkg-query -l "$package" | grep -q ^.i
		c1 := exec.Command("dpkg-query", append([]string{"-Wf'${db:Status-abbrev}'"}, pkgs...)...)
		c2 := exec.Command("grep", "-q", "^i")

		c2.Stdin, _ = c1.StdoutPipe()
		c2.Stdout = os.Stdout

		if err := c2.Start(); err != nil {
			fmt.Println("grep start failed:", err)
			return false, err
		}
		if err := c1.Run(); err != nil {
			// fmt.Println("dpkg-query run failed:", err)
			return false, nil // err
		}
		if err := c2.Wait(); err != nil {
			// if exiterr, ok := err.(*exec.ExitError); ok {
			// 	if status, ok := exiterr.Sys().(syscall.WaitStatus); ok && status.ExitStatus() == 1 {
			// 		return false, nil
			// 	}
			// }
			return false, nil // err
		}
		return true, nil
	},
}
