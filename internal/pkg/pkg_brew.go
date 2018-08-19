package pkg

import (
	"os/exec"
)

// TODO: Brewfile

// https://docs.brew.sh/Manpage
var brew = &Pm{
	Bin: "brew",
	Install: func(m *Pm, pkgs ...string) string {
		// opts := []string{"ls", "--versions", name}
		// err := exec.Command("brew", opts...).Run()
		if Upgrade {
			ok, err := m.Has(pkgs)
			if err == nil && ok {
				return "upgrade"
			}
		}
		return "install"
	},
	Remove: "uninstall",
	Opts:   []string{"--quiet"},
	Env: map[string]string{
		// "HOMEBREW_NO_ANALYTICS": "1",
		"HOMEBREW_NO_AUTO_UPDATE": "1",
		// "HOMEBREW_NO_EMOJI": "1",
	},
	Init: func() error {
		return execCommand("brew", "update", "--quiet")
	},
	Has: func(pkgs []string) (bool, error) {
		// fmt.Printf("brew ls --versions '%s'\n", name)
		cmd := exec.Command("brew", append([]string{"ls", "--versions"}, pkgs...)...)
		// cmd.Stdout = os.Stdout
		// cmd.Stderr = os.Stderr
		// cmd.Stdin = os.Stdin
		err := cmd.Run()
		return err == nil, err
	},
}

var brewCask = &Pm{
	Bin:     "brew",
	Sub:     []string{"cask"},
	Install: "install",
	Remove:  "uninstall",
	Has: func(pkgs []string) (bool, error) {
		// fmt.Printf("brew cask ls --versions '%s'\n", name)
		cmd := exec.Command("brew", append([]string{"cask", "ls", "--versions"}, pkgs...)...)
		// cmd.Stdout = os.Stdout
		// cmd.Stderr = os.Stderr
		// cmd.Stdin = os.Stdin
		err := cmd.Run()
		return err == nil, err
	},
}
