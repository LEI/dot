package pkg

import (
	"os/exec"
)

// https://www.archlinux.org/pacman/pacman.8.html
var pacman = &Pm{
	Sudo: true,
	Bin:  "pacman",
	Acts: map[string]interface{}{
		"install": "--sync",   // -S
		"remove":  "--remove", // -R
	},
	DryRun: []string{"--print"},
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
		// "--unneeded",
	},
	// {
	// 	Args: []string{"--quiet"},
	// 	// If:   []string{"{{eq .Verbose 0}}"},
	// 	HasIf: types.HasIf{If: []string{"{{eq .Verbose 0}}"}},
	// },
	Has: func(name string) (bool, error) {
		//fmt.Printf("pacman -Qqi %s\n", name)
		// Search locally installed packages
		cmd := exec.Command("pacman", "-Qqi", "--noconfirm", name)
		//cmd.Stdout = os.Stdout
		//cmd.Stderr = os.Stderr
		//cmd.Stdin = os.Stdin
		err := cmd.Run()
		if err == nil {
			return true, nil
		}
		// Search installed groups
		//fmt.Printf("pacman -Qqg %s\n", name)
		cmd = exec.Command("pacman", "-Qqg", "--noconfirm", name)
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
	Sudo: false,
	Bin:  "yaourt",
	Acts: map[string]interface{}{
		"install": "--sync",   // -S
		"remove":  "--remove", // -R
	},
	Opts: []string{
		"--noconfirm",
		// "--sysupgrade", // -u
	},
	Has: func(name string) (bool, error) {
		//fmt.Printf("> yaourt -Qqi %s\n", name)
		// Search locally installed packages
		cmd := exec.Command("yaourt", "-Qqi", "--noconfirm", name)
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
