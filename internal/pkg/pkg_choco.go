package pkg

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
		return false, nil // TODO grep -q - Chocolatey vX.X.X
		// opts := []string{"search", "--exact", "--local-only"}
		// opts = append(opts, m.Opts...)
		// opts = append(opts, pkgs...)
		// fmt.Println(m.Bin, opts)
		// cmd := exec.Command(m.Bin, opts...)
		// cmd.Stdout = os.Stdout
		// cmd.Stderr = os.Stderr
		// cmd.Stdin = os.Stdin
		// err := cmd.Run()
		// return err == nil, nil
	},
}
