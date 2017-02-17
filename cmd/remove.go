package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(removeCmd)
}

var removeCmd = &cobra.Command{
	// Hidden: true,
	Use: "remove [flags]",
	// Aliases: []string{"r"},
	// Short: "-R",
	// Long:  `--remove`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("cmd: %+v\n", cmd)
		fmt.Printf("args: %+v\n", args)
		fatal("Not ready yet")
	},
}
