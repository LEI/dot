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
	/* Init: func(m *Pm) error {
		// TODO: check action == "install"
		opts := []string{"update", "--quiet"}
		bin, args, err := getBin(m, opts)
		if err != nil {
			return err
		}
		fmt.Printf("$ %s %s\n", bin, shell.FormatArgs(args))
		return execManagerCommand(m, bin, args...)
		// return execCommand("apt-get", "update", "--quiet")
	}, */
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

// https://wiki.termux.com/wiki/Package_Management
var termux = &Pm{}

func init() {
	*termux = *aptGet
	termux.Sudo = false
	// termux.Bin = "pkg"
	// termux.Init = func(m *Pm) error {
	// 	return nil
	// }
	termux.Has = func(m *Pm, pkgs []string) (bool, error) {
		opts := append([]string{"-Wf'${db:Status-abbrev}'"}, pkgs...)
		cmd := exec.Command("dpkg-query", opts...)
		// cmd.Stdout = os.Stdout
		// cmd.Stderr = os.Stderr
		// cmd.Stdin = os.Stdin
		err := cmd.Run()
		return err == nil, nil
	}
}
