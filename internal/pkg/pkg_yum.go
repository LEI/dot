package pkg

var yum = &Pm{
	Sudo:    true,
	Bin:     "yum",
	Install: "install",
	Remove:  "remove",
	Opts: []string{
		"--assumeyes",
		// "--error=0",
		"--quiet",
	},
}
