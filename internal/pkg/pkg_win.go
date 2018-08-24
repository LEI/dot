package pkg

// Note: do not use the special _windows suffix for now.

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// https://github.com/chocolatey/choco/wiki/CommandsReference
var choco = &Pm{
	// Sudo:    false,
	Bin:     "choco",
	Install: "install",
	Remove:  "uninstall",
	DryRun:  []string{"--noop"}, // --what-if
	Opts: []string{
		"--no-progress",
		"--yes", // --confirm
	},
	// Init: func() error {
	// 	// https://chocolatey.org/docs/installation
	// 	return nil
	// },
	Has: func(m *Pm, pkgs []string) (bool, error) {
		// opts := []string{"info", "--local-only"}
		opts := []string{"search", "--exact", "--local-only"}
		opts = append(opts, m.Opts...)
		opts = append(opts, pkgs...)
		// fmt.Printf("$ %s %s\n", m.Bin, strings.Join(opts, " "))
		var buf bytes.Buffer
		cmd := exec.Command(m.Bin, opts...)
		cmd.Stdout = &buf
		// cmd.Stdout = os.Stdout
		// cmd.Stderr = os.Stderr
		// cmd.Stdin = os.Stdin
		if err := cmd.Run(); err != nil {
			return false, err
		}
		out := buf.String()
		lines := strings.Split(out, "\n")
		if len(lines) != 3 || lines[2] != "0 packages installed.\n" {
			return false, nil
		}
		return true, nil
	},
}

// https://github.com/mobile-shell/mosh/blob/master/appveyor.yml
// C:\cygwin64\...
// c:/cygwin64...
// /cygdrive/c/cygwin64/setup-x86_64.exe
// --quiet-mode --no-shortcuts --upgrade-also --packages
// --download --local-install --packages

// https://github.com/transcode-open/apt-cyg
var aptCyg = &Pm{
	AllowFailure: true,
	// Sudo:    false,
	Bin:     "apt-cyg",
	Install: "install",
	Remove:  "remove",
	// DryRun:  []string{},
	// Opts: []string{},
	Init: func() error {
		// fmt.Println("$ apt-cyg --version")
		cmd := exec.Command("apt-cyg", "--version")
		// cmd.Stdout = os.Stdout
		// cmd.Stderr = os.Stderr
		// cmd.Stdin = os.Stdin
		if err := cmd.Run(); err != nil {
			// Not in %PATH%
			fmt.Fprintf(os.Stderr, "apt-cyg --version: %s", err) // return err
		}
		return nil
	},
	Has: func(m *Pm, pkgs []string) (bool, error) {
		// cygcheck --list-package ...
		opts := []string{"show"}
		// opts = append(opts, m.Opts...)
		opts = append(opts, pkgs...)
		fmt.Printf("$ %s %s\n", m.Bin, strings.Join(opts, " "))
		cmd := exec.Command(m.Bin, opts...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Stdin = os.Stdin
		// if err := cmd.Run(); err != nil {
		// 	return false, err
		// }
		// return true, nil
		err := cmd.Run()
		return err == nil, err
	},
}
