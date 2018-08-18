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
	Use:     "dir [flags]",
	Aliases: []string{"d"},
	Short:   "Dir managed files",
	Long: `
The "dir" task installs or removes directories.
`,
	DisableAutoGenTag: true,
	Args:              cobra.NoArgs,
	RunE:              runDir,
}

func init() {
	cmdInstall.AddCommand(cmdDir)

	cmdRmDir := &cobra.Command{}
	*cmdRmDir = *cmdDir
	cmdRemove.AddCommand(cmdRmDir)

	// flags := cmdDir.Flags()
}

// opts DirOpts, gopts GlobalOptions
func runDir(cmd *cobra.Command, args []string) error {
	action := cmd.Parent().Name()
	for _, r := range globalConfig.Roles {
		for _, d := range r.Dirs {
			if err := runTask(action, d); err != nil {
				return err
			}
		}
	}
	return nil
}
