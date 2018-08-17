package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

// RemoveOptions collects all options for the remove command.
type RemoveOptions struct {
	// Empty bool
}

var removeOptions RemoveOptions

var cmdRemove = &cobra.Command{
	Use:     "remove [flags]",
	Aliases: []string{"rm"},
	Short:   "Remove managed files",
	Long: `
The "remove" command removes roles by executing their tasks.
`,
	DisableAutoGenTag: true,
	Args:              cobra.NoArgs,
	RunE:              runRemove,
}

func init() {
	cmdRoot.AddCommand(cmdRemove)

	// flags := cmdRemove.Flags()
	// flags.BoolVarP(&removeOptions.Empty, "remove-empty", "", false, "remove empty directories and empty files")
}

func runRemove(cmd *cobra.Command, args []string) error {
	action := "remove"
	for _, r := range globalConfig.Roles {
		fmt.Printf("## Removing %s...\n", r.Name)
		for _, c := range r.Files {
			if err := runTask(action, c); err != nil {
				return err
			}
		}
		for _, l := range r.Links {
			if err := runTask(action, l); err != nil {
				return err
			}
		}
		for _, t := range r.Tpls {
			if err := runTask(action, t); err != nil {
				return err
			}
		}
		for _, l := range r.Lines {
			if err := runTask(action, l); err != nil {
				return err
			}
		}
		// Remove directories last
		for _, d := range r.Dirs {
			if err := runTask(action, d); err != nil {
				return err
			}
		}
	}
	return nil
}
