package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

// Options for the install command.
type installOptions struct {
	sync bool

	pkg bool
}

var installOpts installOptions

var cmdInstall = &cobra.Command{
	Use:        "install [flags]",
	Aliases:    []string{"i"},
	SuggestFor: []string{"add"},
	Short:      "Install managed files",
	Long: `
The "install" command installs roles by executing their tasks.
`,
	Args:    cobra.NoArgs,
	PreRunE: preRunInstall,
	RunE:    runInstall,
	// DisableAutoGenTag: true,
}

func init() {
	cmdRoot.AddCommand(cmdInstall)

	pflags := cmdInstall.PersistentFlags()
	pflags.BoolVarP(&installOpts.sync, "sync", "S", false, "synchronize repositories")
	// pflags.BoolVarP(&installOpts.noSync, "no-sync", "", false, "ignore uncommitted changes")

	flags := cmdInstall.Flags()
	flags.BoolVarP(&installOpts.pkg, "packages", "P", false, "manage system packages")
	// flags.BoolVarP(&installOpts.pull, "pull", "p", false, "pull repositories")
}

func preRunInstall(cmd *cobra.Command, args []string) error {
	if installOpts.sync {
		// syncRoles dotConfig.Roles
		if err := runSync(cmd, []string{}); err != nil {
			return err
		}
	}
	return nil
}

// TODO: check nothing to do
func runInstall(cmd *cobra.Command, args []string) error {
	action := "install"
	for _, r := range dotConfig.Roles {
		if dotOpts.verbosity >= 1 {
			fmt.Printf("## Installing %s...\n", r.Name)
		}
		// Pre install hooks
		for _, h := range r.Install {
			h.ExecDir = r.Path
			if err := runTask(action, h); err != nil {
				return err
			}
		}
		// Package management
		if installOpts.pkg {
			for _, p := range r.Pkgs {
				if err := runTask(action, p); err != nil {
					return err
				}
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
		// Post install hooks
		for _, h := range r.PostInstall {
			h.ExecDir = r.Path
			if err := runTask(action, h); err != nil {
				return err
			}
		}
	}
	return nil
}
