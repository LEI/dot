package main

import (
	"github.com/spf13/cobra"
)

// Options for the template command.
type templateOptions struct {
}

var templateOpts templateOptions

var cmdTemplate = &cobra.Command{
	Use:     "template [flags]",
	Aliases: []string{"tpl"},
	Short:   "Template managed files",
	Long: `
The "template" task installs or removes templates.
`,
	DisableAutoGenTag: true,
	Args:              cobra.NoArgs,
	RunE:              runTemplate,
}

func init() {
	cmdInstall.AddCommand(cmdTemplate)
	cmdRemove.AddCommand(cmdTemplate)

	// flags := cmdTemplate.Flags()
}

func runTemplate(cmd *cobra.Command, args []string) error {
	action := cmd.Parent().Name()
	for _, r := range globalConfig.Roles {
		for _, t := range r.Tpls {
			if err := runTask(action, t); err != nil {
				return err
			}
		}
	}
	return nil
}
