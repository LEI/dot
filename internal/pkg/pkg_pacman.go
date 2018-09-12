package pkg

import (
	"os/exec"
)

var pacman, yaourt *Pm

// https://www.archlinux.org/pacman/pacman.8.html
func init() {
	pacman = &Pm{
		Sudo:    true,
		Bin:     "pacman",
		Install: "--sync",   // -S
		Remove:  "--remove", // -R
		Opts: []string{
			"--noconfirm",
			"--noprogressbar",
		},
		InstallOpts: []string{
			"--needed",
			"--quiet",
			"--refresh",    // -y
			"--sysupgrade", // -u
		},
		RemoveOpts: []string{
			"--recursive", // -s
			// "--unneeded",
		},
		DryRunOpts: []string{"--print"},
		// {
		// 	Args: []string{"--quiet"},
		// 	// If:   []string{"{{eq .Verbose 0}}"},
		// 	HasIf: types.HasIf{If: []string{"{{eq .Verbose 0}}"}},
		// },
		Has: func(pkgs []string) (bool, error) {
			//fmt.Printf("pacman -Qqi %s\n", name)
			// Search locally installed packages
			cmd := exec.Command("pacman", append([]string{"-Qqi", "--noconfirm"}, pkgs...)...)
			//cmd.Stdout = Stdout
			//cmd.Stderr = Stderr
			//cmd.Stdin = Stdin
			err := cmd.Run()
			if err == nil {
				return true, nil
			}
			// Search installed groups
			//fmt.Printf("pacman -Qqg %s\n", name)
			cmd = exec.Command("pacman", append([]string{"-Qqg", "--noconfirm"}, pkgs...)...)
			//cmd.Stdout = Stdout
			//cmd.Stderr = Stderr
			//cmd.Stdin = Stdin
			if errg := cmd.Run(); errg == nil {
				// fmt.Fprintf(Stderr, "%s\n", err)
				return false, nil
			}
			// fmt.Fprintf(Stderr, "%s\n", err)
			return false, nil
		},
	}

	// https://archlinux.fr/man/yaourt.8.html
	yaourt = &Pm{
		Sudo:    false,
		Bin:     "yaourt",
		Install: "--sync",   // -S
		Remove:  "--remove", // -R
		Opts: []string{
			"--noconfirm",
			// "--sysupgrade", // -u
		},
		Has: func(pkgs []string) (bool, error) {
			//fmt.Printf("> yaourt -Qqi %s\n", name)
			// Search locally installed packages
			cmd := exec.Command("yaourt", append([]string{"-Qqi", "--noconfirm"}, pkgs...)...)
			//cmd.Stdout = Stdout
			//cmd.Stderr = Stderr
			//cmd.Stdin = Stdin
			if err := cmd.Run(); err != nil {
				// fmt.Fprintf(Stderr, "%s\n", err)
				return false, nil
			}
			return true, nil
		},
	}
}
