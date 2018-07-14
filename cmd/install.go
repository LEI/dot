package cmd

import (
	"fmt"
	// "os"
	// "github.com/jessevdk/go-flags"
)

// InstallCmd ...
type InstallCmd struct {
	BaseTaskCmd
}

// var installCmd InstallCmd

// Execute ...
func (cmd *InstallCmd) Execute(args []string) error {
	if Verbose > 2 {
		fmt.Println("# InstallCmd", args)
	}
	if err := cmd.Copy.Execute([]string{}); err != nil {
		fmt.Println("copy err", err)
		// os.Exit(1)
	}
	if err := cmd.Line.Execute([]string{}); err != nil {
		fmt.Println("line err", err)
	}
	if err := cmd.Link.Execute([]string{}); err != nil {
		fmt.Println("link err", err)
	}
	if err := cmd.Template.Execute([]string{}); err != nil {
		fmt.Println("template err", err)
	}
	return nil
}
