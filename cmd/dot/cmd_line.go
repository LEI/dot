package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

var cmdLine = &cobra.Command{
	Use:     "line [flags]",
	Aliases: []string{"li"},
	Short:   "Line managed files",
	Long: `
The "line" task installs or removes lines.
`,
	DisableAutoGenTag: true,
	Args:              cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runLine(lineOptions, globalOptions, args)
	},
}

// LineOptions collects all options for the line command.
type LineOptions struct {
}

var lineOptions LineOptions

func init() {
	cmdInstall.AddCommand(cmdLine)
	cmdRemove.AddCommand(cmdLine)

	// flags := cmdLine.Flags()
}

func runLine(opts LineOptions, gopts GlobalOptions, args []string) error {
	fmt.Println("runLine", args)
	return nil
}
