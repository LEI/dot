package pkg

import "os/exec"

var gem = &Pm{
	Bin:     "gem",
	Install: "install",
	Remove:  "uninstall",
	Opts: []string{
		// "--no-verbose",
		"--quiet",
	},
	InstallOpts: []string{
		// "--bindir", "/usr/local/bin", // darwin?
		// "--install-dir", "/usr/local",
		// "--no-document", // "rdoc,ri",
		// "--no-post-install-message",
	},
	RemoveOpts: []string{
		// "--install-dir", "/usr/local",
	},
	Init: func(m *Pm) error {
		opts := []string{"update", "--system"}
		// "--bindir", "/usr/local/bin"
		// "--silent"
		bin, opts, err := getBin(m, opts)
		if err != nil {
			return err
		}
		return execManagerCommand(m, bin, opts...)
	},
	Has: func(m *Pm, pkgs []string) (bool, error) {
		opts := []string{"list", "--exact", "--installed"} // --local?
		opts = append(opts, pkgs...)
		cmd := exec.Command(m.Bin, opts...)
		err := cmd.Run()
		return err == nil, nil
	},
}
