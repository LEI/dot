package git

import (
	"os"
	"os/exec"
)

func Exec(args ...string) error {
	// fmt.Printf("git %s\n", strings.Join(args, " "))
	cmd := exec.Command("git", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
