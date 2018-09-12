package pkg

import (
	"os/exec"

	"github.com/LEI/dot/internal/shell"
)

var npm *Pm

// https://docs.npmjs.com/misc/config
// https://docs.npmjs.com/cli
func init() {
	npm = &Pm{
		Bin:     "npm",
		Shell:   shell.Get(),
		Install: "install",   // https://docs.npmjs.com/cli/install
		Remove:  "uninstall", // https://docs.npmjs.com/cli/uninstall
		Opts: []string{
			"--global",
			"--no-progress",
			"--quiet", // --slient
		},
		// InstallOpts: []string{
		// 	"--no-progress",
		// },
		// RemoveOpts: []string{
		// 	"--no-progress",
		// },
		DryRunOpts: []string{"--dry-run"},
		// Init: func() error {
		// 	// TODO: check up to date?
		// 	return npm.Exec([]string{"install", "--global", "npm"}...)
		// },
		Has: func(pkgs []string) (bool, error) {
			// npm info ... --json
			opts := []string{"list", "--global"}
			opts = append(opts, pkgs...)
			cmd := exec.Command(npm.Bin, opts...)
			err := cmd.Run()
			return err == nil, nil
		},
	}
}

// yarn
