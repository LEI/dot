package command

import (
	"fmt"
	"gopkg.in/urfave/cli.v1"
	// "github.com/urfave/cli"
)

var Sync = cli.Command{
	Name: "sync",
	Usage: "Synchronize dotfiles",
	Action: func(c *cli.Context) error {
		fmt.Println("Sync files", c.Args())
		return nil
	},
}
