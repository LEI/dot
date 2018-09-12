package pkg

import (
	"os/exec"
)

// Brewfile: brew bundle
// dot-osx/brew-unbundle.sh

var brew, brewCask *Pm

// https://docs.brew.sh/Manpage
func init() {
	brew = &Pm{
		Bin: "brew",
		Install: func(pkgs ...string) string {
			// opts := []string{"ls", "--versions", name}
			// err := exec.Command("brew", opts...).Run()
			if Upgrade {
				ok, err := brew.Has(pkgs)
				if err == nil && ok {
					return "upgrade"
				}
			}
			return "install"
		},
		Remove: "uninstall",
		// Opts:   []string{"--quiet"},
		Env: map[string]string{
			"HOMEBREW_NO_ANALYTICS":   "1",
			"HOMEBREW_NO_AUTO_UPDATE": "1",
			"HOMEBREW_NO_EMOJI":       "1",
			// "HOMEBREW_VERBOSE": "0",
		},
		Init: func() error {
			opts := []string{"update", "--quiet"}
			return brew.Exec(opts...)
		},
		Has: func(pkgs []string) (bool, error) {
			// fmt.Printf("brew ls --versions %s\n", pkgs)
			cmd := exec.Command(brew.Bin, append([]string{"ls", "--versions"}, pkgs...)...)
			// cmd.Stdout = Stdout
			// cmd.Stderr = Stderr
			// cmd.Stdin = Stdin
			err := cmd.Run()
			return err == nil, nil // err
		},
	}

	brewCask = &Pm{
		Bin:     "brew",
		Sub:     []string{"cask"},
		Install: "install",
		Remove:  "uninstall",
		Has: func(pkgs []string) (bool, error) {
			// fmt.Printf("brew cask ls --versions %s\n", pkgs)
			cmd := exec.Command(brewCask.Bin, append([]string{"cask", "ls", "--versions"}, pkgs...)...)
			// cmd.Stdout = Stdout
			// cmd.Stderr = Stderr
			// cmd.Stdin = Stdin
			err := cmd.Run()
			return err == nil, nil // err
		},
	}
}
