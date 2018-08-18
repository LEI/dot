package pkg

import (
	"os"
	"os/exec"
)

// https://docs.brew.sh/Manpage
var brew = &Pm{
	Bin: "brew",
	Acts: map[string]interface{}{
		"install": func(m *Pm, name string, opts ...string) string {
			// TODO filter strings.HasPrefix(opts, "-")?
			// opts := []string{"ls", "--versions", name}
			// err := exec.Command("brew", opts...).Run()
			if Upgrade {
				ok, err := m.Has(name)
				if err == nil && ok {
					return "upgrade"
				}
			}
			return "install"
		},
		"remove": "uninstall",
	},
	Opts: []*Opt{
		{
			Args: []string{"--quiet"},
		},
	},
	Env: map[string]string{
		// "HOMEBREW_NO_ANALYTICS": "1",
		"HOMEBREW_NO_AUTO_UPDATE": "1",
		// "HOMEBREW_NO_EMOJI": "1",
	},
	Init: func() error {
		return execute("brew", "update", "--quiet")
	},
	Has: func(name string) (bool, error) {
		args := []string{"ls", "--versions", name}
		cmd := exec.Command("brew", args...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Stdin = os.Stdin
		err := cmd.Run()
		return err == nil, err
	},
}

var brewCask = &Pm{
	Bin: "brew",
	Sub: []string{"cask"},
	Acts: map[string]interface{}{
		"install": "install",
		"remove":  "uninstall",
	},
}
