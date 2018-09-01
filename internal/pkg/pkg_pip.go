package pkg

import "os/exec"

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
