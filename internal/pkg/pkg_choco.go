package pkg

// https://github.com/chocolatey/choco/wiki/CommandsInstall
// https://github.com/chocolatey/choco/wiki/CommandsUninstall
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
		return false, nil
	},
}
