package main

import (
	"github.com/spf13/cobra"
)

// LineOptions collects all options for the line command.
type LineOptions struct {
}

var lineOptions LineOptions

var cmdLine = &cobra.Command{
	Use:     "line [flags]",
	Aliases: []string{"li"},
	Short:   "Line managed files",
	Long: `
The "line" task installs or removes lines.
`,
	DisableAutoGenTag: true,
	Args:              cobra.NoArgs,
	RunE:              runLine,
}

func init() {
	cmdInstall.AddCommand(cmdLine)
	cmdRemove.AddCommand(cmdLine)

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
