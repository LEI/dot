package pkg

import "os/exec"

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
	Has: func(pkgs []string) (bool, error) {
		// choco info
		cmd := exec.Command("choco", append([]string{"search", "--exact"}, pkgs...)...)
		// cmd.Stdout = os.Stdout
		// cmd.Stderr = os.Stderr
		// cmd.Stdin = os.Stdin
		err := cmd.Run()
		return err == nil, nil
	},
}
