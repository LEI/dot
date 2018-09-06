package main

import (
	"github.com/spf13/cobra"
)

// Options for the dir command.
type dirOptions struct {
	// empty bool
}

var dirOpts dirOptions

var cmdDir = &cobra.Command{
	Use:     "directory [flags]",
	Aliases: []string{"d", "dir"},
	Short:   "Directory tasks",
	Long:    `The "dir" task installs or removes directories.`,
	// Example: ``,
	Args:    cobra.NoArgs,
	PreRunE: preRunTask,
	RunE:    runDir,
	// DisableAutoGenTag: true,
}

var cmdRmDir = *cmdDir

func init() {
	cmdInstall.AddCommand(cmdDir)
	cmdRemove.AddCommand(&cmdRmDir)

	// flags := cmdDir.Flags()
}

// opts DirOpts, gopts GlobalOptions
func runDir(cmd *cobra.Command, args []string) error {
	// action := cmd.Parent().Name()
	for _, r := range dotConfig.Roles {
		for _, d := range r.Dirs {
			if err := runTask(d); err != nil {
				return err
			}
		}
	}
	return nil
}
