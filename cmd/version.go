package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var Version = "0.0.1"

var versionCmd = &cobra.Command{
	Hidden:  true,
	Use:     "version",
	Aliases: []string{"v"},
	Short:   "Print the version number",
	// Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("%s %s\n", os.Args[0], Version)
	},
}

func init() {
	// DotCmd.AddCommand(versionCmd)
}
