package pkg

var yum = &Pm{
	Sudo: true,
	Bin:  "yum",
	Acts: map[string]interface{}{
		"install": "install",
		"remove":  "remove",
	},
	Opts: []*Opt{
		{
			Args: []string{
				"--assumeyes",
				// "--error=0",
				"--quiet",
			},
		},
	},
}
