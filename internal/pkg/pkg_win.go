package pkg

// Note: do not use the special _windows suffix for now.

import (
	"bytes"
	"fmt"
	"os"
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

var cygwinSetup = []string{
	// Install lynx
	// "/c/cygwin64/setup-x86_64.exe --quiet-mode --no-shortcuts --upgrade-also --packages lynx",
	// "/c/cygwin64/bin/cygcheck -dc cygwin",
	// Install apt-cyg
	"curl -sSL https://rawgit.com/transcode-open/apt-cyg/master/apt-cyg -o apt-cyg",
	"install apt-cyg /bin",
}

// https://github.com/transcode-open/apt-cyg
var aptCyg = &Pm{
	// Sudo:    false,
	Bin:     "apt-cyg",
	Install: "install",
	Remove:  "remove",
	// DryRun:  []string{},
	// Opts: []string{},
	Init: func() error {
		// c := "if ! hash apt-cyg; then ...; fi"
		for _, c := range cygwinSetup {
			fmt.Println("$", c)
			cmd := exec.Command(shell.Get(), "-lc", c)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			cmd.Stdin = os.Stdin
			if err := cmd.Run(); err != nil {
				return err
			}
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
