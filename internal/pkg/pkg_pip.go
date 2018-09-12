package pkg

import (
	"os/exec"
	"runtime"
)

var pip, pip2, pip3 *Pm

// https://pip.pypa.io/en/stable/reference
func init() {
	pip = &Pm{
		Bin:     "pip",
		Install: "install", // "--upgrade",
		Remove:  "uninstall",
		Opts: []string{
			// "--no-cache",
			// "--prefix", "/usr/local",
			"--quiet", // TODO flag
			// "--quiet",
			// "--quiet",
			// "--requirement", "requirements.txt",
		},
		InstallOpts: []string{ // --noinput?
			// "--progress-bar", "off",
		},
		RemoveOpts: []string{
			"--yes",
		},
		// DryRunOpts: []string{},
		/* Init: func() error {
			// TODO: check action == "install" and if pip is up to date
			// FIXME: /!\ sudo is needed on linux unless:
			// curl https://bootstrap.pypa.io/get-pip.py -o get-pip.py
			// python get-pip.py --user
			// bin := "sudo"
			// args := []string{pip.Bin, "install", "--upgrade", "pip"}
			opts := []string{"install", "--upgrade", "pip"}
			if runtime.GOOS == "linux" {
				opts = append(opts, "--user")
			}
			return pip.Exec(opts...)
		}, */
		// FIXME: python2 -c 'import neovim' did not work until
		// pip2 uninstall neovim && pip2 install neovim
		Has: func(pkgs []string) (bool, error) {
			opts := []string{"show"}
			opts = append(opts, pkgs...)
			cmd := exec.Command(pip.Bin, opts...)
			err := cmd.Run()
			return err == nil, nil
		},
	}
	// if runtime.GOOS != "darwin" {
	// 	// brew install python
	// }
	if runtime.GOOS == "linux" {
		pip.InstallOpts = append(pip.InstallOpts, "--user")
	}
	// windows: python -m pip ...

	pip2 = pip
	pip2.Bin = "pip2"

	pip3 = pip
	pip3.Bin = "pip3"
}
