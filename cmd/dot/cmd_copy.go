package main

import (
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
	RunE:              runCopy,
}

func init() {
	cmdInstall.AddCommand(cmdCopy)
	cmdRemove.AddCommand(cmdCopy)

	// flags := cmdCopy.Flags()
}

func runCopy(cmd *cobra.Command, args []string) error {
	action := cmd.Parent().Name()
	for _, r := range globalConfig.Roles {
		for _, c := range r.Files {
			if err := runTask(action, c); err != nil {
				return err
			}
		}
	}
	return nil
}
