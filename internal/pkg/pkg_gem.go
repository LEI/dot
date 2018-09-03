package pkg

import (
	"os/exec"
	"runtime"
)

// https://bundler.io/docs.html

// https://guides.rubygems.org/command-reference
var gem = &Pm{
	Bin:     "gem",
	Install: "install",
	// Install: func(m *Pm, pkgs ...string) string {
	// 	if Upgrade {
	// 		ok, err := m.Has(m, pkgs)
	// 		if err == nil && ok {
	// 			return "update" // upgrade
	// 		}
	// 	}
	// 	return "install"
	// },
	Remove: "uninstall",
	Opts: []string{
		// "--no-verbose",
		"--quiet",
	},
	InstallOpts: []string{
		// "--bindir", "/usr/local/bin", // darwin?
		// "--install-dir", "/usr/local",
		"--no-document", // rdoc,ri
		// "--no-post-install-message",
	},
	RemoveOpts: []string{
		// "--all" // matching versions
		"--executables", // without confirmation
		// "--install-dir", "/usr/local",
	},
	Env: map[string]string{
		// "GEM_HOME": "",
	},
	/* Init: func(m *Pm) error {
		// TODO: check action == "install"

		// // export GEM_HOME="$(ruby -e 'print Gem.user_dir')"
		// cmd := exec.Command("ruby", "-e", "print Gem.user_dir")
		// out, err := cmd.Output()
		// if err != nil {
		// 	return err
		// }
		// m.Env["GEM_HOME"] = string(out)

		opts := []string{"update", "--system"} // local?
		// "--bindir", "/usr/local/bin"
		// "--silent"
		bin, args, err := getBin(m, opts)
		if err != nil {
			return err
		}
		// str := strings.TrimSuffix(string(out), "\n")
		// if os.Getenv("GEM_HOME") == "" {
		// 	err = os.Setenv("GEM_HOME", str)
		// 	if err != nil {
		// 		return err
		// 	}
		// }
		fmt.Printf("$ %s %s\n", bin, shell.FormatArgs(args))
		return execManagerCommand(m, bin, args...)
	}, */
	Has: func(m *Pm, pkgs []string) (bool, error) {
		opts := []string{"list", "--exact", "--installed"} // , "--local"}
		opts = append(opts, pkgs...)
		cmd := exec.Command(m.Bin, opts...)
		err := cmd.Run()
		return err == nil, nil
	},
}

func init() {
	if runtime.GOOS == "linux" {
		// gem.Opts = append(gem.Opts, []string{"--bindir", "/usr/local/bin"}...)

		// (Un)install in user's home directory instead of GEM_HOME.
		gem.Opts = append(gem.Opts, "--user-install") // --local
	}
	// FIXME alpine, debian, arch?
	// if !gem.Sudo && runtime.GOOS == "linux" {
	// 	gem.Sudo = true
	// 	// bin, opts, err := getBin(gem, gem.Opts)
	// 	// if err != nil {
	// 	// 	log.Fatal(err)
	// 	// }
	// 	// gem.Bin = bin
	// 	// gem.Opts = opts
	// }
}
