package pkg

import (
	"fmt"
	"os"
	"os/exec"
)

// https://wiki.alpinelinux.org/wiki/Alpine_Linux_package_management
var apk = &Pm{
	Sudo: true,
	Bin:  "apk",
	Acts: map[string]interface{}{
		"install": "add",
		"remove":  "del",
	},
	Opts: []string{
		"--no-cache",
		"--no-progress",
		"--quiet",
		// --simulate (if DryRun)
		"--update",
		// "--upgrade",
	},
	Has: func(name string) (bool, error) {
		c1 := exec.Command("apk", "search", "--exact", name)
		c2 := exec.Command("grep", "-q", name)

		c2.Stdin, _ = c1.StdoutPipe()
		c2.Stdout = os.Stdout

		if err := c2.Start(); err != nil {
			fmt.Println("grep start failed:", err)
			return false, err
		}
		if err := c1.Run(); err != nil {
			return false, nil // err
		}
		if err := c2.Wait(); err != nil {
			return false, nil // err
		}
		return true, nil
	},
}
