package pkg

import (
	"fmt"
	"os/exec"
)

var aptGet = &Pm{}

// https://manpages.debian.org/stretch/apt/apt-get.8.en.html
func init() {
	aptGet = &Pm{
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
		Init: func() error {
			return aptGet.Exec([]string{"update", "--quiet"}...)
		},
		Has: func(pkgs []string) (bool, error) {
			// dpkg-query -l "$package" | grep -q ^.i
			c1 := exec.Command("dpkg-query", append([]string{"-Wf'${db:Status-abbrev}'"}, pkgs...)...)
			c2 := exec.Command("grep", "-q", "^i")

			c2.Stdin, _ = c1.StdoutPipe()
			c2.Stdout = Stdout

			if err := c2.Start(); err != nil {
				fmt.Println("grep start failed:", err)
				return false, err
			}
			if err := c1.Run(); err != nil {
				// fmt.Println("dpkg-query run failed:", err)
				return false, nil // err
			}
			if err := c2.Wait(); err != nil {
				// fmt.Println("grep run failed:", err)
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
}

// https://wiki.termux.com/wiki/Package_Management
var termux = &Pm{}

func init() {
	*termux = *aptGet
	termux.Sudo = false
	termux.Bin = "pkg"
	termux.Install = "install"
	termux.Remove = "uninstall"
	termux.Init = func() error {
		return termux.Exec([]string{"update", "--quiet"}...)
	}
	// termux.Has = func(pkgs []string) (bool, error) {
	// 	opts := append([]string{"-Wf'${db:Status-abbrev}'"}, pkgs...)
	// 	cmd := exec.Command("dpkg-query", opts...)
	// 	// cmd.Stdout = Stdout
	// 	// cmd.Stderr = Stderr
	// 	// cmd.Stdin = Stdin
	// 	err := cmd.Run()
	// 	return err == nil, nil
	// }
}
