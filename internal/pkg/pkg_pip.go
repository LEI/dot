package pkg

import (
	"fmt"
	"os/exec"
	"runtime"

	"github.com/LEI/dot/internal/shell"
)

// https://pip.pypa.io/en/stable/reference
var pip = &Pm{
	Bin:     "pip",
	Install: "install",
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
		// "--upgrade",
	},
	RemoveOpts: []string{
		"--yes",
	},
	// DryRunOpts: []string{},
	Init: func(m *Pm) error {
		// TODO: check action == "install" and if pip is up to date
		// /!\ sudo is needed on linux unless:
		// curl https://bootstrap.pypa.io/get-pip.py -o get-pip.py
		// python get-pip.py --user
		bin := "sudo"
		args := []string{"pip", "install", "--upgrade", "pip"}
		// if runtime.GOOS == "linux" {
		// 	opts = append(opts, "--user")
		// }
		// bin, args, err := getBin(m, opts)
		// if err != nil {
		// 	return err
		// }
		fmt.Printf("$ %s %s\n", bin, shell.FormatArgs(args))
		return execManagerCommand(m, bin, args...)
	},
	// FIXME: python2 -c 'import neovim' did not work until
	// pip2 uninstall neovim && pip2 install neovim
	Has: func(m *Pm, pkgs []string) (bool, error) {
		opts := []string{"show"}
		opts = append(opts, pkgs...)
		cmd := exec.Command(m.Bin, opts...)
		err := cmd.Run()
		return err == nil, nil
	},
}

var pip2 = &Pm{}
var pip3 = &Pm{}

func init() {
	// if runtime.GOOS != "darwin" {
	// 	// brew install python
	// }
	if runtime.GOOS == "linux" {
		pip.InstallOpts = append(pip.InstallOpts, "--user")
	}
	// windows: python -m pip ...

	*pip2 = *pip
	pip2.Bin = "pip2"

	*pip3 = *pip
	pip3.Bin = "pip3"
}
