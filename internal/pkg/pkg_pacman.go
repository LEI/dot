package pkg

import (
	"os/exec"
)

// https://www.archlinux.org/pacman/pacman.8.html
var pacman = &Pm{
	Sudo:    true,
	Bin:     "pacman",
	Install: "--sync",   // -S
	Remove:  "--remove", // -R
	DryRun:  []string{"--print"},
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
	// {
	// 	Args: []string{"--quiet"},
	// 	// If:   []string{"{{eq .Verbose 0}}"},
	// 	HasIf: types.HasIf{If: []string{"{{eq .Verbose 0}}"}},
	// },
	Has: func(pkgs []string) (bool, error) {
		//fmt.Printf("pacman -Qqi %s\n", name)
		// Search locally installed packages
		cmd := exec.Command("pacman", append([]string{"-Qqi", "--noconfirm"}, pkgs...)...)
		//cmd.Stdout = os.Stdout
		//cmd.Stderr = os.Stderr
		//cmd.Stdin = os.Stdin
		err := cmd.Run()
		if err == nil {
			return true, nil
		}
		// Search installed groups
		//fmt.Printf("pacman -Qqg %s\n", name)
		cmd = exec.Command("pacman", append([]string{"-Qqg", "--noconfirm"}, pkgs...)...)
		//cmd.Stdout = os.Stdout
		//cmd.Stderr = os.Stderr
		//cmd.Stdin = os.Stdin
		if errg := cmd.Run(); errg == nil {
			// fmt.Fprintf(os.Stderr, "%s\n", err)
			return false, nil
		}
		// fmt.Fprintf(os.Stderr, "%s\n", err)
		return false, nil
	},
}

// https://archlinux.fr/man/yaourt.8.html
var yaourt = &Pm{
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
		//cmd.Stdout = os.Stdout
		//cmd.Stderr = os.Stderr
		//cmd.Stdin = os.Stdin
		if err := cmd.Run(); err != nil {
			// fmt.Fprintf(os.Stderr, "%s\n", err)
			return false, nil
		}
		return true, nil
	},
}
