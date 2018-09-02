package pkg

import (
	"fmt"
	"os/exec"

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
		opts := []string{"install", "--upgrade", "pip"}
		bin, args, err := getBin(m, opts)
		if err != nil {
			return err
		}
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
	// if runtime.GOOS == "darwin" {
	// 	pip.Sudo = true
	// }

	*pip2 = *pip
	pip2.Bin = "pip2"

	*pip3 = *pip
	pip3.Bin = "pip3"
}
