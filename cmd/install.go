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

// Install ... var Install *flags.Command

// var installCommands = []*flags.Command{
// 	&flags.Command{
// 		Name: "copy",
// 		Group: &flags.Group{
// 			ShortDescription: "Copy files",
// 			LongDescription: "Copy files",
// 			// data: &copyCmd,
// 		},
// 	},
// }

// Execute ...
func (cmd *InstallCmd) Execute(args []string) error {
	// fmt.Println("exec install cmd", args)
	if err := cmd.Copy.Execute([]string{}); err != nil {
		fmt.Println("copy err", err)
		// os.Exit(1)
	}
	if err := cmd.Link.Execute([]string{}); err != nil {
		fmt.Println("link err", err)
	}
	if err := cmd.Template.Execute([]string{}); err != nil {
		fmt.Println("template err", err)
	}
	return nil
}

func init() {
	// Install, err := parser.AddCommand("install",
	// 	"Install",
	// 	"",
	// 	&installCmd)
	// if err != nil {
	// 	fmt.Println(err)
	// 	os.Exit(1)
	// }
	// Install.Aliases = []string{"i"}
	// Install.SubcommandsOptional = true

	// InstallCopy, err := Install.AddCommand("copy",
	// 	"Install copy",
	// 	"",
	// 	&copyCmd)
	// if err != nil {
	// 	fmt.Println(err)
	// 	os.Exit(1)
	// }
	// InstallCopy.Aliases = []string{"rm"}
}
