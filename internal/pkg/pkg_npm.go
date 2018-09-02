package pkg

import (
	"os/exec"
)

var npm = &Pm{
	Bin: "npm",
	// https://docs.npmjs.com/cli/install
	Install: "install",
	// https://docs.npmjs.com/cli/uninstall
	Remove: "uninstall",
	Opts: []string{
		"--global",
	},
	// InstallOpts: []string{},
	// RemoveOpts:  []string{},
	DryRunOpts: []string{"--dry-run"},
	/* Init: func(m *Pm) error {
		// TODO: check action == "install" and if npm is up to date
		opts := []string{"install", "--global", "npm"}
		bin, args, err := getBin(m, opts)
		if err != nil {
			return err
		}
		fmt.Printf("$ %s %s\n", bin, shell.FormatArgs(args))
		return execManagerCommand(m, bin, args...)
	}, */
	Has: func(m *Pm, pkgs []string) (bool, error) {
		// npm info ... --json
		opts := []string{"list", "--global"}
		opts = append(opts, pkgs...)
		cmd := exec.Command(m.Bin, opts...)
		err := cmd.Run()
		return err == nil, nil
	},
}

// yarn
