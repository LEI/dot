package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

// Options for the remove command.
type removeOptions struct {
	// empty bool
}

var removeOpts removeOptions

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
		// Package management
		for _, p := range r.Pkgs {
			if err := runTask(action, p); err != nil {
				return err
			}
		}
	}
	return nil
}
