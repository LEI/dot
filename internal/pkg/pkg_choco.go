package pkg

import (
	"bytes"
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
		// fmt.Println(m.Bin, opts)
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
