package main

import (
	"github.com/spf13/cobra"
)

// Options for the template command.
type templateOptions struct {
}

var tplOpts templateOptions

var cmdTpl = &cobra.Command{
	Use:     "template [flags]",
	Aliases: []string{"tpl"},
	Short:   "Template tasks",
	Long:    `The "template" task installs or removes templates.`,
	// Example: ``,
	Args:    cobra.NoArgs,
	PreRunE: preRunTask,
	RunE:    runTpl,
	// DisableAutoGenTag: true,
}

var cmdRmTpl = &cobra.Command{}

func init() {
	cmdInstall.AddCommand(cmdTpl)
	*cmdRmTpl = *cmdTpl
	cmdRemove.AddCommand(cmdRmTpl)

	// flags := cmdTpl.Flags()
}

func runTpl(cmd *cobra.Command, args []string) error {
	// action := cmd.Parent().Name()
	for _, r := range dotConfig.Roles {
		for _, t := range r.Tpls {
			if err := runTask(t); err != nil {
				return err
			}
		}
	}
	return nil
}
