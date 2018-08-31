package pkg

import "os/exec"

var npm = &Pm{
	Bin:     "npm",
	Install: "install",
	Remove:  "uninstall",
	Opts: []string{
		"--global",
	},
	InstallOpts: []string{},
	RemoveOpts:  []string{},
	DryRunOpts:  []string{"--dry-run"},
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
