package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

// Options for the install command.
type installOptions struct {
	sync bool
}

var installOpts installOptions

var cmdInstall = &cobra.Command{
	Use:     "install [flags]",
	Aliases: []string{"i"},
	Short:   "Install managed files",
	Long: `
The "install" command installs roles by executing their tasks.
`,
	DisableAutoGenTag: true,
	Args:              cobra.NoArgs,
	PreRunE:           preRunInstall,
	RunE:              runInstall,
}

func init() {
	cmdRoot.AddCommand(cmdInstall)

	f := cmdInstall.PersistentFlags()
	f.BoolVarP(&installOpts.sync, "sync", "S", false, "synchronize repositories")

	// flags := cmdInstall.Flags()
	// flags.BoolVarP(&installOpts.pull, "pull", "p", false, "pull repositories")
	// flags.BoolVarP(&installOpts.force, "force", "f", false, "ignore uncommitted changes")
}

func preRunInstall(cmd *cobra.Command, args []string) error {
	if installOpts.sync {
		for _, r := range globalConfig.Roles {
			// Clone or pull git repository
			if err := r.Sync(); err != nil {
				return err
			}
			// Parse config file again
			if err := r.LoadConfig(); err != nil {
				return err
			}
		}
	}
	return nil
}

// TODO: check nothing to do
func runInstall(cmd *cobra.Command, args []string) error {
	action := "install"
	for _, r := range globalConfig.Roles {
		fmt.Printf("## Installing %s...\n", r.Name)
		// Package management
		for _, p := range r.Pkgs {
			if err := runTask(action, p); err != nil {
				return err
			}
		}
		// Remove directories first
		for _, d := range r.Dirs {
			if err := runTask(action, d); err != nil {
				return err
			}
		}
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
	}
	return nil
}
