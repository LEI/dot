package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

// Options for the remove command.
type removeOptions struct {
	// empty bool
	pkg bool
}

var removeOpts removeOptions

var cmdRemove = &cobra.Command{
	Use:        "remove [flags]",
	Aliases:    []string{"rm"},
	SuggestFor: []string{"delete", "uninstall"},
	Short:      "Remove managed files",
	Long:       `The "remove" command removes roles by executing their tasks.`,
	Args:       cobra.NoArgs,
	PreRunE:    preRunRemove,
	RunE:       runRemove,
	// DisableAutoGenTag: true,
}

func init() {
	cmdRoot.AddCommand(cmdRemove)

	flags := cmdRemove.Flags()
	// flags.BoolVarP(&removeOptions.empty, "remove-empty", "", false, "remove empty directories and empty files")
	flags.BoolVarP(&removeOpts.pkg, "packages", "P", false, "manage system packages")

	addActionFlags(cmdRemove)
}

func preRunRemove(cmd *cobra.Command, args []string) error {
	return nil
}

func runRemove(cmd *cobra.Command, args []string) error {
	action := "remove"
	for _, r := range dotConfig.Roles {
		if dotOpts.verbosity >= 1 {
			fmt.Printf("## Removing %s...\n", r.Name)
		}
		// Pre remove hooks
		for _, h := range r.Remove {
			h.ExecDir = r.Path
			if err := runTask(action, h); err != nil {
				return err
			}
		}
		// Package management
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
		if removeOpts.pkg {
			for _, p := range r.Pkgs {
				if err := runTask(action, p); err != nil {
					return err
				}
			}
		}
		// Post remove hooks
		for _, h := range r.PostRemove {
			h.ExecDir = r.Path
			if err := runTask(action, h); err != nil {
				return err
			}
		}
	}
	return nil
}
