package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

// CopyOptions collects all options for the copy command.
type CopyOptions struct {
}

var copyOptions CopyOptions

var cmdCopy = &cobra.Command{
	Use:     "copy [flags]",
	Aliases: []string{"cp"},
	Short:   "Copy managed files",
	Long: `
The "copy" task installs or removes files.
`,
	DisableAutoGenTag: true,
	Args:              cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runCopy(copyOptions, globalOptions, args)
	},
}

func init() {
	cmdInstall.AddCommand(cmdCopy)
	cmdRemove.AddCommand(cmdCopy)

	// flags := cmdCopy.Flags()
}

func runCopy(opts CopyOptions, gopts GlobalOptions, args []string) error {
	fmt.Println("runCopy", args)
	return nil
}
