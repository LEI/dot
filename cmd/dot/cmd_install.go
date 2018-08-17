package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

// InstallOptions collects all options for the install command.
type InstallOptions struct {
	Sync bool
}

var installOptions InstallOptions

var cmdInstall = &cobra.Command{
	Use:     "install [flags]",
	Aliases: []string{"i"},
	Short:   "Install managed files",
	Long: `
The "install" command installs roles by executing their tasks.
`,
	DisableAutoGenTag: true,
	Args:              cobra.NoArgs,
	RunE:              runInstall,
}

func init() {
	cmdRoot.AddCommand(cmdInstall)

	flags := cmdInstall.Flags()
	flags.BoolVarP(&installOptions.Sync, "sync", "S", false, "synchronize repositories")
	// flags.BoolVarP(&installOptions.Pull, "pull", "p", false, "pull repositories")
	// flags.BoolVarP(&installOptions.Force, "force", "f", false, "ignore uncommitted changes")
}

func runInstall(cmd *cobra.Command, args []string) error {
	// if installOptions.Sync {
	// 	for _, r := range globalConfig.Roles {
	// 		if err := r.RunSync(); err != nil {
	// 			return err
	// 		}
	// 	}
	// }
	action := "install"
	for _, r := range globalConfig.Roles {
		fmt.Printf("## Installing %s...\n", r.Name)
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
		for _, t := range r.Templates {
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
