package pkg

// https://www.archlinux.org/pacman/pacman.8.html
var pacman = &Pm{
	Sudo: true,
	Bin:  "pacman",
	Acts: map[string]interface{}{
		"install": "--sync",   // -S
		"remove":  "--remove", // -R
	},
	Opts: []*Opt{
		{
			Args: []string{
				"--needed",
				"--noconfirm",
				"--noprogressbar",
				"--quiet",
				"--refresh",    // -y
				"--sysupgrade", // -u
			},
		},
		// {
		// 	Args: []string{"--quiet"},
		// 	// If:   []string{"{{eq .Verbose 0}}"},
		// 	HasIf: types.HasIf{If: []string{"{{eq .Verbose 0}}"}},
		// },
	},
}

// https://archlinux.fr/man/yaourt.8.html
var yaourt = &Pm{
	// Sudo: false,
	Bin: "yaourt",
	Acts: map[string]interface{}{
		"install": "--sync",   // -S
		"remove":  "--remove", // -R
	},
	Opts: []*Opt{
		{
			Args: []string{
				"--noconfirm",
				// "--sysupgrade", // -u
			},
		},
	},
}
