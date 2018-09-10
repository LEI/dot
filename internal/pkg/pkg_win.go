package pkg

// Note: do not use the special _windows suffix for now.

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"

	"github.com/LEI/dot/internal/shell"
)

// https://github.com/chocolatey/choco/wiki/CommandsReference
var choco = &Pm{
	// Sudo:    false,
	Bin:     "choco",
	Install: "install",
	Remove:  "uninstall",
	Opts: []string{
		"--no-progress",
		"--yes", // --confirm
	},
	DryRunOpts: []string{"--noop"}, // --what-if
	// Init: func(m *Pm) error {
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
		// cmd.Stdout = Stdout
		// cmd.Stderr = Stderr
		// cmd.Stdin = Stdin
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
	//AllowFailure: true,
	Shell: shell.Get(),
	// Sudo:    false,
	Bin:     "apt-cyg",
	Install: "install",
	Remove:  "remove",
	// Opts: []string{},
	// DryRunOpts:  []string{},
	/* Init: func(m *Pm) error {
		// fmt.Println("$ apt-cyg --version")
		cmd := exec.Command(shell.Get(), "-c", "apt-cyg --version")
		// cmd := exec.Command("apt-cyg", "--version")
		cmd.Stdout = Stdout
		cmd.Stderr = Stderr
		cmd.Stdin = Stdin
		// if err := cmd.Run(); err != nil {
		// 	// Not in %PATH%
		// 	fmt.Printf("apt-cyg --version: error")
		// 	fmt.Fprintln(Stderr, err)
		// 	// return err
		// }
		return cmd.Run()
	}, */
	Has: func(m *Pm, pkgs []string) (bool, error) {
		// cygcheck --list-package ...
		opts := []string{"show"}
		// opts = append(opts, m.Opts...)
		opts = append(opts, pkgs...)
		fmt.Printf("$ %s %s\n", m.Bin, shell.FormatArgs(opts))
		err := execManagerCommand(m, m.Bin, opts...)
		// cmd := exec.Command(shell.Get(), "-c", m.Bin + opts...)
		// cmd.Stdout = Stdout
		// cmd.Stderr = Stderr
		// cmd.Stdin = Stdin
		// if err := cmd.Run(); err != nil {
		// 	fmt.Printf("$ %s %s\n", m.Bin, strings.Join(opts, " "))
		// 	fmt.Fprintln(Stderr, err)
		// 	return false, nil // err
		// // }
		// return true, nil
		return err == nil, err
	},
}
