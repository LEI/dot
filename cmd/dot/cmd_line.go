package main

import (
	"github.com/spf13/cobra"
)

// Options for the line command.
type lineOptions struct {
}

var lineOpts lineOptions

var cmdLine = &cobra.Command{
	Use:     "line [flags]",
	Aliases: []string{"li"},
	Short:   "Line managed files",
	Long: `
The "line" task installs or removes lines.
`,
	DisableAutoGenTag: true,
	Args:              cobra.NoArgs,
	PreRunE:           preRunTask,
	RunE:              runLine,
}

var cmdRmLine = &cobra.Command{}

func init() {
	cmdInstall.AddCommand(cmdLine)
	*cmdRmLine = *cmdLine
	cmdRemove.AddCommand(cmdRmLine)

	// flags := cmdLine.Flags()
}

func runLine(cmd *cobra.Command, args []string) error {
	action := cmd.Parent().Name()
	for _, r := range globalConfig.Roles {
		for _, l := range r.Lines {
			if err := runTask(action, l); err != nil {
				return err
			}
		}
	}
	return nil
}