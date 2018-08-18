package pkg

// https://wiki.alpinelinux.org/wiki/Alpine_Linux_package_management
var apk = &Pm{
	Sudo: true,
	Bin:  "apk",
	Acts: map[string]interface{}{
		"install": "add",
		"remove":  "del",
	},
	Opts: []*Opt{
		{
			Args: []string{
				"--no-cache",
				"--quiet",
				"--update",
			},
		},
	},
}
