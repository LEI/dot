package pkg

import (
	"fmt"
	"os"
	"os/exec"
)

// https://github.com/chocolatey/choco/wiki/CommandsReference
var choco = &Pm{
	Sudo:    true,
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
		fmt.Printf("$ choco info %s\n", pkgs)
		// opts := []string{"search", "--exact", "--local-only"}
		opts := []string{"info", "--local-only"}
		opts = append(opts, m.Opts...)
		opts = append(opts, pkgs...)
		cmd := exec.Command(m.Bin, opts...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Stdin = os.Stdin
		err := cmd.Run()
		return err == nil, nil
	},
}
