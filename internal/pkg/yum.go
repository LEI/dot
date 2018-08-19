package pkg

var yum = &Pm{
	Sudo: true,
	Bin:  "yum",
	Acts: map[string]interface{}{
		"install": "install",
		"remove":  "remove",
	},
	Opts: []string{
		"--assumeyes",
		// "--error=0",
		"--quiet",
	},
}
