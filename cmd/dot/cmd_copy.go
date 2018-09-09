package main

import (
	"github.com/spf13/cobra"
)

// Options for the copy command.
type copyOptions struct {
}

var copyOpts copyOptions

var cmdCopy = &cobra.Command{
	Use:     "copy [flags]",
	Aliases: []string{"cp"},
	Short:   "File tasks",
	Long:    `The "copy" task installs or removes local or remote files.`,
	// Example: ``,
	Args:    cobra.NoArgs,
	PreRunE: preRunTask,
	RunE:    runCopy,
	// DisableAutoGenTag: true,
}

var cmdRmCopy = &cobra.Command{}

func init() {
	cmdInstall.AddCommand(cmdCopy)
	*cmdRmCopy = *cmdCopy
	cmdRemove.AddCommand(cmdRmCopy)

	// flags := cmdCopy.Flags()
}

func runCopy(cmd *cobra.Command, args []string) error {
	// action := cmd.Parent().Name()
	for _, r := range dotConfig.Roles {
		for _, c := range r.Files {
			if err := runTask(c); err != nil {
				return err
			}
		}
	}
	return nil
}
